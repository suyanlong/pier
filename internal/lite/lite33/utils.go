package lite33

import (
	"fmt"
	"sync/atomic"
)

func (lite *Lite33) getDemandHeight() uint64 {
	return atomic.LoadUint64(&lite.height) + 1
}

func (lite *Lite33) updateHeight() {
	atomic.AddUint64(&lite.height, 1)
}

func headerHeightKey() []byte {
	return []byte("lite-height")
}

func headerKey(height uint64) []byte {
	return []byte(fmt.Sprintf("header-%d", height))
}
