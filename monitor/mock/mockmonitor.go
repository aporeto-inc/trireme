// Code generated by MockGen. DO NOT EDIT.
// Source: monitor/interfaces.go

// Package mockmonitor is a generated GoMock package.
package mockmonitor

import (
	reflect "reflect"

	monitor "github.com/aporeto-inc/trireme-lib/monitor"
	policy "github.com/aporeto-inc/trireme-lib/policy"
	gomock "github.com/golang/mock/gomock"
)

// MockMonitor is a mock of Monitor interface
// nolint
type MockMonitor struct {
	ctrl     *gomock.Controller
	recorder *MockMonitorMockRecorder
}

// MockMonitorMockRecorder is the mock recorder for MockMonitor
// nolint
type MockMonitorMockRecorder struct {
	mock *MockMonitor
}

// NewMockMonitor creates a new mock instance
// nolint
func NewMockMonitor(ctrl *gomock.Controller) *MockMonitor {
	mock := &MockMonitor{ctrl: ctrl}
	mock.recorder = &MockMonitorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockMonitor) EXPECT() *MockMonitorMockRecorder {
	return m.recorder
}

// SetupHandlers mocks base method
// nolint
func (m *MockMonitor) SetupHandlers(puHandler monitor.ProcessingUnitsHandler, syncHandler monitor.SynchronizationHandler) {
	m.ctrl.Call(m, "SetupHandlers", puHandler, syncHandler)
}

// SetupHandlers indicates an expected call of SetupHandlers
// nolint
func (mr *MockMonitorMockRecorder) SetupHandlers(puHandler, syncHandler interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetupHandlers", reflect.TypeOf((*MockMonitor)(nil).SetupHandlers), puHandler, syncHandler)
}

// Start mocks base method
// nolint
func (m *MockMonitor) Start() error {
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
// nolint
func (mr *MockMonitorMockRecorder) Start() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockMonitor)(nil).Start))
}

// Stop mocks base method
// nolint
func (m *MockMonitor) Stop() error {
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
// nolint
func (mr *MockMonitorMockRecorder) Stop() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockMonitor)(nil).Stop))
}

// MockProcessingUnitsHandler is a mock of ProcessingUnitsHandler interface
// nolint
type MockProcessingUnitsHandler struct {
	ctrl     *gomock.Controller
	recorder *MockProcessingUnitsHandlerMockRecorder
}

// MockProcessingUnitsHandlerMockRecorder is the mock recorder for MockProcessingUnitsHandler
// nolint
type MockProcessingUnitsHandlerMockRecorder struct {
	mock *MockProcessingUnitsHandler
}

// NewMockProcessingUnitsHandler creates a new mock instance
// nolint
func NewMockProcessingUnitsHandler(ctrl *gomock.Controller) *MockProcessingUnitsHandler {
	mock := &MockProcessingUnitsHandler{ctrl: ctrl}
	mock.recorder = &MockProcessingUnitsHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockProcessingUnitsHandler) EXPECT() *MockProcessingUnitsHandlerMockRecorder {
	return m.recorder
}

// CreatePURuntime mocks base method
// nolint
func (m *MockProcessingUnitsHandler) CreatePURuntime(contextID string, runtimeInfo *policy.PURuntime) error {
	ret := m.ctrl.Call(m, "CreatePURuntime", contextID, runtimeInfo)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePURuntime indicates an expected call of CreatePURuntime
// nolint
func (mr *MockProcessingUnitsHandlerMockRecorder) CreatePURuntime(contextID, runtimeInfo interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePURuntime", reflect.TypeOf((*MockProcessingUnitsHandler)(nil).CreatePURuntime), contextID, runtimeInfo)
}

// HandlePUEvent mocks base method
// nolint
func (m *MockProcessingUnitsHandler) HandlePUEvent(contextID string, event monitor.Event) error {
	ret := m.ctrl.Call(m, "HandlePUEvent", contextID, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandlePUEvent indicates an expected call of HandlePUEvent
// nolint
func (mr *MockProcessingUnitsHandlerMockRecorder) HandlePUEvent(contextID, event interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlePUEvent", reflect.TypeOf((*MockProcessingUnitsHandler)(nil).HandlePUEvent), contextID, event)
}

// MockSynchronizationHandler is a mock of SynchronizationHandler interface
// nolint
type MockSynchronizationHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSynchronizationHandlerMockRecorder
}

// MockSynchronizationHandlerMockRecorder is the mock recorder for MockSynchronizationHandler
// nolint
type MockSynchronizationHandlerMockRecorder struct {
	mock *MockSynchronizationHandler
}

// NewMockSynchronizationHandler creates a new mock instance
// nolint
func NewMockSynchronizationHandler(ctrl *gomock.Controller) *MockSynchronizationHandler {
	mock := &MockSynchronizationHandler{ctrl: ctrl}
	mock.recorder = &MockSynchronizationHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockSynchronizationHandler) EXPECT() *MockSynchronizationHandlerMockRecorder {
	return m.recorder
}

// HandleSynchronization mocks base method
// nolint
func (m *MockSynchronizationHandler) HandleSynchronization(contextID string, state monitor.State, RuntimeReader policy.RuntimeReader, syncType monitor.SynchronizationType) error {
	ret := m.ctrl.Call(m, "HandleSynchronization", contextID, state, RuntimeReader, syncType)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleSynchronization indicates an expected call of HandleSynchronization
// nolint
func (mr *MockSynchronizationHandlerMockRecorder) HandleSynchronization(contextID, state, RuntimeReader, syncType interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSynchronization", reflect.TypeOf((*MockSynchronizationHandler)(nil).HandleSynchronization), contextID, state, RuntimeReader, syncType)
}

// HandleSynchronizationComplete mocks base method
// nolint
func (m *MockSynchronizationHandler) HandleSynchronizationComplete(syncType monitor.SynchronizationType) {
	m.ctrl.Call(m, "HandleSynchronizationComplete", syncType)
}

// HandleSynchronizationComplete indicates an expected call of HandleSynchronizationComplete
// nolint
func (mr *MockSynchronizationHandlerMockRecorder) HandleSynchronizationComplete(syncType interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSynchronizationComplete", reflect.TypeOf((*MockSynchronizationHandler)(nil).HandleSynchronizationComplete), syncType)
}
