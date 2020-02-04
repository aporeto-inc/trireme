// Code generated by MockGen. DO NOT EDIT.
// Source: controller/pkg/remoteenforcer/internal/counterclient/interfaces.go

package mockcounterclient

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCounterClient is a mock of CounterClient interface
// nolint
type MockCounterClient struct {
	ctrl     *gomock.Controller
	recorder *MockCounterClientMockRecorder
}

// MockCounterClientMockRecorder is the mock recorder for MockCounterClient
// nolint
type MockCounterClientMockRecorder struct {
	mock *MockCounterClient
}

// NewMockCounterClient creates a new mock instance
// nolint
func NewMockCounterClient(ctrl *gomock.Controller) *MockCounterClient {
	mock := &MockCounterClient{ctrl: ctrl}
	mock.recorder = &MockCounterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (_m *MockCounterClient) EXPECT() *MockCounterClientMockRecorder {
	return _m.recorder
}

// Run mocks base method
// nolint
func (_m *MockCounterClient) Run(ctx context.Context) error {
	ret := _m.ctrl.Call(_m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
// nolint
func (_mr *MockCounterClientMockRecorder) Run(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "Run", reflect.TypeOf((*MockCounterClient)(nil).Run), arg0)
}
