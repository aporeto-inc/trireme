package dockermonitor

import (
	"go.aporeto.io/enforcerd/trireme-lib/monitor/constants"
	"go.aporeto.io/enforcerd/trireme-lib/monitor/extractors"
)

// Config is the configuration options to start a CNI monitor
type Config struct {
	EventMetadataExtractor   extractors.DockerMetadataExtractor
	SocketType               string
	SocketAddress            string
	SyncAtStart              bool
	DestroyStoppedContainers bool
	ignoreHostModeContainers bool
}

// DefaultConfig provides a default configuration
func DefaultConfig() *Config {
	return &Config{
		EventMetadataExtractor:   extractors.DefaultMetadataExtractor,
		SocketType:               string(constants.DefaultDockerSocketType),
		SocketAddress:            constants.DefaultDockerSocket,
		SyncAtStart:              true,
		ignoreHostModeContainers: true,
	}
}

// SetupDefaultConfig adds defaults to a partial configuration
func SetupDefaultConfig(dockerConfig *Config) *Config {

	defaultConfig := DefaultConfig()

	if dockerConfig.EventMetadataExtractor == nil {
		dockerConfig.EventMetadataExtractor = defaultConfig.EventMetadataExtractor
	}
	if dockerConfig.SocketType == "" {
		dockerConfig.SocketType = defaultConfig.SocketType
	}
	if dockerConfig.SocketAddress == "" {
		dockerConfig.SocketAddress = defaultConfig.SocketAddress
	}
	return dockerConfig
}
