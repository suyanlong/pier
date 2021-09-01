package syncer

import (
	"fmt"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	rpcx "github.com/link33/sidercar/hub/client"
	"github.com/link33/sidercar/model/constant"
	"github.com/link33/sidercar/model/pb"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/types"
)

const (
	srcchainNotAvailable = "current appchain not available"
	invalidIBTP          = "invalid ibtp"
	ibtpIndexExist       = "index already exists"
	ibtpIndexWrong       = "wrong index"
	noBindRule           = "appchain didn't register rule"
)

var (
	ErrIBTPNotFound  = fmt.Errorf("receipt from bitxhub failed")
	ErrMetaOutOfDate = fmt.Errorf("interchain meta is out of date")
)

func (syncer *WrapperSyncer) GetAppchains() ([]*appchainmgr.Appchain, error) {
	panic("implement me")
}

func (syncer *WrapperSyncer) GetInterchainById(from string) *pb.Interchain {
	panic("implement me")
}

func (syncer *WrapperSyncer) QueryInterchainMeta() *pb.Interchain {
	panic("implement me")
}

func (syncer *WrapperSyncer) QueryIBTP(ibtpID string) (*pb.IBTP, bool, error) {
	queryTx, err := syncer.client.GenerateContractTx(pb.TransactionData_BVM, constant.InterchainContractAddr.Address(),
		"GetIBTPByID", rpcx.String(ibtpID))
	if err != nil {
		return nil, false, err
	}
	queryTx.Nonce = 1
	receipt, err := syncer.client.SendView(queryTx)
	if err != nil {
		return nil, false, err
	}

	if !receipt.IsSuccess() {
		return nil, false, fmt.Errorf("%w: %s", ErrIBTPNotFound, string(receipt.Ret))
	}

	hash := types.NewHash(receipt.Ret)
	response, err := syncer.client.GetTransaction(hash.String())
	if err != nil {
		return nil, false, err
	}
	receipt, err = syncer.client.GetReceipt(hash.String())
	if err != nil {
		return nil, false, err
	}
	return response.Tx.GetIBTP(), receipt.Status == pb.Receipt_SUCCESS, nil
}

func (syncer *WrapperSyncer) ListenIBTP() <-chan *pb.IBTPX {
	return syncer.ibtpC
}

func (syncer *WrapperSyncer) SendIBTP(ibtp *pb.IBTP) error {
	panic("implement me")
}

func (syncer *WrapperSyncer) retryFunc(handle func(uint) error) error {
	return retry.Retry(func(attempt uint) error {
		if err := handle(attempt); err != nil {
			syncer.logger.Errorf("retry failed for reason: %s", err.Error())
			return err
		}
		return nil
	}, strategy.Wait(500*time.Millisecond))
}
