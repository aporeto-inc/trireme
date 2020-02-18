// Code generated by MockGen. DO NOT EDIT.
// Source: controller/internal/enforcer/nfqdatapath/tokenaccessor/interfaces.go

// Package mocktokenaccessor is a generated GoMock package.
package mocktokenaccessor

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	claimsheader "go.aporeto.io/trireme-lib/controller/pkg/claimsheader"
	connection "go.aporeto.io/trireme-lib/controller/pkg/connection"
	pucontext "go.aporeto.io/trireme-lib/controller/pkg/pucontext"
	secrets "go.aporeto.io/trireme-lib/controller/pkg/secrets"
	tokens "go.aporeto.io/trireme-lib/controller/pkg/tokens"
)

// MockTokenAccessor is a mock of TokenAccessor interface
// nolint
type MockTokenAccessor struct {
	ctrl     *gomock.Controller
	recorder *MockTokenAccessorMockRecorder
}

// MockTokenAccessorMockRecorder is the mock recorder for MockTokenAccessor
// nolint
type MockTokenAccessorMockRecorder struct {
	mock *MockTokenAccessor
}

// NewMockTokenAccessor creates a new mock instance
// nolint
func NewMockTokenAccessor(ctrl *gomock.Controller) *MockTokenAccessor {
	mock := &MockTokenAccessor{ctrl: ctrl}
	mock.recorder = &MockTokenAccessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockTokenAccessor) EXPECT() *MockTokenAccessorMockRecorder {
	return m.recorder
}

// GetTokenValidity mocks base method
// nolint
func (m *MockTokenAccessor) GetTokenValidity() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenValidity")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// GetTokenValidity indicates an expected call of GetTokenValidity
// nolint
func (mr *MockTokenAccessorMockRecorder) GetTokenValidity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenValidity", reflect.TypeOf((*MockTokenAccessor)(nil).GetTokenValidity))
}

// GetTokenServerID mocks base method
// nolint
func (m *MockTokenAccessor) GetTokenServerID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenServerID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTokenServerID indicates an expected call of GetTokenServerID
// nolint
func (mr *MockTokenAccessorMockRecorder) GetTokenServerID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenServerID", reflect.TypeOf((*MockTokenAccessor)(nil).GetTokenServerID))
}

// CreateAckPacketToken mocks base method
// nolint
func (m *MockTokenAccessor) CreateAckPacketToken(context *pucontext.PUContext, auth *connection.AuthInfo, secrets secrets.Secrets) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAckPacketToken", context, auth, secrets)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAckPacketToken indicates an expected call of CreateAckPacketToken
// nolint
func (mr *MockTokenAccessorMockRecorder) CreateAckPacketToken(context, auth, secrets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAckPacketToken", reflect.TypeOf((*MockTokenAccessor)(nil).CreateAckPacketToken), context, auth, secrets)
}

// CreateSynPacketToken mocks base method
// nolint
func (m *MockTokenAccessor) CreateSynPacketToken(context *pucontext.PUContext, auth *connection.AuthInfo, claimsHeader *claimsheader.ClaimsHeader, secrets secrets.Secrets) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSynPacketToken", context, auth, claimsHeader, secrets)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSynPacketToken indicates an expected call of CreateSynPacketToken
// nolint
func (mr *MockTokenAccessorMockRecorder) CreateSynPacketToken(context, auth, claimsHeader, secrets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSynPacketToken", reflect.TypeOf((*MockTokenAccessor)(nil).CreateSynPacketToken), context, auth, claimsHeader, secrets)
}

// CreateSynAckPacketToken mocks base method
// nolint
func (m *MockTokenAccessor) CreateSynAckPacketToken(context *pucontext.PUContext, auth *connection.AuthInfo, claimsHeader *claimsheader.ClaimsHeader, secrets secrets.Secrets) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSynAckPacketToken", context, auth, claimsHeader, secrets)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSynAckPacketToken indicates an expected call of CreateSynAckPacketToken
// nolint
func (mr *MockTokenAccessorMockRecorder) CreateSynAckPacketToken(context, auth, claimsHeader, secrets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSynAckPacketToken", reflect.TypeOf((*MockTokenAccessor)(nil).CreateSynAckPacketToken), context, auth, claimsHeader, secrets)
}

// ParsePacketToken mocks base method
// nolint
func (m *MockTokenAccessor) ParsePacketToken(auth *connection.AuthInfo, data []byte, secrets secrets.Secrets) (*tokens.ConnectionClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParsePacketToken", auth, data, secrets)
	ret0, _ := ret[0].(*tokens.ConnectionClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParsePacketToken indicates an expected call of ParsePacketToken
// nolint
func (mr *MockTokenAccessorMockRecorder) ParsePacketToken(auth, data, secrets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParsePacketToken", reflect.TypeOf((*MockTokenAccessor)(nil).ParsePacketToken), auth, data, secrets)
}

// ParseAckToken mocks base method
// nolint
func (m *MockTokenAccessor) ParseAckToken(auth *connection.AuthInfo, data []byte, secrets secrets.Secrets) (*tokens.ConnectionClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseAckToken", auth, data, secrets)
	ret0, _ := ret[0].(*tokens.ConnectionClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseAckToken indicates an expected call of ParseAckToken
// nolint
func (mr *MockTokenAccessorMockRecorder) ParseAckToken(auth, data, secrets interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseAckToken", reflect.TypeOf((*MockTokenAccessor)(nil).ParseAckToken), auth, data, secrets)
}
