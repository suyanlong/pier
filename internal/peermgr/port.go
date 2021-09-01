package peermgr

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/link33/sidecar/internal/port"
	"github.com/link33/sidecar/model/pb"
)

type sidecar struct {
	addr  *peer.AddrInfo
	swarm PeerManager
	tag   string
	rev   chan *pb.IBTPX
}

func (s *sidecar) ID() string {
	return s.addr.String()
}

func (s *sidecar) Type() string {
	return port.Sidecar
}

func (s *sidecar) Name() string {
	return s.ID()
}

func (s *sidecar) Tag() string {
	return s.tag
}

func (s *sidecar) Send(msg port.Message) (*pb.Message, error) {
	return s.swarm.Send(s.ID(), msg)
}

func (s *sidecar) AsyncSend(msg port.Message) error {
	return s.swarm.AsyncSend(s.ID(), msg)
}

func (s *sidecar) ListenIBTPX() <-chan *pb.IBTPX {
	return s.rev
}
