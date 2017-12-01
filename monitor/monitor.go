package monitor

import (
	"fmt"

	"github.com/aporeto-inc/trireme-lib/collector"
	"github.com/aporeto-inc/trireme-lib/monitor/impl"
	"github.com/aporeto-inc/trireme-lib/monitor/rpc"
)

// Type specifies the type of monitors supported.
type Type int

// Types supported.
const (
	CNI Type = iota + 1
	Docker
	LinuxProcess
	LinuxHost
	UID
)

// Config specifies the configs for monitors.
type Config struct {
	Collector   collector.EventCollector
	PUHandler   monitorimpl.ProcessingUnitsHandler
	SyncHandler monitorimpl.SynchronizationHandler
	Monitors    map[Type]interface{}
}

type monitors struct {
	config          *Config
	monitors        map[Type]impl.Implementation
	userRPCListener rpcmonitor.Listener
	rootRPCListener rpcmonitor.Listener
}

// New instantiates all/any combination of monitors supported.
func New(c *Config) (Monitor, error) {

	m := &monitors{
		config:   c,
		monitors: make(map[Type]impl.Implementation),
	}

	m.userRPCListener = rpc.New(
		rpcmonitor.DefaultRPCAddress,
		false,
	)

	m.rootRPCListener = rpc.New(
		rpcmonitor.DefaultRootRPCAddress,
		true,
	)

	for k, v := range c.Monitors {
		switch k {
		case CNI:
			monitor := cni.New()
			if err := monitor.SetupHandlers(c.Collector, c.PUHandler, c.SyncHandler); err != nil {
				return nil, err
			}
			if err := monitor.SetupCfg(m.userRPCListener, v); err != nil {
				return nil, err
			}
			m.monitors[CNI] = monitor

		case Docker:
			monitor := docker.New()
			if err := monitor.SetupHandlers(c.Collector, c.PUHandler, c.SyncHandler); err != nil {
				return nil, err
			}
			if err := monitor.SetupCfg(nil, v); err != nil {
				return nil, err
			}
			m.monitors[Docker] = monitor

		case LinuxProcess:
			m.monitors[LinuxProcess] = linux.New()
			if err := monitor.SetupHandlers(c.Collector, c.PUHandler, c.SyncHandler); err != nil {
				return nil, err
			}
			if err := monitor.SetupCfg(m.userRPCListener, v); err != nil {
				return nil, err
			}
			m.monitors[LinuxProcess] = monitor

		case LinuxHost:
			m.monitors[LinuxHost] = linux.New()
			if err := monitor.SetupHandlers(c.Collector, c.PUHandler, c.SyncHandler); err != nil {
				return nil, err
			}
			if err := monitor.SetupCfg(m.rootRPCListener, v); err != nil {
				return nil, err
			}
			m.monitors[LinuxHost] = monitor

		case UID:
			m.monitors[UID] = uid.New()
			if err := monitor.SetupHandlers(c.Collector, c.PUHandler, c.SyncHandler); err != nil {
				return nil, err
			}
			if err := monitor.SetupCfg(m.userRPCListener, v); err != nil {
				return nil, err
			}
			m.monitors[UID] = monitor

		default:
			return nil, fmt.Errorf("Unsupported type %d", k)
		}
	}

	return m, nil
}

func (m *monitors) Start() (err error) {

	for k, v := range m.monitors {
		if err = v.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (m *monitors) Stop() error {

	for k, v := range m.monitors {
		if err = v.Stop(); err != nil {
			return err
		}
	}

	return nil
}
