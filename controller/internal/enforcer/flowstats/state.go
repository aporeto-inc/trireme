package flowstats

import (
	"net"
	"net/http"

	"go.aporeto.io/trireme-lib/v11/controller/internal/enforcer/apiauth"

	"go.aporeto.io/trireme-lib/v11/collector"
	"go.aporeto.io/trireme-lib/v11/controller/pkg/packet"
	"go.aporeto.io/trireme-lib/v11/policy"
)

// ConnectionState captures the connection state. This state
// is passed to the RoundTripper for any last minute adjustments.
type ConnectionState struct {
	Stats  *collector.FlowRecord
	Cookie *http.Cookie
}

// NewAppConnectionState will create the initial connection state object.
func NewAppConnectionState(nativeID string, r *http.Request, authRequest *apiauth.Request, resp *apiauth.AppAuthResponse) *ConnectionState {

	sourceIP := "0.0.0.0/0"
	sourcePort := 0
	if sourceAddress, err := net.ResolveTCPAddr("tcp", r.RemoteAddr); err == nil {
		sourceIP = sourceAddress.IP.String()
		sourcePort = sourceAddress.Port
	}

	return &ConnectionState{
		Stats: &collector.FlowRecord{
			ContextID: nativeID,
			Destination: &collector.EndPoint{
				URI:        r.Method + " " + r.RequestURI,
				HTTPMethod: r.Method,
				Type:       collector.EndPointTypeExternalIP,
				Port:       uint16(authRequest.OriginalDestination.Port),
				IP:         authRequest.OriginalDestination.IP.String(),
				ID:         resp.NetworkServiceID,
			},
			Source: &collector.EndPoint{
				Type:       collector.EnpointTypePU,
				ID:         resp.PUContext.ManagementID(),
				IP:         sourceIP,
				Port:       uint16(sourcePort),
				HTTPMethod: r.Method,
				URI:        r.Method + " " + r.RequestURI,
			},
			Action:      resp.Action,
			L4Protocol:  packet.IPProtocolTCP,
			ServiceType: policy.ServiceHTTP,
			ServiceID:   resp.ServiceID,
			Tags:        resp.PUContext.Annotations(),
			Namespace:   resp.PUContext.ManagementNamespace(),
			PolicyID:    resp.NetworkPolicyID,
			Count:       1,
		},
	}
}

// NewNetworkConnectionState will create the initial connection state object.
func NewNetworkConnectionState(nativeID string, userID string, r *apiauth.Request, d *apiauth.NetworkAuthResponse) *ConnectionState {

	var mgmtID, namespace, serviceID string
	var tags *policy.TagStore

	if d.PUContext != nil {
		mgmtID = d.PUContext.ManagementID()
		namespace = d.PUContext.ManagementNamespace()
		tags = d.PUContext.Annotations()
		serviceID = d.ServiceID
	} else {
		mgmtID = collector.DefaultEndPoint
		namespace = collector.DefaultEndPoint
		tags = policy.NewTagStore()
		serviceID = collector.DefaultEndPoint
	}

	sourceType := collector.EndPointTypeExternalIP
	sourceID := collector.DefaultEndPoint
	networkPolicyID := collector.DefaultEndPoint
	action := policy.Reject

	if d != nil {
		sourceType = d.SourceType
		if sourceType == collector.EndpointTypeClaims {
			sourceType = collector.EndPointTypeExternalIP
		}

		switch d.SourceType {
		case collector.EnpointTypePU:
			sourceID = d.SourcePUID
		case collector.EndpointTypeClaims:
			sourceID = d.NetworkServiceID
		default:
			sourceID = d.NetworkServiceID
		}

		if d.NetworkPolicyID != "" {
			networkPolicyID = d.NetworkPolicyID
		}
		action = d.Action
	}

	c := &ConnectionState{
		Stats: &collector.FlowRecord{
			ContextID: nativeID,
			Destination: &collector.EndPoint{
				ID:         mgmtID,
				Type:       collector.EnpointTypePU,
				IP:         r.OriginalDestination.IP.String(),
				Port:       uint16(r.OriginalDestination.Port),
				URI:        r.Method + " " + r.RequestURI,
				HTTPMethod: r.Method,
				UserID:     userID,
			},
			Source: &collector.EndPoint{
				ID:     sourceID,
				Type:   sourceType,
				IP:     r.SourceAddress.IP.String(),
				Port:   uint16(r.SourceAddress.Port),
				UserID: userID,
			},
			Action:      action,
			L4Protocol:  packet.IPProtocolTCP,
			ServiceType: policy.ServiceHTTP,
			PolicyID:    networkPolicyID,
			ServiceID:   serviceID,
			Tags:        tags,
			Namespace:   namespace,
			Count:       1,
		},
	}

	if d != nil && d.Action.Rejected() {
		c.Stats.DropReason = d.DropReason
	}

	c.Cookie = d.Cookie

	return c
}
