package peermgr

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/link33/sidercar/internal/port"
	"github.com/link33/sidercar/model/pb"
)

type MessageHandler func(port.Port, *pb.Message)
type ConnectHandler func(string)

//go:generate mockgen -destination mock_peermgr/mock_peermgr.go -package mock_peermgr -source peermgr.go
type PeerManager interface {
	DHTManager

	// Start
	Start() error

	// Stop
	Stop() error

	Connect(info *peer.AddrInfo) (string, error)

	// AsyncSend sends message to peer with peer info.
	AsyncSend(string, port.Message) error
	// Send sends message waiting response
	Send(string, port.Message) (*pb.Message, error)

	AsyncSendWithPort(port.Port, port.Message) error

	SendWithPort(s port.Port, msg port.Message) (*pb.Message, error)

	Handler
}

type Handler interface {
	// RegisterMsgHandler
	RegisterMsgHandler(pb.Message_Type, MessageHandler) error

	// RegisterMultiMsgHandler
	RegisterMultiMsgHandler([]pb.Message_Type, MessageHandler) error

	// RegisterConnectHandler
	RegisterConnectHandler(ConnectHandler) error
}

type DHTManager interface {
	// Search for peers who are able to provide a given key
	FindProviders(id string) (string, error)

	// Provide adds the given cid to the content routing system. If 'true' is
	// passed, it also announces it, otherwise it is just kept in the local
	// accounting of which objects are being provided.
	Provider(string, bool) error
}
