// Code generated by MockGen. DO NOT EDIT.
// Source: controller/interfaces.go

// Package mockcontroller is a generated GoMock package.
package mockcontroller

import (
	context "context"
	common "github.com/aporeto-inc/trireme-lib/common"
	secrets "github.com/aporeto-inc/trireme-lib/controller/internal/enforcer/utils/secrets"
	policy "github.com/aporeto-inc/trireme-lib/policy"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTriremeController is a mock of TriremeController interface
type MockTriremeController struct {
	ctrl     *gomock.Controller
	recorder *MockTriremeControllerMockRecorder
}

// MockTriremeControllerMockRecorder is the mock recorder for MockTriremeController
type MockTriremeControllerMockRecorder struct {
	mock *MockTriremeController
}

// NewMockTriremeController creates a new mock instance
func NewMockTriremeController(ctrl *gomock.Controller) *MockTriremeController {
	mock := &MockTriremeController{ctrl: ctrl}
	mock.recorder = &MockTriremeControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTriremeController) EXPECT() *MockTriremeControllerMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockTriremeController) Run(ctx context.Context) error {
	ret := m.ctrl.Call(m, "Run", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run
func (mr *MockTriremeControllerMockRecorder) Run(ctx interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockTriremeController)(nil).Run), ctx)
}

// ProcessEvent mocks base method
func (m *MockTriremeController) ProcessEvent(ctx context.Context, event common.Event, id string, policy *policy.PUPolicy, runtime *policy.PURuntime) error {
	ret := m.ctrl.Call(m, "ProcessEvent", ctx, event, id, policy, runtime)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessEvent indicates an expected call of ProcessEvent
func (mr *MockTriremeControllerMockRecorder) ProcessEvent(ctx, event, id, policy, runtime interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessEvent", reflect.TypeOf((*MockTriremeController)(nil).ProcessEvent), ctx, event, id, policy, runtime)
}

// UpdatePolicy mocks base method
func (m *MockTriremeController) UpdatePolicy(contextID string, policy *policy.PUPolicy, runtime *policy.PURuntime) error {
	ret := m.ctrl.Call(m, "UpdatePolicy", contextID, policy, runtime)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePolicy indicates an expected call of UpdatePolicy
func (mr *MockTriremeControllerMockRecorder) UpdatePolicy(contextID, policy, runtime interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePolicy", reflect.TypeOf((*MockTriremeController)(nil).UpdatePolicy), contextID, policy, runtime)
}

// UpdateSecrets mocks base method
func (m *MockTriremeController) UpdateSecrets(secrets secrets.Secrets) error {
	ret := m.ctrl.Call(m, "UpdateSecrets", secrets)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSecrets indicates an expected call of UpdateSecrets
func (mr *MockTriremeControllerMockRecorder) UpdateSecrets(secrets interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecrets", reflect.TypeOf((*MockTriremeController)(nil).UpdateSecrets), secrets)
}
