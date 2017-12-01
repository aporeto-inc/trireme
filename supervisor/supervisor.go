package supervisor

import (
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/aporeto-inc/trireme-lib/cache"
	"github.com/aporeto-inc/trireme-lib/collector"
	"github.com/aporeto-inc/trireme-lib/constants"
	"github.com/aporeto-inc/trireme-lib/enforcer/policyenforcer"
	"github.com/aporeto-inc/trireme-lib/enforcer/utils/fqconfig"
	"github.com/aporeto-inc/trireme-lib/policy"
	"github.com/aporeto-inc/trireme-lib/portset"
	"github.com/aporeto-inc/trireme-lib/supervisor/ipsetctrl"
	"github.com/aporeto-inc/trireme-lib/supervisor/iptablesctrl"
)

type cacheData struct {
	version       int
	ips           policy.ExtendedMap
	mark          string
	port          string
	uid           string
	containerInfo *policy.PUInfo
}

// Config is the structure holding all information about the supervisor
type Config struct {
	//implType ImplementationType
	mode constants.ModeType

	versionTracker cache.DataStore
	collector      collector.EventCollector
	filterQueue    *fqconfig.FilterQueue
	excludedIPs    []string
	impl           Implementor

	triremeNetworks []string

	portSetInstance portset.PortSet

	sync.Mutex
}

// NewSupervisor will create a new connection supervisor that uses IPTables
// to redirect specific packets to userspace. It instantiates multiple data stores
// to maintain efficient mappings between contextID, policy and IP addresses. This
// simplifies the lookup operations at the expense of memory.
func NewSupervisor(collector collector.EventCollector, enforcerInstance policyenforcer.Enforcer, mode constants.ModeType, implementation constants.ImplementationType, networks []string) (*Config, error) {

	if collector == nil {
		return nil, fmt.Errorf("Collector cannot be nil")
	}

	if enforcerInstance == nil {
		return nil, fmt.Errorf("Enforcer cannot be nil")
	}

	filterQueue := enforcerInstance.GetFilterQueue()

	if filterQueue == nil {
		return nil, fmt.Errorf("Enforcer FilterQueues cannot be nil")
	}

	portSetInstance := enforcerInstance.GetPortSetInstance()

	if portSetInstance == nil {
		return nil, fmt.Errorf("Enforcer portset instance cannot be nil")
	}

	s := &Config{
		mode:            mode,
		impl:            nil,
		versionTracker:  cache.NewCache("SupVersionTracker"),
		collector:       collector,
		filterQueue:     filterQueue,
		excludedIPs:     []string{},
		triremeNetworks: networks,
		portSetInstance: portSetInstance,
	}

	var err error
	switch implementation {
	case constants.IPSets:
		s.impl, err = ipsetctrl.NewInstance(s.filterQueue, false, mode)
	default:
		s.impl, err = iptablesctrl.NewInstance(s.filterQueue, mode, portSetInstance)
	}
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize supervisor controllers: %s", err)
	}

	return s, nil
}

// Supervise creates a mapping between an IP address and the corresponding labels.
// it invokes the various handlers that process the parameter policy.
func (s *Config) Supervise(contextID string, containerInfo *policy.PUInfo) error {

	if containerInfo == nil {
		return fmt.Errorf("ContainerInfo must not be nil")
	}
	if containerInfo.Policy == nil {
		return fmt.Errorf("containerInfo.Policy must not be nil")
	}
	if containerInfo.Runtime == nil {
		return fmt.Errorf("containerInfo.Runtime must not be nil")
	}

	_, err := s.versionTracker.Get(contextID)

	if err != nil {
		// ContextID is not found in Cache, New PU: Do create.
		return s.doCreatePU(contextID, containerInfo)
	}

	// Context already in the cache. Just run update
	return s.doUpdatePU(contextID, containerInfo)
}

// Unsupervise removes the mapping from cache and cleans up the iptable rules. ALL
// remove operations will print errors by they don't return error. We want to force
// as much cleanup as possible to avoid stale state
func (s *Config) Unsupervise(contextID string) error {

	version, err := s.versionTracker.Get(contextID)

	if err != nil {
		return fmt.Errorf("Cannot find policy version: %s", err)
	}

	cacheEntry := version.(*cacheData)
	port := cacheEntry.containerInfo.Runtime.Options().ProxyPort
	proxyPortSetName := iptablesctrl.PuPortSetName(contextID, cacheEntry.mark, "Proxy-")
	if err := s.impl.DeleteRules(cacheEntry.version, contextID, cacheEntry.ips, cacheEntry.port, cacheEntry.mark, cacheEntry.uid, port, proxyPortSetName); err != nil {
		zap.L().Warn("Some rules were not deleted during unsupervise", zap.Error(err))
	}

	if err := s.versionTracker.Remove(contextID); err != nil {
		zap.L().Warn("Failed to clean the rule version cache", zap.Error(err))
	}

	return nil
}

// Start starts the supervisor
func (s *Config) Start() error {

	if err := s.impl.Start(); err != nil {
		return fmt.Errorf("Unable to start the implementer: %s", err)
	}

	s.Lock()
	if err := s.impl.SetTargetNetworks([]string{}, s.triremeNetworks); err != nil {
		return err
	}
	s.Unlock()

	zap.L().Debug("Started the supervisor")

	return nil
}

// Stop stops the supervisor
func (s *Config) Stop() error {

	if err := s.impl.Stop(); err != nil {
		return fmt.Errorf("Unable to stop the implementer: %s", err)
	}

	return nil
}

// SetTargetNetworks sets the target networks of the supervisor
func (s *Config) SetTargetNetworks(networks []string) error {

	s.Lock()
	defer s.Unlock()

	// If there are no target networks, capture all traffic
	if len(networks) == 0 {
		networks = []string{"0.0.0.0/1", "128.0.0.0/1"}
	}

	if err := s.impl.SetTargetNetworks(s.triremeNetworks, networks); err != nil {
		return err
	}

	s.triremeNetworks = networks

	return nil
}

func (s *Config) doCreatePU(contextID string, containerInfo *policy.PUInfo) error {

	zap.L().Debug("IPTables update for the creation of a pu", zap.String("contextID", contextID))

	version := 0
	mark := containerInfo.Runtime.Options().CgroupMark
	port := policy.ConvertServicesToPortList(containerInfo.Runtime.Options().Services)
	uid := containerInfo.Runtime.Options().UserID

	cacheEntry := &cacheData{
		version:       version,
		ips:           containerInfo.Policy.IPAddresses(),
		mark:          mark,
		port:          port,
		uid:           uid,
		containerInfo: containerInfo,
	}

	// Version the policy so that we can do hitless policy changes
	s.versionTracker.AddOrUpdate(contextID, cacheEntry)

	if err := s.impl.ConfigureRules(version, contextID, containerInfo); err != nil {
		if uerr := s.Unsupervise(contextID); uerr != nil {
			zap.L().Warn("Failed to clean up state while creating the PU",
				zap.String("contextID", contextID),
				zap.Error(uerr),
			)
		}
		return err
	}

	return nil
}

// UpdatePU creates a mapping between an IP address and the corresponding labels
//and the invokes the various handlers that process all policies.
func (s *Config) doUpdatePU(contextID string, containerInfo *policy.PUInfo) error {

	cacheEntry, err := s.versionTracker.LockedModify(contextID, add, 1)

	if err != nil {
		return fmt.Errorf("Unable to find PU %s in cache: %s", contextID, err)
	}

	cachedEntry := cacheEntry.(*cacheData)
	if err := s.impl.UpdateRules(cachedEntry.version, contextID, containerInfo, cachedEntry.containerInfo); err != nil {
		if uerr := s.Unsupervise(contextID); uerr != nil {
			zap.L().Warn("Failed to clean up state while updating the PU",
				zap.String("contextID", contextID),
				zap.Error(uerr),
			)
		}
		return err
	}

	return nil
}

func add(a, b interface{}) interface{} {
	entry := a.(*cacheData)
	entry.version = entry.version ^ 1
	return entry
}
