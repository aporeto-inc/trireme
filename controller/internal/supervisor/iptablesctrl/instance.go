package iptablesctrl

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.aporeto.io/enforcerd/trireme-lib/controller/constants"
	provider "go.aporeto.io/enforcerd/trireme-lib/controller/pkg/aclprovider"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/ebpf"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/fqconfig"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/ipsetmanager"
	"go.aporeto.io/enforcerd/trireme-lib/controller/runtime"
	"go.aporeto.io/enforcerd/trireme-lib/policy"
	"go.uber.org/zap"
)

const (
	//IPV4 version for ipv4
	IPV4 = iota
	//IPV6 version for ipv6
	IPV6
)

//Instance is the structure holding the ipv4 and ipv6 handles
type Instance struct {
	iptv4 *iptables
	iptv6 *iptables
}

// SetTargetNetworks updates ths target networks. There are three different
// types of target networks:
//   - TCPTargetNetworks for TCP traffic (by default 0.0.0.0/0)
//   - UDPTargetNetworks for UDP traffic (by default empty)
//   - ExcludedNetworks that are always ignored (by default empty)
func (i *Instance) SetTargetNetworks(c *runtime.Configuration) error {

	if err := i.iptv4.SetTargetNetworks(c); err != nil {
		return err
	}

	if err := i.iptv6.SetTargetNetworks(c); err != nil {
		return err
	}

	return nil
}

// Run starts the iptables controller
func (i *Instance) Run(ctx context.Context) error {

	if err := i.iptv4.Run(ctx); err != nil {
		return err
	}

	if err := i.iptv6.Run(ctx); err != nil {
		return err
	}

	return nil
}

// ConfigureRules implments the ConfigureRules interface. It will create the
// port sets and then it will call install rules to create all the ACLs for
// the given chains. PortSets are only created here. Updates will use the
// exact same logic.
func (i *Instance) ConfigureRules(version int, contextID string, pu *policy.PUInfo) error {
	if err := i.iptv4.ConfigureRules(version, contextID, pu); err != nil {
		return err
	}

	if err := i.iptv6.ConfigureRules(version, contextID, pu); err != nil {
		return err
	}

	return nil
}

// DeleteRules implements the DeleteRules interface. This is responsible
// for cleaning all ACLs and associated chains, as well as ll the sets
// that we have created. Note, that this only clears up the state
// for a given processing unit.
func (i *Instance) DeleteRules(version int, contextID string, tcpPorts, udpPorts string, mark string, username string, containerInfo *policy.PUInfo) error {

	if err := i.iptv4.DeleteRules(version, contextID, tcpPorts, udpPorts, mark, username, containerInfo); err != nil {
		zap.L().Warn("Delete rules for iptables v4 returned error")
	}

	if err := i.iptv6.DeleteRules(version, contextID, tcpPorts, udpPorts, mark, username, containerInfo); err != nil {
		zap.L().Warn("Delete rules for iptables v6 returned error")
	}

	return nil
}

// UpdateRules implements the update part of the interface. Update will call
// installrules to install the new rules and then it will delete the old rules.
// For installations that do not have latests iptables-restore we time
// the operations so that the switch is almost atomic, by creating the new rules
// first. For latest kernel versions iptables-restorce will update all the rules
// in one shot.
func (i *Instance) UpdateRules(version int, contextID string, containerInfo *policy.PUInfo, oldContainerInfo *policy.PUInfo) error {

	if err := i.iptv4.UpdateRules(version, contextID, containerInfo, oldContainerInfo); err != nil {
		return err
	}

	if err := i.iptv6.UpdateRules(version, contextID, containerInfo, oldContainerInfo); err != nil {
		return err
	}

	return nil
}

// CleanUp requires the implementor to clean up all ACLs and destroy all
// the IP sets.
func (i *Instance) CleanUp() error {

	if err := i.iptv4.CleanUp(); err != nil {
		zap.L().Error("Failed to cleanup ipv4 rules")
	}

	if err := i.iptv6.CleanUp(); err != nil {
		zap.L().Error("Failed to cleanup ipv6 rules")
	}

	return nil
}

// CreateCustomRulesChain creates a custom rules chain if it doesnt exist
func (i *Instance) CreateCustomRulesChain() error {
	nonbatchedv4tableprovider, _ := provider.NewGoIPTablesProviderV4([]string{}, CustomQOSChain)
	nonbatchedv6tableprovider, _ := provider.NewGoIPTablesProviderV6([]string{}, CustomQOSChain)
	err := nonbatchedv4tableprovider.NewChain(customQOSChainTable, CustomQOSChain)
	if err != nil {
		zap.L().Debug("Chain already exists", zap.Error(err))

	}
	postroutingchainrulesv4, err := nonbatchedv4tableprovider.ListRules(customQOSChainTable, customQOSChainNFHook)
	if err != nil {
		zap.L().Error("ListRules returned error", zap.Error(err))
		return err
	}
	checkCustomRulesv4 := func() bool {
		for _, rule := range postroutingchainrulesv4 {
			if strings.Contains(rule, CustomQOSChain) {
				return true
			}
		}
		return false
	}
	if !checkCustomRulesv4() {
		if err := nonbatchedv4tableprovider.Insert(customQOSChainTable, customQOSChainNFHook, 1,
			"-m", "addrtype",
			"--src-type", "LOCAL",
			"-j", CustomQOSChain,
		); err != nil {
			zap.L().Debug("Unable to create ipv4 custom rule", zap.Error(err))
		}
	}

	err = nonbatchedv6tableprovider.NewChain(customQOSChainTable, CustomQOSChain)
	if err != nil {
		zap.L().Debug("Chain already exists", zap.Error(err))
	}
	postroutingchainrulesv6, err := nonbatchedv6tableprovider.ListRules(customQOSChainTable, customQOSChainNFHook)
	if err != nil {
		return err
	}
	checkCustomRulesv6 := func() bool {
		for _, rule := range postroutingchainrulesv6 {
			if strings.Contains(rule, CustomQOSChain) {
				return true
			}
		}
		return false
	}
	if !checkCustomRulesv6() {
		if err := nonbatchedv6tableprovider.Append(customQOSChainTable, customQOSChainNFHook,
			"-m", "addrtype",
			"--src-type", "LOCAL",
			"-j", CustomQOSChain,
		); err != nil {
			zap.L().Debug("Unable to create ipv6 custom rule", zap.Error(err))
		}
	}

	return nil
}

// NewInstance creates a new iptables controller instance
func NewInstance(fqc fqconfig.FilterQueue, mode constants.ModeType, ipv6Enabled bool, ebpf ebpf.BPFModule, iptablesLockfile string, serviceMeshType policy.ServiceMesh) (*Instance, error) {

	// our iptables binary `aporeto-iptables` uses the environment variable XT_LOCK_NAME
	// to set the iptables lockfile. Standard iptables does not look at this environment variable
	if iptablesLockfile != "" {
		if err := os.Setenv("XT_LOCK_NAME", iptablesLockfile); err != nil {
			return nil, fmt.Errorf("unable to set XT_LOCK_NAME: %s", err)
		}
	}

	ipv4Impl, err := GetIPv4Impl()
	if err != nil {
		return nil, fmt.Errorf("unable to create ipv4 instance: %s", err)
	}

	ipsetV4 := ipsetmanager.V4()
	iptInstanceV4 := createIPInstance(ipv4Impl, ipsetV4, fqc, mode, ebpf, serviceMeshType)

	ipv6Impl, err := GetIPv6Impl(ipv6Enabled)
	if err != nil {
		return nil, fmt.Errorf("unable to create ipv6 instance: %s", err)
	}

	ipsetV6 := ipsetmanager.V6()
	iptInstanceV6 := createIPInstance(ipv6Impl, ipsetV6, fqc, mode, ebpf, serviceMeshType)

	return newInstanceWithProviders(iptInstanceV4, iptInstanceV6)
}

// newInstanceWithProviders is called after ipt and ips have been created. This helps
// with all the unit testing to be able to mock the providers.
func newInstanceWithProviders(iptv4 *iptables, iptv6 *iptables) (*Instance, error) {

	i := &Instance{
		iptv4: iptv4,
		iptv6: iptv6,
	}

	return i, nil
}

// ACLProvider returns the current ACL provider that can be re-used by other entities.
func (i *Instance) ACLProvider() []provider.IptablesProvider {
	return []provider.IptablesProvider{i.iptv4.impl, i.iptv6.impl}
}
