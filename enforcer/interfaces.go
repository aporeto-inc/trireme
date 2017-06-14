package enforcer

import (
	"sync"
	"time"

	"github.com/aporeto-inc/trireme/constants"
	"github.com/aporeto-inc/trireme/enforcer/lookup"
	"github.com/aporeto-inc/trireme/enforcer/utils/fqconfig"
	"github.com/aporeto-inc/trireme/enforcer/utils/packet"
	"github.com/aporeto-inc/trireme/enforcer/utils/secrets"
	"github.com/aporeto-inc/trireme/enforcer/utils/tokens"
	"github.com/aporeto-inc/trireme/policy"
)

// A PolicyEnforcer is implementing the enforcer that will modify//analyze the capture packets
type PolicyEnforcer interface {

	// Enforce starts enforcing policies for the given policy.PUInfo.
	Enforce(contextID string, puInfo *policy.PUInfo) error

	// Unenforce stops enforcing policy for the given IP.
	Unenforce(contextID string) error

	// GetFilterQueue returns the current FilterQueueConfig.
	GetFilterQueue() *fqconfig.FilterQueue

	// Start starts the PolicyEnforcer.
	Start() error

	// Stop stops the PolicyEnforcer.
	Stop() error
}

// PublicKeyAdder register a publicKey for a Node.
type PublicKeyAdder interface {

	// PublicKeyAdd adds the given cert for the given host.
	PublicKeyAdd(host string, cert []byte) error
}

// PacketProcessor is an interface implemented to stitch into our enforcer
type PacketProcessor interface {
	// Initialize  initializes the secrets of the processor
	Initialize(s secrets.Secrets, fq *fqconfig.FilterQueue)

	// PreProcessTCPAppPacket will be called for application packets and return value of false means drop packet.
	PreProcessTCPAppPacket(p *packet.Packet, context *PUContext, conn *TCPConnection) bool

	// PostProcessTCPAppPacket will be called for application packets and return value of false means drop packet.
	PostProcessTCPAppPacket(p *packet.Packet, action interface{}, context *PUContext, conn *TCPConnection) bool

	// PreProcessTCPNetPacket will be called for network packets and return value of false means drop packet
	PreProcessTCPNetPacket(p *packet.Packet, context *PUContext, conn *TCPConnection) bool

	// PostProcessTCPNetPacket will be called for network packets and return value of false means drop packet
	PostProcessTCPNetPacket(p *packet.Packet, action interface{}, claims *tokens.ConnectionClaims, context *PUContext, conn *TCPConnection) bool
}

//Reporting is an interface for using mock
type Reporting interface {
	//reportFlow will be called when the flow is either accepted or rejected
	reportFlow(p *packet.Packet, connection *TCPConnection, sourceID string, destID string, context *PUContext, action string, mode string)
	//reportAcceptedFlow will be called when the flow is accepted
	reportAcceptedFlow(p *packet.Packet, conn *TCPConnection, sourceID string, destID string, context *PUContext)
	//reportRejectedFlow will be called when the flow is rejected
	reportRejectedFlow(p *packet.Packet, conn *TCPConnection, sourceID string, destID string, context *PUContext, mode string)
}

// PUContext holds data indexed by the PU ID
type PUContext struct {
	ID             string
	ManagementID   string
	Identity       *policy.TagsMap
	Annotations    *policy.TagsMap
	AcceptTxtRules *lookup.PolicyDB
	RejectTxtRules *lookup.PolicyDB
	AcceptRcvRules *lookup.PolicyDB
	RejectRcvRules *lookup.PolicyDB
	Extension      interface{}
	IP             string
	Mark           string
	Ports          []string
	PUType         constants.PUType
	synToken       []byte
	synExpiration  time.Time
	sync.Mutex
}
