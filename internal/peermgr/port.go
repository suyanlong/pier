package peermgr

import (
	"github.com/link33/sidercar/internal/port"
	"github.com/link33/sidercar/model/pb"
)

type sidercar struct {
	peerIDs string
	swarm   *Swarm
	tag     string
	rev     chan *pb.Message
}

func (s *sidercar) ID() string {
	return s.peerIDs
}

func (s *sidercar) Type() string {
	return "sidercar"
}

func (s *sidercar) Name() string {
	return s.peerIDs
}

func (s *sidercar) Tag() string {
	return s.tag
}

func (s *sidercar) Send(msg port.Message) (*pb.Message, error) {
	panic("implement me")
	//s.swarm.Send(s.peerIDs,)
}

func (s *sidercar) AsyncSend(msg port.Message) error {
	panic("implement me")
	//s.swarm.Send(s.peerIDs,)
}

func (s *sidercar) ListenIBTPX() <-chan *pb.IBTPX {
	return s.rev
}
