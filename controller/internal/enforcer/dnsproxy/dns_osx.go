// +build !linux darwin

package dnsproxy

import (
	"go.aporeto.io/trireme-lib/collector"
	"go.aporeto.io/trireme-lib/controller/pkg/flowtracking"
	"go.aporeto.io/trireme-lib/utils/cache"
)

// Proxy struct represents the object for dns proxy
type Proxy struct {
}

// New creates an instance of the dns proxy
func New(puFromID cache.DataStore, conntrack flowtracking.FlowClient, c collector.EventCollector) *Proxy {
	return &Proxy{}
}

// ShutdownDNS shuts down the dns server for contextID
func (p *Proxy) ShutdownDNS(contextID string) {

}

// StartDNSServer starts the dns server on the port provided for contextID
func (p *Proxy) StartDNSServer(contextID, port string) error {
	return nil
}
