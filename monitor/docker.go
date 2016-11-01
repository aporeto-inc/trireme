package monitor

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aporeto-inc/trireme/collector"
	"github.com/aporeto-inc/trireme/policy"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/golang/glog"

	dockerClient "github.com/docker/docker/client"
)

// DockerEvent is the type of various docker events.
type DockerEvent string

const (
	// DockerEventStart represents the Docker "start" event.
	DockerEventStart DockerEvent = "start"

	// DockerEventDie represents the Docker "die" event.
	DockerEventDie DockerEvent = "die"

	// DockerEventDestroy represents the Docker "destroy" event.
	DockerEventDestroy DockerEvent = "destroy"

	// DockerEventConnect represents the Docker "connect" event.
	DockerEventConnect DockerEvent = "connect"

	// DockerClientVersion is the version sent out as the client
	DockerClientVersion = "v1.23"
)

// A DockerEventHandler is type of docker event handler functions.
type DockerEventHandler func(event *events.Message) error

// A DockerMetadataExtractor is a function used to extract a *policy.PURuntime from a given
// docker ContainerJSON.
type DockerMetadataExtractor func(*types.ContainerJSON) (*policy.PURuntime, error)

func contextIDFromDockerID(dockerID string) (string, error) {

	if dockerID == "" {
		return "", fmt.Errorf("Empty DockerID String")
	}
	if len(dockerID) < 12 {
		return "", fmt.Errorf("dockerID smaller than 12 characters")
	}
	return dockerID[:12], nil
}

func initDockerClient(socketType string, socketAddress string) (*dockerClient.Client, error) {

	var socket string

	switch socketType {
	case "tcp":
		socket = "https://" + socketAddress

	case "unix":
		// Sanity check that this path exists
		if _, oserr := os.Stat(socketAddress); os.IsNotExist(oserr) {
			return nil, oserr
		}
		socket = "unix://" + socketAddress

	default:
		return nil, fmt.Errorf("Bad socket type %s", socketType)
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-dockerClient-1.0"}
	dockerClient, err := dockerClient.NewClient(socket, DockerClientVersion, nil, defaultHeaders)
	if err != nil {
		return nil, fmt.Errorf("Error creating Docker Client %s", err)
	}

	return dockerClient, nil
}

func defaultDockerMetadataExtractor(info *types.ContainerJSON) (*policy.PURuntime, error) {

	runtimeInfo := policy.NewPURuntime()

	tags := policy.TagsMap{}
	tags["image"] = info.Config.Image
	tags["name"] = info.Name

	for k, v := range info.Config.Labels {
		tags[k] = v
	}

	ipa := map[string]string{}
	ipa["bridge"] = info.NetworkSettings.IPAddress

	runtimeInfo.SetName(info.Name)
	runtimeInfo.SetPid(info.State.Pid)
	runtimeInfo.SetIPAddresses(ipa)
	runtimeInfo.SetTags(tags)

	return runtimeInfo, nil
}

// dockerMonitor implements the connection to Docker and monitoring based on events
type dockerMonitor struct {
	dockerClient       *dockerClient.Client
	metadataExtractor  DockerMetadataExtractor
	handlers           map[DockerEvent]func(event *events.Message) error
	eventnotifications chan *events.Message
	stopprocessor      chan bool
	stoplistener       chan bool
	syncAtStart        bool

	collector collector.EventCollector
	puHandler ProcessingUnitsHandler
}

// NewDockerMonitor returns a pointer to a DockerMonitor initialized with the given
// socketType ('tcp' or 'unix') and socketAddress (a port for 'tcp' or
// a socket file for 'unix').
//
// After creating a new DockerMonitor, call addHandler to install one
// or more callback handlers for the events to monitor. Then call Start.
func NewDockerMonitor(
	socketType string,
	socketAddress string,
	p ProcessingUnitsHandler,
	m DockerMetadataExtractor,
	l collector.EventCollector, syncAtStart bool,
) Monitor {

	cli, err := initDockerClient(socketType, socketAddress)
	if err != nil {
		panic(fmt.Sprintf("Unable to initialize Docker client: %s", err))
	}

	d := &dockerMonitor{
		puHandler:          p,
		collector:          l,
		syncAtStart:        syncAtStart,
		eventnotifications: make(chan *events.Message, 1000),
		handlers:           make(map[DockerEvent]func(event *events.Message) error),
		stoplistener:       make(chan bool),
		stopprocessor:      make(chan bool),
		metadataExtractor:  m,
		dockerClient:       cli,
	}

	// Add handlers for the events that we know how to process
	d.addHandler(DockerEventStart, d.handleStartEvent)
	d.addHandler(DockerEventDie, d.handleDieEvent)
	d.addHandler(DockerEventDestroy, d.handleDestroyEvent)
	d.addHandler(DockerEventConnect, d.handleNetworkConnectEvent)

	return d
}

// addHandler adds a callback handler for the given docker event.
// Interesting event names include 'start' and 'die'. For more on events see
// https://docs.docker.com/engine/reference/api/docker_remote_api/
// under the section 'Docker Events'.
func (d *dockerMonitor) addHandler(event DockerEvent, handler DockerEventHandler) {
	d.handlers[event] = handler
}

// Start will start the DockerPolicy Enforcement.
// It applies a policy to each Container already Up and Running.
// It listens to all ContainerEvents
func (d *dockerMonitor) Start() error {

	glog.Infoln("Starting the docker monitor ...")

	// Starting the eventListener First.
	go d.eventListener()

	//Syncing all Existing containers depending on MonitorSetting
	if d.syncAtStart {
		err := d.syncContainers()
		if err != nil {
			glog.V(1).Infoln("Error Syncing existingContainers: %s", err)
		}
	}

	// Processing the events received duringthe time of Sync.
	go d.eventProcessor()

	return nil
}

// Stop monitoring docker events.
func (d *dockerMonitor) Stop() error {

	glog.Infoln("Stopping the docker monitor ...")
	d.stoplistener <- true
	d.stopprocessor <- true

	return nil
}

// eventProcessor processes docker events
func (d *dockerMonitor) eventProcessor() {

	for {
		select {
		case event := <-d.eventnotifications:
			if event.Action != "" {
				f, present := d.handlers[DockerEvent(event.Action)]
				if present {
					glog.V(6).Infof("Handling docker event [%s].", event.Action)
					err := f(event)
					if err != nil {
						glog.V(1).Infof("Error while handling event [%s]. : %s", event.Action, err)
					}
				} else {
					glog.V(10).Infof("Docker event [%s] not handled.", event.Action)
				}
			}
		case <-d.stopprocessor:
			return
		}
	}
}

// eventListener listens to Docker events from the daemon and passes to
// to the processor through a buffered channel. This minimizes the chances
// that we will miss events because the processor is delayed
func (d *dockerMonitor) eventListener() {

	messages, errs := d.dockerClient.Events(context.Background(), types.EventsOptions{})

	for {
		select {
		case message := <-messages:
			d.eventnotifications <- &message
		case err := <-errs:
			if err != nil && err != io.EOF {
				glog.V(1).Infoln("Received docker event error ", err)
			}
		case stop := <-d.stoplistener:
			if stop {
				return
			}
		}
	}
}

// syncContainers resyncs all the existing containers on the Host, using the
// same process as when a container is initially spawn up
func (d *dockerMonitor) syncContainers() error {

	glog.V(2).Infoln("Syncing all existing containers")

	options := types.ContainerListOptions{All: true}
	containers, err := d.dockerClient.ContainerList(context.Background(), options)
	if err != nil {
		return fmt.Errorf("Error Getting ContainerList: %s", err)
	}

	for _, c := range containers {
		container, err := d.dockerClient.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			glog.V(1).Infof("Error Syncing existing Container: %s", err)
		}
		if err := d.addOrUpdateDockerContainer(&container); err != nil {
			glog.V(1).Infof("Error Syncing existing Container: %s", err)
		}
	}
	return nil
}

func (d *dockerMonitor) addOrUpdateDockerContainer(dockerInfo *types.ContainerJSON) error {

	timeout := time.Second * 0

	if !dockerInfo.State.Running {
		glog.V(2).Infoln("Container is not running - Activation not needed.")
		return nil
	}

	contextID, err := contextIDFromDockerID(dockerInfo.ID)
	if err != nil {
		return fmt.Errorf("Couldn't generate ContextID: %s", err)
	}

	runtimeInfo, err := d.extractMetadata(dockerInfo)

	if err != nil {
		return fmt.Errorf("Error getting some of the Docker primitives")
	}

	ip, ok := runtimeInfo.DefaultIPAddress()
	if !ok || ip == "" {
		glog.V(2).Infof("IP Not present in container, not policing")
		return nil
	}

	returnChan := d.puHandler.HandleCreate(contextID, runtimeInfo)
	if err := <-returnChan; err != nil {
		glog.V(2).Infoln("Setting policy failed. Stopping the container")
		d.dockerClient.ContainerStop(context.Background(), dockerInfo.ID, &timeout)
		d.collector.CollectContainerEvent(contextID, ip, nil, collector.ContainerFailed)
		return fmt.Errorf("Policy cound't be set - container was killed")
	}

	d.collector.CollectContainerEvent(contextID, ip, runtimeInfo.Tags(), collector.ContainerStart)

	return nil
}

func (d *dockerMonitor) removeDockerContainer(dockerID string) error {

	contextID, err := contextIDFromDockerID(dockerID)

	if err != nil {
		return fmt.Errorf("Couldn't generate ContextID: %s", err)
	}

	errchan := d.puHandler.HandleDelete(contextID)
	return <-errchan
}

// ExtractMetadata generates the RuntimeInfo based on Docker primitive
func (d *dockerMonitor) extractMetadata(dockerInfo *types.ContainerJSON) (*policy.PURuntime, error) {

	if dockerInfo == nil {
		return nil, fmt.Errorf("DockerInfo is empty.")
	}

	if d.metadataExtractor != nil {
		return d.metadataExtractor(dockerInfo)
	}

	return defaultDockerMetadataExtractor(dockerInfo)
}

// handleStartEvent will notify the agent immediately about the event in order
//to start the implementation of the functions. The agent must query
//the policy engine for details on what to do with this container.
func (d *dockerMonitor) handleStartEvent(event *events.Message) error {

	timeout := time.Second * 0
	dockerID := event.ID
	contextID, err := contextIDFromDockerID(dockerID)
	if err != nil {
		return fmt.Errorf("Error Generating ContextID: %s", err)
	}

	info, err := d.dockerClient.ContainerInspect(context.Background(), dockerID)
	if err != nil {
		glog.V(2).Infoln("Killing container because inspect returned error")
		//If we see errors, we will kill the container for security reasons.
		d.dockerClient.ContainerStop(context.Background(), dockerID, &timeout)
		d.collector.CollectContainerEvent(contextID, "", nil, collector.ContainerFailed)
		return fmt.Errorf("Cannot read container information. Killing container. ")
	}

	if err := d.addOrUpdateDockerContainer(&info); err != nil {
		glog.V(2).Infof("Error while trying to add container: %s", err)
		return err
	}

	return nil
}

//handleDie event is called when a container dies. It updates the agent
//data structures and stops enforcement.
func (d *dockerMonitor) handleDieEvent(event *events.Message) error {

	dockerID := event.ID
	contextID, err := contextIDFromDockerID(dockerID)
	if err != nil {
		return fmt.Errorf("Error Generating ContextID: %s", err)
	}
	d.collector.CollectContainerEvent(contextID, "", nil, collector.ContainerStop)

	return d.removeDockerContainer(dockerID)
}

// handleDestroyEvent handles destroy events from Docker. This means that the Policy
// can be safely deleted.
func (d *dockerMonitor) handleDestroyEvent(event *events.Message) error {

	dockerID := event.ID
	contextID, err := contextIDFromDockerID(dockerID)
	if err != nil {
		return fmt.Errorf("Error Generating ContextID: %s", err)
	}

	d.collector.CollectContainerEvent(contextID, "", nil, collector.UnknownContainerDelete)
	// Clear the policy cache
	errChan := d.puHandler.HandleDestroy(contextID)
	return <-errChan
}

func (d *dockerMonitor) handleNetworkConnectEvent(event *events.Message) error {

	id := event.Actor.Attributes["container"]

	_, err := d.dockerClient.ContainerInspect(context.Background(), id)
	if err != nil {
		glog.V(2).Infoln("Failed to read the affected container. %s ", err)
		return err
	}

	return nil
}
