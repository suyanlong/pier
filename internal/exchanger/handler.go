package exchanger

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	network "github.com/meshplus/go-lightp2p"
	"github.com/meshplus/pier/internal/peermgr"
	peerMsg "github.com/meshplus/pier/internal/peermgr/proto"
	"github.com/meshplus/pier/internal/syncer"
	"github.com/meshplus/pier/pkg/model"
	"github.com/sirupsen/logrus"
)

// handleIBTP handle ibtps from bitxhub
func (ex *Exchanger) handleIBTP(wIbtp *model.WrappedIBTP, entry logrus.FieldLogger) {
	ibtp := wIbtp.Ibtp
	err := ex.checker.Check(ibtp)
	if err != nil {
		// todo: send receipt back to bitxhub
		return
	}
	entry.Debugf("IBTP pass check")
	if pb.IBTP_ASSET_EXCHANGE_REDEEM == ibtp.Type || pb.IBTP_ASSET_EXCHANGE_REFUND == ibtp.Type {
		if err := retry.Retry(func(attempt uint) error {
			if err := ex.fetchSignsToIBTP(ibtp); err != nil {
				return err
			}
			return nil
		}, strategy.Wait(1*time.Second)); err != nil {
			ex.logger.Panic(err)
		}
	}

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
		err = ex.syncer.SendIBTP(receipt)
		if err != nil {
			ex.logger.Errorf("send ibtp error: %s", err.Error())
			if errors.Is(err, syncer.ErrMetaOutOfDate) {
				ex.updateSourceReceiptMeta()
				return
			}
			// if sending receipt failed, try to get new receipt from appchain and retry
		queryLoop:
			for {
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

// handleIBTP handle ibtps from bitxhub
func (ex *Exchanger) handleUnionIBTP(wIbtp *model.WrappedIBTP) {
	ibtp := wIbtp.Ibtp
	if ibtp.To == ex.appchainDID {
		ex.exec.ExecuteIBTP(wIbtp)
		ex.logger.WithFields(logrus.Fields{
			"index": ibtp.Index,
			"type":  ibtp.Type,
			"from":  ibtp.From,
			"id":    ibtp.ID(),
		}).Infof("Handle union ibtp sent to executor")
		return
	}

	ibtp.From = ex.appchainDID + "-" + ibtp.From // for inter-relay they're the same
	var signs []byte
	if err := retry.Retry(func(attempt uint) error {
		var err error
		signs, err = ex.syncer.GetIBTPSigns(ibtp)
		if err != nil {
			return err
		}
		return nil
	}, strategy.Wait(1*time.Second)); err != nil {
		ex.logger.Panic(err)
	}

	ibtp.Proof = signs

	if err := retry.Retry(func(attempt uint) error {
		err := ex.router.Route(ibtp)
		if err != nil {
			return err
		}
		return nil
	}, strategy.Wait(1*time.Second)); err != nil {
		ex.logger.Panic(err)
	}
	ex.logger.WithFields(logrus.Fields{"index": ibtp.Index, "type": ibtp.Type, "from": ibtp.From, "id": ibtp.ID()}).Info("Route tx successful")
}

//联盟模式: 主要是从bithub上同步appchain信息。
func (ex *Exchanger) handleProviderAppchains() error {
	appchains, err := ex.syncer.GetAppchains()
	if err != nil {
		return fmt.Errorf("get appchains:%w", err)
	}
	return ex.router.AddAppchains(appchains)
}

//联盟模式: handleRouterSendIBTPMessage handles IBTP from union interchain network
func (ex *Exchanger) handleRouterSendIBTPMessage(stream network.Stream, msg *peerMsg.Message) {
	handle := func() error {
		wIbtp := &model.WrappedIBTP{}
		if err := json.Unmarshal(msg.Payload.Data, wIbtp); err != nil {
			return fmt.Errorf("unmarshal ibtp: %w", err)
		}

		ibtp := wIbtp.Ibtp
		entry := ex.logger.WithFields(logrus.Fields{
			"index": ibtp.Index,
			"type":  ibtp.Type,
			"from":  ibtp.From,
			"to":    ibtp.To,
			"id":    ibtp.ID(),
		})

		// 这个可以放到Router里面路由。
		retMsg := peermgr.Message(peerMsg.Message_ACK, true, nil)
		if !ex.router.ExistAppchain(ibtp.To) {
			entry.WithField("appchain", ibtp.To).Errorf("cannot found appchain in relay network")
			retMsg.Payload.Ok = false
		}

		err := ex.peerMgr.AsyncSendWithStream(stream, retMsg)
		if err != nil {
			return fmt.Errorf("send back ibtp: %w", err)
		}
		entry.Infof("receive ibtp successfully")

		if err := ex.sendIBTP(ibtp); err != nil {
			return fmt.Errorf("send ibtp: %v", err)
		}
		return nil
	}

	if err := handle(); err != nil {
		ex.logger.Error(err)
		return
	}
}

func (ex *Exchanger) timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		ex.sendIBTPTimer.Add(tc)
	}
}

//联盟模式:
func (ex *Exchanger) handleRecover(ibtp *pb.IBTP) (*rpcx.Interchain, error) {
	pierId, err := ex.peerMgr.FindProviders(ibtp.To)
	if err != nil {
		return nil, err
	}
	msg := peermgr.Message(peerMsg.Message_ROUTER_INTERCHAIN_SEND, true, []byte(ibtp.From))
	res, err := ex.peerMgr.Send(pierId, msg)
	if err != nil {
		return nil, fmt.Errorf("router interchain:%v", err)
	}
	interchain := &pb.Interchain{}
	err = json.Unmarshal(res.Payload.Data, interchain)
	if err != nil {
		return nil, err
	}
	ic := &rpcx.Interchain{
		ID:                   interchain.ID,
		InterchainCounter:    interchain.InterchainCounter,
		ReceiptCounter:       interchain.ReceiptCounter,
		SourceReceiptCounter: interchain.SourceReceiptCounter,
	}
	return ic, nil
}

//联盟模式:
func (ex *Exchanger) handleRouterInterchain(s network.Stream, msg *peerMsg.Message) {
	ic := ex.syncer.GetInterchainById(string(msg.Payload.Data))
	data, err := ic.Marshal()
	if err != nil {
		panic(err)
	}
	retMsg := peermgr.Message(peerMsg.Message_ACK, true, data)
	err = ex.peerMgr.AsyncSendWithStream(s, retMsg)
	if err != nil {
		ex.logger.Error(err)
	}
}

func (ex *Exchanger) fetchSignsToIBTP(ibtp *pb.IBTP) error {
	signs, err := ex.syncer.GetAssetExchangeSigns(string(ibtp.Extra))
	if err != nil {
		return fmt.Errorf("get asset exchange signs: %w", err)
	}

	ibtp.Extra = signs

	return nil
}
