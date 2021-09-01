package syncer

import (
	rpcx "github.com/link33/sidercar/hub/client"
	"github.com/link33/sidercar/model/pb"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
)

type IBTPHandler func(ibtp *pb.IBTP)

type AppchainHandler func() error

type RecoverUnionHandler func(ibtp *pb.IBTP) (*rpcx.Interchain, error)

type RollbackHandler func(ibtp *pb.IBTP)

//go:generate mockgen -destination mock_syncer/mock_syncer.go -package mock_syncer -source interface.go
type Syncer interface {
	// QueryInterchainMeta queries meta including interchain and receipt related meta from bitxhub
	QueryInterchainMeta() *pb.Interchain

	// QueryIBTP query ibtp from bitxhub by its id.
	// if error occurs, it means this ibtp is not existed on bitxhub
	QueryIBTP(ibtpID string) (*pb.IBTP, bool, error)

	// ListenIBTP listen on the ibtps destined for this sidercar from bitxhub
	ListenIBTP() <-chan *pb.IBTPX

	// SendIBTP sends interchain or receipt type of ibtp to bitxhub
	// if error occurs, user need to reconstruct this ibtp cause it means ibtp is invalid on bitxhub
	SendIBTP(ibtp *pb.IBTP) error

	//GetAppchains gets appchains from bitxhub node
	//GetAppchains() ([]*appchainmgr.Appchain, error)

	////GetInterchainById gets interchain meta by appchain id
	//GetInterchainById(from string) *pb.Interchain

	// RegisterRecoverHandler registers handler that recover ibtps from bitxhub
	//RegisterRecoverHandler(RecoverUnionHandler) error

	// RegisterAppchainHandler registers handler that fetch appchains information
	//RegisterAppchainHandler(handler AppchainHandler) error

	RegisterRollbackHandler(handler RollbackHandler) error
}

type Handler interface {
	//GetAppchains gets appchains from bitxhub node
	GetAppchains() ([]*appchainmgr.Appchain, error)

	//GetInterchainById gets interchain meta by appchain id
	GetInterchainById(from string) *pb.Interchain

	//RegisterRecoverHandler registers handler that recover ibtps from bitxhub
	RegisterRecoverHandler(RecoverUnionHandler) error

	//RegisterAppchainHandler registers handler that fetch appchains information
	RegisterAppchainHandler(handler AppchainHandler) error
}
