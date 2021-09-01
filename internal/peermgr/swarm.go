package peermgr

import (
	"context"
	"fmt"
	"github.com/link33/sidecar/internal/port"
	"github.com/link33/sidecar/model/pb"
	"strings"
	"sync"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/ipfs/go-cid"
	crypto2 "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/link33/sidecar/internal/repo"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym/ecdsa"
	network "github.com/meshplus/go-lightp2p"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

const (
	protocolID          = "/sidecar/1.0.0" // magic protocol
	defaultProvidersNum = 1
)

var _ PeerManager = (*Swarm)(nil)

type Swarm struct {
	p2p            network.Network
	logger         logrus.FieldLogger
	peers          map[string]*peer.AddrInfo
	connectedPeers *port.PortMap

	providers       uint64
	privKey         crypto.PrivateKey
	msgHandlers     sync.Map
	connectHandlers []ConnectHandler

	lock   sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func New(config *repo.Config, portMap *port.PortMap, nodePrivKey crypto.PrivateKey, privKey crypto.PrivateKey, providers uint64, logger logrus.FieldLogger) (*Swarm, error) {
	libp2pPrivKey, err := convertToLibp2pPrivKey(nodePrivKey)
	if err != nil {
		return nil, fmt.Errorf("convert private key: %w", err)
	}
	var local string
	var remotes map[string]*peer.AddrInfo

	local, remotes, err = loadPeers(config.Peer.Peers, libp2pPrivKey)

	var protocolIDs = []string{protocolID}

	p2p, err := network.New(
		network.WithLocalAddr(local),
		network.WithPrivateKey(libp2pPrivKey),
		network.WithProtocolIDs(protocolIDs),
		network.WithLogger(logger),
	)
	if err != nil {
		return nil, fmt.Errorf("create p2p: %w", err)
	}

	if providers == 0 {
		providers = defaultProvidersNum
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Swarm{
		providers: providers,
		p2p:       p2p,
		logger:    logger,
		peers:     remotes,
		privKey:   privKey,
		ctx:       ctx,
		cancel:    cancel,
	}, nil
}

func (swarm *Swarm) Start() error {
	swarm.p2p.SetMessageHandler(swarm.handleMessage)

	if err := swarm.RegisterMsgHandler(pb.Message_ADDRESS_GET, swarm.handleGetAddressMessage); err != nil {
		return fmt.Errorf("register get address msg handler: %w", err)
	}

	if err := swarm.p2p.Start(); err != nil {
		return fmt.Errorf("p2p module start: %w", err)
	}

	//need to connect one other sidecar at least
	wg := &sync.WaitGroup{}
	wg.Add(1)

	for id, addr := range swarm.peers {
		go func(id string, addr *peer.AddrInfo) {
			if err := retry.Retry(func(attempt uint) error {
				if err := swarm.p2p.Connect(*addr); err != nil {
					if attempt != 0 && attempt%5 == 0 {
						swarm.logger.WithFields(logrus.Fields{
							"node":  id,
							"error": err,
						}).Error("Connect failed")
					}
					return err
				}

				address, err := swarm.getRemoteAddress(addr.ID)
				if err != nil {
					swarm.logger.WithFields(logrus.Fields{
						"node":  id,
						"error": err,
					}).Error("Get remote address failed")
					return err
				}

				swarm.logger.WithFields(logrus.Fields{
					"node":     id,
					"address:": address,
				}).Info("Connect successfully")

				rec := make(chan *pb.IBTPX)
				p := &sidecar{
					addr:  addr,
					swarm: swarm,
					tag:   "",
					rev:   rec,
				}
				swarm.connectedPeers.Store(id, p)
				swarm.lock.RLock()
				defer swarm.lock.RUnlock()
				for _, handler := range swarm.connectHandlers {
					go func(connectHandler ConnectHandler, address string) {
						connectHandler(address)
					}(handler, address)
				}
				wg.Done()
				return nil
			},
				strategy.Wait(1*time.Second),
			); err != nil {
				swarm.logger.Error(err)
			}
		}(id, addr)
	}

	wg.Wait()

	return nil
}

//注册异步处理数据的方法
func (swarm *Swarm) handleMessage(s network.Stream, data []byte) {
	m := &pb.Message{}
	if err := m.Unmarshal(data); err != nil {
		swarm.logger.Error(err)
		return
	}

	pack := m.Payload.Data
	t := m.Type
	switch {
	case t == pb.Message_IBTP_SEND || t == pb.Message_IBTP_GET || t == pb.Message_IBTP_RECEIPT_SEND:
		ibtpx := &pb.IBTPX{}
		if err := m.Unmarshal(pack); err != nil {
			swarm.logger.Error(err)
			return
		}
		p, is := swarm.connectedPeers.Load(s.RemotePeerID())
		if is {
			ps, iss := p.(*sidecar)
			if iss {
				ps.rev <- ibtpx
				return
			}
		}
		addr, _ := peer.AddrInfoFromP2pAddr(s.RemotePeerAddr())
		rec := make(chan *pb.IBTPX)
		newPort := &sidecar{
			addr:  addr,
			swarm: swarm,
			tag:   "",
			rev:   rec,
		}
		swarm.connectedPeers.Store(addr.ID.String(), newPort)
		rec <- ibtpx
	default:

	}
}

func (swarm *Swarm) Stop() error {
	if err := swarm.p2p.Stop(); err != nil {
		return err
	}
	swarm.cancel()
	return nil
}

func (swarm *Swarm) Connect(addrInfo *peer.AddrInfo) (string, error) {
	err := swarm.p2p.Connect(*addrInfo)
	if err != nil {
		return "", err
	}
	address, err := swarm.getRemoteAddress(addrInfo.ID)
	if err != nil {
		return "", err
	}
	swarm.logger.WithFields(logrus.Fields{
		"address":  address,
		"addrInfo": addrInfo,
	}).Info("Connect peer")

	rec := make(chan *pb.IBTPX)
	p := &sidecar{
		addr:  addrInfo,
		swarm: swarm,
		tag:   "",
		rev:   rec,
	}
	swarm.connectedPeers.Add(p)
	return addrInfo.ID.String(), nil
}

func (swarm *Swarm) AsyncSendWithPort(s port.Port, msg port.Message) error {
	return s.AsyncSend(msg)
}

func (swarm *Swarm) SendWithPort(s port.Port, msg port.Message) (*pb.Message, error) {
	return s.Send(msg)
}

func (swarm *Swarm) AsyncSend(id string, msg port.Message) error {
	data, err := msg.Marshal()
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	return swarm.p2p.AsyncSend(id, data)
}

func (swarm *Swarm) Send(id string, msg port.Message) (*pb.Message, error) {
	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	ret, err := swarm.p2p.Send(id, data)
	if err != nil {
		return nil, fmt.Errorf("sync send: %w", err)
	}

	m := &pb.Message{}
	if err := m.Unmarshal(ret); err != nil {
		return nil, err
	}

	return m, nil
}

func (swarm *Swarm) Peers() map[string]*peer.AddrInfo {
	m := make(map[string]*peer.AddrInfo)
	for id, addr := range swarm.peers {
		m[id] = addr
	}
	return m
}

func convertToLibp2pPrivKey(privateKey crypto.PrivateKey) (crypto2.PrivKey, error) {
	ecdsaPrivKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("convert to libp2p private key: not ecdsa private key")
	}

	libp2pPrivKey, _, err := crypto2.ECDSAKeyPairFromKey(ecdsaPrivKey.K)
	if err != nil {
		return nil, err
	}

	return libp2pPrivKey, nil
}

func loadPeers(peers []string, privateKey crypto2.PrivKey) (string, map[string]*peer.AddrInfo, error) {
	var local string
	remotes := make(map[string]*peer.AddrInfo)

	id, err := peer.IDFromPrivateKey(privateKey)
	if err != nil {
		return "", nil, err
	}

	for _, p := range peers {
		if strings.HasSuffix(p, id.String()) {
			idx := strings.LastIndex(p, "/p2p/")
			if idx == -1 {
				return "", nil, fmt.Errorf("pid is not existed in bootstrap")
			}

			local = p[:idx]
		} else {
			addr, err := AddrToPeerInfo(p)
			if err != nil {
				return "", nil, fmt.Errorf("wrong network addr: %w", err)
			}
			remotes[addr.ID.String()] = addr
		}
	}

	if local == "" {
		return "", nil, fmt.Errorf("get local addr: no local addr is configured")
	}

	return local, remotes, nil
}

// AddrToPeerInfo transfer addr to PeerInfo
// addr example: "/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64"
func AddrToPeerInfo(multiAddr string) (*peer.AddrInfo, error) {
	maddr, err := ma.NewMultiaddr(multiAddr)
	if err != nil {
		return nil, err
	}

	return peer.AddrInfoFromP2pAddr(maddr)
}

func (swarm *Swarm) getRemoteAddress(id peer.ID) (string, error) {
	msg := Message(pb.Message_ADDRESS_GET, true, nil)
	reqData, err := msg.Marshal()
	if err != nil {
		return "", err
	}
	retData, err := swarm.p2p.Send(id.String(), reqData) //同步获取数据
	if err != nil {
		return "", fmt.Errorf("sync send: %w", err)
	}
	ret := &pb.Message{}
	if err := ret.Unmarshal(retData); err != nil {
		return "", err
	}

	return string(ret.Payload.Data), nil
}

func (swarm *Swarm) FindProviders(id string) (string, error) {
	format := cid.V0Builder{}
	toCid, err := format.Sum([]byte(id))
	if err != nil {
		return "", err
	}
	providers, err := swarm.p2p.FindProvidersAsync(toCid.String(), int(swarm.providers))
	if err != nil {
		swarm.logger.WithFields(logrus.Fields{
			"id": id,
		}).Error("Not find providers")
		return "", err
	}

	for provider := range providers {
		swarm.logger.WithFields(logrus.Fields{
			"id":          id,
			"cid":         toCid.String(),
			"provider_id": provider.ID.String(),
		}).Info("Find provider")

		sidecarId, err := swarm.Connect(&provider)
		if err != nil {
			swarm.logger.WithFields(logrus.Fields{"peerId": sidecarId,
				"cid": provider.ID.String()}).Error("connect error: ", err)
			continue
		}
		return sidecarId, nil
	}

	swarm.logger.WithFields(logrus.Fields{
		"id":  id,
		"cid": toCid.String(),
	}).Warning("No providers found") // TODO add error
	return "", nil
}

func (swarm *Swarm) Provider(key string, passed bool) error {
	return swarm.p2p.Provider(key, passed)
}

func (swarm *Swarm) handleGetAddressMessage(p port.Port, message *pb.Message) {
	addr, err := swarm.privKey.PublicKey().Address()
	if err != nil {
		swarm.logger.Error(err)
		return
	}
	retMsg := Message(pb.Message_ACK, true, []byte(addr.String()))
	err = swarm.AsyncSendWithPort(p, retMsg)
	if err != nil {
		swarm.logger.Error(err)
	}
}
