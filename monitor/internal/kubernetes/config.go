package kubernetesmonitor

import (
	"github.com/aporeto-inc/trireme-lib/monitor/extractors"
	dockerMonitor "github.com/aporeto-inc/trireme-lib/monitor/internal/docker"
)

// Config is the config for the Kubernetes monitor
type Config struct { // nolint
	DockerConfig dockerMonitor.Config

	Kubeconfig     string
	Nodename       string
	EnableHostPods bool

	KubernetesExtractor extractors.KubernetesMetadataExtractorType
	DockerExtractor     extractors.DockerMetadataExtractor
}

// DefaultConfig provides a default configuration
func DefaultConfig() *Config {
	return &Config{
		KubernetesExtractor: extractors.DefaultKubernetesMetadataExtractor,
		DockerExtractor:     extractors.DefaultMetadataExtractor,
		EnableHostPods:      false,
		Kubeconfig:          "",
		Nodename:            "",
	}
}

// SetupDefaultConfig adds defaults to a partial configuration
func SetupDefaultConfig(kubernetesConfig *Config) *Config {
	return kubernetesConfig
}
