// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mocktrireme is a generated GoMock package.
package mocktrireme

import (
	reflect "reflect"

	constants "github.com/aporeto-inc/trireme-lib/constants"
	secrets "github.com/aporeto-inc/trireme-lib/enforcer/utils/secrets"
	supervisor "github.com/aporeto-inc/trireme-lib/internal/supervisor"
	policy "github.com/aporeto-inc/trireme-lib/policy"
	events "github.com/aporeto-inc/trireme-lib/rpc/events"
	gomock "github.com/golang/mock/gomock"
)

// MockTrireme is a mock of Trireme interface
// nolint
type MockTrireme struct {
	ctrl     *gomock.Controller
	recorder *MockTriremeMockRecorder
}

// MockTriremeMockRecorder is the mock recorder for MockTrireme
// nolint
type MockTriremeMockRecorder struct {
	mock *MockTrireme
}

// NewMockTrireme creates a new mock instance
// nolint
func NewMockTrireme(ctrl *gomock.Controller) *MockTrireme {
	mock := &MockTrireme{ctrl: ctrl}
	mock.recorder = &MockTriremeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockTrireme) EXPECT() *MockTriremeMockRecorder {
	return m.recorder
}

// PURuntime mocks base method
// nolint
func (m *MockTrireme) PURuntime(contextID string) (policy.RuntimeReader, error) {
	ret := m.ctrl.Call(m, "PURuntime", contextID)
	ret0, _ := ret[0].(policy.RuntimeReader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PURuntime indicates an expected call of PURuntime
// nolint
func (mr *MockTriremeMockRecorder) PURuntime(contextID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PURuntime", reflect.TypeOf((*MockTrireme)(nil).PURuntime), contextID)
}

// Start mocks base method
// nolint
func (m *MockTrireme) Start() error {
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
// nolint
func (mr *MockTriremeMockRecorder) Start() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTrireme)(nil).Start))
}

// Stop mocks base method
// nolint
func (m *MockTrireme) Stop() error {
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
// nolint
func (mr *MockTriremeMockRecorder) Stop() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTrireme)(nil).Stop))
}

// Supervisor mocks base method
// nolint
func (m *MockTrireme) Supervisor(kind constants.PUType) supervisor.Supervisor {
	ret := m.ctrl.Call(m, "Supervisor", kind)
	ret0, _ := ret[0].(supervisor.Supervisor)
	return ret0
}

// Supervisor indicates an expected call of Supervisor
// nolint
func (mr *MockTriremeMockRecorder) Supervisor(kind interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Supervisor", reflect.TypeOf((*MockTrireme)(nil).Supervisor), kind)
}

// CreatePURuntime mocks base method
// nolint
func (m *MockTrireme) CreatePURuntime(contextID string, runtimeInfo *policy.PURuntime) error {
	ret := m.ctrl.Call(m, "CreatePURuntime", contextID, runtimeInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePURuntime indicates an expected call of CreatePURuntime
// nolint
func (mr *MockTriremeMockRecorder) CreatePURuntime(contextID, runtimeInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePURuntime", reflect.TypeOf((*MockTrireme)(nil).CreatePURuntime), contextID, runtimeInfo)
}

// HandlePUEvent mocks base method
// nolint
func (m *MockTrireme) HandlePUEvent(contextID string, event events.Event) error {
	ret := m.ctrl.Call(m, "HandlePUEvent", contextID, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandlePUEvent indicates an expected call of HandlePUEvent
// nolint
func (mr *MockTriremeMockRecorder) HandlePUEvent(contextID, event interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlePUEvent", reflect.TypeOf((*MockTrireme)(nil).HandlePUEvent), contextID, event)
}

// UpdatePolicy mocks base method
// nolint
func (m *MockTrireme) UpdatePolicy(contextID string, policy *policy.PUPolicy) error {
	ret := m.ctrl.Call(m, "UpdatePolicy", contextID, policy)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePolicy indicates an expected call of UpdatePolicy
// nolint
func (mr *MockTriremeMockRecorder) UpdatePolicy(contextID, policy interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePolicy", reflect.TypeOf((*MockTrireme)(nil).UpdatePolicy), contextID, policy)
}

// UpdateSecrets mocks base method
// nolint
func (m *MockTrireme) UpdateSecrets(secrets secrets.Secrets) error {
	ret := m.ctrl.Call(m, "UpdateSecrets", secrets)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecrets indicates an expected call of UpdateSecrets
// nolint
func (mr *MockTriremeMockRecorder) UpdateSecrets(secrets interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecrets", reflect.TypeOf((*MockTrireme)(nil).UpdateSecrets), secrets)
}

// MockPolicyUpdater is a mock of PolicyUpdater interface
// nolint
type MockPolicyUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyUpdaterMockRecorder
}

// MockPolicyUpdaterMockRecorder is the mock recorder for MockPolicyUpdater
// nolint
type MockPolicyUpdaterMockRecorder struct {
	mock *MockPolicyUpdater
}

// NewMockPolicyUpdater creates a new mock instance
// nolint
func NewMockPolicyUpdater(ctrl *gomock.Controller) *MockPolicyUpdater {
	mock := &MockPolicyUpdater{ctrl: ctrl}
	mock.recorder = &MockPolicyUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockPolicyUpdater) EXPECT() *MockPolicyUpdaterMockRecorder {
	return m.recorder
}

// UpdatePolicy mocks base method
// nolint
func (m *MockPolicyUpdater) UpdatePolicy(contextID string, policy *policy.PUPolicy) error {
	ret := m.ctrl.Call(m, "UpdatePolicy", contextID, policy)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePolicy indicates an expected call of UpdatePolicy
// nolint
func (mr *MockPolicyUpdaterMockRecorder) UpdatePolicy(contextID, policy interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePolicy", reflect.TypeOf((*MockPolicyUpdater)(nil).UpdatePolicy), contextID, policy)
}

// MockPolicyResolver is a mock of PolicyResolver interface
// nolint
type MockPolicyResolver struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyResolverMockRecorder
}

// MockPolicyResolverMockRecorder is the mock recorder for MockPolicyResolver
// nolint
type MockPolicyResolverMockRecorder struct {
	mock *MockPolicyResolver
}

// NewMockPolicyResolver creates a new mock instance
// nolint
func NewMockPolicyResolver(ctrl *gomock.Controller) *MockPolicyResolver {
	mock := &MockPolicyResolver{ctrl: ctrl}
	mock.recorder = &MockPolicyResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockPolicyResolver) EXPECT() *MockPolicyResolverMockRecorder {
	return m.recorder
}

// ResolvePolicy mocks base method
// nolint
func (m *MockPolicyResolver) ResolvePolicy(contextID string, RuntimeReader policy.RuntimeReader) (*policy.PUPolicy, error) {
	ret := m.ctrl.Call(m, "ResolvePolicy", contextID, RuntimeReader)
	ret0, _ := ret[0].(*policy.PUPolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolvePolicy indicates an expected call of ResolvePolicy
// nolint
func (mr *MockPolicyResolverMockRecorder) ResolvePolicy(contextID, RuntimeReader interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolvePolicy", reflect.TypeOf((*MockPolicyResolver)(nil).ResolvePolicy), contextID, RuntimeReader)
}

// HandlePUEvent mocks base method
// nolint
func (m *MockPolicyResolver) HandlePUEvent(contextID string, eventType events.Event) {
	m.ctrl.Call(m, "HandlePUEvent", contextID, eventType)
}

// HandlePUEvent indicates an expected call of HandlePUEvent
// nolint
func (mr *MockPolicyResolverMockRecorder) HandlePUEvent(contextID, eventType interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlePUEvent", reflect.TypeOf((*MockPolicyResolver)(nil).HandlePUEvent), contextID, eventType)
}

// MockSecretsUpdater is a mock of SecretsUpdater interface
// nolint
type MockSecretsUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsUpdaterMockRecorder
}

// MockSecretsUpdaterMockRecorder is the mock recorder for MockSecretsUpdater
// nolint
type MockSecretsUpdaterMockRecorder struct {
	mock *MockSecretsUpdater
}

// NewMockSecretsUpdater creates a new mock instance
// nolint
func NewMockSecretsUpdater(ctrl *gomock.Controller) *MockSecretsUpdater {
	mock := &MockSecretsUpdater{ctrl: ctrl}
	mock.recorder = &MockSecretsUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockSecretsUpdater) EXPECT() *MockSecretsUpdaterMockRecorder {
	return m.recorder
}

// UpdateSecrets mocks base method
// nolint
func (m *MockSecretsUpdater) UpdateSecrets(secrets secrets.Secrets) error {
	ret := m.ctrl.Call(m, "UpdateSecrets", secrets)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecrets indicates an expected call of UpdateSecrets
// nolint
func (mr *MockSecretsUpdaterMockRecorder) UpdateSecrets(secrets interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecrets", reflect.TypeOf((*MockSecretsUpdater)(nil).UpdateSecrets), secrets)
}
