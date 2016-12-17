// Automatically generated by MockGen. DO NOT EDIT!
// Source: supervisor/iptablesutils/iptablesutils.go

package mockuptablesutils

import (
	gomock "github.com/aporeto-inc/mock/gomock"
	policy "github.com/aporeto-inc/trireme/policy"
)

// Mock of IptableCommon interface
type MockIptableCommon struct {
	ctrl     *gomock.Controller
	recorder *_MockIptableCommonRecorder
}

// Recorder for MockIptableCommon (not exported)
type _MockIptableCommonRecorder struct {
	mock *MockIptableCommon
}

func NewMockIptableCommon(ctrl *gomock.Controller) *MockIptableCommon {
	mock := &MockIptableCommon{ctrl: ctrl}
	mock.recorder = &_MockIptableCommonRecorder{mock}
	return mock
}

func (_m *MockIptableCommon) EXPECT() *_MockIptableCommonRecorder {
	return _m.recorder
}

func (_m *MockIptableCommon) AppChainPrefix(contextID string, index int) string {
	ret := _m.ctrl.Call(_m, "AppChainPrefix", contextID, index)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockIptableCommonRecorder) AppChainPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AppChainPrefix", arg0, arg1)
}

func (_m *MockIptableCommon) NetChainPrefix(contextID string, index int) string {
	ret := _m.ctrl.Call(_m, "NetChainPrefix", contextID, index)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockIptableCommonRecorder) NetChainPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NetChainPrefix", arg0, arg1)
}

func (_m *MockIptableCommon) DefaultCacheIP(ips []string) (string, error) {
	ret := _m.ctrl.Call(_m, "DefaultCacheIP", ips)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIptableCommonRecorder) DefaultCacheIP(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DefaultCacheIP", arg0)
}

func (_m *MockIptableCommon) chainRules(appChain string, netChain string, ip string) [][]string {
	ret := _m.ctrl.Call(_m, "chainRules", appChain, netChain, ip)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIptableCommonRecorder) chainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "chainRules", arg0, arg1, arg2)
}

func (_m *MockIptableCommon) trapRules(appChain string, netChain string, network string, appQueue string, netQueue string) [][]string {
	ret := _m.ctrl.Call(_m, "trapRules", appChain, netChain, network, appQueue, netQueue)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIptableCommonRecorder) trapRules(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "trapRules", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockIptableCommon) CleanACLs() error {
	ret := _m.ctrl.Call(_m, "CleanACLs")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableCommonRecorder) CleanACLs() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CleanACLs")
}

// Mock of IptableProviderUtils interface
type MockIptableProviderUtils struct {
	ctrl     *gomock.Controller
	recorder *_MockIptableProviderUtilsRecorder
}

// Recorder for MockIptableProviderUtils (not exported)
type _MockIptableProviderUtilsRecorder struct {
	mock *MockIptableProviderUtils
}

func NewMockIptableProviderUtils(ctrl *gomock.Controller) *MockIptableProviderUtils {
	mock := &MockIptableProviderUtils{ctrl: ctrl}
	mock.recorder = &_MockIptableProviderUtilsRecorder{mock}
	return mock
}

func (_m *MockIptableProviderUtils) EXPECT() *_MockIptableProviderUtilsRecorder {
	return _m.recorder
}

func (_m *MockIptableProviderUtils) FilterMarkedPackets(mark int) error {
	ret := _m.ctrl.Call(_m, "FilterMarkedPackets", mark)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) FilterMarkedPackets(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FilterMarkedPackets", arg0)
}

func (_m *MockIptableProviderUtils) AddContainerChain(appChain string, netChain string) error {
	ret := _m.ctrl.Call(_m, "AddContainerChain", appChain, netChain)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) AddContainerChain(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddContainerChain", arg0, arg1)
}

func (_m *MockIptableProviderUtils) deleteChain(context string, chain string) error {
	ret := _m.ctrl.Call(_m, "deleteChain", context, chain)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) deleteChain(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "deleteChain", arg0, arg1)
}

func (_m *MockIptableProviderUtils) DeleteAllContainerChains(appChain string, netChain string) error {
	ret := _m.ctrl.Call(_m, "DeleteAllContainerChains", appChain, netChain)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) DeleteAllContainerChains(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAllContainerChains", arg0, arg1)
}

func (_m *MockIptableProviderUtils) AddChainRules(appChain string, netChain string, ip string) error {
	ret := _m.ctrl.Call(_m, "AddChainRules", appChain, netChain, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) AddChainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddChainRules", arg0, arg1, arg2)
}

func (_m *MockIptableProviderUtils) DeleteChainRules(appChain string, netChain string, ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteChainRules", appChain, netChain, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) DeleteChainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteChainRules", arg0, arg1, arg2)
}

func (_m *MockIptableProviderUtils) AddPacketTrap(appChain string, netChain string, ip string, targetNetworks []string, appQueue string, netQueue string) error {
	ret := _m.ctrl.Call(_m, "AddPacketTrap", appChain, netChain, ip, targetNetworks, appQueue, netQueue)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) AddPacketTrap(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddPacketTrap", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockIptableProviderUtils) DeletePacketTrap(appChain string, netChain string, ip string, targetNetworks []string, appQueue string, netQueue string) error {
	ret := _m.ctrl.Call(_m, "DeletePacketTrap", appChain, netChain, ip, targetNetworks, appQueue, netQueue)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) DeletePacketTrap(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeletePacketTrap", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockIptableProviderUtils) AddAppACLs(chain string, ip string, rules []policy.IPRule) error {
	ret := _m.ctrl.Call(_m, "AddAppACLs", chain, ip, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) AddAppACLs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddAppACLs", arg0, arg1, arg2)
}

func (_m *MockIptableProviderUtils) AddNetACLs(chain string, ip string, rules []policy.IPRule) error {
	ret := _m.ctrl.Call(_m, "AddNetACLs", chain, ip, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) AddNetACLs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddNetACLs", arg0, arg1, arg2)
}

func (_m *MockIptableProviderUtils) cleanACLSection(context string, section string, chainPrefix string) {
	_m.ctrl.Call(_m, "cleanACLSection", context, section, chainPrefix)
}

func (_mr *_MockIptableProviderUtilsRecorder) cleanACLSection(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "cleanACLSection", arg0, arg1, arg2)
}

func (_m *MockIptableProviderUtils) exclusionChainRules(ip string) [][]string {
	ret := _m.ctrl.Call(_m, "exclusionChainRules", ip)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) exclusionChainRules(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "exclusionChainRules", arg0)
}

func (_m *MockIptableProviderUtils) AddExclusionChainRules(ip string) error {
	ret := _m.ctrl.Call(_m, "AddExclusionChainRules", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) AddExclusionChainRules(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddExclusionChainRules", arg0)
}

func (_m *MockIptableProviderUtils) DeleteExclusionChainRules(ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteExclusionChainRules", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableProviderUtilsRecorder) DeleteExclusionChainRules(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteExclusionChainRules", arg0)
}

// Mock of IpsetProviderUtils interface
type MockIpsetProviderUtils struct {
	ctrl     *gomock.Controller
	recorder *_MockIpsetProviderUtilsRecorder
}

// Recorder for MockIpsetProviderUtils (not exported)
type _MockIpsetProviderUtilsRecorder struct {
	mock *MockIpsetProviderUtils
}

func NewMockIpsetProviderUtils(ctrl *gomock.Controller) *MockIpsetProviderUtils {
	mock := &MockIpsetProviderUtils{ctrl: ctrl}
	mock.recorder = &_MockIpsetProviderUtilsRecorder{mock}
	return mock
}

func (_m *MockIpsetProviderUtils) EXPECT() *_MockIpsetProviderUtilsRecorder {
	return _m.recorder
}

func (_m *MockIpsetProviderUtils) SetupIpset(name string, ips []string) error {
	ret := _m.ctrl.Call(_m, "SetupIpset", name, ips)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) SetupIpset(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetupIpset", arg0, arg1)
}

func (_m *MockIpsetProviderUtils) AddIpsetOption(ip string) error {
	ret := _m.ctrl.Call(_m, "AddIpsetOption", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) AddIpsetOption(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddIpsetOption", arg0)
}

func (_m *MockIpsetProviderUtils) DeleteIpsetOption(ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteIpsetOption", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) DeleteIpsetOption(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteIpsetOption", arg0)
}

func (_m *MockIpsetProviderUtils) AddAppSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "AddAppSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) AddAppSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddAppSetRule", arg0, arg1)
}

func (_m *MockIpsetProviderUtils) DeleteAppSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteAppSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) DeleteAppSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAppSetRule", arg0, arg1)
}

func (_m *MockIpsetProviderUtils) AddNetSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "AddNetSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) AddNetSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddNetSetRule", arg0, arg1)
}

func (_m *MockIpsetProviderUtils) DeleteNetSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteNetSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) DeleteNetSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteNetSetRule", arg0, arg1)
}

func (_m *MockIpsetProviderUtils) SetupTrapRules(set string, networkQueues string, applicationQueues string) error {
	ret := _m.ctrl.Call(_m, "SetupTrapRules", set, networkQueues, applicationQueues)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) SetupTrapRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetupTrapRules", arg0, arg1, arg2)
}

func (_m *MockIpsetProviderUtils) CreateACLSets(set string, rules []policy.IPRule) error {
	ret := _m.ctrl.Call(_m, "CreateACLSets", set, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) CreateACLSets(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateACLSets", arg0, arg1)
}

func (_m *MockIpsetProviderUtils) DeleteSet(set string) error {
	ret := _m.ctrl.Call(_m, "DeleteSet", set)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) DeleteSet(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteSet", arg0)
}

func (_m *MockIpsetProviderUtils) CleanIPSets() error {
	ret := _m.ctrl.Call(_m, "CleanIPSets")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetProviderUtilsRecorder) CleanIPSets() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CleanIPSets")
}

// Mock of IptableUtils interface
type MockIptableUtils struct {
	ctrl     *gomock.Controller
	recorder *_MockIptableUtilsRecorder
}

// Recorder for MockIptableUtils (not exported)
type _MockIptableUtilsRecorder struct {
	mock *MockIptableUtils
}

func NewMockIptableUtils(ctrl *gomock.Controller) *MockIptableUtils {
	mock := &MockIptableUtils{ctrl: ctrl}
	mock.recorder = &_MockIptableUtilsRecorder{mock}
	return mock
}

func (_m *MockIptableUtils) EXPECT() *_MockIptableUtilsRecorder {
	return _m.recorder
}

func (_m *MockIptableUtils) AppChainPrefix(contextID string, index int) string {
	ret := _m.ctrl.Call(_m, "AppChainPrefix", contextID, index)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AppChainPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AppChainPrefix", arg0, arg1)
}

func (_m *MockIptableUtils) NetChainPrefix(contextID string, index int) string {
	ret := _m.ctrl.Call(_m, "NetChainPrefix", contextID, index)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) NetChainPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NetChainPrefix", arg0, arg1)
}

func (_m *MockIptableUtils) DefaultCacheIP(ips []string) (string, error) {
	ret := _m.ctrl.Call(_m, "DefaultCacheIP", ips)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIptableUtilsRecorder) DefaultCacheIP(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DefaultCacheIP", arg0)
}

func (_m *MockIptableUtils) chainRules(appChain string, netChain string, ip string) [][]string {
	ret := _m.ctrl.Call(_m, "chainRules", appChain, netChain, ip)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) chainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "chainRules", arg0, arg1, arg2)
}

func (_m *MockIptableUtils) trapRules(appChain string, netChain string, network string, appQueue string, netQueue string) [][]string {
	ret := _m.ctrl.Call(_m, "trapRules", appChain, netChain, network, appQueue, netQueue)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) trapRules(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "trapRules", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockIptableUtils) CleanACLs() error {
	ret := _m.ctrl.Call(_m, "CleanACLs")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) CleanACLs() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CleanACLs")
}

func (_m *MockIptableUtils) FilterMarkedPackets(mark int) error {
	ret := _m.ctrl.Call(_m, "FilterMarkedPackets", mark)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) FilterMarkedPackets(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FilterMarkedPackets", arg0)
}

func (_m *MockIptableUtils) AddContainerChain(appChain string, netChain string) error {
	ret := _m.ctrl.Call(_m, "AddContainerChain", appChain, netChain)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AddContainerChain(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddContainerChain", arg0, arg1)
}

func (_m *MockIptableUtils) deleteChain(context string, chain string) error {
	ret := _m.ctrl.Call(_m, "deleteChain", context, chain)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) deleteChain(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "deleteChain", arg0, arg1)
}

func (_m *MockIptableUtils) DeleteAllContainerChains(appChain string, netChain string) error {
	ret := _m.ctrl.Call(_m, "DeleteAllContainerChains", appChain, netChain)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) DeleteAllContainerChains(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAllContainerChains", arg0, arg1)
}

func (_m *MockIptableUtils) AddChainRules(appChain string, netChain string, ip string) error {
	ret := _m.ctrl.Call(_m, "AddChainRules", appChain, netChain, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AddChainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddChainRules", arg0, arg1, arg2)
}

func (_m *MockIptableUtils) DeleteChainRules(appChain string, netChain string, ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteChainRules", appChain, netChain, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) DeleteChainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteChainRules", arg0, arg1, arg2)
}

func (_m *MockIptableUtils) AddPacketTrap(appChain string, netChain string, ip string, targetNetworks []string, appQueue string, netQueue string) error {
	ret := _m.ctrl.Call(_m, "AddPacketTrap", appChain, netChain, ip, targetNetworks, appQueue, netQueue)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AddPacketTrap(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddPacketTrap", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockIptableUtils) DeletePacketTrap(appChain string, netChain string, ip string, targetNetworks []string, appQueue string, netQueue string) error {
	ret := _m.ctrl.Call(_m, "DeletePacketTrap", appChain, netChain, ip, targetNetworks, appQueue, netQueue)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) DeletePacketTrap(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeletePacketTrap", arg0, arg1, arg2, arg3, arg4, arg5)
}

func (_m *MockIptableUtils) AddAppACLs(chain string, ip string, rules []policy.IPRule) error {
	ret := _m.ctrl.Call(_m, "AddAppACLs", chain, ip, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AddAppACLs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddAppACLs", arg0, arg1, arg2)
}

func (_m *MockIptableUtils) AddNetACLs(chain string, ip string, rules []policy.IPRule) error {
	ret := _m.ctrl.Call(_m, "AddNetACLs", chain, ip, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AddNetACLs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddNetACLs", arg0, arg1, arg2)
}

func (_m *MockIptableUtils) cleanACLSection(context string, section string, chainPrefix string) {
	_m.ctrl.Call(_m, "cleanACLSection", context, section, chainPrefix)
}

func (_mr *_MockIptableUtilsRecorder) cleanACLSection(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "cleanACLSection", arg0, arg1, arg2)
}

func (_m *MockIptableUtils) exclusionChainRules(ip string) [][]string {
	ret := _m.ctrl.Call(_m, "exclusionChainRules", ip)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) exclusionChainRules(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "exclusionChainRules", arg0)
}

func (_m *MockIptableUtils) AddExclusionChainRules(ip string) error {
	ret := _m.ctrl.Call(_m, "AddExclusionChainRules", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) AddExclusionChainRules(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddExclusionChainRules", arg0)
}

func (_m *MockIptableUtils) DeleteExclusionChainRules(ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteExclusionChainRules", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIptableUtilsRecorder) DeleteExclusionChainRules(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteExclusionChainRules", arg0)
}

// Mock of IpsetUtils interface
type MockIpsetUtils struct {
	ctrl     *gomock.Controller
	recorder *_MockIpsetUtilsRecorder
}

// Recorder for MockIpsetUtils (not exported)
type _MockIpsetUtilsRecorder struct {
	mock *MockIpsetUtils
}

func NewMockIpsetUtils(ctrl *gomock.Controller) *MockIpsetUtils {
	mock := &MockIpsetUtils{ctrl: ctrl}
	mock.recorder = &_MockIpsetUtilsRecorder{mock}
	return mock
}

func (_m *MockIpsetUtils) EXPECT() *_MockIpsetUtilsRecorder {
	return _m.recorder
}

func (_m *MockIpsetUtils) AppChainPrefix(contextID string, index int) string {
	ret := _m.ctrl.Call(_m, "AppChainPrefix", contextID, index)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) AppChainPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AppChainPrefix", arg0, arg1)
}

func (_m *MockIpsetUtils) NetChainPrefix(contextID string, index int) string {
	ret := _m.ctrl.Call(_m, "NetChainPrefix", contextID, index)
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) NetChainPrefix(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NetChainPrefix", arg0, arg1)
}

func (_m *MockIpsetUtils) DefaultCacheIP(ips []string) (string, error) {
	ret := _m.ctrl.Call(_m, "DefaultCacheIP", ips)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIpsetUtilsRecorder) DefaultCacheIP(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DefaultCacheIP", arg0)
}

func (_m *MockIpsetUtils) chainRules(appChain string, netChain string, ip string) [][]string {
	ret := _m.ctrl.Call(_m, "chainRules", appChain, netChain, ip)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) chainRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "chainRules", arg0, arg1, arg2)
}

func (_m *MockIpsetUtils) trapRules(appChain string, netChain string, network string, appQueue string, netQueue string) [][]string {
	ret := _m.ctrl.Call(_m, "trapRules", appChain, netChain, network, appQueue, netQueue)
	ret0, _ := ret[0].([][]string)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) trapRules(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "trapRules", arg0, arg1, arg2, arg3, arg4)
}

func (_m *MockIpsetUtils) CleanACLs() error {
	ret := _m.ctrl.Call(_m, "CleanACLs")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) CleanACLs() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CleanACLs")
}

func (_m *MockIpsetUtils) SetupIpset(name string, ips []string) error {
	ret := _m.ctrl.Call(_m, "SetupIpset", name, ips)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) SetupIpset(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetupIpset", arg0, arg1)
}

func (_m *MockIpsetUtils) AddIpsetOption(ip string) error {
	ret := _m.ctrl.Call(_m, "AddIpsetOption", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) AddIpsetOption(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddIpsetOption", arg0)
}

func (_m *MockIpsetUtils) DeleteIpsetOption(ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteIpsetOption", ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) DeleteIpsetOption(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteIpsetOption", arg0)
}

func (_m *MockIpsetUtils) AddAppSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "AddAppSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) AddAppSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddAppSetRule", arg0, arg1)
}

func (_m *MockIpsetUtils) DeleteAppSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteAppSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) DeleteAppSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAppSetRule", arg0, arg1)
}

func (_m *MockIpsetUtils) AddNetSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "AddNetSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) AddNetSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddNetSetRule", arg0, arg1)
}

func (_m *MockIpsetUtils) DeleteNetSetRule(set string, ip string) error {
	ret := _m.ctrl.Call(_m, "DeleteNetSetRule", set, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) DeleteNetSetRule(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteNetSetRule", arg0, arg1)
}

func (_m *MockIpsetUtils) SetupTrapRules(set string, networkQueues string, applicationQueues string) error {
	ret := _m.ctrl.Call(_m, "SetupTrapRules", set, networkQueues, applicationQueues)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) SetupTrapRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetupTrapRules", arg0, arg1, arg2)
}

func (_m *MockIpsetUtils) CreateACLSets(set string, rules []policy.IPRule) error {
	ret := _m.ctrl.Call(_m, "CreateACLSets", set, rules)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) CreateACLSets(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateACLSets", arg0, arg1)
}

func (_m *MockIpsetUtils) DeleteSet(set string) error {
	ret := _m.ctrl.Call(_m, "DeleteSet", set)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) DeleteSet(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteSet", arg0)
}

func (_m *MockIpsetUtils) CleanIPSets() error {
	ret := _m.ctrl.Call(_m, "CleanIPSets")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIpsetUtilsRecorder) CleanIPSets() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CleanIPSets")
}
