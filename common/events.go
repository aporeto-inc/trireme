package common

import (
	"context"
)

// TriremeSocket is the standard API server Trireme socket path
// it is set via ConfigureTriremeSocketPath() and canonicalized with
// utils.GetPathOnHostViaProcRoot() at point of use
var TriremeSocket = "/var/run/trireme.sock"

// ConfigureTriremeSocketPath updates the TriremeSocket path
func ConfigureTriremeSocketPath(path string) {
	TriremeSocket = path
}

// PUType defines the PU type
type PUType int

const (
	// ContainerPU indicates that this PU is a container
	ContainerPU PUType = iota
	// LinuxProcessPU indicates that this is Linux process
	LinuxProcessPU
	// WindowsProcessPU indicates that this is Windows process
	WindowsProcessPU
	// HostPU is a host wrapping PU
	HostPU
	// HostNetworkPU is a PU for a network service in a host
	HostNetworkPU
	// KubernetesPU indicates that this is KubernetesPod
	KubernetesPU
	// TransientPU PU -- placeholder to run processing. This should not
	// be inserted in any cache. This is valid only for processing a packet
	TransientPU
)

const (
	// TriremeCgroupPath is the standard Trireme cgroup path
	TriremeCgroupPath = "/trireme/"

	// TriremeDockerHostNetwork is the path for Docker HostNetwork container based activations
	TriremeDockerHostNetwork = "/trireme_docker_hostnet/"
)

// EventInfo is a generic structure that defines all the information related to a PU event.
// EventInfo should be used as a normalized struct container that
type EventInfo struct {

	// EventType refers to one of the standard events that Trireme handles.
	EventType Event `json:"eventtype,omitempty"`

	// PUType is the the type of the PU
	PUType PUType `json:"putype,omitempty"`

	// The PUID is a unique value for the Processing Unit. Ideally this should be the UUID.
	PUID string `json:"puid,omitempty"`

	// The Name is a user-friendly name for the Processing Unit.
	Name string `json:"name,omitempty"`

	// The Executable is the executable name  for the Processing Unit.
	Executable string `json:"executable,omitempty"`

	// Tags represents the set of MetadataTags associated with this PUID.
	Tags []string `json:"tags,omitempty"`

	// The path for the Network Namespace.
	NS string `json:"namespace,omitempty"`

	// Cgroup is the path to the cgroup - used for deletes
	Cgroup string `json:"cgroup,omitempty"`

	// IPs is a map of all the IPs that fully belong to this processing Unit.
	IPs map[string]string `json:"ipaddressesutype,omitempty"`

	// Services is a list of services of interest - for host control
	Services []Service `json:"services,omitempty"`

	// The PID is the PID on the system where this Processing Unit is running.
	PID int32 `json:"pid,omitempty"`

	// HostService indicates that the request is for the root namespace
	HostService bool `json:"hostservice,omitempty"`

	// AutoPort indicates that the PU will have auto port feature enabled
	AutoPort bool `json:"autoport,omitempty"`

	// NetworkOnlyTraffic indicates that traffic towards the applications must be controlled.
	NetworkOnlyTraffic bool `json:"networktrafficonly,omitempty"`

	// Root indicates that this request is coming from a roor user. Its overwritten by the enforcer
	Root bool `json:"root,omitempty"`
}

// Event represents the event picked up by the monitor.
type Event string

// Values of the events
const (
	EventStart   Event = "start"
	EventStop    Event = "stop"
	EventUpdate  Event = "update"
	EventCreate  Event = "create"
	EventDestroy Event = "destroy"
	EventPause   Event = "pause"
	EventUnpause Event = "unpause"
	EventResync  Event = "resync"
)

var (
	// EventMap used for validations
	EventMap = map[Event]*struct{}{
		"start":   nil,
		"stop":    nil,
		"update":  nil,
		"create":  nil,
		"destroy": nil,
		"pause":   nil,
		"unpause": nil,
		"resync":  nil,
	}
)

// EventResponse encapsulate the error response if any.
type EventResponse struct {
	Error string
}

// A EventHandler is type of event handler functions.
type EventHandler func(ctx context.Context, event *EventInfo) error

// A State describes the state of the PU.
type State int

const (
	// StateStarted is the state of a started PU.
	StateStarted State = iota + 1

	// StateStopped is the state of stopped PU.
	StateStopped

	// StatePaused is the state of a paused PU.
	StatePaused

	// StateDestroyed is the state of destroyed PU.
	StateDestroyed

	// StateUnknwown is the state of PU in an unknown state.
	StateUnknwown
)
