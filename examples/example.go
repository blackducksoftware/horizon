package examples

import (
	"log"
	"strings"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/deployer"
	"k8s.io/client-go/rest"

	"github.com/blackducksoftware/horizon/pkg/components"
)

// Container defines the configuration for a container
type Container struct {
	ContainerConfig       *horizonapi.ContainerConfig
	EnvConfigs            []*horizonapi.EnvConfig
	VolumeMounts          []*horizonapi.VolumeMountConfig
	PortConfig            *horizonapi.PortConfig
	ActionConfig          *horizonapi.ActionConfig
	ReadinessProbeConfigs []*horizonapi.ProbeConfig
	LivenessProbeConfigs  []*horizonapi.ProbeConfig
}

// CreatePod will create a pod to be used later in an RC or deployment, syntactic sugar.
func CreatePod(name string, serviceAccount string, volumes []*components.Volume, containers []*Container, initContainers []*Container, affinityConfigs []horizonapi.AffinityConfig) *components.Pod {
	pod := components.NewPod(horizonapi.PodConfig{
		Name: name,
	})

	if !strings.EqualFold(serviceAccount, "") {
		pod.GetObj().Account = serviceAccount
	}

	for _, volume := range volumes {
		pod.AddVolume(volume)
	}

	pod.AddLabels(map[string]string{
		"app":  name,
		"tier": name,
	})

	for _, affinityConfig := range affinityConfigs {
		pod.AddAffinity(affinityConfig)
	}

	for _, containerConfig := range containers {
		container := CreateContainer(
			containerConfig.ContainerConfig,
			containerConfig.EnvConfigs,
			containerConfig.VolumeMounts,
			containerConfig.PortConfig,
			containerConfig.ActionConfig, containerConfig.LivenessProbeConfigs, containerConfig.ReadinessProbeConfigs)
		pod.AddContainer(container)
	}

	for _, initContainerConfig := range initContainers {
		initContainer := CreateContainer(initContainerConfig.ContainerConfig, initContainerConfig.EnvConfigs, initContainerConfig.VolumeMounts,
			initContainerConfig.PortConfig, initContainerConfig.ActionConfig, initContainerConfig.LivenessProbeConfigs, initContainerConfig.ReadinessProbeConfigs)
		err := pod.AddInitContainer(initContainer)
		if err != nil {
			log.Printf("failed to create the init container because %+v", err)
		}
	}

	return pod
}

// CreateContainer will create a container to be used later in a deployment or rc.  Syntactic sugar.
func CreateContainer(config *horizonapi.ContainerConfig, envs []*horizonapi.EnvConfig, volumeMounts []*horizonapi.VolumeMountConfig, port *horizonapi.PortConfig,
	actionConfig *horizonapi.ActionConfig, livenessProbeConfigs []*horizonapi.ProbeConfig, readinessProbeConfigs []*horizonapi.ProbeConfig) *components.Container {

	container := components.NewContainer(*config)

	for _, env := range envs {
		container.AddEnv(*env)
	}

	for _, volumeMount := range volumeMounts {
		container.AddVolumeMount(*volumeMount)
	}

	container.AddPort(*port)
	if actionConfig != nil {
		container.AddPostStartAction(*actionConfig)
	}

	for _, livenessProbe := range livenessProbeConfigs {
		container.AddLivenessProbe(*livenessProbe)
	}

	for _, readinessProbe := range readinessProbeConfigs {
		container.AddReadinessProbe(*readinessProbe)
	}

	return container
}

// CreateReplicationController will create a replication controller with minimal syntax
func CreateReplicationController(replicationControllerConfig *horizonapi.ReplicationControllerConfig, pod *components.Pod) *components.ReplicationController {
	rc := components.NewReplicationController(*replicationControllerConfig)
	rc.AddLabelSelectors(map[string]string{
		"app":  replicationControllerConfig.Name,
		"tier": replicationControllerConfig.Name,
	})
	rc.AddPod(pod)
	return rc
}

// CreateReplicationControllerFromContainer will create a replication controller with multiple containers inside a pod
func CreateReplicationControllerFromContainer(replicationControllerConfig *horizonapi.ReplicationControllerConfig, serviceAccount string, containers []*Container, volumes []*components.Volume, initContainers []*Container, affinityConfigs []horizonapi.AffinityConfig) *components.ReplicationController {
	pod := CreatePod(replicationControllerConfig.Name, serviceAccount, volumes, containers, initContainers, affinityConfigs)
	rc := CreateReplicationController(replicationControllerConfig, pod)
	return rc
}

// NewDeployerWithDefaultKubeConfig: So you dont have to make a rest config if 
// you just want to compile some horizon examples.
func NewDeployerWithDefaultKubeconfig() (*deployer.Deployer, error) {
	rc := &rest.Config{}
	return deployer.NewDeployer(rc)
}

// NewRestConfig: So you dont have to figure out the right way
// to import the kube client.
func NewRestConfig() *rest.Config {
	rc := &rest.Config{}
	return rc
}
