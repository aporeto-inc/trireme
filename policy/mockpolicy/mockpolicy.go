// Code generated by MockGen. DO NOT EDIT.
// Source: policy/interfaces.go

// Package mockpolicy is a generated GoMock package.
package mockpolicy

import (
	context "context"
	reflect "reflect"

	nat "github.com/docker/go-connections/nat"
	gomock "github.com/golang/mock/gomock"
	common "go.aporeto.io/trireme-lib/common"
	policy "go.aporeto.io/trireme-lib/policy"
)

// MockRuntimeReader is a mock of RuntimeReader interface
// nolint
type MockRuntimeReader struct {
	ctrl     *gomock.Controller
	recorder *MockRuntimeReaderMockRecorder
}

// MockRuntimeReaderMockRecorder is the mock recorder for MockRuntimeReader
// nolint
type MockRuntimeReaderMockRecorder struct {
	mock *MockRuntimeReader
}

// NewMockRuntimeReader creates a new mock instance
// nolint
func NewMockRuntimeReader(ctrl *gomock.Controller) *MockRuntimeReader {
	mock := &MockRuntimeReader{ctrl: ctrl}
	mock.recorder = &MockRuntimeReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockRuntimeReader) EXPECT() *MockRuntimeReaderMockRecorder {
	return m.recorder
}

// Pid mocks base method
// nolint
func (m *MockRuntimeReader) Pid() int {
	ret := m.ctrl.Call(m, "Pid")
	ret0, _ := ret[0].(int)
	return ret0
}

// Pid indicates an expected call of Pid
// nolint
func (mr *MockRuntimeReaderMockRecorder) Pid() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pid", reflect.TypeOf((*MockRuntimeReader)(nil).Pid))
}

// Name mocks base method
// nolint
func (m *MockRuntimeReader) Name() string {
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
// nolint
func (mr *MockRuntimeReaderMockRecorder) Name() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockRuntimeReader)(nil).Name))
}

// NSPath mocks base method
// nolint
func (m *MockRuntimeReader) NSPath() string {
	ret := m.ctrl.Call(m, "NSPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// NSPath indicates an expected call of NSPath
// nolint
func (mr *MockRuntimeReaderMockRecorder) NSPath() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NSPath", reflect.TypeOf((*MockRuntimeReader)(nil).NSPath))
}

// Tag mocks base method
// nolint
func (m *MockRuntimeReader) Tag(arg0 string) (string, bool) {
	ret := m.ctrl.Call(m, "Tag", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Tag indicates an expected call of Tag
// nolint
func (mr *MockRuntimeReaderMockRecorder) Tag(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockRuntimeReader)(nil).Tag), arg0)
}

// Tags mocks base method
// nolint
func (m *MockRuntimeReader) Tags() *policy.TagStore {
	ret := m.ctrl.Call(m, "Tags")
	ret0, _ := ret[0].(*policy.TagStore)
	return ret0
}

// Tags indicates an expected call of Tags
// nolint
func (mr *MockRuntimeReaderMockRecorder) Tags() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tags", reflect.TypeOf((*MockRuntimeReader)(nil).Tags))
}

// Options mocks base method
// nolint
func (m *MockRuntimeReader) Options() policy.OptionsType {
	ret := m.ctrl.Call(m, "Options")
	ret0, _ := ret[0].(policy.OptionsType)
	return ret0
}

// Options indicates an expected call of Options
// nolint
func (mr *MockRuntimeReaderMockRecorder) Options() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Options", reflect.TypeOf((*MockRuntimeReader)(nil).Options))
}

// IPAddresses mocks base method
// nolint
func (m *MockRuntimeReader) IPAddresses() policy.ExtendedMap {
	ret := m.ctrl.Call(m, "IPAddresses")
	ret0, _ := ret[0].(policy.ExtendedMap)
	return ret0
}

// IPAddresses indicates an expected call of IPAddresses
// nolint
func (mr *MockRuntimeReaderMockRecorder) IPAddresses() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IPAddresses", reflect.TypeOf((*MockRuntimeReader)(nil).IPAddresses))
}

// PUType mocks base method
// nolint
func (m *MockRuntimeReader) PUType() common.PUType {
	ret := m.ctrl.Call(m, "PUType")
	ret0, _ := ret[0].(common.PUType)
	return ret0
}

// PUType indicates an expected call of PUType
// nolint
func (mr *MockRuntimeReaderMockRecorder) PUType() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PUType", reflect.TypeOf((*MockRuntimeReader)(nil).PUType))
}

// SetServices mocks base method
// nolint
func (m *MockRuntimeReader) SetServices(services []common.Service) {
	m.ctrl.Call(m, "SetServices", services)
}

// SetServices indicates an expected call of SetServices
// nolint
func (mr *MockRuntimeReaderMockRecorder) SetServices(services interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetServices", reflect.TypeOf((*MockRuntimeReader)(nil).SetServices), services)
}

// PortMap mocks base method
// nolint
func (m *MockRuntimeReader) PortMap() map[nat.Port][]string {
	ret := m.ctrl.Call(m, "PortMap")
	ret0, _ := ret[0].(map[nat.Port][]string)
	return ret0
}

// PortMap indicates an expected call of PortMap
// nolint
func (mr *MockRuntimeReaderMockRecorder) PortMap() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PortMap", reflect.TypeOf((*MockRuntimeReader)(nil).PortMap))
}

// MockResolver is a mock of Resolver interface
// nolint
type MockResolver struct {
	ctrl     *gomock.Controller
	recorder *MockResolverMockRecorder
}

// MockResolverMockRecorder is the mock recorder for MockResolver
// nolint
type MockResolverMockRecorder struct {
	mock *MockResolver
}

// NewMockResolver creates a new mock instance
// nolint
func NewMockResolver(ctrl *gomock.Controller) *MockResolver {
	mock := &MockResolver{ctrl: ctrl}
	mock.recorder = &MockResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
// nolint
func (m *MockResolver) EXPECT() *MockResolverMockRecorder {
	return m.recorder
}

// HandlePUEvent mocks base method
// nolint
func (m *MockResolver) HandlePUEvent(ctx context.Context, puID string, event common.Event, runtime policy.RuntimeReader) error {
	ret := m.ctrl.Call(m, "HandlePUEvent", ctx, puID, event, runtime)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandlePUEvent indicates an expected call of HandlePUEvent
// nolint
func (mr *MockResolverMockRecorder) HandlePUEvent(ctx, puID, event, runtime interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlePUEvent", reflect.TypeOf((*MockResolver)(nil).HandlePUEvent), ctx, puID, event, runtime)
}
