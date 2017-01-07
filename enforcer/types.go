package enforcer

import (
	"github.com/aporeto-inc/trireme/enforcer/lookup"
	"github.com/aporeto-inc/trireme/enforcer/utils/packet"
	"github.com/aporeto-inc/trireme/policy"
)

// FlowState identifies the constants of the state of a connectioncon
type FlowState int

const (

	// SynSend is the state where the Syn packets has been send, but no response has been received
	SynSend FlowState = iota

	//SynReceived indicates that the syn packet has been received
	SynReceived

	//SynAckSend indicates that the SynAck packet has been send
	SynAckSend

	// SynAckReceived is the state where the SynAck has been received
	SynAckReceived

	// AckSend indicates that the ack packets has been send
	AckSend

	// AckProcessed is the state that the negotiation has been completed
	AckProcessed
)

const (
	// DefaultNetwork is the default IP address used when we don't care about IP addresses
	DefaultNetwork = "0.0.0.0/0"
)

var (
	// TransmitterLabel is the name of the label used to identify the Transmitter Context
	TransmitterLabel = "AporetoContextID"
)

// InterfaceStats for interface
type InterfaceStats struct {
	IncomingPackets     uint32
	OutgoingPackets     uint32
	ProtocolDropPackets uint32
	CreateDropPackets   uint32
}

// PacketStats for interface
type PacketStats struct {
	IncomingPackets        uint32
	OutgoingPackets        uint32
	AuthDropPackets        uint32
	ServicePreDropPackets  uint32
	ServicePostDropPackets uint32
}

// FilterQueue captures all the configuration parameters of the NFQUEUEs
type FilterQueue struct {
	// Network Queue is the queue number of the base queue for network packets
	NetworkQueue uint16
	// NetworkQueueSize is the size of the network queue
	NetworkQueueSize uint32
	// NumberOfNetworkQueues is the number of network queues allocated
	NumberOfNetworkQueues uint16
	// ApplicationQueue is the queue number of the first application queue
	ApplicationQueue uint16
	// ApplicationQueueSize is the size of the application queue
	ApplicationQueueSize uint32
	// NumberOfApplicationQueues is the number of queues that must be allocated
	NumberOfApplicationQueues uint16
	// MarkValue is the default mark to set in packets in the RAW chain
	MarkValue int
}

// PUContext holds data indexed by the docker ID
type PUContext struct {
	ID             string
	Identity       *policy.TagsMap
	Annotations    *policy.TagsMap
	acceptTxtRules *lookup.PolicyDB
	rejectTxtRules *lookup.PolicyDB
	acceptRcvRules *lookup.PolicyDB
	rejectRcvRules *lookup.PolicyDB
	Extension      interface{}
}

// StatsPayload holds the payload for statistics
type StatsPayload struct {
	ContextID string
	Tags      *policy.TagsMap
	Action    string
	Mode      string
	Source    string
	Packet    *packet.Packet
}
