package exchanger

import (
	"github.com/link33/sidercar/model/pb"
	"sync"
)

type Pool struct {
	ibtps *sync.Map
	ch    chan *pb.IBTPX
}

func NewPool() *Pool {
	return &Pool{
		ibtps: &sync.Map{},
		ch:    make(chan *pb.IBTPX, 40960),
	}
}

func (pool *Pool) feed(ibtp *pb.IBTPX) {
	pool.ch <- ibtp
}

func (pool *Pool) put(ibtp *pb.IBTPX) {
	pool.ibtps.Store(ibtp.Ibtp.Index, ibtp)
}

func (pool *Pool) delete(idx uint64) {
	pool.ibtps.Delete(idx)
}

func (pool *Pool) get(index uint64) *pb.IBTPX {
	ibtp, ok := pool.ibtps.Load(index)
	if !ok {
		return nil
	}

	return ibtp.(*pb.IBTPX)
}
