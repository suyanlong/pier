package plugins

import (
	"github.com/link33/sidercar/internal"
	"github.com/link33/sidercar/model/pb"
	//"github.com/link33/sidercar/internal/port"
)

type Kernel interface {
	Kill()
	Exited() bool
}

// Client defines the interface that interacts with appchain 与 appchain交互。
//go:generate mockgen -destination mock_client/mock_client.go -package mock_client -source interface.go
type Client interface { //业务实现委托接口。需要实现的那有些。
	internal.Launcher
	Kernel
	Bind(kern Kernel)
	// Initialize initialize plugin client
	Initialize(configPath string, ID string, extra []byte) error

	// GetIBTP gets an interchain ibtp channel generated by client
	GetIBTP() chan *pb.IBTP

	// SubmitIBTP submits the interchain ibtp to appchain
	SubmitIBTP(*pb.IBTP) (*pb.SubmitIBTPResponse, error)

	// RollbackIBTP rollbacks the interchain ibtp to appchain
	RollbackIBTP(ibtp *pb.IBTP, isSrcChain bool) (*pb.RollbackIBTPResponse, error)

	// increase in meta without actually executing it
	IncreaseInMeta(ibtp *pb.IBTP) (*pb.IBTP, error)

	// GetOutMessage gets interchain ibtp by index and target chain_id from broker contract
	GetOutMessage(to string, idx uint64) (*pb.IBTP, error)

	// GetInMessage gets receipt by index and source chain_id
	GetInMessage(from string, idx uint64) ([][]byte, error)

	// GetOutMeta gets an index map, which implicates the greatest index of
	// ingoing interchain txs for each source chain
	GetInMeta() (map[string]uint64, error)

	// GetOutMeta gets an index map, which implicates the greatest index of
	// outgoing interchain txs for each receiving chain
	GetOutMeta() (map[string]uint64, error)

	// GetReceiptMeta gets an index map, which implicates the greatest index of
	// executed callback txs for each receiving chain
	GetCallbackMeta() (map[string]uint64, error)

	// CommitCallback is a callback function when get receipt from bitxhub success
	CommitCallback(ibtp *pb.IBTP) error

	// GetReceipt gets receipt of an executed IBTP
	GetReceipt(ibtp *pb.IBTP) (*pb.IBTP, error)

	// Name gets name of blockchain from plugin
	Name() string

	// Type gets type of blockchain from plugin
	Type() string

	// ID
	ID() string
}

//type ClientX interface {
//	Client
//
//	ID() string
//	Send(ibtp *pb.IBTP) error
//}
