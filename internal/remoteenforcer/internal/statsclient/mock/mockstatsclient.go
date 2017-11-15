// Automatically generated by MockGen. DO NOT EDIT!
// Source: internal/remoteenforcer/internal/statsclient/interfaces.go

package mockstatsclient

import (
	gomock "github.com/aporeto-inc/mock/gomock"
)

// Mock of StatsClient interface
// nolint
type MockStatsClient struct {
	ctrl     *gomock.Controller
	recorder *_MockStatsClientRecorder
}

// Recorder for MockStatsClient (not exported)
// nolint
type _MockStatsClientRecorder struct {
	mock *MockStatsClient
}

// nolint
func NewMockStatsClient(ctrl *gomock.Controller) *MockStatsClient {
	mock := &MockStatsClient{ctrl: ctrl}
	mock.recorder = &_MockStatsClientRecorder{mock}
	return mock
}

// nolint
func (_m *MockStatsClient) EXPECT() *_MockStatsClientRecorder {
	return _m.recorder
}

// nolint
func (_m *MockStatsClient) Start() error {
	ret := _m.ctrl.Call(_m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// nolint
func (_mr *_MockStatsClientRecorder) Start() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Start")
}

// nolint
func (_m *MockStatsClient) Stop() {
	_m.ctrl.Call(_m, "Stop")
}

// nolint
func (_mr *_MockStatsClientRecorder) Stop() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Stop")
}
