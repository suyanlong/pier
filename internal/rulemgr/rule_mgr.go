package rulemgr

import (
	"sync"

	"github.com/meshplus/bitxhub-core/validator"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/pier/internal/peermgr"
	peerproto "github.com/meshplus/pier/internal/peermgr/proto"
	"github.com/sirupsen/logrus"
)

type Rule struct {
	Code    []byte `json:"code"`
	Address string `json:"address"`
}

type RuleResponse struct {
	Ok      bool   `json:"ok"`
	Content string `json:"content"`
}

type RuleMgr struct {
	Ledger      *CodeLedger
	PeerManager peermgr.PeerManager
	Ve          *validator.ValidationEngine
	logger      logrus.FieldLogger
}

func New(storage storage.Storage, pm peermgr.PeerManager, logger logrus.FieldLogger) (*RuleMgr, error) {
	ledger := &CodeLedger{
		storage: storage,
	}
	instances := &sync.Map{}
	ve := validator.NewValidationEngine(ledger, instances, logger)
	rm := &RuleMgr{
		Ledger:      ledger,
		PeerManager: pm,
		Ve:          ve,
		logger:      logger,
	}
	if err := pm.RegisterMsgHandler(peerproto.Message_RULE_DEPLOY, rm.handleRule); err != nil {
		return nil, err
	}
	return rm, nil
}

func (rm *RuleMgr) Validate(address, from string, proof, payload []byte, validators string) (bool, error) {
	return rm.Ve.Validate(address, from, proof, payload, validators)
}

//规则验证，主要是验证交易的背书。当前阶段可以不需要加入。主要是通过WASM执行验证。
