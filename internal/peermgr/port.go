package peermgr

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/link33/sidercar/internal/port"
	"github.com/link33/sidercar/model/pb"
)

type sidercar struct {
	addr  *peer.AddrInfo
	swarm PeerManager
	tag   string
	rev   chan *pb.IBTPX
}

func (s *sidercar) ID() string {
	return s.addr.String()
}

func (s *sidercar) Type() string {
	return port.Sidercar
}

func (s *sidercar) Name() string {
	return s.ID()
}

func (s *sidercar) Tag() string {
	return s.tag
}

func (s *sidercar) Send(msg port.Message) (*pb.Message, error) {
	return s.swarm.Send(s.ID(), msg)
}

func (s *sidercar) AsyncSend(msg port.Message) error {
	return s.swarm.AsyncSend(s.ID(), msg)
}

func (s *sidercar) ListenIBTPX() <-chan *pb.IBTPX {
	return s.rev
}
