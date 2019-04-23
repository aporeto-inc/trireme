// Code generated by MockGen. DO NOT EDIT.
// Source: controller/internal/supervisor/interfaces.go

// Package mocksupervisor is a generated GoMock package.
package mocksupervisor

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	common "go.aporeto.io/trireme-lib/common"
	aclprovider "go.aporeto.io/trireme-lib/controller/pkg/aclprovider"
	runtime "go.aporeto.io/trireme-lib/controller/runtime"
	policy "go.aporeto.io/trireme-lib/policy"
)

// MockSupervisor is a mock of Supervisor interface
// nolint
type MockSupervisor struct {
	ctrl     *gomock.Controller
	recorder *MockSupervisorMockRecorder
}

// MockSupervisorMockRecorder is the mock recorder for MockSupervisor
// nolint
type MockSupervisorMockRecorder struct {
	mock *MockSupervisor
}

// NewMockSupervisor creates a new mock instance
// nolint
func NewMockSupervisor(ctrl *gomock.Controller) *MockSupervisor {
	mock := &MockSupervisor{ctrl: ctrl}
	mock.recorder = &MockSupervisorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockSupervisor) EXPECT() *MockSupervisorMockRecorder {
	return m.recorder
}

// Supervise mocks base method
// nolint
func (m *MockSupervisor) Supervise(contextID string, puInfo *policy.PUInfo) error {
	ret := m.ctrl.Call(m, "Supervise", contextID, puInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Supervise indicates an expected call of Supervise
// nolint
func (mr *MockSupervisorMockRecorder) Supervise(contextID, puInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Supervise", reflect.TypeOf((*MockSupervisor)(nil).Supervise), contextID, puInfo)
}

// Unsupervise mocks base method
// nolint
func (m *MockSupervisor) Unsupervise(contextID string) error {
	ret := m.ctrl.Call(m, "Unsupervise", contextID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsupervise indicates an expected call of Unsupervise
// nolint
func (mr *MockSupervisorMockRecorder) Unsupervise(contextID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsupervise", reflect.TypeOf((*MockSupervisor)(nil).Unsupervise), contextID)
}

// Run mocks base method
// nolint
func (m *MockSupervisor) Run(ctx context.Context) error {
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
// nolint
func (mr *MockSupervisorMockRecorder) Run(ctx interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockSupervisor)(nil).Run), ctx)
}

// SetTargetNetworks mocks base method
// nolint
func (m *MockSupervisor) SetTargetNetworks(cfg *runtime.Configuration) error {
	ret := m.ctrl.Call(m, "SetTargetNetworks", cfg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTargetNetworks indicates an expected call of SetTargetNetworks
// nolint
func (mr *MockSupervisorMockRecorder) SetTargetNetworks(cfg interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTargetNetworks", reflect.TypeOf((*MockSupervisor)(nil).SetTargetNetworks), cfg)
}

// CleanUp mocks base method
// nolint
func (m *MockSupervisor) CleanUp() error {
	ret := m.ctrl.Call(m, "CleanUp")
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanUp indicates an expected call of CleanUp
// nolint
func (mr *MockSupervisorMockRecorder) CleanUp() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanUp", reflect.TypeOf((*MockSupervisor)(nil).CleanUp))
}

// EnableIPTablesPacketTracing mocks base method
// nolint
func (m *MockSupervisor) EnableIPTablesPacketTracing(ctx context.Context, contextID string, interval time.Duration) error {
	ret := m.ctrl.Call(m, "EnableIPTablesPacketTracing", ctx, contextID, interval)
	ret0, _ := ret[0].(error)
	return ret0
}

// EnableIPTablesPacketTracing indicates an expected call of EnableIPTablesPacketTracing
// nolint
func (mr *MockSupervisorMockRecorder) EnableIPTablesPacketTracing(ctx, contextID, interval interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableIPTablesPacketTracing", reflect.TypeOf((*MockSupervisor)(nil).EnableIPTablesPacketTracing), ctx, contextID, interval)
}

// MockImplementor is a mock of Implementor interface
// nolint
type MockImplementor struct {
	ctrl     *gomock.Controller
	recorder *MockImplementorMockRecorder
}

// MockImplementorMockRecorder is the mock recorder for MockImplementor
// nolint
type MockImplementorMockRecorder struct {
	mock *MockImplementor
}

// NewMockImplementor creates a new mock instance
// nolint
func NewMockImplementor(ctrl *gomock.Controller) *MockImplementor {
	mock := &MockImplementor{ctrl: ctrl}
	mock.recorder = &MockImplementorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockImplementor) EXPECT() *MockImplementorMockRecorder {
	return m.recorder
}

// ConfigureRules mocks base method
// nolint
func (m *MockImplementor) ConfigureRules(version int, contextID string, containerInfo *policy.PUInfo) error {
	ret := m.ctrl.Call(m, "ConfigureRules", version, contextID, containerInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfigureRules indicates an expected call of ConfigureRules
// nolint
func (mr *MockImplementorMockRecorder) ConfigureRules(version, contextID, containerInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigureRules", reflect.TypeOf((*MockImplementor)(nil).ConfigureRules), version, contextID, containerInfo)
}

// UpdateRules mocks base method
// nolint
func (m *MockImplementor) UpdateRules(version int, contextID string, containerInfo, oldContainerInfo *policy.PUInfo) error {
	ret := m.ctrl.Call(m, "UpdateRules", version, contextID, containerInfo, oldContainerInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRules indicates an expected call of UpdateRules
// nolint
func (mr *MockImplementorMockRecorder) UpdateRules(version, contextID, containerInfo, oldContainerInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRules", reflect.TypeOf((*MockImplementor)(nil).UpdateRules), version, contextID, containerInfo, oldContainerInfo)
}

// DeleteRules mocks base method
// nolint
func (m *MockImplementor) DeleteRules(version int, context, tcpPorts, udpPorts, mark, uid, proxyPort string, puType common.PUType) error {
	ret := m.ctrl.Call(m, "DeleteRules", version, context, tcpPorts, udpPorts, mark, uid, proxyPort, puType)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRules indicates an expected call of DeleteRules
// nolint
func (mr *MockImplementorMockRecorder) DeleteRules(version, context, tcpPorts, udpPorts, mark, uid, proxyPort, puType interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRules", reflect.TypeOf((*MockImplementor)(nil).DeleteRules), version, context, tcpPorts, udpPorts, mark, uid, proxyPort, puType)
}

// SetTargetNetworks mocks base method
// nolint
func (m *MockImplementor) SetTargetNetworks(cfg *runtime.Configuration) error {
	ret := m.ctrl.Call(m, "SetTargetNetworks", cfg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTargetNetworks indicates an expected call of SetTargetNetworks
// nolint
func (mr *MockImplementorMockRecorder) SetTargetNetworks(cfg interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTargetNetworks", reflect.TypeOf((*MockImplementor)(nil).SetTargetNetworks), cfg)
}

// Run mocks base method
// nolint
func (m *MockImplementor) Run(ctx context.Context) error {
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
// nolint
func (mr *MockImplementorMockRecorder) Run(ctx interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockImplementor)(nil).Run), ctx)
}

// CleanUp mocks base method
// nolint
func (m *MockImplementor) CleanUp() error {
	ret := m.ctrl.Call(m, "CleanUp")
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanUp indicates an expected call of CleanUp
// nolint
func (mr *MockImplementorMockRecorder) CleanUp() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanUp", reflect.TypeOf((*MockImplementor)(nil).CleanUp))
}

// ACLProvider mocks base method
// nolint
func (m *MockImplementor) ACLProvider() []aclprovider.IptablesProvider {
	ret := m.ctrl.Call(m, "ACLProvider")
	ret0, _ := ret[0].(aclprovider.IptablesProvider)

	ret = m.ctrl.Call(m, "ACLProvider")
	ret1, _ := ret[0].(aclprovider.IptablesProvider)

	return []aclprovider.IptablesProvider{ret0, ret1}
}

// ACLProvider indicates an expected call of ACLProvider
// nolint
func (mr *MockImplementorMockRecorder) ACLProvider() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ACLProvider", reflect.TypeOf((*MockImplementor)(nil).ACLProvider))
}
