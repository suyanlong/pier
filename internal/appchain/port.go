package appchain

import (
	"github.com/link33/sidecar/internal/port"
	"github.com/link33/sidecar/model/pb"
)

func (a *appChain) ID() string {
	return a.client.ID()
}

func (a *appChain) Type() string {
	return a.client.Type()
}

func (a *appChain) Tag() string {
	return a.client.Type()
}

func (a *appChain) Name() string {
	return a.client.Name()
}

func (a *appChain) Send(msg port.Message) (*pb.Message, error) {
	panic("implement me")
}

func (a *appChain) AsyncSend(msg port.Message) error {
	panic("implement me")
}

func (a *appChain) ListenIBTPX() <-chan *pb.IBTPX {
	panic("implement me")
}
