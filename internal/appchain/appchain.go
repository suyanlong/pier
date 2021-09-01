package appchain

import (
	"context"
	"fmt"
	"github.com/link33/sidercar/internal/port"
	"strconv"
	"strings"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/link33/sidercar/internal/txcrypto"
	"github.com/link33/sidercar/model/pb"
	"github.com/link33/sidercar/pkg/plugins"
	"github.com/sirupsen/logrus"
)

type AppChain interface {
	Executor
	Monitor
}

type appChain struct {
	client  plugins.Client
	recvCh  chan *pb.IBTP
	logger  logrus.FieldLogger
	ctx     context.Context
	cancel  context.CancelFunc
	cryptor txcrypto.Cryptor
}

func NewPort(client plugins.Client, cryptor txcrypto.Cryptor, logger logrus.FieldLogger) (port.Port, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &appChain{
		client:  client,
		cryptor: cryptor,
		logger:  logger,
		recvCh:  make(chan *pb.IBTP, 1024),
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func NewPorts(clients []plugins.Client, cryptor txcrypto.Cryptor, logger logrus.FieldLogger) []port.Port {
	var ps []port.Port
	for _, c := range clients {
		p, err := NewPort(c, cryptor, logger)
		if err != nil {
			panic(err)
		}
		ps = append(ps, p)
	}

	return ps
}

func (a *appChain) QueryInterchainMeta() map[string]uint64 {
	execMeta, err := a.client.GetInMeta()
	if err != nil {
		return map[string]uint64{}
	}
	return execMeta
}

func (a *appChain) QueryCallbackMeta() map[string]uint64 {
	callbackMeta, err := a.client.GetCallbackMeta()
	if err != nil {
		return map[string]uint64{}
	}
	return callbackMeta
}

// getReceipt only generates one receipt given source chain id and interchain tx index
func (a *appChain) QueryIBTPReceipt(originalIBTP *pb.IBTP) (*pb.IBTP, error) {
	if originalIBTP == nil {
		return nil, fmt.Errorf("empty original ibtp")
	}
	return a.client.GetReceipt(originalIBTP)
}

func (a *appChain) ExecuteIBTP(ibtp *pb.IBTP) (*pb.IBTP, error) {
	if ibtp == nil {
		a.logger.Error("empty ibtp structure")
		return nil, fmt.Errorf("nil ibtp structure")
	}
	a.logger.WithFields(logrus.Fields{
		"index": ibtp.Index,
		"type":  ibtp.Type,
		"from":  ibtp.From,
		"id":    ibtp.ID(),
	}).Info("Apply tx")

	switch ibtp.Type {
	case pb.IBTP_INTERCHAIN:
		return a.applyInterchainIBTP(ibtp)
	case pb.IBTP_RECEIPT_SUCCESS, pb.IBTP_RECEIPT_FAILURE:
		err := a.applyReceiptIBTP(ibtp)
		return nil, err
	default:
		return nil, fmt.Errorf("wrong ibtp type")
	}
}

func (a *appChain) applyInterchainIBTP(ibtp *pb.IBTP) (*pb.IBTP, error) {
	entry := a.logger.WithFields(logrus.Fields{
		"from":  ibtp.From,
		"type":  ibtp.Type,
		"index": ibtp.Index,
	})

	// todo: deal with plugin returned error
	// execute interchain tx, and if execution failed, try to rollback
	response, err := a.client.SubmitIBTP(ibtp)
	if err != nil {
		entry.WithField("error", err).Panic("Submit ibtp")
	}

	if response == nil || response.Result == nil {
		entry.WithField("error", err).Panic("empty response")
	}

	if !response.Status {
		pd := &pb.Payload{}
		if err := pd.Unmarshal(response.Result.Payload); err != nil {
			entry.Panic("Unmarshal payload")
		}

		entry.WithFields(logrus.Fields{
			"result":  response.Message,
			"payload": pd,
		}).Warn("Get wrong response, need rollback on source chain")
	}

	return response.Result, nil
}

func (a *appChain) applyReceiptIBTP(ibtp *pb.IBTP) error {
	pd := &pb.Payload{}
	if err := pd.Unmarshal(ibtp.Payload); err != nil {
		return fmt.Errorf("unmarshal receipt type ibtp payload: %w", err)
	}

	ct := &pb.Content{}
	contentByte := pd.Content

	var err error
	if pd.Encrypted {
		contentByte, err = a.cryptor.Decrypt(contentByte, ibtp.To)
		if err != nil {
			return fmt.Errorf("decrypt ibtp payload content: %w", err)
		}
	}

	if err := ct.Unmarshal(contentByte); err != nil {
		return fmt.Errorf("unmarshal payload content: %w", err)
	}

	if err := retry.Retry(func(attempt uint) error {
		if err := a.execCallback(ibtp); err != nil {
			a.logger.Errorf("Execute callback tx: %s, retry sending tx", err.Error())
			return fmt.Errorf("execute callback tx: %w", err)
		}
		return nil
	}, strategy.Wait(1*time.Second)); err != nil {
		a.logger.Errorf("Execution of callback function failed: %s", err.Error())
	}
	return nil
}

func (a *appChain) execCallback(ibtp *pb.IBTP) error {
	ibtp.From, ibtp.To = ibtp.To, ibtp.From

	// no need to send receipt for callback
	resp, err := a.client.SubmitIBTP(ibtp)
	if err != nil {
		return fmt.Errorf("handle ibtp of callback %w", err)
	}

	// executor should not change the content of ibtp
	ibtp.From, ibtp.To = ibtp.To, ibtp.From
	a.logger.WithFields(logrus.Fields{
		"index":  ibtp.Index,
		"type":   ibtp.Type,
		"status": resp.Status,
		"msg":    resp.Message,
	}).Info("Execute callback")

	return nil
}

func (a *appChain) Rollback(ibtp *pb.IBTP, isSrcChain bool) {
	if err := retry.Retry(func(attempt uint) error {
		err := a.execRollback(ibtp, isSrcChain)
		if err != nil {
			a.logger.Errorf("Execute callback tx: %s, retry sending tx", err.Error())
			return fmt.Errorf("execute callback tx: %w", err)
		}
		return nil
	}, strategy.Wait(1*time.Second)); err != nil {
		a.logger.Errorf("Execution of callback function failed: %s", err.Error())
	}
}

func (a *appChain) execRollback(ibtp *pb.IBTP, isSrcChain bool) error {
	// no need to send receipt for callback
	resp, err := a.client.RollbackIBTP(ibtp, isSrcChain)
	if err != nil {
		return fmt.Errorf("rollback ibtp on source appchain %w", err)
	}

	a.logger.WithFields(logrus.Fields{
		"index":  ibtp.Index,
		"type":   ibtp.Type,
		"status": resp.Status,
		"msg":    resp.Message,
	}).Info("Executed rollbcak")
	return nil
}

//-------------------------------------------------------------------------------

// Start implements Monitor
func (a *appChain) Start() error {
	if err := a.client.Start(); err != nil {
		return err
	}

	ch := a.client.GetIBTP()
	go func() {
		for {
			select {
			case e := <-ch:
				a.logger.Debugf("Receive ibtp %s from plugin", e.ID())
				a.handleIBTP(e)
			case <-a.ctx.Done():
				return
			}
		}
	}()
	a.logger.Info("Monitor started")
	return nil
}

// Stop implements Monitor
func (a *appChain) Stop() error {
	a.cancel()
	a.logger.Info("Monitor stopped")
	return nil
}

func (a *appChain) ListenIBTP() <-chan *pb.IBTP {
	return a.recvCh
}

// QueryIBTP queries interchain tx recorded in appchain given ibtp id
func (a *appChain) QueryIBTP(id string) (*pb.IBTP, error) {
	// TODO(xcc): Encapsulate as a function
	args := strings.Split(id, "-")
	if len(args) != 3 {
		return nil, fmt.Errorf("invalid ibtp id %s", id)
	}

	idx, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid ibtp index")
	}

	c := make(chan *pb.IBTP, 1)
	if err := retry.Retry(func(attempt uint) error {
		// TODO(xcc): Need to distinguish error types
		e, err := a.client.GetOutMessage(args[1], idx)
		if err != nil {
			a.logger.WithFields(logrus.Fields{
				"error":   err,
				"ibtp_id": id,
			}).Error("Query ibtp")
			return err
		}
		c <- e
		return nil
	}, strategy.Wait(2*time.Second)); err != nil {
		panic(err)
	}

	ibtp := <-c

	if err := a.encryption(ibtp); err != nil {
		return nil, err
	}

	return ibtp, nil
}

// QueryOuterMeta queries outer meta from appchain.
// It will loop until the result is returned or panic.
func (a *appChain) QueryOuterMeta() map[string]uint64 {
	var (
		meta map[string]uint64
		err  error
	)
	if err := retry.Retry(func(attempt uint) error {
		meta, err = a.client.GetOutMeta()
		if err != nil {
			a.logger.WithField("error", err).Error("Get outer meta from appchain")
			return err
		}
		return nil
	}, strategy.Wait(2*time.Second)); err != nil {
		panic(err)
	}

	return meta
}

// handleIBTP handle the ibtp package captured by monitor.
func (a *appChain) handleIBTP(ibtp *pb.IBTP) {
	if err := a.encryption(ibtp); err != nil {
		a.logger.WithFields(logrus.Fields{
			"index": ibtp.Index,
			"to":    ibtp.To,
		}).Error("check encryption")
		return
	}

	a.recvCh <- ibtp
}

func (a *appChain) encryption(ibtp *pb.IBTP) error {
	pld := &pb.Payload{}
	if err := pld.Unmarshal(ibtp.Payload); err != nil {
		return err
	}
	if !pld.Encrypted {
		return nil
	}

	ctb, err := a.cryptor.Encrypt(pld.Content, ibtp.To)
	if err != nil {
		return err
	}
	pld.Content = ctb
	payload, err := pld.Marshal()
	if err != nil {
		return err
	}
	ibtp.Payload = payload

	return nil
}
