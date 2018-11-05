package examples

import (
	"fmt"
	"testing"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/horizon/pkg/deployer"
)

// IntToInt32 will convert from int to int32
func IntToInt32(i int) *int32 {
	j := int32(i)
	return &j
}

// IntToInt64 will convert from int to int64
func IntToInt64(i int) *int64 {
	j := int64(i)
	return &j
}

func TestExport(t *testing.T) {
	d := deployer.NewDeployerExporter()
	webServerContainerConfig := &Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "webserver", Image: fmt.Sprintf("%s/%s/%s-nginx:%s", "reg",
			"docker.io", "pre", "1.0"),
			PullPolicy: horizonapi.PullAlways, MinMem: "1G", MaxMem: "1G", MinCPU: "", MaxCPU: "", UID: IntToInt64(1000)},
		VolumeMounts: []*horizonapi.VolumeMountConfig{
			{Name: "dir-webserver", MountPath: "/opt/blackduck/hub/webserver/security"},
			{Name: "certificate", MountPath: "/tmp/secrets"},
		},
		PortConfig: &horizonapi.PortConfig{ContainerPort: "8080"},
	}
	webserver := CreateReplicationControllerFromContainer(
		&horizonapi.ReplicationControllerConfig{
			Namespace: "ns",
			Name:      "webserver",
			Replicas:  IntToInt32(1)},

		"ns",
		[]*Container{
			webServerContainerConfig},
		[]*components.Volume{},
		[]*Container{}, []horizonapi.AffinityConfig{})
	// log.Infof("webserver : %v\n", webserver.GetObj())
	d.AddReplicationController(webserver)
	for _, t := range d.Export() {
		fmt.Println(t)
	}
}
