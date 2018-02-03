// Code generated by MockGen. DO NOT EDIT.
// Source: controller/pkg/remoteenforcer/interfaces.go

// Package mockremoteenforcer is a generated GoMock package.
package mockremoteenforcer

import (
	reflect "reflect"

	rpcwrapper "github.com/aporeto-inc/trireme-lib/controller/internal/enforcer/utils/rpcwrapper"
	gomock "github.com/golang/mock/gomock"
)

// MockRemoteIntf is a mock of RemoteIntf interface
// nolint
type MockRemoteIntf struct {
	ctrl     *gomock.Controller
	recorder *MockRemoteIntfMockRecorder
}

// MockRemoteIntfMockRecorder is the mock recorder for MockRemoteIntf
// nolint
type MockRemoteIntfMockRecorder struct {
	mock *MockRemoteIntf
}

// NewMockRemoteIntf creates a new mock instance
// nolint
func NewMockRemoteIntf(ctrl *gomock.Controller) *MockRemoteIntf {
	mock := &MockRemoteIntf{ctrl: ctrl}
	mock.recorder = &MockRemoteIntfMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockRemoteIntf) EXPECT() *MockRemoteIntfMockRecorder {
	return m.recorder
}

// InitEnforcer mocks base method
// nolint
func (m *MockRemoteIntf) InitEnforcer(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "InitEnforcer", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitEnforcer indicates an expected call of InitEnforcer
// nolint
func (mr *MockRemoteIntfMockRecorder) InitEnforcer(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitEnforcer", reflect.TypeOf((*MockRemoteIntf)(nil).InitEnforcer), req, resp)
}

// InitSupervisor mocks base method
// nolint
func (m *MockRemoteIntf) InitSupervisor(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "InitSupervisor", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// InitSupervisor indicates an expected call of InitSupervisor
// nolint
func (mr *MockRemoteIntfMockRecorder) InitSupervisor(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitSupervisor", reflect.TypeOf((*MockRemoteIntf)(nil).InitSupervisor), req, resp)
}

// Supervise mocks base method
// nolint
func (m *MockRemoteIntf) Supervise(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "Supervise", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Supervise indicates an expected call of Supervise
// nolint
func (mr *MockRemoteIntfMockRecorder) Supervise(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Supervise", reflect.TypeOf((*MockRemoteIntf)(nil).Supervise), req, resp)
}

// Unenforce mocks base method
// nolint
func (m *MockRemoteIntf) Unenforce(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "Unenforce", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unenforce indicates an expected call of Unenforce
// nolint
func (mr *MockRemoteIntfMockRecorder) Unenforce(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unenforce", reflect.TypeOf((*MockRemoteIntf)(nil).Unenforce), req, resp)
}

// Unsupervise mocks base method
// nolint
func (m *MockRemoteIntf) Unsupervise(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "Unsupervise", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsupervise indicates an expected call of Unsupervise
// nolint
func (mr *MockRemoteIntfMockRecorder) Unsupervise(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsupervise", reflect.TypeOf((*MockRemoteIntf)(nil).Unsupervise), req, resp)
}

// Enforce mocks base method
// nolint
func (m *MockRemoteIntf) Enforce(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "Enforce", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enforce indicates an expected call of Enforce
// nolint
func (mr *MockRemoteIntfMockRecorder) Enforce(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enforce", reflect.TypeOf((*MockRemoteIntf)(nil).Enforce), req, resp)
}

// EnforcerExit mocks base method
// nolint
func (m *MockRemoteIntf) EnforcerExit(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	ret := m.ctrl.Call(m, "EnforcerExit", req, resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnforcerExit indicates an expected call of EnforcerExit
// nolint
func (mr *MockRemoteIntfMockRecorder) EnforcerExit(req, resp interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnforcerExit", reflect.TypeOf((*MockRemoteIntf)(nil).EnforcerExit), req, resp)
}
