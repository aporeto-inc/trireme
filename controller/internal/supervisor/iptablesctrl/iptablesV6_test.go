package iptablesctrl

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aporeto-inc/go-ipset/ipset"
	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/trireme-lib/common"
	"go.aporeto.io/trireme-lib/controller/constants"
	provider "go.aporeto.io/trireme-lib/controller/pkg/aclprovider"
	"go.aporeto.io/trireme-lib/controller/pkg/fqconfig"
	"go.aporeto.io/trireme-lib/controller/runtime"
	"go.aporeto.io/trireme-lib/policy"
	"go.aporeto.io/trireme-lib/utils/portspec"
)

func createTestInstance(ipsv4 provider.IpsetProvider, ipsv6 provider.IpsetProvider, iptv4 provider.IptablesProvider, iptv6 provider.IptablesProvider, mode constants.ModeType) (*Instance, error) {

	ipv4Impl := &ipv4{ipt: iptv4}
	ipv6Impl := &ipv6{ipt: iptv6}

	iptInstanceV4, err := createIPInstance(ipv4Impl, ipsv4, fqconfig.NewFilterQueueWithDefaults(), mode)
	if err != nil {
		return nil, fmt.Errorf("unable to create v4 instance: %s", err)
	}

	iptInstanceV6, err := createIPInstance(ipv6Impl, ipsv6, fqconfig.NewFilterQueueWithDefaults(), mode)
	if err != nil {
		return nil, fmt.Errorf("unable to create ipv6 instance: %s", err)
	}

	iptInstanceV4.conntrackCmd = func([]string) {}
	iptInstanceV6.conntrackCmd = func([]string) {}

	return newInstanceWithProviders(iptInstanceV4, iptInstanceV6)
}

// Fake iptables controller that always returns succes.
type baseIpt struct{}

// Append apends a rule to chain of table
func (b *baseIpt) Append(table, chain string, rulespec ...string) error { return nil }

// Insert inserts a rule to a chain of table at the required pos
func (b *baseIpt) Insert(table, chain string, pos int, rulespec ...string) error { return nil }

// Delete deletes a rule of a chain in the given table
func (b *baseIpt) Delete(table, chain string, rulespec ...string) error { return nil }

// ListChains lists all the chains associated with a table
func (b *baseIpt) ListChains(table string) ([]string, error) { return nil, nil }

// ClearChain clears a chain in a table
func (b *baseIpt) ClearChain(table, chain string) error { return nil }

// DeleteChain deletes a chain in the table. There should be no references to this chain
func (b *baseIpt) DeleteChain(table, chain string) error { return nil }

// NewChain creates a new chain
func (b *baseIpt) NewChain(table, chain string) error { return nil }

// Fake memory IPset that will tell us if we are deleting or installing
// bad things.
type memoryIPSet struct {
	set map[string]struct{}
}

func (m *memoryIPSet) Add(entry string, timeout int) error {
	m.set[entry] = struct{}{}
	return nil
}

func (m *memoryIPSet) AddOption(entry string, option string, timeout int) error {
	return nil
}

func (m *memoryIPSet) Del(entry string) error {
	if _, ok := m.set[entry]; !ok {
		return fmt.Errorf("not found")
	}
	delete(m.set, entry)
	return nil
}

func (m *memoryIPSet) Destroy() error {
	m.set = map[string]struct{}{}
	return nil
}

func (m *memoryIPSet) Flush() error {
	m.set = map[string]struct{}{}
	return nil
}

func (m *memoryIPSet) Test(entry string) (bool, error) {
	_, ok := m.set[entry]
	return ok, nil
}

// Fake IpSetProvider that will use memory and allow us to
// to simulate the system.
type memoryIPSetProvider struct {
	sets map[string]*memoryIPSet
}

func (m *memoryIPSetProvider) NewIpset(name string, hasht string, p *ipset.Params) (provider.Ipset, error) {
	_, ok := m.sets[name]
	if ok {
		return nil, fmt.Errorf("set exists")
	}

	newSet := &memoryIPSet{set: map[string]struct{}{}}
	m.sets[name] = newSet

	return newSet, nil
}

func (m *memoryIPSetProvider) GetIpset(name string) provider.Ipset {
	return m.sets[name]
}

func (m *memoryIPSetProvider) DestroyAll(prefix string) error {
	for set := range m.sets {
		if strings.HasPrefix(set, prefix) {
			delete(m.sets, set)
		}
	}
	return nil
}

func (m *memoryIPSetProvider) ListIPSets() ([]string, error) {
	allSets := []string{}
	for set := range m.sets {
		allSets = append(allSets, set)
	}
	return allSets, nil
}

func TestNewInstanceV4(t *testing.T) {
	Convey("When I create a new iptables instance", t, func() {
		Convey("If I create a remote implemenetation and iptables exists", func() {
			ipsv4 := provider.NewTestIpsetProvider()
			ipsv6 := provider.NewTestIpsetProvider()
			iptv4 := provider.NewTestIptablesProvider()
			iptv6 := provider.NewTestIptablesProvider()

			i, err := createTestInstance(ipsv4, ipsv6, iptv4, iptv6, constants.RemoteContainer)
			Convey("It should succeed", func() {
				So(i, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When I create a new iptables instance", t, func() {
		Convey("If I create a Linux server implemenetation and iptables exists", func() {
			ipsv4 := provider.NewTestIpsetProvider()
			ipsv6 := provider.NewTestIpsetProvider()
			iptv4 := provider.NewTestIptablesProvider()
			iptv6 := provider.NewTestIptablesProvider()

			i, err := createTestInstance(ipsv4, ipsv6, iptv4, iptv6, constants.LocalServer)
			Convey("It should succeed", func() {
				So(i, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})
}

func Test_NegativeConfigureRules(t *testing.T) {
	Convey("Given a valid instance", t, func() {
		ipsv4 := provider.NewTestIpsetProvider()
		ipsv6 := provider.NewTestIpsetProvider()
		iptv4 := provider.NewTestIptablesProvider()
		iptv6 := provider.NewTestIptablesProvider()

		i, err := createTestInstance(ipsv4, ipsv6, iptv4, iptv6, constants.LocalServer)
		So(err, ShouldBeNil)
		cfg := &runtime.Configuration{}
		i.SetTargetNetworks(cfg)

		ipl := policy.ExtendedMap{}
		policyrules := policy.NewPUPolicy("Context",
			policy.Police,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			ipl,
			0,
			nil,
			nil,
			[]string{},
		)
		containerinfo := policy.NewPUInfo("Context", common.ContainerPU)
		containerinfo.Policy = policyrules
		containerinfo.Runtime = policy.NewPURuntimeWithDefaults()
		containerinfo.Runtime.SetOptions(policy.OptionsType{
			CgroupMark: "10",
		})

		Convey("When I configure the rules with no errors, it should succeed", func() {
			i.iptv4.createPUPortSet = func(string) error { return nil }
			i.iptv6.createPUPortSet = func(string) error { return nil }
			err := i.ConfigureRules(1, "ID", containerinfo)
			So(err, ShouldBeNil)
		})

		Convey("When I configure the rules and the port set fails, it should error ", func() {
			i.iptv4.createPUPortSet = func(string) error { return fmt.Errorf("error") }
			i.iptv6.createPUPortSet = func(string) error { return fmt.Errorf("error") }
			err := i.ConfigureRules(1, "ID", containerinfo)
			So(err, ShouldNotBeNil)
		})

		Convey("When I configure the rules and the proxy set fails, it should error", func() {
			ipsv4.MockNewIpset(t, func(name, hash string, p *ipset.Params) (provider.Ipset, error) {
				return nil, fmt.Errorf("error")
			})
			err := i.ConfigureRules(1, "ID", containerinfo)
			So(err, ShouldNotBeNil)
		})

		Convey("When I configure the rules and acls fail, it should error", func() {
			iptv4.MockAppend(t, func(table, chain string, rulespec ...string) error {
				return fmt.Errorf("error")
			})
			err := i.ConfigureRules(1, "ID", containerinfo)
			So(err, ShouldNotBeNil)
		})

		Convey("When I configure the rules and commit fails, it should error", func() {
			iptv4.MockCommit(t, func() error {
				return fmt.Errorf("error")
			})
			err := i.ConfigureRules(1, "ID", containerinfo)
			So(err, ShouldNotBeNil)
		})
	})
}

var (
	expectedGlobalMangleChains = map[string][]string{
		"INPUT": {
			"-m set ! --match-set TRI-v4-Excluded src -j TRI-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-App",
		},
		"TRI-App": {
			"-j TRI-Prx-App",
			"-m mark --mark 1073741922 -j ACCEPT",
			"-m connmark --mark 61166 -j ACCEPT",
			"-j TRI-UID-App",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j MARK --set-mark 99",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 8:11 --queue-bypass",
			"-j TRI-Pid-App",
			"-j TRI-Svc-App",
			"-j TRI-Hst-App",
		},
		"TRI-Net": {
			"-j TRI-Prx-Net",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m string --string n30njxq7bmiwr6dtxq --algo bm --to 65535 -j NFQUEUE --queue-bypass --queue-balance 24:27",
			"-m connmark --mark 61166 -j ACCEPT",
			"-j TRI-UID-Net",
			"-m set --match-set TRI-v4-TargetTCP src -p tcp -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 24:27 --queue-bypass",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-option 34 --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19 --queue-bypass",
			"-j TRI-Pid-Net",
			"-j TRI-Svc-Net",
			"-j TRI-Hst-Net",
		},
		"TRI-Pid-App": {},
		"TRI-Pid-Net": {},
		"TRI-Prx-App": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
		"TRI-Prx-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
		"TRI-Hst-App": {},
		"TRI-Hst-Net": {},
		"TRI-Svc-App": {},
		"TRI-Svc-Net": {},
		"TRI-UID-App": {},
		"TRI-UID-Net": {},
	}

	expectedGlobalNATChains = map[string][]string{
		"PREROUTING": {
			"-p tcp -m addrtype --dst-type LOCAL -m set ! --match-set TRI-v4-Excluded src -j TRI-Redir-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-Redir-App",
		},
		"TRI-Redir-App": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
		"TRI-Redir-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
	}

	expectedGlobalIPSets = map[string][]string{
		"TRI" + "-v4-" + targetTCPNetworkSet: {"0.0.0.0/1", "128.0.0.0/1"},
		"TRI" + "-v4-" + targetUDPNetworkSet: {"10.0.0.0/8"},
		"TRI" + "-v4-" + excludedNetworkSet:  {"127.0.0.1"},
	}

	expectedMangleAfterPUInsert = map[string][]string{
		"INPUT": {
			"-m set ! --match-set TRI-v4-Excluded src -j TRI-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-App",
		},
		"TRI-App": {
			"-j TRI-Prx-App",
			"-m mark --mark 1073741922 -j ACCEPT",
			"-m connmark --mark 61166 -j ACCEPT",
			"-j TRI-UID-App",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j MARK --set-mark 99",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 8:11 --queue-bypass",
			"-j TRI-Pid-App",
			"-j TRI-Svc-App",
			"-j TRI-Hst-App",
		},
		"TRI-Net": {
			"-j TRI-Prx-Net",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m string --string n30njxq7bmiwr6dtxq --algo bm --to 65535 -j NFQUEUE --queue-bypass --queue-balance 24:27",
			"-m connmark --mark 61166 -j ACCEPT",
			"-j TRI-UID-Net",
			"-m set --match-set TRI-v4-TargetTCP src -p tcp -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 24:27 --queue-bypass",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-option 34 --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19 --queue-bypass",
			"-j TRI-Pid-Net",
			"-j TRI-Svc-Net",
			"-j TRI-Hst-Net",
		},
		"TRI-Pid-App": {
			"-m cgroup --cgroup 10 -m comment --comment PU-Chain -j MARK --set-mark 10",
			"-m mark --mark 10 -m comment --comment PU-Chain -j TRI-App-pu1N7uS6--0",
		},
		"TRI-Pid-Net": {
			"-p tcp -m multiport --destination-ports 9000 -m comment --comment PU-Chain -j TRI-Net-pu1N7uS6--0", "-p udp -m multiport --destination-ports 5000 -m comment --comment PU-Chain -j TRI-Net-pu1N7uS6--0",
		},
		"TRI-Prx-App": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m tcp --sport 0 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv src -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst dst,dst -m mark ! --mark 0x40 -j ACCEPT",
		},
		"TRI-Prx-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst src,src -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv src -m addrtype --src-type LOCAL -j ACCEPT",
			"-p tcp -m tcp --dport 0 -j ACCEPT",
		},
		"TRI-Hst-App": {},
		"TRI-Hst-Net": {},
		"TRI-Svc-App": {},
		"TRI-Svc-Net": {},
		"TRI-UID-App": {},
		"TRI-UID-Net": {},

		"TRI-Net-pu1N7uS6--0": {
			"-p UDP -m set --match-set TRI-v4-ext-6zlJIpu19gtV src -m state --state ESTABLISHED -j ACCEPT",
			"-p TCP -m set --match-set TRI-v4-ext-w5frVpu19gtV src -m state --state NEW -m set ! --match-set TRI-v4-TargetTCP src --match multiport --dports 80 -j DROP",
			"-p UDP -m set --match-set TRI-v4-ext-IuSLspu19gtV src --match multiport --dports 443 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-flags SYN,ACK ACK -j NFQUEUE --queue-balance 20:23",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m state --state ESTABLISHED -j NFQUEUE --queue-balance 16:19",
			"-p tcp -m state --state ESTABLISHED -m comment --comment TCP-Established-Connections -j ACCEPT",
			"-s 0.0.0.0/0 -m state --state NEW -j NFLOG --nflog-group 11 --nflog-prefix pu1:default:default6",
			"-s 0.0.0.0/0 -j DROP",
		},

		"TRI-App-pu1N7uS6--0": {
			"-p TCP -m set --match-set TRI-v4-ext-uNdc0pu19gtV dst -m state --state NEW -m set ! --match-set TRI-v4-TargetTCP dst --match multiport --dports 80 -j DROP",
			"-p UDP -m set --match-set TRI-v4-ext-6zlJIpu19gtV dst --match multiport --dports 443 -j ACCEPT",
			"-p icmp -m set --match-set TRI-v4-ext-w5frVpu19gtV dst -j ACCEPT",
			"-p UDP -m set --match-set TRI-v4-ext-IuSLspu19gtV dst -m state --state ESTABLISHED -j ACCEPT",
			"-p tcp -m tcp --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 0:3",
			"-p tcp -m tcp --tcp-flags SYN,ACK ACK -j NFQUEUE --queue-balance 4:7",
			"-p udp -m set --match-set TRI-v4-TargetUDP dst -j NFQUEUE --queue-balance 0:3",
			"-p udp -m set --match-set TRI-v4-TargetUDP dst -m state --state ESTABLISHED -m comment --comment UDP-Established-Connections -j ACCEPT",
			"-p tcp -m state --state ESTABLISHED -m comment --comment TCP-Established-Connections -j ACCEPT",
			"-d 0.0.0.0/0 -m state --state NEW -j NFLOG --nflog-group 10 --nflog-prefix pu1:default:default6",
			"-d 0.0.0.0/0 -j DROP",
		},
	}

	expectedNATAfterPUInsert = map[string][]string{
		"PREROUTING": {
			"-p tcp -m addrtype --dst-type LOCAL -m set ! --match-set TRI-v4-Excluded src -j TRI-Redir-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-Redir-App",
		},
		"TRI-Redir-App": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst dst,dst -m mark ! --mark 0x40 -m cgroup --cgroup 10 -j REDIRECT --to-ports 0",
		},
		"TRI-Redir-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv dst -m mark ! --mark 0x40 -j REDIRECT --to-ports 0",
		},
		"POSTROUTING": {
			"-p udp -m addrtype --src-type LOCAL -m multiport --source-ports 5000 -j ACCEPT",
		},
	}

	expectedIPSetsAfterPUInsert = map[string][]string{
		"TRI" + "-v4-" + targetTCPNetworkSet: {"0.0.0.0/1", "128.0.0.0/1"},
		"TRI" + "-v4-" + targetUDPNetworkSet: {"10.0.0.0/8"},
		"TRI" + "-v4-" + excludedNetworkSet:  {"127.0.0.1"},
		"TRI-v4-ProcPort-pu19gtV":            {},
		"TRI-v4-ext-6zlJIpu19gtV":            {"30.0.0.0/24"},
		"TRI-v4-ext-uNdc0pu19gtV":            {"30.0.0.0/24"},
		"TRI-v4-ext-w5frVpu19gtV":            {"40.0.0.0/24"},
		"TRI-v4-ext-IuSLspu19gtV":            {"40.0.0.0/24"},
		"TRI-v4-Proxy-pu19gtV-dst":           {},
		"TRI-v4-Proxy-pu19gtV-srv":           {},
	}

	expectedMangleAfterPUUpdate = map[string][]string{
		"INPUT": {
			"-m set ! --match-set TRI-v4-Excluded src -j TRI-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-App",
		},
		"TRI-App": {
			"-j TRI-Prx-App",
			"-m mark --mark 1073741922 -j ACCEPT",
			"-m connmark --mark 61166 -j ACCEPT",
			"-j TRI-UID-App",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j MARK --set-mark 99",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 8:11 --queue-bypass",
			"-j TRI-Pid-App",
			"-j TRI-Svc-App",
			"-j TRI-Hst-App",
		},
		"TRI-Net": {
			"-j TRI-Prx-Net",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m string --string n30njxq7bmiwr6dtxq --algo bm --to 65535 -j NFQUEUE --queue-bypass --queue-balance 24:27",
			"-m connmark --mark 61166 -j ACCEPT",
			"-j TRI-UID-Net",
			"-m set --match-set TRI-v4-TargetTCP src -p tcp -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 24:27 --queue-bypass",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-option 34 --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19 --queue-bypass",
			"-j TRI-Pid-Net",
			"-j TRI-Svc-Net",
			"-j TRI-Hst-Net",
		},
		"TRI-Pid-App": {
			"-m cgroup --cgroup 10 -m comment --comment PU-Chain -j MARK --set-mark 10",
			"-m mark --mark 10 -m comment --comment PU-Chain -j TRI-App-pu1N7uS6--1",
		},
		"TRI-Pid-Net": {
			"-p tcp -m set --match-set TRI-v4-ProcPort-pu19gtV dst -m comment --comment PU-Chain -j TRI-Net-pu1N7uS6--1",
		},
		"TRI-Prx-App": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m tcp --sport 0 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv src -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst dst,dst -m mark ! --mark 0x40 -j ACCEPT",
		},
		"TRI-Prx-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst src,src -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv src -m addrtype --src-type LOCAL -j ACCEPT",
			"-p tcp -m tcp --dport 0 -j ACCEPT",
		},
		"TRI-Hst-App": {},
		"TRI-Hst-Net": {},
		"TRI-Svc-App": {},
		"TRI-Svc-Net": {},
		"TRI-UID-App": {},
		"TRI-UID-Net": {},

		"TRI-Net-pu1N7uS6--1": {
			"-p TCP -m set --match-set TRI-v4-ext-w5frVpu19gtV src -m state --state NEW -m set ! --match-set TRI-v4-TargetTCP src --match multiport --dports 80 -j DROP",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-flags SYN,ACK ACK -j NFQUEUE --queue-balance 20:23",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m state --state ESTABLISHED -j NFQUEUE --queue-balance 16:19",
			"-p tcp -m state --state ESTABLISHED -m comment --comment TCP-Established-Connections -j ACCEPT",
			"-s 0.0.0.0/0 -m state --state NEW -j NFLOG --nflog-group 11 --nflog-prefix pu1:default:default6",
			"-s 0.0.0.0/0 -j DROP",
		},

		"TRI-App-pu1N7uS6--1": {
			"-p TCP -m set --match-set TRI-v4-ext-uNdc0pu19gtV dst -m state --state NEW -m set ! --match-set TRI-v4-TargetTCP dst --match multiport --dports 80 -j DROP",
			"-p tcp -m tcp --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 0:3",
			"-p tcp -m tcp --tcp-flags SYN,ACK ACK -j NFQUEUE --queue-balance 4:7",
			"-p udp -m set --match-set TRI-v4-TargetUDP dst -j NFQUEUE --queue-balance 0:3",
			"-p udp -m set --match-set TRI-v4-TargetUDP dst -m state --state ESTABLISHED -m comment --comment UDP-Established-Connections -j ACCEPT",
			"-p tcp -m state --state ESTABLISHED -m comment --comment TCP-Established-Connections -j ACCEPT",
			"-d 0.0.0.0/0 -m state --state NEW -j NFLOG --nflog-group 10 --nflog-prefix pu1:default:default6",
			"-d 0.0.0.0/0 -j DROP",
		},
	}
)

func Test_OperationWithLinuxServices(t *testing.T) {
	Convey("Given an iptables controller with a memory backend ", t, func() {
		cfg := &runtime.Configuration{
			TCPTargetNetworks: []string{"0.0.0.0/0"},
			UDPTargetNetworks: []string{"10.0.0.0/8"},
			ExcludedNetworks:  []string{"127.0.0.1"},
		}

		commitFunc := func(buf *bytes.Buffer) error {
			return nil
		}

		iptv4 := provider.NewCustomBatchProvider(&baseIpt{}, commitFunc, []string{"nat", "mangle"})
		So(iptv4, ShouldNotBeNil)

		iptv6 := provider.NewCustomBatchProvider(&baseIpt{}, commitFunc, []string{"nat", "mangle"})
		So(iptv6, ShouldNotBeNil)

		ipsv4 := &memoryIPSetProvider{sets: map[string]*memoryIPSet{}}
		ipsv6 := &memoryIPSetProvider{sets: map[string]*memoryIPSet{}}

		i, err := createTestInstance(ipsv4, ipsv6, iptv4, iptv6, constants.LocalServer)
		So(err, ShouldBeNil)
		So(i, ShouldNotBeNil)

		i.SetTargetNetworks(cfg)

		Convey("When I start the controller, I should get the right global chains and ipsets", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err := i.Run(ctx)
			So(err, ShouldBeNil)

			for set, targets := range ipsv4.sets {
				So(expectedGlobalIPSets, ShouldContainKey, set)
				for target := range targets.set {
					So(expectedGlobalIPSets[set], ShouldContain, target)
				}
			}

			t := i.iptv4.impl.RetrieveTable()
			So(t, ShouldNotBeNil)
			So(len(t), ShouldEqual, 2)
			So(t["mangle"], ShouldNotBeNil)
			So(t["nat"], ShouldNotBeNil)

			for chain, rules := range t["mangle"] {
				So(expectedGlobalMangleChains, ShouldContainKey, chain)
				So(rules, ShouldResemble, expectedGlobalMangleChains[chain])
			}

			for chain, rules := range t["nat"] {
				So(expectedGlobalNATChains, ShouldContainKey, chain)
				So(rules, ShouldResemble, expectedGlobalNATChains[chain])
			}

			Convey("When I configure a new set of rules, the ACLs must be correct", func() {
				// Mock the exec commands
				i.iptv4.createPUPortSet = func(setName string) error {
					_, err := ipsv4.NewIpset(setName, "bitmap:port", &ipset.Params{})
					return err
				}
				i.iptv6.createPUPortSet = func(setName string) error {
					_, err := ipsv6.NewIpset(setName, "bitmap:port", &ipset.Params{})
					return err
				}

				appACLs := policy.IPRuleList{
					policy.IPRule{
						Addresses: []string{"30.0.0.0/24"},
						Ports:     []string{"80"},
						Protocols: []string{"TCP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Reject,
							ServiceID: "s1",
							PolicyID:  "1",
						},
					},
					policy.IPRule{
						Addresses: []string{"30.0.0.0/24"},
						Ports:     []string{"443"},
						Protocols: []string{"UDP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Accept,
							ServiceID: "s2",
							PolicyID:  "2",
						},
					},
					policy.IPRule{
						Addresses: []string{"50.0.0.0/24"},
						Ports:     []string{"443"},
						Protocols: []string{"icmp"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Accept,
							ServiceID: "s3",
							PolicyID:  "3",
						},
					},
				}
				netACLs := policy.IPRuleList{
					policy.IPRule{
						Addresses: []string{"40.0.0.0/24"},
						Ports:     []string{"80"},
						Protocols: []string{"TCP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Reject,
							ServiceID: "s3",
							PolicyID:  "1",
						},
					},
					policy.IPRule{
						Addresses: []string{"40.0.0.0/24"},
						Ports:     []string{"443"},
						Protocols: []string{"UDP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Accept,
							ServiceID: "s4",
							PolicyID:  "2",
						},
					},
				}
				ipl := policy.ExtendedMap{}
				policyrules := policy.NewPUPolicy("Context",
					policy.Police,
					appACLs,
					netACLs,
					nil,
					nil,
					nil,
					nil,
					nil,
					ipl,
					0,
					nil,
					nil,
					[]string{},
				)
				puInfo := policy.NewPUInfo("Context", common.LinuxProcessPU)
				puInfo.Policy = policyrules
				puInfo.Runtime.SetOptions(policy.OptionsType{
					CgroupMark: "10",
				})

				udpPortSpec, err := portspec.NewPortSpecFromString("5000", nil)
				So(err, ShouldBeNil)
				tcpPortSpec, err := portspec.NewPortSpecFromString("9000", nil)
				So(err, ShouldBeNil)

				puInfo.Runtime.SetServices([]common.Service{
					{
						Ports:    udpPortSpec,
						Protocol: 17,
					},
					{
						Ports:    tcpPortSpec,
						Protocol: 6,
					},
				})

				err = i.ConfigureRules(0, "pu1", puInfo)
				So(err, ShouldBeNil)
				t := i.iptv4.impl.RetrieveTable()

				for chain, rules := range t["mangle"] {
					So(expectedMangleAfterPUInsert, ShouldContainKey, chain)
					So(rules, ShouldResemble, expectedMangleAfterPUInsert[chain])
				}

				for chain, rules := range t["nat"] {
					So(expectedNATAfterPUInsert, ShouldContainKey, chain)
					So(rules, ShouldResemble, expectedNATAfterPUInsert[chain])
				}

				for set, targets := range ipsv4.sets {
					So(expectedIPSetsAfterPUInsert, ShouldContainKey, set)
					for target := range targets.set {
						So(expectedIPSetsAfterPUInsert[set], ShouldContain, target)
					}
				}

				Convey("When I update the policy, the update must result in correct state", func() {
					appACLs := policy.IPRuleList{
						policy.IPRule{
							Addresses: []string{"30.0.0.0/24"},
							Ports:     []string{"80"},
							Protocols: []string{"TCP"},
							Policy: &policy.FlowPolicy{
								Action:    policy.Reject,
								ServiceID: "s1",
								PolicyID:  "1",
							},
						},
					}
					netACLs := policy.IPRuleList{
						policy.IPRule{
							Addresses: []string{"40.0.0.0/24"},
							Ports:     []string{"80"},
							Protocols: []string{"TCP"},
							Policy: &policy.FlowPolicy{
								Action:    policy.Reject,
								ServiceID: "s3",
								PolicyID:  "1",
							},
						},
					}
					ipl := policy.ExtendedMap{}
					policyrules := policy.NewPUPolicy("Context",
						policy.Police,
						appACLs,
						netACLs,
						nil,
						nil,
						nil,
						nil,
						nil,
						ipl,
						0,
						nil,
						nil,
						[]string{},
					)
					puInfoUpdated := policy.NewPUInfo("Context", common.LinuxProcessPU)
					puInfoUpdated.Policy = policyrules
					puInfoUpdated.Runtime.SetOptions(policy.OptionsType{
						CgroupMark: "10",
					})

					err := i.UpdateRules(1, "pu1", puInfoUpdated, puInfo)
					So(err, ShouldBeNil)

					t := i.iptv4.impl.RetrieveTable()
					for chain, rules := range t["mangle"] {
						So(expectedMangleAfterPUUpdate, ShouldContainKey, chain)
						So(rules, ShouldResemble, expectedMangleAfterPUUpdate[chain])
					}

					Convey("When I delete the same rule, the chains must be restored in the global state", func() {
						err := i.DeleteRules(1, "pu1", "0", "5000", "10", "", "0", common.LinuxProcessPU)
						So(err, ShouldBeNil)

						t := i.iptv4.impl.RetrieveTable()

						So(t["mangle"], ShouldNotBeNil)
						So(t["nat"], ShouldNotBeNil)

						for chain, rules := range t["mangle"] {
							So(expectedGlobalMangleChains, ShouldContainKey, chain)
							So(rules, ShouldResemble, expectedGlobalMangleChains[chain])
						}

						for chain, rules := range t["nat"] {
							if len(rules) > 0 {
								So(expectedGlobalNATChains, ShouldContainKey, chain)
								So(rules, ShouldResemble, expectedGlobalNATChains[chain])
							}
						}

						Convey("When I cancel the context, it should cleanup", func() {
							cancel()
							time.Sleep(1 * time.Second)
							t := i.iptv4.impl.RetrieveTable()
							So(len(t["mangle"]), ShouldEqual, 2)
						})
					})
				})
			})
		})
	})
}

var (
	expectedContainerGlobalMangleChains = map[string][]string{
		"INPUT": {
			"-m set ! --match-set TRI-v4-Excluded src -j TRI-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-App",
		},
		"TRI-App": {
			"-j TRI-Prx-App",
			"-m mark --mark 1073741922 -j ACCEPT",
			"-m connmark --mark 61166 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j MARK --set-mark 99",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 8:11 --queue-bypass",
		},
		"TRI-Net": {
			"-j TRI-Prx-Net",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m string --string n30njxq7bmiwr6dtxq --algo bm --to 65535 -j NFQUEUE --queue-bypass --queue-balance 24:27",
			"-m connmark --mark 61166 -j ACCEPT",
			"-m set --match-set TRI-v4-TargetTCP src -p tcp -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 24:27 --queue-bypass",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-option 34 --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19 --queue-bypass",
		},
		"TRI-Prx-App": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
		"TRI-Prx-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
	}

	expectedContainerGlobalNATChains = map[string][]string{
		"PREROUTING": {
			"-p tcp -m addrtype --dst-type LOCAL -m set ! --match-set TRI-v4-Excluded src -j TRI-Redir-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-Redir-App",
		},
		"TRI-Redir-App": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
		"TRI-Redir-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
		},
	}

	expectedContainerGlobalIPSets = map[string][]string{
		"TRI" + "-v4-" + targetTCPNetworkSet: {"0.0.0.0/1", "128.0.0.0/1"},
		"TRI" + "-v4-" + targetUDPNetworkSet: {"10.0.0.0/8"},
		"TRI" + "-v4-" + excludedNetworkSet:  {"127.0.0.1"},
	}

	expectedContainerMangleAfterPUInsert = map[string][]string{
		"INPUT": {
			"-m set ! --match-set TRI-v4-Excluded src -j TRI-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-App",
		},
		"TRI-App": {
			"-j TRI-Prx-App",
			"-m mark --mark 1073741922 -j ACCEPT",
			"-m connmark --mark 61166 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j MARK --set-mark 99",
			"-p tcp -m set --match-set TRI-v4-TargetTCP dst -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 8:11 --queue-bypass",
			"-m comment --comment Container-specific-chain -j TRI-App-pu1N7uS6--0",
		},
		"TRI-Net": {
			"-j TRI-Prx-Net",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m string --string n30njxq7bmiwr6dtxq --algo bm --to 65535 -j NFQUEUE --queue-bypass --queue-balance 24:27",
			"-m connmark --mark 61166 -j ACCEPT",
			"-m set --match-set TRI-v4-TargetTCP src -p tcp -m tcp --tcp-flags SYN,ACK SYN,ACK -j NFQUEUE --queue-balance 24:27 --queue-bypass",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-option 34 --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19 --queue-bypass",
			"-m comment --comment Container-specific-chain -j TRI-Net-pu1N7uS6--0",
		},
		"TRI-Prx-App": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m tcp --sport 0 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv src -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst dst,dst -m mark ! --mark 0x40 -j ACCEPT",
		},
		"TRI-Prx-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst src,src -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv src -m addrtype --src-type LOCAL -j ACCEPT",
			"-p tcp -m tcp --dport 0 -j ACCEPT",
		},
		"TRI-Net-pu1N7uS6--0": {
			"-p UDP -m set --match-set TRI-v4-ext-6zlJIpu19gtV src -m state --state ESTABLISHED -j ACCEPT",
			"-p TCP -m set --match-set TRI-v4-ext-w5frVpu19gtV src -m state --state NEW -m set ! --match-set TRI-v4-TargetTCP src --match multiport --dports 80 -j DROP",
			"-p UDP -m set --match-set TRI-v4-ext-IuSLspu19gtV src --match multiport --dports 443 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 16:19",
			"-p tcp -m set --match-set TRI-v4-TargetTCP src -m tcp --tcp-flags SYN,ACK ACK -j NFQUEUE --queue-balance 20:23",
			"-p udp -m set --match-set TRI-v4-TargetUDP src -m state --state ESTABLISHED -j NFQUEUE --queue-balance 16:19",
			"-p tcp -m state --state ESTABLISHED -m comment --comment TCP-Established-Connections -j ACCEPT",
			"-s 0.0.0.0/0 -m state --state NEW -j NFLOG --nflog-group 11 --nflog-prefix pu1:default:default6",
			"-s 0.0.0.0/0 -j DROP",
		},

		"TRI-App-pu1N7uS6--0": {
			"-p TCP -m set --match-set TRI-v4-ext-uNdc0pu19gtV dst -m state --state NEW -m set ! --match-set TRI-v4-TargetTCP dst --match multiport --dports 80 -j DROP",
			"-p UDP -m set --match-set TRI-v4-ext-6zlJIpu19gtV dst --match multiport --dports 443 -j ACCEPT",
			"-p UDP -m set --match-set TRI-v4-ext-IuSLspu19gtV dst -m state --state ESTABLISHED -j ACCEPT",
			"-p tcp -m tcp --tcp-flags SYN,ACK SYN -j NFQUEUE --queue-balance 0:3",
			"-p tcp -m tcp --tcp-flags SYN,ACK ACK -j NFQUEUE --queue-balance 4:7",
			"-p udp -m set --match-set TRI-v4-TargetUDP dst -j NFQUEUE --queue-balance 0:3",
			"-p udp -m set --match-set TRI-v4-TargetUDP dst -m state --state ESTABLISHED -m comment --comment UDP-Established-Connections -j ACCEPT",
			"-p tcp -m state --state ESTABLISHED -m comment --comment TCP-Established-Connections -j ACCEPT",
			"-d 0.0.0.0/0 -m state --state NEW -j NFLOG --nflog-group 10 --nflog-prefix pu1:default:default6",
			"-d 0.0.0.0/0 -j DROP",
		},
	}

	expectedContainerNATAfterPUInsert = map[string][]string{
		"PREROUTING": {
			"-p tcp -m addrtype --dst-type LOCAL -m set ! --match-set TRI-v4-Excluded src -j TRI-Redir-Net",
		},
		"OUTPUT": {
			"-m set ! --match-set TRI-v4-Excluded dst -j TRI-Redir-App",
		},
		"TRI-Redir-App": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-dst dst,dst -m mark ! --mark 0x40 -m cgroup --cgroup 10 -j REDIRECT --to-ports 0",
		},
		"TRI-Redir-Net": {
			"-m mark --mark 0x40 -j ACCEPT",
			"-p tcp -m set --match-set TRI-v4-Proxy-pu19gtV-srv dst -m mark ! --mark 0x40 -j REDIRECT --to-ports 0",
		},
	}

	expectedContainerIPSetsAfterPUInsert = map[string][]string{
		"TRI-v4-" + targetTCPNetworkSet: {"0.0.0.0/1", "128.0.0.0/1"},
		"TRI-v4-" + targetUDPNetworkSet: {"10.0.0.0/8"},
		"TRI-v4-" + excludedNetworkSet:  {"127.0.0.1"},
		"TRI-v4-ProcPort-pu19gtV":       {},
		"TRI-v4-ext-6zlJIpu19gtV":       {"30.0.0.0/24"},
		"TRI-v4-ext-uNdc0pu19gtV":       {"30.0.0.0/24"},
		"TRI-v4-ext-w5frVpu19gtV":       {"40.0.0.0/24"},
		"TRI-v4-ext-IuSLspu19gtV":       {"40.0.0.0/24"},
		"TRI-v4-Proxy-pu19gtV-dst":      {},
		"TRI-v4-Proxy-pu19gtV-srv":      {},
	}
)

func Test_OperationWithContainers(t *testing.T) {
	Convey("Given an iptables controller with a memory backend for containers ", t, func() {
		cfg := &runtime.Configuration{
			TCPTargetNetworks: []string{"0.0.0.0/0"},
			UDPTargetNetworks: []string{"10.0.0.0/8"},
			ExcludedNetworks:  []string{"127.0.0.1"},
		}

		commitFunc := func(buf *bytes.Buffer) error {
			return nil
		}

		iptv4 := provider.NewCustomBatchProvider(&baseIpt{}, commitFunc, []string{"nat", "mangle"})
		So(iptv4, ShouldNotBeNil)

		iptv6 := provider.NewCustomBatchProvider(&baseIpt{}, commitFunc, []string{"nat", "mangle"})
		So(iptv6, ShouldNotBeNil)

		ipsv4 := &memoryIPSetProvider{sets: map[string]*memoryIPSet{}}
		ipsv6 := &memoryIPSetProvider{sets: map[string]*memoryIPSet{}}

		i, err := createTestInstance(ipsv4, ipsv6, iptv4, iptv6, constants.RemoteContainer)
		So(err, ShouldBeNil)
		So(i, ShouldNotBeNil)

		i.SetTargetNetworks(cfg)

		Convey("When I start the controller, I should get the right global chains and sets", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			err := i.Run(ctx)
			So(err, ShouldBeNil)

			for set, targets := range ipsv4.sets {
				So(expectedContainerGlobalIPSets, ShouldContainKey, set)
				for target := range targets.set {
					So(expectedContainerGlobalIPSets[set], ShouldContain, target)
				}
			}

			t := i.iptv4.impl.RetrieveTable()
			So(t, ShouldNotBeNil)
			So(len(t), ShouldEqual, 2)
			So(t["mangle"], ShouldNotBeNil)
			So(t["nat"], ShouldNotBeNil)

			for chain, rules := range t["mangle"] {
				So(expectedContainerGlobalMangleChains, ShouldContainKey, chain)
				So(rules, ShouldResemble, expectedContainerGlobalMangleChains[chain])
			}

			for chain, rules := range t["nat"] {
				So(expectedContainerGlobalNATChains, ShouldContainKey, chain)
				So(rules, ShouldResemble, expectedContainerGlobalNATChains[chain])
			}

			Convey("When I configure a new set of rules, the ACLs must be correct", func() {
				// Mock the exec commands
				i.iptv4.createPUPortSet = func(setName string) error {
					_, err := ipsv4.NewIpset(setName, "bitmap:port", &ipset.Params{})
					return err
				}

				i.iptv6.createPUPortSet = func(setName string) error {
					_, err := ipsv6.NewIpset(setName, "bitmap:port", &ipset.Params{})
					return err
				}

				appACLs := policy.IPRuleList{
					policy.IPRule{
						Addresses: []string{"30.0.0.0/24"},
						Ports:     []string{"80"},
						Protocols: []string{"TCP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Reject,
							ServiceID: "s1",
							PolicyID:  "1",
						},
					},
					policy.IPRule{
						Addresses: []string{"30.0.0.0/24"},
						Ports:     []string{"443"},
						Protocols: []string{"UDP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Accept,
							ServiceID: "s2",
							PolicyID:  "2",
						},
					},
				}
				netACLs := policy.IPRuleList{
					policy.IPRule{
						Addresses: []string{"40.0.0.0/24"},
						Ports:     []string{"80"},
						Protocols: []string{"TCP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Reject,
							ServiceID: "s3",
							PolicyID:  "1",
						},
					},
					policy.IPRule{
						Addresses: []string{"40.0.0.0/24"},
						Ports:     []string{"443"},
						Protocols: []string{"UDP"},
						Policy: &policy.FlowPolicy{
							Action:    policy.Accept,
							ServiceID: "s4",
							PolicyID:  "2",
						},
					},
				}
				ipl := policy.ExtendedMap{}
				policyrules := policy.NewPUPolicy("Context",
					policy.Police,
					appACLs,
					netACLs,
					nil,
					nil,
					nil,
					nil,
					nil,
					ipl,
					0,
					nil,
					nil,
					[]string{},
				)
				puInfo := policy.NewPUInfo("Context", common.ContainerPU)
				puInfo.Policy = policyrules
				puInfo.Runtime.SetOptions(policy.OptionsType{
					CgroupMark: "10",
				})
				err := i.ConfigureRules(0, "pu1", puInfo)
				So(err, ShouldBeNil)
				t := i.iptv4.impl.RetrieveTable()

				for chain, rules := range t["mangle"] {
					So(expectedContainerMangleAfterPUInsert, ShouldContainKey, chain)
					So(rules, ShouldResemble, expectedContainerMangleAfterPUInsert[chain])
				}

				for chain, rules := range t["nat"] {
					So(expectedContainerNATAfterPUInsert, ShouldContainKey, chain)
					So(rules, ShouldResemble, expectedContainerNATAfterPUInsert[chain])
				}

				for set, targets := range ipsv4.sets {
					So(expectedContainerIPSetsAfterPUInsert, ShouldContainKey, set)
					for target := range targets.set {
						So(expectedContainerIPSetsAfterPUInsert[set], ShouldContain, target)
					}
				}

				Convey("When I delete the same rule, the chains must be restored in the global state", func() {
					err := i.DeleteRules(0, "pu1", "0", "0", "10", "", "0", common.ContainerPU)
					So(err, ShouldBeNil)

					t := i.iptv4.impl.RetrieveTable()
					if err != nil {
						printTable(t)
					}

					So(t["mangle"], ShouldNotBeNil)
					So(t["nat"], ShouldNotBeNil)

					for chain, rules := range t["mangle"] {
						So(expectedContainerGlobalMangleChains, ShouldContainKey, chain)
						So(rules, ShouldResemble, expectedContainerGlobalMangleChains[chain])
					}

					for chain, rules := range t["nat"] {
						So(expectedContainerGlobalNATChains, ShouldContainKey, chain)
						So(rules, ShouldResemble, expectedContainerGlobalNATChains[chain])
					}

					Convey("When I cancel the context, it should cleanup", func() {
						cancel()
						time.Sleep(1 * time.Second)
						t := i.iptv4.impl.RetrieveTable()
						So(len(t["mangle"]), ShouldEqual, 2)
					})
				})

			})
		})
	})
}

func printTable(t map[string]map[string][]string) {
	fmt.Printf("\n")
	for table, chains := range t {
		fmt.Println(table)
		for chain, rules := range chains {
			fmt.Println(chain)
			for _, rule := range rules {
				fmt.Println(rule)
			}
		}
	}
}
