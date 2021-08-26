package exchanger

import (
	"errors"
	"time"

	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/pier/internal/syncer"
	"github.com/meshplus/pier/pkg/model"
	"github.com/sirupsen/logrus"
)

//中继模式 handleIBTP handle ibtps from bitxhub
func (ex *Exchanger) handleIBTP(wIbtp *model.WrappedIBTP, entry logrus.FieldLogger) {
	ibtp := wIbtp.Ibtp
	err := ex.checker.Check(ibtp)
	if err != nil {
		// todo: send receipt back to bitxhub
		return
	}
	entry.Debugf("IBTP pass check")
	receipt, err := ex.exec.ExecuteIBTP(wIbtp)
	if err != nil {
		ex.logger.Errorf("execute ibtp error:%s", err.Error())
	}
	if receipt == nil {
		ex.logger.WithFields(logrus.Fields{"type": ibtp.Type, "id": ibtp.ID()}).Info("Handle ibtp receipt success")
		return
	}

sendReceiptLoop:
	for {
		err = ex.syncer.SendIBTP(receipt) // 回执一定是要到中继链上的，作为数据凭证。
		if err != nil {
			ex.logger.Errorf("send ibtp error: %s", err.Error())
			if errors.Is(err, syncer.ErrMetaOutOfDate) {
				ex.updateSourceReceiptMeta()
				return
			}
			// if sending receipt failed, try to get new receipt from appchain and retry
		queryLoop:
			for {
				// 死循环，直到成功。
				receipt, err = ex.exec.QueryIBTPReceipt(ibtp)
				if err != nil {
					ex.logger.Errorf("Query ibtp receipt for %s error: %s", ibtp.ID(), err.Error())
					time.Sleep(1 * time.Second)
					continue queryLoop
				}
				time.Sleep(1 * time.Second)
				continue sendReceiptLoop
			}
		}
		break
	}
	ex.logger.WithFields(logrus.Fields{"type": ibtp.Type, "id": ibtp.ID()}).Info("Handle ibtp success")
}

// 中继模式：处理链间交易回执
func (ex *Exchanger) applyReceipt(wIbtp *model.WrappedIBTP, entry logrus.FieldLogger) {
	ibtp := wIbtp.Ibtp
	index := ex.callbackCounter[ibtp.To]
	if index >= ibtp.Index {
		entry.Infof("Ignore ibtp callback, expected index %d", index+1)
		return
	}

	if index+1 < ibtp.Index {
		entry.Infof("Get missing ibtp receipt, expected index %d", index+1)
		// todo: need to handle missing ibtp receipt or not?
		return
	}
	ex.handleIBTP(wIbtp, entry)
	ex.callbackCounter[ibtp.To] = ibtp.Index
}

// 中继链架构,处理链间交易，
func (ex *Exchanger) applyInterchain(wIbtp *model.WrappedIBTP, entry logrus.FieldLogger) {
	ibtp := wIbtp.Ibtp
	index := ex.executorCounter[ibtp.From]
	if index >= ibtp.Index {
		entry.Infof("Ignore ibtp, expected %d", index+1)
		return
	}

	if index+1 < ibtp.Index {
		entry.Info("Get missing ibtp")
		if err := ex.handleMissingIBTPFromSyncer(ibtp.From, index+1, ibtp.Index); err != nil {
			entry.WithField("err", err).Error("Handle missing ibtp")
			return
		}
	}
	ex.handleIBTP(wIbtp, entry)
	ex.executorCounter[ibtp.From] = ibtp.Index
}

//中继模式
func (ex *Exchanger) handleRollback(ibtp *pb.IBTP) {
	if ibtp.Category() == pb.IBTP_RESPONSE {
		// if this is receipt type of ibtp, no need to rollback
		return
	}
	ex.feedIBTPReceipt(&model.WrappedIBTP{Ibtp: ibtp, IsValid: false})
}

func (ex *Exchanger) timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		ex.sendIBTPTimer.Add(tc)
	}
}
