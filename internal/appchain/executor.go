package appchain

import (
	"github.com/link33/sidecar/model/pb"
)

//go:generate mockgen -destination mock_executor/mock_executor.go -package mock_executor -source interface.go
type Executor interface {
	// ExecuteIBTP handles interchain ibtps from other appchains
	// and return the receipt ibtp for ack or callback
	ExecuteIBTP(Ibtp *pb.IBTP) (*pb.IBTP, error)

	// Rollback rollbacks ibtp on appchain
	Rollback(ibtp *pb.IBTP, isSrcChain bool)

	// QueryInterchainMeta queries latest index map of ibtps executed on appchain
	// For the returned map, key is the source chain ID,
	// and value is the latest index of tx executed on appchain
	QueryInterchainMeta() map[string]uint64

	// QueryCallbackMeta queries latest index map of ibtps callbacks executed on appchain
	// For the returned map, key is the destination chain ID,
	// and value is the latest index of callback executed on appchain
	QueryCallbackMeta() map[string]uint64

	// QueryIBTPReceipt query receipt for original interchain ibtp
	QueryIBTPReceipt(originalIBTP *pb.IBTP) (*pb.IBTP, error)
}
