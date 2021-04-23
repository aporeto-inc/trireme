package connection

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"go.aporeto.io/enforcerd/trireme-lib/controller/internal/enforcer/nfqdatapath/afinetrawsocket"
	"go.aporeto.io/enforcerd/trireme-lib/controller/internal/enforcer/utils/ephemeralkeys"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/counters"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/packet"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/pingconfig"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/pucontext"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/secrets"
	"go.aporeto.io/enforcerd/trireme-lib/controller/pkg/tokens"
	"go.aporeto.io/enforcerd/trireme-lib/policy"
	"go.aporeto.io/enforcerd/trireme-lib/utils/cache"
	"go.aporeto.io/enforcerd/trireme-lib/utils/crypto"
	"go.uber.org/zap"
)

// TCPFlowState identifies the constants of the state of a TCP connectioncon
type TCPFlowState int

// UDPFlowState identifies the constants of the state of a UDP connection.
type UDPFlowState int

// ProxyConnState identifies the constants of the state of a proxied connection
type ProxyConnState int

const (

	// TCPSynSend is the state where the Syn packets has been send, but no response has been received
	TCPSynSend TCPFlowState = iota

	// TCPSynReceived indicates that the syn packet has been received
	TCPSynReceived

	// TCPSynAckSend indicates that the SynAck packet has been send
	TCPSynAckSend

	// TCPSynAckReceived is the state where the SynAck has been received
	TCPSynAckReceived

	// TCPAckSend indicates that the ack packets has been sent
	TCPAckSend

	// TCPAckProcessed is the state that the negotiation has been completed
	TCPAckProcessed

	// TCPData indicates that the packets are now data packets
	TCPData

	// UnknownState indicates that this an existing connection in the unknown state.
	UnknownState
)

const (
	// ClientTokenSend Init token send for client
	ClientTokenSend ProxyConnState = iota

	// ServerReceivePeerToken -- waiting to receive peer token
	ServerReceivePeerToken

	// ServerSendToken -- Send our own token and the client tokens
	ServerSendToken

	// ClientPeerTokenReceive -- Receive signed tokens from server
	ClientPeerTokenReceive

	// ClientSendSignedPair -- Sign the (token/nonce pair) and send
	ClientSendSignedPair

	// ServerAuthenticatePair -- Authenticate pair of tokens
	ServerAuthenticatePair
)

const (
	// UDPStart is the state where a syn will be sent.
	UDPStart UDPFlowState = iota

	// UDPClientSendSyn is the state where a syn has been sent.
	UDPClientSendSyn

	// UDPClientSendAck  is the state where application side has send the ACK.
	UDPClientSendAck

	// UDPReceiverSendSynAck is the state where syn ack packet has been sent.
	UDPReceiverSendSynAck

	// UDPReceiverProcessedAck is the state that the negotiation has been completed.
	UDPReceiverProcessedAck

	// UDPData is the state where data is being transmitted.
	UDPData

	// UDPRST is the state when we received rst from peer. This connection is dead
	UDPRST
)

// MaximumUDPQueueLen is the maximum number of UDP packets buffered.
const MaximumUDPQueueLen = 50

// AuthInfo keeps authentication information about a connection
type AuthInfo struct {
	Nonce                        [tokens.NonceLength]byte
	RemoteNonce                  []byte
	RemoteContextID              string
	RemoteIP                     string
	RemotePort                   string
	LocalDatapathPrivateKey      *ephemeralkeys.PrivateKey
	SecretKey                    []byte
	LocalDatapathPublicKeyV1     []byte
	LocalDatapathPublicKeySignV1 []byte
	LocalDatapathPublicKeyV2     []byte
	LocalDatapathPublicKeySignV2 []byte
	ConnectionClaims             tokens.ConnectionClaims
	SynAckToken                  []byte
	AckToken                     []byte
	Proto314                     bool
}

//TCPTuple contains the 4 tuple for tcp connection
type TCPTuple struct {
	SourceAddress      net.IP
	DestinationAddress net.IP
	SourcePort         uint16
	DestinationPort    uint16
}

// TCPConnection is information regarding TCP Connection
type TCPConnection struct {
	sync.RWMutex

	state TCPFlowState
	Auth  AuthInfo

	// ServiceData allows services to associate state with a connection
	ServiceData interface{}

	// Context is the pucontext.PUContext that is associated with this connection
	// Minimizes the number of caches and lookups
	Context *pucontext.PUContext

	// TimeOut signals the timeout to be used by the state machines
	TimeOut time.Duration

	// ServiceConnection indicates that this connection is handled by a service
	ServiceConnection bool

	// LoopbackConnection indicates that this connections is within the same pu context.
	loopbackConnection bool

	// ReportFlowPolicy holds the last matched observed policy
	ReportFlowPolicy *policy.FlowPolicy

	// PacketFlowPolicy holds the last matched actual policy
	PacketFlowPolicy *policy.FlowPolicy

	// MarkForDeletion -- this is is used only in conjunction with serviceconnection. Its a hint for us if we have a fin for an earlier connection
	// and this is reused port flow.
	MarkForDeletion bool

	RetransmittedSynAck bool

	expiredConnection bool

	// TCPtuple is tcp tuple
	TCPtuple *TCPTuple

	// PingConfig is the config that holds ping related information.
	PingConfig *pingconfig.PingConfig

	Secrets secrets.Secrets

	SourceController      string
	DestinationController string
	initialSequenceNumber uint32
	timer                 *time.Timer
	counter               uint32
	reportReason          string
	connectionTimeout     time.Duration
	EncodedBuf            [tokens.ClaimsEncodedBufSize]byte
}

//DefaultConnectionTimeout is used as the timeout for connection in the cache.
var DefaultConnectionTimeout = 24 * time.Second

//StartTimer starts the timer for 24 seconds and
//on expiry will call the function passed in the argument.
func (c *TCPConnection) StartTimer(f func()) {
	if c.timer == nil {
		c.timer = time.AfterFunc(c.connectionTimeout, f)
	} else {
		c.timer.Reset(c.connectionTimeout)
	}
}

//StopTimer will stop the timer in the connection object.
func (c *TCPConnection) StopTimer() {
	if c.timer != nil {
		c.timer.Stop()
	}
}

//ResetTimer resets the timer
func (c *TCPConnection) ResetTimer(newTimeout time.Duration) {
	if c.timer != nil {
		c.timer.Reset(newTimeout)
	}
}

func (tcpTuple *TCPTuple) String() string {
	return "sip: " + tcpTuple.SourceAddress.String() + " dip: " + tcpTuple.DestinationAddress.String() + " sport: " + strconv.Itoa(int(tcpTuple.SourcePort)) + " dport: " + strconv.Itoa(int(tcpTuple.DestinationPort))
}

// String returns a printable version of connection
func (c *TCPConnection) String() string {
	return fmt.Sprintf("state:%d auth: %+v", c.state, c.Auth)
}

// PingEnabled returns true if ping is enabled for this connection
func (c *TCPConnection) PingEnabled() bool {
	return c.PingConfig != nil
}

// GetState is used to return the state
func (c *TCPConnection) GetState() TCPFlowState {
	return c.state
}

// GetStateString is used to return the state as string
func (c *TCPConnection) GetStateString() string {

	switch c.state {
	case TCPSynSend:
		return "TCPSynSend"

	case TCPSynReceived:
		return "TCPSynReceived"

	case TCPSynAckSend:
		return "TCPSynAckSend"

	case TCPSynAckReceived:
		return "TCPSynAckReceived"

	case TCPAckSend:
		return "TCPAckSend"

	case TCPAckProcessed:
		return "TCPAckProcessed"

	case TCPData:
		return "TCPData"

	default:
		return "UnknownState"
	}
}

// GetInitialSequenceNumber returns the initial sequence number that was found on the syn packet
// corresponding to this TCP Connection
func (c *TCPConnection) GetInitialSequenceNumber() uint32 {
	return c.initialSequenceNumber
}

// GetMarkForDeletion returns the state of markForDeletion flag
func (c *TCPConnection) GetMarkForDeletion() bool {
	c.RLock()
	defer c.RUnlock()
	return c.MarkForDeletion
}

// IncrementCounter increments counter for this connection
func (c *TCPConnection) IncrementCounter() {
	atomic.AddUint32(&c.counter, 1)
}

// GetCounterAndReset returns the counter and resets it to zero
func (c *TCPConnection) GetCounterAndReset() uint32 {
	return atomic.SwapUint32(&c.counter, 0)
}

// GetReportReason returns the reason for reporting this connection
func (c *TCPConnection) GetReportReason() string {

	c.RLock()
	defer c.RUnlock()

	return c.reportReason
}

// SetReportReason sets the reason for reporting this connection
func (c *TCPConnection) SetReportReason(reason string) {

	c.Lock()
	c.reportReason = reason
	c.Unlock()
}

// SetState is used to setup the state for the TCP connection
func (c *TCPConnection) SetState(state TCPFlowState) {
	c.state = state
}

// Cleanup will provide information when a connection is removed by a timer.
func (c *TCPConnection) Cleanup() {

	c.Lock()
	if !c.expiredConnection && c.state != TCPData {
		c.expiredConnection = true
		if c.Context != nil {
			c.Context.Counters().IncrementCounter(counters.ErrTCPConnectionsExpired)
		}
	}
	c.Unlock()
}

// SetLoopbackConnection sets LoopbackConnection field.
func (c *TCPConnection) SetLoopbackConnection(isLoopback bool) {
	// Logging information
	c.loopbackConnection = isLoopback
}

// IsLoopbackConnection sets LoopbackConnection field.
func (c *TCPConnection) IsLoopbackConnection() bool {
	// Logging information
	return c.loopbackConnection
}

// SetLoopbackConnection sets LoopbackConnection field.
func (c *UDPConnection) SetLoopbackConnection(isLoopback bool) {
	// Logging information
	c.loopbackConnection = isLoopback
}

// IsLoopbackConnection sets LoopbackConnection field.
func (c *UDPConnection) IsLoopbackConnection() bool {
	// Logging information
	return c.loopbackConnection
}

// NewTCPConnection returns a TCPConnection information struct
func NewTCPConnection(context *pucontext.PUContext, p *packet.Packet) *TCPConnection {

	var initialSeqNumber uint32

	// Default tuple in case the packet is nil.
	tuple := &TCPTuple{}

	if p != nil {
		tuple.SourceAddress = p.SourceAddress()
		tuple.DestinationAddress = p.DestinationAddress()
		tuple.SourcePort = p.SourcePort()
		tuple.DestinationPort = p.DestPort()
		initialSeqNumber = p.TCPSequenceNumber()
	}

	tcp := &TCPConnection{
		state:                 TCPSynSend,
		Context:               context,
		Auth:                  AuthInfo{},
		initialSequenceNumber: initialSeqNumber,
		TCPtuple:              tuple,
		connectionTimeout:     DefaultConnectionTimeout,
	}

	crypto.Nonce().GenerateNonce16Bytes(tcp.Auth.Nonce[:])
	return tcp

}

// ProxyConnection is a record to keep state of proxy auth
type ProxyConnection struct {
	sync.Mutex

	state            ProxyConnState
	Auth             AuthInfo
	ReportFlowPolicy *policy.FlowPolicy
	PacketFlowPolicy *policy.FlowPolicy
	reported         bool
	Secrets          secrets.Secrets
}

// NewProxyConnection returns a new Proxy Connection
func NewProxyConnection(keyPair ephemeralkeys.KeyAccessor) *ProxyConnection {

	p := &ProxyConnection{
		state: ClientTokenSend,
		Auth: AuthInfo{
			LocalDatapathPublicKeyV1: keyPair.DecodingKeyV1(),
			LocalDatapathPrivateKey:  keyPair.PrivateKey(),
		},
	}

	crypto.Nonce().GenerateNonce16Bytes(p.Auth.Nonce[:])

	return p
}

// GetState returns the state of a proxy connection
func (c *ProxyConnection) GetState() ProxyConnState {

	return c.state
}

// SetState is used to setup the state for the Proxy Connection
func (c *ProxyConnection) SetState(state ProxyConnState) {

	c.state = state
}

// SetReported sets the flag to reported when the conn is reported
func (c *ProxyConnection) SetReported(reported bool) {
	c.reported = reported
}

// UDPConnection is information regarding UDP connection.
type UDPConnection struct {
	sync.RWMutex

	state   UDPFlowState
	Context *pucontext.PUContext
	Auth    AuthInfo

	ReportFlowPolicy *policy.FlowPolicy
	PacketFlowPolicy *policy.FlowPolicy
	// ServiceData allows services to associate state with a connection
	ServiceData interface{}

	// PacketQueue indicates app UDP packets queued while authorization is in progress.
	PacketQueue chan *packet.Packet
	Writer      afinetrawsocket.SocketWriter
	// ServiceConnection indicates that this connection is handled by a service
	ServiceConnection bool
	// LoopbackConnection indicates that this connections is within the same pu context.
	loopbackConnection bool
	// Stop channels for restransmissions
	synStop    chan bool
	synAckStop chan bool
	ackStop    chan bool

	TestIgnore           bool
	udpQueueFullDropCntr uint64
	expiredConnection    bool

	Secrets secrets.Secrets

	SourceController      string
	DestinationController string
	EncodedBuf            [tokens.ClaimsEncodedBufSize]byte
}

// NewUDPConnection returns UDPConnection struct.
func NewUDPConnection(context *pucontext.PUContext, writer afinetrawsocket.SocketWriter) *UDPConnection {

	u := &UDPConnection{
		state:       UDPStart,
		Context:     context,
		PacketQueue: make(chan *packet.Packet, MaximumUDPQueueLen),
		Writer:      writer,
		Auth:        AuthInfo{},
		synStop:     make(chan bool),
		synAckStop:  make(chan bool),
		ackStop:     make(chan bool),
		TestIgnore:  true,
	}

	crypto.Nonce().GenerateNonce16Bytes(u.Auth.Nonce[:])
	return u
}

// SynStop issues a stop on the synStop channel.
func (c *UDPConnection) SynStop() {
	select {
	case c.synStop <- true:
	default:
		zap.L().Debug("Packet loss - channel was already done")
	}

}

// SynAckStop issues a stop in the synAckStop channel.
func (c *UDPConnection) SynAckStop() {
	select {
	case c.synAckStop <- true:
	default:
		zap.L().Debug("Packet loss - channel was already done")
	}
}

// AckStop issues a stop in the Ack channel.
func (c *UDPConnection) AckStop() {
	select {
	case c.ackStop <- true:
	default:
		zap.L().Debug("Packet loss - channel was already done")
	}

}

// SynChannel returns the SynStop channel.
func (c *UDPConnection) SynChannel() chan bool {
	return c.synStop

}

// SynAckChannel returns the SynAck stop channel.
func (c *UDPConnection) SynAckChannel() chan bool {
	return c.synAckStop
}

// AckChannel returns the Ack stop channel.
func (c *UDPConnection) AckChannel() chan bool {
	return c.ackStop
}

// GetState is used to get state of UDP Connection.
func (c *UDPConnection) GetState() UDPFlowState {
	return c.state
}

// SetState is used to setup the state for the UDP Connection.
func (c *UDPConnection) SetState(state UDPFlowState) {
	c.state = state
}

// QueuePackets queues UDP packets till the flow is authenticated.
func (c *UDPConnection) QueuePackets(udpPacket *packet.Packet) (err error) {
	buffer := make([]byte, len(udpPacket.GetBuffer(0)))
	copy(buffer, udpPacket.GetBuffer(0))

	copyPacket, err := packet.New(packet.PacketTypeApplication, buffer, udpPacket.Mark, true)
	if err != nil {
		return fmt.Errorf("Unable to copy packets to queue:%s", err)
	}

	if udpPacket.PlatformMetadata != nil {
		copyPacket.PlatformMetadata = udpPacket.PlatformMetadata.Clone()
	}

	select {
	case c.PacketQueue <- copyPacket:
	default:
		// connection object is always locked.
		c.udpQueueFullDropCntr++
		return fmt.Errorf("Queue is full")
	}

	return nil
}

// DropPackets drops packets on errors during Authorization.
func (c *UDPConnection) DropPackets() {
	close(c.PacketQueue)
	c.PacketQueue = make(chan *packet.Packet, MaximumUDPQueueLen)
}

// ReadPacket reads a packet from the queue.
func (c *UDPConnection) ReadPacket() *packet.Packet {
	select {
	case p := <-c.PacketQueue:
		return p
	default:
		return nil
	}
}

// Cleanup is called on cache expiry of the connection to record incomplete connections
func (c *UDPConnection) Cleanup() {

	c.Lock()
	if !c.expiredConnection && c.state != UDPData {
		c.expiredConnection = true
		if c.Context != nil {
			c.Context.Counters().IncrementCounter(counters.ErrUDPConnectionsExpired)
		}
	}
	c.Unlock()
}

// String returns a printable version of connection
func (c *UDPConnection) String() string {

	return fmt.Sprintf("udp-conn state:%d auth: %+v", c.state, c.Auth)
}

// UDPConnectionExpirationNotifier expiration notifier when cache entry expires
func UDPConnectionExpirationNotifier(c cache.DataStore, id interface{}, item interface{}) {

	if conn, ok := item.(*UDPConnection); ok {
		conn.Cleanup()
	}
}

// ChangeConnectionTimeout is used by test code to change the default
// connection timeout
func (c *TCPConnection) ChangeConnectionTimeout(t time.Duration) {
	// Logging information
	c.connectionTimeout = t
}
