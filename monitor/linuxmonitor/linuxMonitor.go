package linuxmonitor

import (
	"fmt"
	"strconv"

	"github.com/aporeto-inc/trireme/monitor/linuxmonitor/cgnetcls"
	"github.com/aporeto-inc/trireme/monitor/rpcmonitor"
	"github.com/aporeto-inc/trireme/policy"
)

// SystemdRPCMetadataExtractor is a systemd based metadata extractor
func SystemdRPCMetadataExtractor(event *rpcmonitor.EventInfo) (*policy.PURuntime, error) {

	if event.Name == "" {
		return nil, fmt.Errorf("EventInfo PU Name is empty")
	}

	if event.PID == "" {
		return nil, fmt.Errorf("EventInfo PID is empty")
	}

	if event.PUID == "" {
		return nil, fmt.Errorf("EventInfo PUID is empty")
	}

	runtimeTags := policy.NewTagsMap(event.Tags)

	options := policy.NewTagsMap(map[string]string{cgnetcls.PortTag: "0"})

	if _, ok := runtimeTags.Tags[cgnetcls.PortTag]; ok {
		options.Tags[cgnetcls.PortTag] = runtimeTags.Tags[cgnetcls.PortTag]
	}

	options.Tags[cgnetcls.CgroupNameTag] = event.PUID
	options.Tags[cgnetcls.CgroupMarkTag] = <-cgnetcls.MarkVal()

	runtimeIps := policy.NewIPMap(map[string]string{"bridge": "0.0.0.0/0"})

	runtimePID, err := strconv.Atoi(event.PID)

	if err != nil {
		return nil, fmt.Errorf("PID is invalid: %s", err)
	}

	return policy.NewPURuntime(event.Name, runtimePID, runtimeTags, runtimeIps, policy.LinuxProcessPU, options), nil
}
