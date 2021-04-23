package iptablesctrl

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"text/template"

	"go.aporeto.io/enforcerd/trireme-lib/controller/constants"
	provider "go.aporeto.io/enforcerd/trireme-lib/controller/pkg/aclprovider"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/ebpf"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/fqconfig"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/ipsetmanager"
	"go.aporeto.io/enforcerd/trireme-lib/controller/runtime"
	"go.aporeto.io/enforcerd/trireme-lib/monitor/extractors"
	"go.aporeto.io/enforcerd/trireme-lib/policy"
	"go.uber.org/zap"
)

const (
	mainAppChain        = constants.ChainPrefix + "App"
	mainNetChain        = constants.ChainPrefix + "Net"
	appChainPrefix      = constants.ChainPrefix + "App-"
	netChainPrefix      = constants.ChainPrefix + "Net-"
	natProxyOutputChain = constants.ChainPrefix + "Redir-App"
	natProxyInputChain  = constants.ChainPrefix + "Redir-Net"
	proxyOutputChain    = constants.ChainPrefix + "Prx-App"
	proxyInputChain     = constants.ChainPrefix + "Prx-Net"
	istioChain          = constants.ChainPrefix + "Istio"

	// TriremeInput represent the chain that contains pu input rules.
	TriremeInput = constants.ChainPrefix + "Pid-Net"
	// TriremeOutput represent the chain that contains pu output rules.
	TriremeOutput = constants.ChainPrefix + "Pid-App"

	// NetworkSvcInput represent the chain that contains NetworkSvc input rules.
	NetworkSvcInput = constants.ChainPrefix + "Svc-Net"

	// NetworkSvcOutput represent the chain that contains NetworkSvc output rules.
	NetworkSvcOutput = constants.ChainPrefix + "Svc-App"

	// HostModeInput represent the chain that contains Hostmode input rules.
	HostModeInput = constants.ChainPrefix + "Hst-Net"

	// HostModeOutput represent the chain that contains Hostmode output rules.
	HostModeOutput = constants.ChainPrefix + "Hst-App"
	// NfqueueOutput represents the chain that contains the nfqueue output rules
	NfqueueOutput = constants.ChainPrefix + "Nfq-OUT"
	// NfqueueInput represents the chain that contains the nfqueue input rules
	NfqueueInput = constants.ChainPrefix + "Nfq-IN"
	// IstioUID is the UID of the istio-proxy(envoy) that is used in the iptables to identify the
	// envoy generated traffic
	IstioUID = "1337"
	// IstioRedirPort is the port where the App traffic from the output chain
	// is redirected into Istio-proxy, we need to accept this traffic as we don't to come in between
	// APP --> Envoy traffic.
	IstioRedirPort = "15001"
)

type iptables struct {
	impl            IPImpl
	fqc             fqconfig.FilterQueue
	mode            constants.ModeType
	ipsetmanager    ipsetmanager.IPSetManager
	bpf             ebpf.BPFModule
	serviceMeshType policy.ServiceMesh
}

// IPImpl interface is to be used by the iptable implentors like ipv4 and ipv6.
type IPImpl interface {
	provider.IptablesProvider
	IPVersion() int
	ProtocolAllowed(proto string) bool
	IPFilter() func(net.IP) bool
	GetDefaultIP() string
	NeedICMP() bool
}

type ipFilter func(net.IP) bool

func createIPInstance(impl IPImpl, ipsetmanager ipsetmanager.IPSetManager, fqc fqconfig.FilterQueue, mode constants.ModeType, ebpf ebpf.BPFModule, ServiceMeshType policy.ServiceMesh) *iptables {

	return &iptables{
		impl:            impl,
		fqc:             fqc,
		mode:            mode,
		ipsetmanager:    ipsetmanager,
		bpf:             ebpf,
		serviceMeshType: ServiceMeshType,
	}
}

func (i *iptables) SetTargetNetworks(c *runtime.Configuration) error {
	if c == nil {
		return nil
	}

	tcp := c.TCPTargetNetworks
	udp := c.UDPTargetNetworks
	excluded := c.ExcludedNetworks

	// If there are no target networks, capture all traffic
	if len(tcp) == 0 {
		tcp = []string{IPv4DefaultIP, IPv6DefaultIP}
	}

	return i.ipsetmanager.UpdateIPsetsForTargetAndExcludedNetworks(tcp, udp, excluded)
}

func (i *iptables) Run(ctx context.Context) error {

	// Clean any previous ACLs. This is needed in case we crashed at some
	// earlier point or there are other ACLs that create conflicts. We
	// try to clean only ACLs related to Trireme.
	if err := i.cleanACLs(); err != nil {
		return fmt.Errorf("Unable to clean previous acls while starting the supervisor: %s", err)
	}

	if err := i.ipsetmanager.DestroyAllIPsets(); err != nil {
		zap.L().Debug("ipset destroy all ipset returned error", zap.Error(err))
	}

	if err := i.ipsetmanager.CreateIPsetsForTargetAndExcludedNetworks(); err != nil {
		if err1 := i.ipsetmanager.DestroyAllIPsets(); err1 != nil {
			zap.L().Debug("ipset destroy all ipset returned error", zap.Error(err1))
		}
		return fmt.Errorf("unable to create target network ipsets: %s", err)
	}

	// Windows needs to initialize some ipsets
	if err := i.platformInit(); err != nil {
		return err
	}

	// Initialize all the global Trireme chains. There are several global chaims
	// that apply to all PUs:
	// Tri-App/Tri-Net are the main chains for the egress/ingress directions
	// UID related chains for any UID PUs.
	// Host, Service, Pid chains for the different modes of operation (host mode, pu mode, host service).
	// The priority is explicit (Pid activations take precedence of Service activations and Host Services)
	if err := i.initializeChains(); err != nil {
		return fmt.Errorf("Unable to initialize chains: %s", err)
	}

	// Insert the global ACLS. These are the main ACLs that will direct traffic from
	// the INPUT/OUTPUT chains to the Trireme chains. They also includes the main
	// rules of the main chains. These rules are never touched again, unless
	// if we gracefully terminate.
	if err := i.setGlobalRules(); err != nil {
		return fmt.Errorf("failed to update synack networks: %s", err)
	}

	if err := i.impl.Commit(); err != nil {
		return err
	}

	return nil
}

func (i *iptables) ConfigureRules(version int, contextID string, pu *policy.PUInfo) error {
	var err error
	var cfg *ACLInfo

	// First we create an IPSet for destination matching ports. This only
	// applies to Linux type PUs. A port set is associated with every PU,
	// and packets matching this destination get associated with the context
	// of the PU.
	if i.mode != constants.RemoteContainer {
		if err = i.ipsetmanager.CreateServerPortSet(contextID); err != nil {
			return err
		}
	}

	// Create the proxy sets. These are the target sets that will match
	// traffic towards the L4 and L4 services. There are two sets created
	// for every PU in this context (for outgoing and incoming traffic).
	// The outgoing sets capture all traffic towards specific destinations
	// as proxied traffic. Incoming sets correspond to the listening
	// services.
	// create proxySets only if there is no serviceMesh.
	if i.serviceMeshType == policy.None {
		if err := i.ipsetmanager.CreateProxySets(contextID); err != nil {
			return err
		}
	}

	// We create the generic ACL object that is used for all the templates.
	cfg, err = i.newACLInfo(version, contextID, pu, pu.Runtime.PUType())
	if err != nil {
		return err
	}

	// At this point we can install all the ACL rules that will direct
	// traffic to user space, allow for external access or direct
	// traffic towards the proxies
	if err = i.installRules(cfg, pu); err != nil {
		return err
	}

	// We commit the ACLs at the end. Note, that some of the ACLs in the
	// NAT table are not committed as a group. The commit function only
	// applies when newer versions of tables are installed (1.6.2 and above).
	if err = i.impl.Commit(); err != nil {
		zap.L().Error("unable to configure rules", zap.Error(err))
		return err
	}

	return nil
}

func (i *iptables) DeleteRules(version int, contextID string, tcpPorts, udpPorts string, mark string, username string, containerInfo *policy.PUInfo) error {
	cfg, err := i.newACLInfo(version, contextID, nil, containerInfo.Runtime.PUType())
	if err != nil {
		zap.L().Error("unable to create cleanup configuration", zap.Error(err))
		return err
	}
	if i.mode == constants.LocalServer {
		cfg.PacketMark = mark
	}
	cfg.UDPPorts = udpPorts
	cfg.TCPPorts = tcpPorts
	cfg.CgroupMark = mark
	cfg.Mark = mark

	cfg.PUType = containerInfo.Runtime.PUType()
	cfg.ProxyPort = containerInfo.Policy.ServicesListeningPort()
	cfg.DNSProxyPort = containerInfo.Policy.DNSProxyPort()
	// We clean up the chain rules first, so that we can delete the chains.
	// If any rule is not deleted, then the chain will show as busy.
	if err := i.deleteChainRules(cfg); err != nil {
		zap.L().Warn("Failed to clean rules", zap.Error(err))
	}

	// We can now delete the chains we have created for this PU. Note that
	// in every case we only create two chains for every PU. All other
	// chains are global.
	if err = i.deletePUChains(cfg); err != nil {
		zap.L().Warn("Failed to clean container chains while deleting the rules", zap.Error(err))
	}

	// We call commit to update all the changes, before destroying the ipsets.
	// References must be deleted for ipset deletion to succeed.
	if err := i.impl.Commit(); err != nil {
		zap.L().Warn("Failed to commit ACL changes", zap.Error(err))
	}

	if i.mode != constants.RemoteContainer {
		// We delete the set that captures all destination ports of the
		// PU. This only holds for Linux PUs.
		if err := i.ipsetmanager.DestroyServerPortSet(contextID); err != nil {
			zap.L().Warn("Failed to remove port set")
		}
	}

	// if serviceMesh is enabled then don't detroy the proxySets as we have not create them.
	if i.serviceMeshType == policy.None {
		// We delete the proxy port sets that were created for this PU.
		i.ipsetmanager.DestroyProxySets(contextID)
	}
	return nil
}

func (i *iptables) UpdateRules(version int, contextID string, containerInfo *policy.PUInfo, oldContainerInfo *policy.PUInfo) error {
	policyrules := containerInfo.Policy
	if policyrules == nil {
		return errors.New("policy rules cannot be nil")
	}

	// We cache the old config and we use it to delete the previous
	// rules. Every time we update the policy the version changes to
	// its binary complement.
	newCfg, err := i.newACLInfo(version, contextID, containerInfo, containerInfo.Runtime.PUType())
	if err != nil {
		return err
	}

	oldCfg, err := i.newACLInfo(version^1, contextID, oldContainerInfo, containerInfo.Runtime.PUType())
	if err != nil {
		return err
	}

	// Install all the new rules. The hooks to the new chains are appended
	// and do not take effect yet.
	if err := i.installRules(newCfg, containerInfo); err != nil {
		return err
	}

	// Remove mapping from old chain. By removing the old hooks the new
	// hooks take priority.
	if err := i.deleteChainRules(oldCfg); err != nil {
		return err
	}

	// Delete the old chains, since there are not references any more.
	if err := i.deletePUChains(oldCfg); err != nil {
		return err
	}

	// Commit all actions in on iptables-restore function.
	if err := i.impl.Commit(); err != nil {
		return err
	}

	return nil
}

func (i *iptables) CleanUp() error {

	if err := i.cleanACLs(); err != nil {
		zap.L().Error("Failed to clean acls while stopping the supervisor", zap.Error(err))
	}

	if err := i.ipsetmanager.DestroyAllIPsets(); err != nil {
		zap.L().Error("Failed to clean up ipsets", zap.Error(err))
	}

	i.ipsetmanager.Reset()

	return nil
}

// InitializeChains initializes the chains.
func (i *iptables) initializeChains() error {

	cfg, err := i.newACLInfo(0, "", nil, 0)
	if err != nil {
		return err
	}
	tmpl := template.Must(template.New(triremChains).Funcs(template.FuncMap{
		"isLocalServer": func() bool {
			return i.mode == constants.LocalServer
		},
		"isIstioEnabled": func() bool {
			return i.serviceMeshType == policy.Istio
		},
	}).Parse(triremChains))

	rules, err := extractRulesFromTemplate(tmpl, cfg)
	if err != nil {
		return fmt.Errorf("unable to create trireme chains:%s", err)
	}
	for _, rule := range rules {
		if len(rule) != 4 {
			continue
		}
		if err := i.impl.NewChain(rule[1], rule[3]); err != nil {
			return err
		}
	}

	return nil
}

// configureContainerRules adds the chain rules for a container.
// We separate in different methods to keep track of the changes
// independently.
func (i *iptables) configureContainerRules(cfg *ACLInfo) error {
	return i.addChainRules(cfg)
}

// configureLinuxRules adds the chain rules for a linux process or a UID process.
func (i *iptables) configureLinuxRules(cfg *ACLInfo) error {

	// These checks are for rather unusal error scenarios. We should
	// never see errors here. But better safe than sorry.
	if cfg.CgroupMark == "" {
		return errors.New("no mark value found")
	}

	if cfg.TCPPortSet == "" {
		return fmt.Errorf("port set was not found for the contextID. This should not happen")
	}

	return i.addChainRules(cfg)
}

type aclIPset struct {
	ipset string
	*policy.IPRule
}

func (i *iptables) getACLIPSets(ipRules policy.IPRuleList) []aclIPset {

	ipsets := i.ipsetmanager.GetACLIPsetsNames(ipRules)

	aclIPsets := make([]aclIPset, 0)

	for i, ipset := range ipsets {
		if len(ipset) > 0 {
			aclIPsets = append(aclIPsets, aclIPset{ipset, &ipRules[i]})
		}
	}

	return aclIPsets
}

// Install rules will install all the rules and update the port sets.
func (i *iptables) installRules(cfg *ACLInfo, containerInfo *policy.PUInfo) error {

	policyrules := containerInfo.Policy

	// update the proxy set only if there is no serviceMesh enabled.
	if i.serviceMeshType == policy.None {
		if err := i.updateProxySet(cfg.ContextID, containerInfo.Policy); err != nil {
			return err
		}
	}

	appACLIPset := i.getACLIPSets(policyrules.ApplicationACLs())
	netACLIPset := i.getACLIPSets(policyrules.NetworkACLs())

	// Install the PU specific chain first.
	if err := i.addContainerChain(cfg); err != nil {
		return err
	}

	// If its a remote and thus container, configure container rules.
	if i.mode == constants.RemoteContainer {
		if err := i.configureContainerRules(cfg); err != nil {
			return err
		}
	}

	// If its a Linux process configure the Linux rules.
	if i.mode == constants.LocalServer {
		if err := i.configureLinuxRules(cfg); err != nil {
			return err
		}
	}

	isHostPU := extractors.IsHostPU(containerInfo.Runtime, i.mode)

	if err := i.addPreNetworkACLRules(cfg); err != nil {
		return err
	}

	if err := i.addExternalACLs(cfg, cfg.AppChain, cfg.NetChain, appACLIPset, true); err != nil {
		return err
	}

	if err := i.addExternalACLs(cfg, cfg.NetChain, cfg.AppChain, netACLIPset, false); err != nil {
		return err
	}

	appAnyRules, netAnyRules, err := i.getProtocolAnyRules(cfg, appACLIPset, netACLIPset)
	if err != nil {
		return err
	}

	return i.addPacketTrap(cfg, isHostPU, appAnyRules, netAnyRules)
}

func (i *iptables) updateProxySet(contextID string, policy *policy.PUPolicy) error {
	i.ipsetmanager.FlushProxySets(contextID)

	for _, dependentService := range policy.DependentServices() {
		addresses := dependentService.NetworkInfo.Addresses
		min, max := dependentService.NetworkInfo.Ports.Range()

		for addrS := range addresses {
			_, addr, _ := net.ParseCIDR(addrS)
			for port := int(min); port <= int(max); port++ {
				if err := i.ipsetmanager.AddIPPortToDependentService(contextID, addr, strconv.Itoa(port)); err != nil {
					return fmt.Errorf("unable to add dependent ip %v to dependent networks ipset: %v", port, err)
				}
			}
		}
	}

	for _, exposedService := range policy.ExposedServices() {
		min, max := exposedService.PrivateNetworkInfo.Ports.Range()
		for port := int(min); port <= int(max); port++ {
			if err := i.ipsetmanager.AddPortToExposedService(contextID, strconv.Itoa(port)); err != nil {
				zap.L().Error("Failed to add vip", zap.Error(err))
				return fmt.Errorf("unable to add port %d to exposed ports ipset: %s", port, err)
			}
		}

		if exposedService.PublicNetworkInfo != nil {
			min, max := exposedService.PublicNetworkInfo.Ports.Range()
			for port := int(min); port <= int(max); port++ {
				if err := i.ipsetmanager.AddPortToExposedService(contextID, strconv.Itoa(port)); err != nil {
					zap.L().Error("Failed to VIP for public network", zap.Error(err))
					return fmt.Errorf("Failed to program VIP: %s", err)
				}
			}
		}
	}

	return nil
}
