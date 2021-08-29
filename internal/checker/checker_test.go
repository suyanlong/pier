package checker

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/pier/model/pb"
	"github.com/stretchr/testify/require"
)

const (
	from           = "0xe02d8fdacd59020d7f292ab3278d13674f5c404d"
	to             = "0x0915fdfc96232c95fb9c62d27cc9dc0f13f50161"
	from2          = "0x0915fdfc96232c95fb9c62d27cc9dc0f13f50162"
	rulePrefix     = "validation-rule-"
	proofPath      = "./testdata/proof_1.0.0_rc"
	proofPath2     = "./testdata/proof_1.0.0_rc_complex"
	validatorsPath = "./testdata/single_validator"
)

func TestMockChecker_Check(t *testing.T) {
	checker := &MockChecker{}
	require.Nil(t, checker.Check(nil))
}

func getAppchain(id, chainType string) (*appchainmgr.Appchain, error) {
	validators, err := ioutil.ReadFile(validatorsPath)
	if err != nil {
		return nil, err
	}

	app := &appchainmgr.Appchain{
		ID:            id,
		Name:          "chainA",
		Validators:    string(validators),
		ConsensusType: "rbft",
		ChainType:     chainType,
		Desc:          "appchain",
		Version:       "1.4.3",
		PublicKey:     "",
	}

	return app, nil
}

func getIBTP(t *testing.T, index uint64, typ pb.IBTP_Type, fid, tid, proofPath string) *pb.IBTP {
	ct := &pb.Content{
		SrcContractId: "mychannel&transfer",
		DstContractId: "mychannel&transfer",
		Func:          "get",
		Args:          [][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		Callback:      "interchainConfirm",
	}
	c, err := ct.Marshal()
	require.Nil(t, err)

	pd := pb.Payload{
		Encrypted: false,
		Content:   c,
	}
	ibtppd, err := pd.Marshal()
	require.Nil(t, err)

	proof, err := ioutil.ReadFile(proofPath)
	require.Nil(t, err)

	return &pb.IBTP{
		From:      fid,
		To:        tid,
		Payload:   ibtppd,
		Index:     index,
		Type:      typ,
		Timestamp: time.Now().UnixNano(),
		Proof:     proof,
	}
}

// MockAppchainMgr================================================
type MockAppchainMgr struct {
}

func (m MockAppchainMgr) ChangeStatus(id, trigger string, extra []byte) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) CountAvailable(extra []byte) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) CountAll(extra []byte) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) All(extra []byte) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) QueryById(id string, extra []byte) (bool, []byte) {
	if id == from || id == from2 {
		app, err := getAppchain(id, "fabric")
		if err != nil {
			return false, nil
		}
		data, err := json.Marshal(app)
		if err != nil {
			return false, nil
		}
		return true, data
	} else if id == to {
		app, err := getAppchain(id, "ethereum")
		data, err := json.Marshal(app)
		if err != nil {
			return false, nil
		}
		return true, data
	} else if id == "10" {
		return true, []byte("10")
	} else {
		return false, nil
	}
}

func (m MockAppchainMgr) Register(info []byte) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) Update(info []byte) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) CountAvailableAppchains() (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) UpdateAppchain(id, appchainOwner, docAddr, docHash, validators string, consensusType, chainType, name, desc, version, pubkey string) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) Audit(proposer string, isApproved int32, desc string) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) FetchAuditRecords(id string) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) CountApprovedAppchains() (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) CountAppchains() (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) Appchains() (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) DeleteAppchain(id string) (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) Appchain() (bool, []byte) {
	return true, nil
}

func (m MockAppchainMgr) GetPubKeyByChainID(id string) (bool, []byte) {
	return true, nil
}

var _ appchainmgr.AppchainMgr = &MockAppchainMgr{}
