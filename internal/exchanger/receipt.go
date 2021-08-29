package exchanger

import (
	"github.com/meshplus/pier/pkg/model"
)

//中继模式
func (ex *Exchanger) feedIBTPReceipt(receipt *model.WrappedIBTP) {
	var pool *Pool
	act, loaded := ex.receipts.Load(receipt.Ibtp.To)
	if !loaded {
		pool = NewPool()
		ex.receipts.Store(receipt.Ibtp.To, pool)
	} else {
		pool = act.(*Pool)
	}
	pool.feed(receipt)
}
