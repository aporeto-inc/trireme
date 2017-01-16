// Package configurator provides some helper functions to helpe
// you create default Trireme and Monitor configurations.
package configurator

import (
	"crypto/ecdsa"

	log "github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/trireme"
	"github.com/aporeto-inc/trireme/collector"
	"github.com/aporeto-inc/trireme/constants"
	"github.com/aporeto-inc/trireme/enforcer"
	"github.com/aporeto-inc/trireme/monitor"
	"github.com/aporeto-inc/trireme/monitor/contextstore"
	"github.com/aporeto-inc/trireme/monitor/dockermonitor"
	"github.com/aporeto-inc/trireme/monitor/linuxmonitor"
	"github.com/aporeto-inc/trireme/monitor/linuxmonitor/cgnetcls"
	"github.com/aporeto-inc/trireme/monitor/rpcmonitor"
	"github.com/aporeto-inc/trireme/policy"

	"github.com/aporeto-inc/trireme/enforcer/utils/tokens"

	"github.com/aporeto-inc/trireme/enforcer/proxy"
	"github.com/aporeto-inc/trireme/enforcer/utils/rpcwrapper"
	"github.com/aporeto-inc/trireme/supervisor"
	"github.com/aporeto-inc/trireme/supervisor/proxy"
)

// NewTriremeLinuxProcess instantiates Trireme for a Linux process implementation
func NewTriremeLinuxProcess(
	serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	secrets tokens.Secrets) trireme.Trireme {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	enforcers := map[policy.PUType]enforcer.PolicyEnforcer{
		policy.LinuxProcessPU: enforcer.NewDefaultDatapathEnforcer(serverID,
			eventCollector,
			nil,
			secrets,
			constants.LocalServer,
		)}

	s, err := supervisor.NewSupervisor(
		eventCollector,
		enforcers[policy.LinuxProcessPU],
		networks,
		constants.LocalServer,
		constants.IPTables,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"package": "configurator",
			"error":   err.Error(),
		}).Fatal("Failed to load Supervisor")
	}

	supervisors := map[policy.PUType]supervisor.Supervisor{policy.ContainerPU: s}

	return trireme.NewTrireme(serverID, resolver, supervisors, enforcers, eventCollector)
}

// NewLocalTriremeDocker instantiates Trireme for Docker using enforcement on the
// main namespace
func NewLocalTriremeDocker(
	serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	secrets tokens.Secrets,
	impl constants.ImplementationType) trireme.Trireme {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	enforcers := map[policy.PUType]enforcer.PolicyEnforcer{
		policy.ContainerPU: enforcer.NewDefaultDatapathEnforcer(serverID,
			eventCollector,
			nil,
			secrets,
			constants.LocalContainer,
		)}

	s, err := supervisor.NewSupervisor(
		eventCollector,
		enforcers[policy.ContainerPU],
		networks,
		constants.LocalContainer,
		impl,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"package": "configurator",
			"error":   err.Error(),
		}).Fatal("Failed to load Supervisor")
	}

	supervisors := map[policy.PUType]supervisor.Supervisor{policy.ContainerPU: s}

	return trireme.NewTrireme(serverID, resolver, supervisors, enforcers, eventCollector)
}

// NewDistributedTriremeDocker instantiates Trireme using remote enforcers on
// the container namespaces
func NewDistributedTriremeDocker(serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	secrets tokens.Secrets,
	impl constants.ImplementationType) trireme.Trireme {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	rpcwrapper := rpcwrapper.NewRPCWrapper()

	enforcers := map[policy.PUType]enforcer.PolicyEnforcer{
		policy.ContainerPU: enforcerproxy.NewDefaultProxyEnforcer(
			serverID,
			eventCollector,
			secrets,
			rpcwrapper),
	}

	s, err := supervisorproxy.NewProxySupervisor(eventCollector, enforcers[0], networks, rpcwrapper)

	if err != nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Fatal("Cannot initialize proxy supervisor")
	}

	supervisors := map[policy.PUType]supervisor.Supervisor{policy.ContainerPU: s}

	return trireme.NewTrireme(serverID, resolver, supervisors, enforcers, eventCollector)
}

// NewHybridTrireme instantiates Trireme with both Linux and Docker enforcers.
// The Docker enforcers are remote
func NewHybridTrireme(
	serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	secrets tokens.Secrets,
) trireme.Trireme {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	rpcwrapper := rpcwrapper.NewRPCWrapper()
	containerEnforcer := enforcerproxy.NewDefaultProxyEnforcer(
		serverID,
		eventCollector,
		secrets,
		rpcwrapper)

	containerSupervisor, cerr := supervisorproxy.NewProxySupervisor(
		eventCollector,
		containerEnforcer,
		networks,
		rpcwrapper)

	if cerr != nil {
		log.WithFields(log.Fields{
			"package": "configurator",
			"error":   cerr.Error(),
		}).Fatal("Failed to load Supervisor")

	}

	processEnforcer := enforcer.NewDefaultDatapathEnforcer(serverID,
		eventCollector,
		processor,
		secrets,
		constants.LocalServer,
	)

	processSupervisor, perr := supervisor.NewSupervisor(
		eventCollector,
		processEnforcer,
		networks,
		constants.LocalServer,
		constants.IPTables,
	)

	if perr != nil {
		log.WithFields(log.Fields{
			"package": "configurator",
			"error":   perr.Error(),
		}).Fatal("Failed to load Supervisor")

	}

	enforcers := map[policy.PUType]enforcer.PolicyEnforcer{
		policy.ContainerPU:    containerEnforcer,
		policy.LinuxProcessPU: processEnforcer,
	}

	supervisors := map[policy.PUType]supervisor.Supervisor{
		policy.ContainerPU:    containerSupervisor,
		policy.LinuxProcessPU: processSupervisor,
	}

	trireme := trireme.NewTrireme(serverID, resolver, supervisors, enforcers, eventCollector)

	return trireme
}

// NewSecretsFromPSK creates secrets from a pre-shared key
func NewSecretsFromPSK(key []byte) tokens.Secrets {
	return tokens.NewPSKSecrets(key)
}

// NewSecretsFromPKI creates secrets from a PKI
func NewSecretsFromPKI(keyPEM, certPEM, caCertPEM []byte) tokens.Secrets {
	return tokens.NewPKISecrets(keyPEM, certPEM, caCertPEM, map[string]*ecdsa.PublicKey{})
}

// NewPSKTriremeWithDockerMonitor creates a new network isolator. The calling module must provide
// a policy engine implementation and a pre-shared secret. This is for backward
// compatibility. Will be removed
func NewPSKTriremeWithDockerMonitor(
	serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	syncAtStart bool,
	key []byte,
	dockerMetadataExtractor dockermonitor.DockerMetadataExtractor,
	remoteEnforcer bool,
) (trireme.Trireme, monitor.Monitor, supervisor.Excluder) {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	secrets := NewSecretsFromPSK(key)

	var triremeInstance trireme.Trireme

	if remoteEnforcer {
		triremeInstance = NewDistributedTriremeDocker(
			serverID,
			networks,
			resolver,
			processor,
			eventCollector,
			secrets,
			constants.IPTables)
	} else {
		triremeInstance = NewLocalTriremeDocker(
			serverID,
			networks,
			resolver,
			processor,
			eventCollector,
			secrets,
			constants.IPTables)
	}

	monitorInstance := dockermonitor.NewDockerMonitor(
		constants.DefaultDockerSocketType,
		constants.DefaultDockerSocket,
		triremeInstance,
		dockerMetadataExtractor,
		eventCollector,
		syncAtStart,
		nil)

	return triremeInstance, monitorInstance, triremeInstance.Supervisor(policy.ContainerPU).(supervisor.Excluder)

}

// NewPKITriremeWithDockerMonitor creates a new network isolator. The calling module must provide
// a policy engine implementation and private/public key pair and parent certificate.
// All certificates are passed in PEM format. If a certificate pool is provided
// certificates will not be transmitted on the wire
func NewPKITriremeWithDockerMonitor(
	serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	syncAtStart bool,
	keyPEM []byte,
	certPEM []byte,
	caCertPEM []byte,
	dockerMetadataExtractor dockermonitor.DockerMetadataExtractor,
	remoteEnforcer bool,
) (trireme.Trireme, monitor.Monitor, supervisor.Excluder, enforcer.PublicKeyAdder) {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	publicKeyAdder := tokens.NewPKISecrets(keyPEM, certPEM, caCertPEM, map[string]*ecdsa.PublicKey{})

	var triremeInstance trireme.Trireme

	if remoteEnforcer {
		triremeInstance = NewDistributedTriremeDocker(
			serverID,
			networks,
			resolver,
			processor,
			eventCollector,
			publicKeyAdder,
			constants.IPTables)
	} else {
		triremeInstance = NewLocalTriremeDocker(
			serverID,
			networks,
			resolver,
			processor,
			eventCollector,
			publicKeyAdder,
			constants.IPTables)
	}

	monitorInstance := dockermonitor.NewDockerMonitor(
		constants.DefaultDockerSocketType,
		constants.DefaultDockerSocket,
		triremeInstance,
		dockerMetadataExtractor,
		eventCollector,
		syncAtStart,
		nil)

	return triremeInstance, monitorInstance, triremeInstance.Supervisor(policy.ContainerPU).(supervisor.Excluder), publicKeyAdder

}

// NewPSKHybridTriremeWithMonitor creates a new network isolator. The calling module must provide
// a policy engine implementation and a pre-shared secret. This is for backward
// compatibility. Will be removed
func NewPSKHybridTriremeWithMonitor(
	serverID string,
	networks []string,
	resolver trireme.PolicyResolver,
	processor enforcer.PacketProcessor,
	eventCollector collector.EventCollector,
	syncAtStart bool,
	key []byte,
	dockerMetadataExtractor dockermonitor.DockerMetadataExtractor,
) (trireme.Trireme, monitor.Monitor, monitor.Monitor, supervisor.Excluder) {

	if eventCollector == nil {
		log.WithFields(log.Fields{
			"package": "configurator",
		}).Warn("Using a default collector for events")
		eventCollector = &collector.DefaultCollector{}
	}

	secrets := NewSecretsFromPSK(key)

	triremeInstance := NewHybridTrireme(
		serverID,
		networks,
		resolver,
		processor,
		eventCollector,
		secrets,
	)

	monitorDocker := dockermonitor.NewDockerMonitor(
		constants.DefaultDockerSocketType,
		constants.DefaultDockerSocket,
		triremeInstance,
		dockerMetadataExtractor,
		eventCollector,
		syncAtStart,
		nil,
	)

	//use rpcmonitor no need to return it since no other consumer for it
	netcls := cgnetcls.NewCgroupNetController()
	contextstorehdl := contextstore.NewContextStore()
	rpcmonitor, _ := rpcmonitor.NewRPCMonitor(
		rpcmonitor.Rpcaddress,
		linuxmonitor.SystemdRPCMetadataExtractor,
		triremeInstance,
		eventCollector,
		netcls,
		contextstorehdl,
	)

	return triremeInstance, monitorDocker, rpcmonitor, triremeInstance.Supervisor(policy.ContainerPU).(supervisor.Excluder)

}
