package monitor

import (
	"context"
	"fmt"

	"github.com/aporeto-inc/trireme-lib/common"
	"github.com/aporeto-inc/trireme-lib/monitor/config"
	"github.com/aporeto-inc/trireme-lib/monitor/internal/cni"
	"github.com/aporeto-inc/trireme-lib/monitor/internal/docker"
	"github.com/aporeto-inc/trireme-lib/monitor/internal/linux"
	"github.com/aporeto-inc/trireme-lib/monitor/internal/uid"
	"github.com/aporeto-inc/trireme-lib/monitor/registerer"
	"github.com/aporeto-inc/trireme-lib/monitor/remoteapi/server"
	"github.com/aporeto-inc/trireme-lib/policy"
	"go.uber.org/zap"
)

type monitors struct {
	config     *config.MonitorConfig
	monitors   map[config.Type]Implementation
	registerer registerer.Registerer
	server     server.APIServer
	policy     policy.Resolver
}

// NewMonitors instantiates all/any combination of monitors supported.
func NewMonitors(opts ...Options) (Monitor, error) {

	var err error

	c := &config.MonitorConfig{
		MergeTags: []string{},
		Common:    &config.ProcessorConfig{},
		Monitors:  map[config.Type]interface{}{},
	}

	for _, opt := range opts {
		opt(c)
	}

	if err = c.Common.IsComplete(); err != nil {
		return nil, err
	}

	m := &monitors{
		config:   c,
		monitors: make(map[config.Type]Implementation),
	}

	m.registerer = registerer.New()

	m.server, err = server.NewEventServer(common.TriremeSocket, m.registerer)
	if err != nil {
		return nil, err
	}

	for k, v := range c.Monitors {
		switch k {
		case config.CNI:
			mon := cnimonitor.New()
			mon.SetupHandlers(c.Common)
			if err := mon.SetupConfig(m.registerer, v); err != nil {
				return nil, fmt.Errorf("CNI: %s", err.Error())
			}
			m.monitors[config.CNI] = mon

		case config.Docker:
			mon := dockermonitor.New()
			mon.SetupHandlers(c.Common)
			if err := mon.SetupConfig(nil, v); err != nil {
				return nil, fmt.Errorf("Docker: %s", err.Error())
			}
			m.monitors[config.Docker] = mon

		case config.LinuxProcess:
			mon := linuxmonitor.New()
			mon.SetupHandlers(c.Common)
			if err := mon.SetupConfig(m.registerer, v); err != nil {
				return nil, fmt.Errorf("Process: %s", err.Error())
			}
			m.monitors[config.LinuxProcess] = mon

		case config.LinuxHost:
			mon := linuxmonitor.New()
			mon.SetupHandlers(c.Common)
			if err := mon.SetupConfig(m.registerer, v); err != nil {
				return nil, fmt.Errorf("Host: %s", err.Error())
			}
			m.monitors[config.LinuxHost] = mon

		case config.UID:
			mon := uidmonitor.New()
			mon.SetupHandlers(c.Common)
			if err := mon.SetupConfig(m.registerer, v); err != nil {
				return nil, fmt.Errorf("UID: %s", err.Error())
			}
			m.monitors[config.UID] = mon

		default:
			return nil, fmt.Errorf("Unsupported type %d", k)
		}
	}

	zap.L().Debug("Monitor configuration", zap.String("conf", m.config.String()))

	return m, nil
}

func (m *monitors) Run(ctx context.Context) (err error) {

	if err := m.server.Run(ctx); err != nil {
		return err
	}

	for _, v := range m.monitors {
		if err = v.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

// UpdateConfiguration updates the configuration of the monitors.
func (m *monitors) UpdateConfiguration(ctx context.Context, config *config.MonitorConfig) error {
	// Monitor configuration cannot change at this time.
	// TODO:
	return nil
}

// Resync resyncs the monitor
func (m *monitors) Resync(ctx context.Context) error {

	failure := false
	var errs string

	for _, i := range m.monitors {
		if err := i.ReSync(ctx); err != nil {
			errs = errs + err.Error()
			failure = true
		}
	}

	if failure {
		return fmt.Errorf("Monitor resync failed: %s", errs)
	}

	return nil
}
