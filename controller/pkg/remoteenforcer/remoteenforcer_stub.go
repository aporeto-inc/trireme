// +build !linux

package remoteenforcer

import (
	"context"

	"github.com/blang/semver"
	"go.aporeto.io/trireme-lib/v11/controller/pkg/remoteenforcer/internal/client"
	"go.aporeto.io/trireme-lib/v11/controller/pkg/remoteenforcer/internal/tokenissuer"
	"go.aporeto.io/trireme-lib/v11/policy"

	"go.aporeto.io/trireme-lib/v11/controller/internal/enforcer"
	"go.aporeto.io/trireme-lib/v11/controller/internal/enforcer/utils/rpcwrapper"
	"go.aporeto.io/trireme-lib/v11/controller/internal/supervisor"
	"go.aporeto.io/trireme-lib/v11/controller/pkg/packetprocessor"
	"go.aporeto.io/trireme-lib/v11/controller/pkg/remoteenforcer/internal/statscollector"
	"go.uber.org/zap"
)

//nolint:unused // seem to be used by the test
var (
	createEnforcer   = enforcer.New
	createSupervisor = supervisor.NewSupervisor
)

// newServer is a fake implementation for building on darwin/windows.
func newRemoteEnforcer(
	ctx context.Context,
	cancel context.CancelFunc,
	service packetprocessor.PacketProcessor,
	rpcHandle rpcwrapper.RPCServer,
	secret string,
	statsClient client.Reporter,
	collector statscollector.Collector,
	reportsClient client.Reporter,
	tokenIssuer tokenissuer.TokenClient,
	zapConfig zap.Config,
	enforcerType policy.EnforcerType,
	agentVersion semver.Version,
) (*RemoteEnforcer, error) {
	return nil, nil
}

// LaunchRemoteEnforcer is a fake implementation for building on darwin.
func LaunchRemoteEnforcer(service packetprocessor.PacketProcessor, zapConfig zap.Config, agentVersion semver.Version) error {
	return nil
}

// InitEnforcer is a function called from the controller using RPC. It intializes data structure required by the
// remote enforcer
func (s *RemoteEnforcer) InitEnforcer(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	return nil
}

// Enforce this method calls the enforce method on the enforcer created during initenforcer
func (s *RemoteEnforcer) Enforce(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	return nil
}

// Unenforce this method calls the unenforce method on the enforcer created from initenforcer
func (s *RemoteEnforcer) Unenforce(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	return nil
}

// EnforcerExit this method is called when  we received a killrpocess message from the controller
// This allows a graceful exit of the enforcer
func (s *RemoteEnforcer) EnforcerExit(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	return nil
}

// EnableDatapathPacketTracing enable nfq datapath packet tracing
func (s *RemoteEnforcer) EnableDatapathPacketTracing(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	return nil
}

// EnableIPTablesPacketTracing enables iptables trace packet tracing
func (s *RemoteEnforcer) EnableIPTablesPacketTracing(req rpcwrapper.Request, resp *rpcwrapper.Response) error {
	return nil
}
