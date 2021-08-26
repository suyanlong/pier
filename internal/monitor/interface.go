package monitor

import "github.com/meshplus/bitxhub-model/pb"

//go:generate mockgen -destination mock_monitor/mock_monitor.go -package mock_monitor -source interface.go
type Monitor interface {
	// Start starts the service of monitor
	Start() error
	// Stop stops the service of monitor
	Stop() error
	// listen on interchain ibtp from appchain
	ListenIBTP() <-chan *pb.IBTP
	// query historical ibtp by its id
	QueryIBTP(id string) (*pb.IBTP, error)
	// QueryLatestMeta queries latest index map of ibtps threw on appchain
	QueryOuterMeta() map[string]uint64
}

// 是否可以重构，意义不是多大。视乎与Executor意义一样。可以抽象为同一接口。
// 注册moniter
// 注册Executor
//
//from: 这一方方需要实现moniter接口，Syncer这个接口和moniter一样。不过适用于同步to一方数据的。
//to: 这一方方需要实现Executor接口。
