package iptablesctrl

import (
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/aporeto-inc/trireme/constants"
	"github.com/aporeto-inc/trireme/enforcer/utils/fqconfig"
	"github.com/aporeto-inc/trireme/policy"
	"github.com/bvandewalle/go-ipset/ipset"

	"github.com/aporeto-inc/trireme/supervisor/provider"
)

const (
	uidchain         = "UIDCHAIN"
	chainPrefix      = "TRIREME-"
	appChainPrefix   = chainPrefix + "App-"
	netChainPrefix   = chainPrefix + "Net-"
	targetNetworkSet = "TargetNetSet"
	//PuPortSet The prefix for portset names
	PuPortSet                 = "PUPortSet-"
	ipTableSectionOutput      = "OUTPUT"
	ipTableSectionInput       = "INPUT"
	ipTableSectionPreRouting  = "PREROUTING"
	ipTableSectionPostRouting = "POSTROUTING"
	natProxyOutputChain       = "RedirProxy-App"
	natProxyInputChain        = "RedirProxy-Net"
	proxyOutputChain          = "Proxy-App"
	proxyInputChain           = "Proxy-Net"
	destProxyServiceSet       = "Service-VIP"
	srcProxyServiceSet        = "Service-PIP"
	proxyMark                 = "0x40"
	//ProxyPort DefaultProxyPort
	ProxyPort = "5000"
)

// Instance  is the structure holding all information about a implementation
type Instance struct {
	fqc                        *fqconfig.FilterQueue
	ipt                        provider.IptablesProvider
	ipset                      provider.IpsetProvider
	vipTargetSet               provider.Ipset
	pipTargetSet               provider.Ipset
	targetSet                  provider.Ipset
	appPacketIPTableContext    string
	appAckPacketIPTableContext string
	appProxyIPTableContext     string
	appPacketIPTableSection    string
	netPacketIPTableContext    string
	netPacketIPTableSection    string
	appCgroupIPTableSection    string
	appSynAckIPTableSection    string
	mode                       constants.ModeType
}

// NewInstance creates a new iptables controller instance
func NewInstance(fqc *fqconfig.FilterQueue, mode constants.ModeType) (*Instance, error) {

	ipt, err := provider.NewGoIPTablesProvider()
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize IPtables provider: %s", err)
	}

	ips := provider.NewGoIPsetProvider()
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize ipsets")
	}

	i := &Instance{
		fqc:   fqc,
		ipt:   ipt,
		ipset: ips,
		appPacketIPTableContext:    "raw",
		appAckPacketIPTableContext: "mangle",
		netPacketIPTableContext:    "mangle",
		appProxyIPTableContext:     "nat",
		mode: mode,
	}

	if mode == constants.LocalServer || mode == constants.RemoteContainer {
		i.appPacketIPTableSection = ipTableSectionOutput
		i.appCgroupIPTableSection = ipTableSectionOutput
		i.netPacketIPTableSection = ipTableSectionInput
		i.appSynAckIPTableSection = ipTableSectionOutput
	} else {
		i.appPacketIPTableSection = ipTableSectionPreRouting
		i.appCgroupIPTableSection = ipTableSectionOutput
		i.netPacketIPTableSection = ipTableSectionPostRouting
		i.appSynAckIPTableSection = ipTableSectionInput
	}

	return i, nil

}

// chainPrefix returns the chain name for the specific PU
func (i *Instance) chainName(contextID string, version int) (app, net string) {
	app = appChainPrefix + contextID + "-" + strconv.Itoa(version)
	net = netChainPrefix + contextID + "-" + strconv.Itoa(version)
	return app, net
}

//PuPortSetName returns the name of the pu portset
func PuPortSetName(contextID string, mark string) string {

	return (PuPortSet + contextID + "-" + mark)
}

// DefaultIPAddress returns the default IP address for the processing unit
func (i *Instance) defaultIP(addresslist map[string]string) (string, bool) {

	if ip, ok := addresslist[policy.DefaultNamespace]; ok && len(ip) > 0 {
		return ip, true
	}

	if i.mode == constants.LocalContainer {
		return "0.0.0.0/0", false
	}

	return "0.0.0.0/0", true
}

// ConfigureRules implmenets the ConfigureRules interface
func (i *Instance) ConfigureRules(version int, contextID string, containerInfo *policy.PUInfo) error {
	policyrules := containerInfo.Policy

	appChain, netChain := i.chainName(contextID, version)
	// policyrules.DefaultIPAddress()

	// Supporting only one ip
	ipAddress, ok := i.defaultIP(policyrules.IPAddresses())

	if !ok {
		return fmt.Errorf("No ip address found ")
	}
	proxyPort := containerInfo.Runtime.Options().ProxyPort
	zap.L().Error("COnfigureRules", zap.String("proxyPort", proxyPort))
	// Configure all the ACLs
	if err := i.addContainerChain(appChain, netChain); err != nil {
		return err
	}

	if i.mode != constants.LocalServer {

		if err := i.addChainRules("", appChain, netChain, ipAddress, "", "", "", proxyPort); err != nil {
			return err
		}

	} else {
		mark := containerInfo.Runtime.Options().CgroupMark
		if mark == "" {
			return fmt.Errorf("No Mark value found")
		}

		port := policy.ConvertServicesToPortList(containerInfo.Runtime.Options().Services)

		uid := containerInfo.Runtime.Options().UserID
		if uid != "" {
			//We are about to create a uid login pu
			//This set will be empty and we will only fill it when we find a port for it
			//The reason to use contextID here is to ensure that we don't need to talk between supervisor and enforcer to share names the id is derivable from information available in the enforcer
			if puseterr := i.createPUPortSet(PuPortSetName(contextID, mark)); puseterr != nil {
				return puseterr
			}
		}

		portSetName := PuPortSetName(contextID, mark)
		if err := i.addChainRules(portSetName, appChain, netChain, ipAddress, port, mark, uid, proxyPort); err != nil {

			return err
		}
	}

	if err := i.addPacketTrap(appChain, netChain, ipAddress, containerInfo.Policy.TriremeNetworks()); err != nil {
		return err
	}

	if err := i.addAppACLs(contextID, appChain, ipAddress, policyrules.ApplicationACLs()); err != nil {
		return err
	}

	if err := i.addNetACLs(contextID, netChain, ipAddress, policyrules.NetworkACLs()); err != nil {
		return err
	}

	if err := i.addExclusionACLs(appChain, netChain, ipAddress, policyrules.ExcludedNetworks()); err != nil {
		return err
	}

	return nil
}

// DeleteRules implements the DeleteRules interface
func (i *Instance) DeleteRules(version int, contextID string, ipAddresses policy.ExtendedMap, port string, mark string, uid string, proxyPort string) error {
	var ipAddress string
	var ok bool

	// Supporting only one ip
	if i.mode != constants.LocalServer {
		if ipAddresses == nil {
			return fmt.Errorf("Provided map of IP addresses is nil")
		}

		ipAddress, ok = i.defaultIP(ipAddresses)
		if !ok {
			return fmt.Errorf("No ip address found ")
		}
	}

	appChain, netChain := i.chainName(contextID, version)
	portSetName := PuPortSetName(contextID, mark)
	if derr := i.deleteChainRules(portSetName, appChain, netChain, ipAddress, port, mark, uid, proxyPort); derr != nil {
		zap.L().Warn("Failed to clean rules", zap.Error(derr))
	}

	if err := i.deleteAllContainerChains(appChain, netChain); err != nil {
		zap.L().Warn("Failed to clean container chains while deleting the rules", zap.Error(err))
	}
	if uid != "" {
		ips := ipset.IPSet{
			Name: PuPortSetName(contextID, mark),
		}
		if err := ips.Destroy(); err != nil {
			zap.L().Warn("Failed to clear puport set", zap.Error(err))
		}
	}
	return nil
}

// UpdateRules implements the update part of the interface
func (i *Instance) UpdateRules(version int, contextID string, containerInfo *policy.PUInfo) error {

	if containerInfo == nil {
		return fmt.Errorf("Container info cannot be nil")
	}

	policyrules := containerInfo.Policy
	if policyrules == nil {
		return fmt.Errorf("Policy rules cannot be nil")
	}

	// Supporting only one ip
	ipAddress, ok := i.defaultIP(policyrules.IPAddresses())
	if !ok {
		return fmt.Errorf("No ip address found ")
	}
	proxyPort := containerInfo.Runtime.Options().ProxyPort

	appChain, netChain := i.chainName(contextID, version)

	oldAppChain, oldNetChain := i.chainName(contextID, version^1)

	//Add a new chain for this update and map all rules there
	if err := i.addContainerChain(appChain, netChain); err != nil {
		return err
	}

	if err := i.addPacketTrap(appChain, netChain, ipAddress, containerInfo.Policy.TriremeNetworks()); err != nil {
		return err
	}

	if err := i.addAppACLs(contextID, appChain, ipAddress, policyrules.ApplicationACLs()); err != nil {
		return err
	}

	if err := i.addNetACLs(contextID, netChain, ipAddress, policyrules.NetworkACLs()); err != nil {
		return err
	}

	if err := i.addExclusionACLs(appChain, netChain, ipAddress, policyrules.ExcludedNetworks()); err != nil {
		return err
	}

	// Add mapping to new chain
	if i.mode != constants.LocalServer {
		if err := i.addChainRules("", appChain, netChain, ipAddress, "", "", "", proxyPort); err != nil {
			return err
		}
	} else {
		mark := containerInfo.Runtime.Options().CgroupMark
		if mark == "" {
			return fmt.Errorf("No Mark value found")
		}
		portlist := policy.ConvertServicesToPortList(containerInfo.Runtime.Options().Services)
		uid := containerInfo.Runtime.Options().UserID

		portSetName := PuPortSetName(contextID, mark)
		if err := i.addChainRules(portSetName, appChain, netChain, ipAddress, portlist, mark, uid, proxyPort); err != nil {

			return err
		}
	}

	//Remove mapping from old chain
	if i.mode != constants.LocalServer {
		if err := i.deleteChainRules("", oldAppChain, oldNetChain, ipAddress, "", "", "", proxyPort); err != nil {

			return err
		}
	} else {
		mark := containerInfo.Runtime.Options().CgroupMark
		port := policy.ConvertServicesToPortList(containerInfo.Runtime.Options().Services)
		uid := containerInfo.Runtime.Options().UserID
		portSetName := PuPortSetName(contextID, mark)
		if err := i.deleteChainRules(portSetName, oldAppChain, oldNetChain, ipAddress, port, mark, uid, proxyPort); err != nil {
			return err
		}
	}

	// Delete the old chain to clean up
	if err := i.deleteAllContainerChains(oldAppChain, oldNetChain); err != nil {
		return err
	}

	return nil
}

// Start starts the iptables controller
func (i *Instance) Start() error {

	// Clean any previous ACLs
	if err := i.cleanACLs(); err != nil {
		zap.L().Warn("Failed to clean previous acls while starting the supervisor", zap.Error(err))
	}

	if i.mode == constants.LocalContainer {
		if i.acceptMarkedPackets() != nil {
			return fmt.Errorf("Filter of marked packets was not set")
		}
	}

	zap.L().Debug("Started the iptables controller")

	return nil
}

// SetTargetNetworks updates ths target networks for SynAck packets
func (i *Instance) SetTargetNetworks(current, networks []string) error {

	if len(networks) == 0 {
		networks = []string{"0.0.0.0/1", "128.0.0.0/1"}
	}

	// Cleanup old ACLs
	if len(current) > 0 {
		return i.updateTargetNetworks(current, networks)
	}

	// Create the target network set
	if err := i.createTargetSet(networks); err != nil {
		return err
	}

	if err := i.createProxySets([]string{}, []string{}); err != nil {
		return err
	}
	i.ipt.NewChain(i.appAckPacketIPTableContext, uidchain)
	i.ipt.NewChain(i.appProxyIPTableContext, natProxyInputChain)
	i.ipt.NewChain(i.appProxyIPTableContext, natProxyOutputChain)
	i.ipt.NewChain(i.appAckPacketIPTableContext, proxyOutputChain)
	i.ipt.NewChain(i.appAckPacketIPTableContext, proxyInputChain)
	i.ipt.Insert(i.appAckPacketIPTableContext, i.appPacketIPTableSection, 1, "-j", uidchain)

	// Insert the ACLS that point to the target networks
	if err := i.setGlobalRules(i.appPacketIPTableSection, i.netPacketIPTableSection); err != nil {
		return fmt.Errorf("Failed to update synack networks")
	}

	return nil
}

// Stop stops the supervisor
func (i *Instance) Stop() error {

	zap.L().Debug("Stop the supervisor")

	// Clean any previous ACLs that we have installed
	if err := i.cleanACLs(); err != nil {
		zap.L().Error("Failed to clean acls while stopping the supervisor", zap.Error(err))
	}

	if err := i.ipset.DestroyAll(); err != nil {
		zap.L().Error("Failed to clean up ipsets", zap.Error(err))
	}

	return nil
}
