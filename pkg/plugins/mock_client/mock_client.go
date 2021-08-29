// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pb "github.com/meshplus/pier/model/pb"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CommitCallback mocks base method.
func (m *MockClient) CommitCallback(ibtp *pb.IBTP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitCallback", ibtp)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitCallback indicates an expected call of CommitCallback.
func (mr *MockClientMockRecorder) CommitCallback(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitCallback", reflect.TypeOf((*MockClient)(nil).CommitCallback), ibtp)
}

// GetCallbackMeta mocks base method.
func (m *MockClient) GetCallbackMeta() (map[string]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCallbackMeta")
	ret0, _ := ret[0].(map[string]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCallbackMeta indicates an expected call of GetCallbackMeta.
func (mr *MockClientMockRecorder) GetCallbackMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCallbackMeta", reflect.TypeOf((*MockClient)(nil).GetCallbackMeta))
}

// GetIBTP mocks base method.
func (m *MockClient) GetIBTP() chan *pb.IBTP {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIBTP")
	ret0, _ := ret[0].(chan *pb.IBTP)
	return ret0
}

// GetIBTP indicates an expected call of GetIBTP.
func (mr *MockClientMockRecorder) GetIBTP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIBTP", reflect.TypeOf((*MockClient)(nil).GetIBTP))
}

// GetInMessage mocks base method.
func (m *MockClient) GetInMessage(from string, idx uint64) ([][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInMessage", from, idx)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInMessage indicates an expected call of GetInMessage.
func (mr *MockClientMockRecorder) GetInMessage(from, idx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInMessage", reflect.TypeOf((*MockClient)(nil).GetInMessage), from, idx)
}

// GetInMeta mocks base method.
func (m *MockClient) GetInMeta() (map[string]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInMeta")
	ret0, _ := ret[0].(map[string]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInMeta indicates an expected call of GetInMeta.
func (mr *MockClientMockRecorder) GetInMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInMeta", reflect.TypeOf((*MockClient)(nil).GetInMeta))
}

// GetOutMessage mocks base method.
func (m *MockClient) GetOutMessage(to string, idx uint64) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutMessage", to, idx)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutMessage indicates an expected call of GetOutMessage.
func (mr *MockClientMockRecorder) GetOutMessage(to, idx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutMessage", reflect.TypeOf((*MockClient)(nil).GetOutMessage), to, idx)
}

// GetOutMeta mocks base method.
func (m *MockClient) GetOutMeta() (map[string]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutMeta")
	ret0, _ := ret[0].(map[string]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutMeta indicates an expected call of GetOutMeta.
func (mr *MockClientMockRecorder) GetOutMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutMeta", reflect.TypeOf((*MockClient)(nil).GetOutMeta))
}

// GetReceipt mocks base method.
func (m *MockClient) GetReceipt(ibtp *pb.IBTP) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReceipt", ibtp)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceipt indicates an expected call of GetReceipt.
func (mr *MockClientMockRecorder) GetReceipt(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceipt", reflect.TypeOf((*MockClient)(nil).GetReceipt), ibtp)
}

// IncreaseInMeta mocks base method.
func (m *MockClient) IncreaseInMeta(ibtp *pb.IBTP) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseInMeta", ibtp)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncreaseInMeta indicates an expected call of IncreaseInMeta.
func (mr *MockClientMockRecorder) IncreaseInMeta(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseInMeta", reflect.TypeOf((*MockClient)(nil).IncreaseInMeta), ibtp)
}

// Initialize mocks base method.
func (m *MockClient) Initialize(configPath, pierID string, extra []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", configPath, pierID, extra)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockClientMockRecorder) Initialize(configPath, pierID, extra interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockClient)(nil).Initialize), configPath, pierID, extra)
}

// Name mocks base method.
func (m *MockClient) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockClientMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockClient)(nil).Name))
}

// RollbackIBTP mocks base method.
func (m *MockClient) RollbackIBTP(ibtp *pb.IBTP, isSrcChain bool) (*pb.RollbackIBTPResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackIBTP", ibtp, isSrcChain)
	ret0, _ := ret[0].(*pb.RollbackIBTPResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RollbackIBTP indicates an expected call of RollbackIBTP.
func (mr *MockClientMockRecorder) RollbackIBTP(ibtp, isSrcChain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackIBTP", reflect.TypeOf((*MockClient)(nil).RollbackIBTP), ibtp, isSrcChain)
}

// Start mocks base method.
func (m *MockClient) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockClientMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockClient)(nil).Start))
}

// Stop mocks base method.
func (m *MockClient) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockClientMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockClient)(nil).Stop))
}

// SubmitIBTP mocks base method.
func (m *MockClient) SubmitIBTP(arg0 *pb.IBTP) (*pb.SubmitIBTPResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitIBTP", arg0)
	ret0, _ := ret[0].(*pb.SubmitIBTPResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitIBTP indicates an expected call of SubmitIBTP.
func (mr *MockClientMockRecorder) SubmitIBTP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitIBTP", reflect.TypeOf((*MockClient)(nil).SubmitIBTP), arg0)
}

// Type mocks base method.
func (m *MockClient) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockClientMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockClient)(nil).Type))
}

// MockClientX is a mock of ClientX interface.
type MockClientX struct {
	ctrl     *gomock.Controller
	recorder *MockClientXMockRecorder
}

// MockClientXMockRecorder is the mock recorder for MockClientX.
type MockClientXMockRecorder struct {
	mock *MockClientX
}

// NewMockClientX creates a new mock instance.
func NewMockClientX(ctrl *gomock.Controller) *MockClientX {
	mock := &MockClientX{ctrl: ctrl}
	mock.recorder = &MockClientXMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientX) EXPECT() *MockClientXMockRecorder {
	return m.recorder
}

// CommitCallback mocks base method.
func (m *MockClientX) CommitCallback(ibtp *pb.IBTP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitCallback", ibtp)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitCallback indicates an expected call of CommitCallback.
func (mr *MockClientXMockRecorder) CommitCallback(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitCallback", reflect.TypeOf((*MockClientX)(nil).CommitCallback), ibtp)
}

// GetCallbackMeta mocks base method.
func (m *MockClientX) GetCallbackMeta() (map[string]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCallbackMeta")
	ret0, _ := ret[0].(map[string]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCallbackMeta indicates an expected call of GetCallbackMeta.
func (mr *MockClientXMockRecorder) GetCallbackMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCallbackMeta", reflect.TypeOf((*MockClientX)(nil).GetCallbackMeta))
}

// GetIBTP mocks base method.
func (m *MockClientX) GetIBTP() chan *pb.IBTP {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIBTP")
	ret0, _ := ret[0].(chan *pb.IBTP)
	return ret0
}

// GetIBTP indicates an expected call of GetIBTP.
func (mr *MockClientXMockRecorder) GetIBTP() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIBTP", reflect.TypeOf((*MockClientX)(nil).GetIBTP))
}

// GetInMessage mocks base method.
func (m *MockClientX) GetInMessage(from string, idx uint64) ([][]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInMessage", from, idx)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInMessage indicates an expected call of GetInMessage.
func (mr *MockClientXMockRecorder) GetInMessage(from, idx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInMessage", reflect.TypeOf((*MockClientX)(nil).GetInMessage), from, idx)
}

// GetInMeta mocks base method.
func (m *MockClientX) GetInMeta() (map[string]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInMeta")
	ret0, _ := ret[0].(map[string]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInMeta indicates an expected call of GetInMeta.
func (mr *MockClientXMockRecorder) GetInMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInMeta", reflect.TypeOf((*MockClientX)(nil).GetInMeta))
}

// GetOutMessage mocks base method.
func (m *MockClientX) GetOutMessage(to string, idx uint64) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutMessage", to, idx)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutMessage indicates an expected call of GetOutMessage.
func (mr *MockClientXMockRecorder) GetOutMessage(to, idx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutMessage", reflect.TypeOf((*MockClientX)(nil).GetOutMessage), to, idx)
}

// GetOutMeta mocks base method.
func (m *MockClientX) GetOutMeta() (map[string]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutMeta")
	ret0, _ := ret[0].(map[string]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOutMeta indicates an expected call of GetOutMeta.
func (mr *MockClientXMockRecorder) GetOutMeta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutMeta", reflect.TypeOf((*MockClientX)(nil).GetOutMeta))
}

// GetReceipt mocks base method.
func (m *MockClientX) GetReceipt(ibtp *pb.IBTP) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReceipt", ibtp)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceipt indicates an expected call of GetReceipt.
func (mr *MockClientXMockRecorder) GetReceipt(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceipt", reflect.TypeOf((*MockClientX)(nil).GetReceipt), ibtp)
}

// ID mocks base method.
func (m *MockClientX) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockClientXMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockClientX)(nil).ID))
}

// IncreaseInMeta mocks base method.
func (m *MockClientX) IncreaseInMeta(ibtp *pb.IBTP) (*pb.IBTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseInMeta", ibtp)
	ret0, _ := ret[0].(*pb.IBTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncreaseInMeta indicates an expected call of IncreaseInMeta.
func (mr *MockClientXMockRecorder) IncreaseInMeta(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseInMeta", reflect.TypeOf((*MockClientX)(nil).IncreaseInMeta), ibtp)
}

// Initialize mocks base method.
func (m *MockClientX) Initialize(configPath, pierID string, extra []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", configPath, pierID, extra)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockClientXMockRecorder) Initialize(configPath, pierID, extra interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockClientX)(nil).Initialize), configPath, pierID, extra)
}

// Name mocks base method.
func (m *MockClientX) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockClientXMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockClientX)(nil).Name))
}

// RollbackIBTP mocks base method.
func (m *MockClientX) RollbackIBTP(ibtp *pb.IBTP, isSrcChain bool) (*pb.RollbackIBTPResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackIBTP", ibtp, isSrcChain)
	ret0, _ := ret[0].(*pb.RollbackIBTPResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RollbackIBTP indicates an expected call of RollbackIBTP.
func (mr *MockClientXMockRecorder) RollbackIBTP(ibtp, isSrcChain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackIBTP", reflect.TypeOf((*MockClientX)(nil).RollbackIBTP), ibtp, isSrcChain)
}

// Send mocks base method.
func (m *MockClientX) Send(ibtp *pb.IBTP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ibtp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockClientXMockRecorder) Send(ibtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockClientX)(nil).Send), ibtp)
}

// Start mocks base method.
func (m *MockClientX) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockClientXMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockClientX)(nil).Start))
}

// Stop mocks base method.
func (m *MockClientX) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockClientXMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockClientX)(nil).Stop))
}

// SubmitIBTP mocks base method.
func (m *MockClientX) SubmitIBTP(arg0 *pb.IBTP) (*pb.SubmitIBTPResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitIBTP", arg0)
	ret0, _ := ret[0].(*pb.SubmitIBTPResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitIBTP indicates an expected call of SubmitIBTP.
func (mr *MockClientXMockRecorder) SubmitIBTP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitIBTP", reflect.TypeOf((*MockClientX)(nil).SubmitIBTP), arg0)
}

// Type mocks base method.
func (m *MockClientX) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockClientXMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockClientX)(nil).Type))
}
