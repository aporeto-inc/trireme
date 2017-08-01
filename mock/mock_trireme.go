// Automatically generated by MockGen. DO NOT EDIT!
// Source: ./monitor/monitor.go

// nolint
package mock_trireme

import (
	"github.com/aporeto-inc/trireme/monitor"
	policy "github.com/aporeto-inc/trireme/policy"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Monitor interface
type MockMonitor struct {
	ctrl     *gomock.Controller
	recorder *_MockMonitorRecorder
}

// Recorder for MockMonitor (not exported)
type _MockMonitorRecorder struct {
	mock *MockMonitor
}

func NewMockMonitor(ctrl *gomock.Controller) *MockMonitor {
	mock := &MockMonitor{ctrl: ctrl}
	mock.recorder = &_MockMonitorRecorder{mock}
	return mock
}

func (_m *MockMonitor) EXPECT() *_MockMonitorRecorder {
	return _m.recorder
}

func (_m *MockMonitor) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockMonitorRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

func (_m *MockMonitor) Stop() error {
	ret := _m.ctrl.Call(_m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockMonitorRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}

// Mock of ProcessingUnitsHandler interface
type MockProcessingUnitsHandler struct {
	ctrl     *gomock.Controller
	recorder *_MockProcessingUnitsHandlerRecorder
}

// Recorder for MockProcessingUnitsHandler (not exported)
type _MockProcessingUnitsHandlerRecorder struct {
	mock *MockProcessingUnitsHandler
}

func NewMockProcessingUnitsHandler(ctrl *gomock.Controller) *MockProcessingUnitsHandler {
	mock := &MockProcessingUnitsHandler{ctrl: ctrl}
	mock.recorder = &_MockProcessingUnitsHandlerRecorder{mock}
	return mock
}

func (_m *MockProcessingUnitsHandler) EXPECT() *_MockProcessingUnitsHandlerRecorder {
	return _m.recorder
}

func (_m *MockProcessingUnitsHandler) SetPURuntime(contextID string, runtimeInfo *policy.PURuntime) error {

	ret := _m.ctrl.Call(_m, "SetPURuntime", contextID, runtimeInfo)

	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProcessingUnitsHandlerRecorder) SetPURuntime(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetPURuntime", arg0, arg1)
}

func (_m *MockProcessingUnitsHandler) HandlePUEvent(contextID string, event monitor.Event) error {

	ret := _m.ctrl.Call(_m, "HandlePUEvent", contextID, event)

	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockProcessingUnitsHandlerRecorder) HandlePUEvent(arg0, arg1 interface{}) *gomock.Call {

	return _mr.mock.ctrl.RecordCall(_mr.mock, "HandlePUEvent", arg0, arg1)
}
