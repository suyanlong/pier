package manger

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/link33/sidercar/internal"
	"github.com/link33/sidercar/internal/port"
	"github.com/link33/sidercar/internal/repo"
	"github.com/link33/sidercar/internal/router"
	"github.com/link33/sidercar/model/pb"

	"github.com/link33/sidercar/internal/peermgr"
	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	logger logrus.FieldLogger
	Mgr    appchainmgr.AppchainMgr
}

func NewManager(addr string, logger logrus.FieldLogger) (*Manager, error) {
	appchainMgr := appchainmgr.New(&Persister{addr: addr, logger: logger})
	am := &Manager{
		Mgr:    appchainMgr,
		logger: logger,
	}

	err := pm.RegisterMultiMsgHandler([]pb.Message_Type{
		pb.Message_APPCHAIN_REGISTER,
		pb.Message_APPCHAIN_UPDATE,
		pb.Message_APPCHAIN_GET,
	}, am.handleMessage)
	if err != nil {
		return nil, err
	}

	return am, nil
}

func (mgr *Manager) handleMessage(s port.Port, msg *pb.Message) {
	var res []byte
	var ok bool
	switch msg.Type {
	case pb.Message_APPCHAIN_REGISTER:
		ok, res = mgr.Mgr.Register(msg.Payload.Data)
	case pb.Message_APPCHAIN_UPDATE:
		ok, res = mgr.Mgr.Update(msg.Payload.Data)
	case pb.Message_APPCHAIN_GET:
		app := &appchainmgr.Appchain{}
		if err := json.Unmarshal(msg.Payload.Data, app); err != nil {
			mgr.logger.Error(err)
			return
		}
		ok, res = mgr.Mgr.QueryById(app.ID, nil)
	default:
		m := "wrong appchain message type"
		res = []byte(m)
		mgr.logger.Error(m)
	}

	ackMsg := peermgr.Message(msg.Type, ok, res)
	err := s.AsyncSend(ackMsg)
	if err != nil {
		mgr.logger.Error(err)
	}

	appchainRes := &appchainmgr.Appchain{}
	if err := json.Unmarshal(res, appchainRes); err != nil {
		mgr.logger.Error(err)
		return
	}

	mgr.logger.WithFields(logrus.Fields{
		"type":           msg.Type,
		"from_id":        appchainRes.ID,
		"name":           appchainRes.Name,
		"desc":           appchainRes.Desc,
		"chain_type":     appchainRes.ChainType,
		"consensus_type": appchainRes.ConsensusType,
	}).Info("Handle appchain message")
}

type Manger interface {
	internal.Launcher
	Remove()
	Add()
	Query()
}

type MangerPort struct {
	// peer manger
	// sidercar manger
	logger  logrus.FieldLogger
	ctx     context.Context
	cancel  context.CancelFunc
	router  router.Router
	portMap *port.PortMap

	methodMap map[string]routeMethod
}

type routeMethod func([]string) []port.Port

func NewMangerPort(logger logrus.FieldLogger) *MangerPort {
	ctx, cancel := context.WithCancel(context.Background())
	return &MangerPort{
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
		router:  nil,
		portMap: port.NewPortMap(),
	}
}

func (m *MangerPort) Start() error {
	if err := m.router.Start(); err != nil {
		return err
	}
	m.methodMap["single"] = m.Single
	m.methodMap["multicast"] = m.Multicast
	m.methodMap["broadcast"] = m.Broadcast
	m.methodMap["official"] = m.Official
	return nil
}

func (m *MangerPort) Stop() error {
	if err := m.router.Stop(); err != nil {
		return err
	}
	return nil
}

//TODO
func (m *MangerPort) Add(p port.Port) error {
	m.portMap.Add(p)
	go func() {
		c := p.ListenIBTPX()
		for {
			select {
			case ibtpx := <-c:
				err := m.Route(ibtpx)
				if err != nil {
					m.logger.Error(err)
				}
			}
		}
	}()

	return nil
}

func (m *MangerPort) Remove(p port.Port) error {
	m.portMap.Remove(p)
	return nil
}

//TODO 本机找到的appchain是自己的appchain
func (m *MangerPort) Route(ibtpx *pb.IBTPX) error {
	mode := ibtpx.Mode
	//本网关已签名、中继链已背书、to是本网关内部的appchain，即顺利通过并转发，否则打断。
	if !((m.isSign(ibtpx) && mode == repo.RelayMode && m.isEndorse(ibtpx)) || !m.isSign(ibtpx)) {
		return nil
	}
	//本网关签名
	if !m.isSign(ibtpx) {
		m.Sign(ibtpx)
	}
	ibtp := ibtpx.Ibtp
	_, to := ibtp.From, ibtp.To
	if pp, is := m.portMap.Port(to); is {
		switch {
		case pp.Type() == port.Hub || pp.Type() == port.Sidercar:
			return pp.AsyncSend(ibtpx)
		case pp.Type() == port.Appchain:
			switch mode {
			case repo.RelayMode:
				hub, is := m.getHub()
				if is && !m.isEndorse(ibtpx) {
					return hub.AsyncSend(ibtpx)
				} else {
					return pp.AsyncSend(ibtpx)
				}
			case repo.DirectMode:
				return pp.AsyncSend(ibtpx)
			default:
				//TODO  跳过
				return nil
			}
		default:
			//TODO 跳过
			return nil
		}
	}
	//规则判断
	method := ibtpx.RouteMethod
	if md, is := m.methodMap[method]; is {
		ports := md(ibtpx.RouteMethodArg)
		if len(ports) == 0 {
			m.firstRoute(ibtpx)
		}
		for _, p := range ports {
			_ := p.AsyncSend(ibtpx)
		}
	}
	m.firstRoute(ibtpx)
	return nil
}

func (m *MangerPort) firstRoute(ibtp *pb.IBTPX) {
	panic("implement me")
}

func (m *MangerPort) getHub() (port.Port, bool) {

	return nil, false
}

func (m *MangerPort) isSign(ibtpx *pb.IBTPX) bool {
	panic("implement me")
}

func (m *MangerPort) Sign(ibtpx *pb.IBTPX) {
	panic("implement me")
}

// hub endorse
func (m *MangerPort) isEndorse(ibtpx *pb.IBTPX) bool {
	panic("implement me")
}

func (m *MangerPort) HandlerMethod() {}

func (m *MangerPort) Single([]string) []port.Port {
	panic("implement me")
}

func (m *MangerPort) Multicast([]string) []port.Port {

	panic("implement me")
}

func (m *MangerPort) Broadcast([]string) []port.Port {
	panic("implement me")
}

func (m *MangerPort) Official([]string) []port.Port {
	panic("implement me")
}

func (m *MangerPort) Send(id string, msg port.Message) (*pb.Message, error) {
	if p, is := m.portMap.Port(id); is {
		return p.Send(msg)
	}

	return nil, errors.New("id error!")
}

func (m *MangerPort) AsyncSend(id string, msg port.Message) error {
	if p, is := m.portMap.Port(id); is {
		return p.AsyncSend(msg)
	}

	return errors.New("id error!")
}
