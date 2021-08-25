package peermgr

import (
	"github.com/libp2p/go-libp2p-core/peer"
	peermgr "github.com/meshplus/pier/internal/peermgr/proto"
	"github.com/meshplus/pier/internal/port"
)

type MessageHandler func(port.Port, *peermgr.Message)
type ConnectHandler func(string)

//go:generate mockgen -destination mock_peermgr/mock_peermgr.go -package mock_peermgr -source peermgr.go
type PeerManager interface {
	DHTManager

	// Start
	Start() error

	// Stop
	Stop() error

	// AsyncSend sends message to peer with peer info.
	AsyncSend(string, port.Message) error

	Connect(info *peer.AddrInfo) (string, error)

	// AsyncSendWithStream sends message using existed stream
	AsyncSendWithPort(port.Port, port.Message) error

	// SendWithStream sends message using existed stream
	//SendWithStreamX(port.Port, port.Message) (port.Message, error)
	//
	//// AsyncSendWithPort sends message using existed stream
	//AsyncSendWithStreamX(port.Port,, port.Message) error

	// Send sends message waiting response
	Send(string, port.Message) (*peermgr.Message, error)

	// Peers
	Peers() map[string]*peer.AddrInfo

	// RegisterMsgHandler
	RegisterMsgHandler(peermgr.Message_Type, MessageHandler) error

	// RegisterMultiMsgHandler
	RegisterMultiMsgHandler([]peermgr.Message_Type, MessageHandler) error

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
