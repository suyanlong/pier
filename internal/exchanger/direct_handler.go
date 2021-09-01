package exchanger

import (
	"encoding/json"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/link33/sidecar/internal/peermgr"
	"github.com/link33/sidecar/internal/port"
	"github.com/link33/sidecar/model/pb"
	"github.com/sirupsen/logrus"
	"time"
)

func (ex *Exchanger) feedIBTP(wIbtp *pb.IBTPX) {
	var pool *Pool
	ibtp := wIbtp.Ibtp
	act, loaded := ex.ibtps.Load(ibtp.From)
	if !loaded {
		pool = NewPool()
		ex.ibtps.Store(ibtp.From, pool)
	} else {
		pool = act.(*Pool)
	}
	pool.feed(wIbtp)

	if !loaded {
		go func(pool *Pool) {
			defer func() {
				if e := recover(); e != nil {
					ex.logger.Error(fmt.Errorf("%v", e))
				}
			}()
			inMeta := ex.exec.QueryInterchainMeta()
			for wIbtp := range pool.ch {
				ibtp := wIbtp.Ibtp
				idx := inMeta[ibtp.From]
				if ibtp.Index <= idx {
					pool.delete(ibtp.Index)
					ex.logger.Warnf("ignore ibtp with invalid index: %d", ibtp.Index)
					continue
				}
				if idx+1 == ibtp.Index {
					ex.processIBTP(wIbtp)
					pool.delete(ibtp.Index)
					index := ibtp.Index + 1
					wIbtp := pool.get(index)
					for wIbtp != nil {
						ex.processIBTP(wIbtp)
						pool.delete(wIbtp.Ibtp.Index)
						index++
						wIbtp = pool.get(index)
					}
				} else {
					pool.put(wIbtp)
				}
			}
		}(pool)
	}
}

// 直连
func (ex *Exchanger) processIBTP(wIbtp *pb.IBTPX) {
	receipt, err := ex.exec.ExecuteIBTP(wIbtp)
	if err != nil {
		ex.logger.Errorf("Execute ibtp error: %s", err.Error())
		return
	}
	ex.postHandleIBTP(wIbtp.Ibtp.From, receipt)
	ex.sendIBTPCounter.Inc()
}

// 直连
func (ex *Exchanger) feedReceipt(receipt *pb.IBTP) {
	var pool *Pool
	act, loaded := ex.ibtps.Load(receipt.To)
	if !loaded {
		pool = NewPool()
		ex.ibtps.Store(receipt.To, pool)
	} else {
		pool = act.(*Pool)
	}
	pool.feed(&pb.IBTPX{Ibtp: receipt, IsValid: true})

	if !loaded {
		go func(pool *Pool) {
			defer func() {
				if e := recover(); e != nil {
					ex.logger.Error(fmt.Errorf("%v", e))
				}
			}()
			callbackMeta := ex.exec.QueryCallbackMeta()
			for wIbtp := range pool.ch {
				ibtp := wIbtp.Ibtp
				if ibtp.Index <= callbackMeta[ibtp.To] {
					pool.delete(ibtp.Index)
					ex.logger.Warn("ignore ibtp with invalid index")
					continue
				}
				if callbackMeta[ibtp.To]+1 == ibtp.Index {
					ex.processIBTP(wIbtp)
					pool.delete(ibtp.Index)
					index := ibtp.Index + 1
					wIbtp := pool.get(index)
					for wIbtp != nil {
						ibtp := wIbtp.Ibtp
						receipt, _ := ex.exec.ExecuteIBTP(wIbtp)
						ex.postHandleIBTP(ibtp.From, receipt)
						pool.delete(ibtp.Index)
						index++
						wIbtp = pool.get(index)
					}
				} else {
					pool.put(wIbtp)
				}
			}
		}(pool)
	}
}

func (ex *Exchanger) postHandleIBTP(from string, receipt *pb.IBTP) {
	if receipt == nil {
		retMsg := peermgr.Message(pb.Message_IBTP_RECEIPT_SEND, true, nil)
		err := ex.peerMgr.AsyncSend(from, retMsg)
		if err != nil {
			ex.logger.Errorf("Send back empty ibtp receipt: %s", err.Error())
		}
		return
	}

	data, _ := receipt.Marshal()
	retMsg := peermgr.Message(pb.Message_IBTP_RECEIPT_SEND, true, data)
	if err := ex.peerMgr.AsyncSend(from, retMsg); err != nil {
		ex.logger.Errorf("Send back ibtp receipt: %s", err.Error())
	}
}

//直链模式
func (ex *Exchanger) handleSendIBTPMessage(p port.Port, msg *pb.Message) {
	ex.ch <- struct{}{}
	go func(msg *pb.Message) {
		wIbtp := &pb.IBTPX{}
		if err := json.Unmarshal(msg.Payload.Data, wIbtp); err != nil {
			ex.logger.Errorf("Unmarshal ibtp: %s", err.Error())
			return
		}
		defer ex.timeCost()()
		err := ex.checker.Check(wIbtp.Ibtp)
		if err != nil {
			ex.logger.Error("check ibtp: %w", err)
			return
		}
		ex.feedIBTP(wIbtp)
		<-ex.ch
	}(msg)
}

//直链模式
func (ex *Exchanger) handleSendIBTPReceiptMessage(p port.Port, msg *pb.Message) {
	if msg.Payload.Data == nil {
		return
	}
	receipt := &pb.IBTP{}
	if err := receipt.Unmarshal(msg.Payload.Data); err != nil {
		ex.logger.Error("unmarshal ibtp: %w", err)
		return
	}

	// ignore msg for receipt type
	if receipt.Type == pb.IBTP_RECEIPT_SUCCESS || receipt.Type == pb.IBTP_RECEIPT_FAILURE {
		//ex.logger.Warn("ignore receipt ibtp")
		return
	}

	err := ex.checker.Check(receipt)
	if err != nil {
		ex.logger.Error("check ibtp: %w", err)
		return
	}

	ex.feedReceipt(receipt)

	ex.logger.Info("Receive ibtp receipt from other sidecar")
}

// 直连
func (ex *Exchanger) handleGetIBTPMessage(p port.Port, msg *pb.Message) {
	ibtpID := string(msg.Payload.Data)
	ibtp, err := ex.mnt.QueryIBTP(ibtpID)
	if err != nil {
		ex.logger.Error("Get wrong ibtp id")
		return
	}

	data, err := ibtp.Marshal()
	if err != nil {
		return
	}

	retMsg := peermgr.Message(pb.Message_ACK, true, data)

	err = ex.peerMgr.AsyncSendWithPort(p, retMsg)
	if err != nil {
		ex.logger.Error(err)
	}
}

// 直连
func (ex *Exchanger) handleNewConnection(dstSidecarID string) {
	appchainMethod := []byte(ex.appchainDID)
	msg := peermgr.Message(pb.Message_INTERCHAIN_META_GET, true, appchainMethod)

	indices := &struct {
		InterchainIndex uint64 `json:"interchain_index"`
		ReceiptIndex    uint64 `json:"receipt_index"`
	}{}

	loop := func() error {
		interchainMeta, err := ex.peerMgr.Send(dstSidecarID, msg)
		if err != nil {
			return err
		}

		if !interchainMeta.Payload.Ok {
			return fmt.Errorf("interchain meta message payload is false")
		}

		if err = json.Unmarshal(interchainMeta.Payload.Data, indices); err != nil {
			return err
		}

		return nil
	}

	if err := retry.Retry(func(attempt uint) error {
		return loop()
	}, strategy.Wait(1*time.Second)); err != nil {
		ex.logger.Panic(err)
	}
}

//直链模式
func (ex *Exchanger) handleGetInterchainMessage(p port.Port, msg *pb.Message) {
	mntMeta := ex.mnt.QueryOuterMeta()
	execMeta := ex.exec.QueryInterchainMeta()

	indices := &struct {
		InterchainIndex uint64 `json:"interchain_index"`
		ReceiptIndex    uint64 `json:"receipt_index"`
	}{}

	execLoad, ok := execMeta[string(msg.Payload.Data)]
	if ok {
		indices.InterchainIndex = execLoad
	}

	mntLoad, ok := mntMeta[string(msg.Payload.Data)]
	if ok {
		indices.InterchainIndex = mntLoad
	}

	data, err := json.Marshal(indices)
	if err != nil {
		panic(err)
	}

	retMsg := peermgr.Message(pb.Message_ACK, true, data)
	if err := ex.peerMgr.AsyncSendWithPort(p, retMsg); err != nil {
		ex.logger.Error(err)
		return
	}
}

//直链模式
func (ex *Exchanger) analysisDirectTPS() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	current := time.Now()
	counter := ex.sendIBTPCounter.Load()
	for {
		select {
		case <-ticker.C:
			tps := ex.sendIBTPCounter.Load() - counter
			counter = ex.sendIBTPCounter.Load()
			totalTimer := ex.sendIBTPTimer.Load()

			if tps != 0 {
				ex.logger.WithFields(logrus.Fields{
					"tps":      tps,
					"tps_sum":  counter,
					"tps_time": totalTimer.Milliseconds() / int64(counter),
					"tps_avg":  float64(counter) / time.Since(current).Seconds(),
				}).Info("analysis")
			}

		case <-ex.ctx.Done():
			return
		}
	}
}
