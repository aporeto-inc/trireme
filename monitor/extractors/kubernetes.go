package extractors

import (
	"fmt"

	"go.aporeto.io/trireme-lib/policy"
	"go.uber.org/zap"
	api "k8s.io/api/core/v1"
)

// KubernetesPodNameIdentifier is the label used by Docker for the K8S pod name.
const KubernetesPodNameIdentifier = "@usr:io.kubernetes.pod.name"

// KubernetesPodNamespaceIdentifier is the label used by Docker for the K8S namespace.
const KubernetesPodNamespaceIdentifier = "@usr:io.kubernetes.pod.namespace"

// KubernetesContainerNameIdentifier is the label used by Docker for the K8S container name.
const KubernetesContainerNameIdentifier = "@usr:io.kubernetes.container.name"

// KubernetesInfraContainerName is the name of the infra POD.
const KubernetesInfraContainerName = "POD"

// UpstreamNameIdentifier is the identifier used to identify the nane on the resulting PU
const UpstreamNameIdentifier = "@k8s:name"

// UpstreamNamespaceIdentifier is the identifier used to identify the nanespace on the resulting PU
const UpstreamNamespaceIdentifier = "@k8s:namespace"

// UserLabelPrefix is the label prefix for all user defined labels
const UserLabelPrefix = "@usr:"

// KubernetesMetadataExtractorType is an extractor function for Kubernetes.
// It takes as parameter a standard Docker runtime and a Pod Kubernetes definition and return a PolicyRuntime
// This extractor also provides an extra boolean parameter that is used as a token to decide if activation is required.
type KubernetesMetadataExtractorType func(runtime policy.RuntimeReader, pod *api.Pod) (*policy.PURuntime, bool, error)

// DefaultKubernetesMetadataExtractor is a default implementation for the medatadata extractor for Kubernetes
// It only activates the POD//INFRA containers and strips all the labels from docker to only keep the ones from Kubernetes
func DefaultKubernetesMetadataExtractor(runtime policy.RuntimeReader, pod *api.Pod) (*policy.PURuntime, bool, error) {

	if runtime == nil {
		return nil, false, fmt.Errorf("empty runtime")
	}

	if pod == nil {
		return nil, false, fmt.Errorf("empty pod")
	}

	// In this specific metadataExtractor we only want to activate the Infra Container for each pod.
	process, err := isPodInfraContainer(runtime)
	if err != nil {
		return nil, false, fmt.Errorf("Error while processing Kubernetes pod %s", err)
	}

	if !process {
		return nil, false, nil
	}

	podLabels := pod.GetLabels()
	if podLabels == nil {
		zap.L().Debug("couldn't get labels.")
		return nil, false, nil
	}

	tags := policy.NewTagStoreFromMap(podLabels)
	tags.AppendKeyValue(UpstreamNameIdentifier, pod.GetName())
	tags.AppendKeyValue(UpstreamNamespaceIdentifier, pod.GetNamespace())

	originalRuntime, ok := runtime.(*policy.PURuntime)
	if !ok {
		return nil, false, fmt.Errorf("Error casting puruntime")
	}

	newRuntime := originalRuntime.Clone()
	newRuntime.SetTags(tags)

	zap.L().Debug("kubernetes runtime tags", zap.String("name", pod.GetName()), zap.String("namespace", pod.GetNamespace()), zap.Strings("tags", newRuntime.Tags().GetSlice()))

	return newRuntime, true, nil
}

// isPodInfraContainer returns true if the runtime represents the infra container for the POD
func isPodInfraContainer(runtime policy.RuntimeReader) (bool, error) {
	// The Infra container can be found by checking env. variable.
	tagContent, ok := runtime.Tag(KubernetesContainerNameIdentifier)
	if !ok || tagContent != KubernetesInfraContainerName {
		return false, nil
	}

	return true, nil
}
