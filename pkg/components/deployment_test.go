/*
Copyright (C) 2018 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package components

import (
	"reflect"
	"strings"
	"testing"

	"github.com/blackducksoftware/horizon/pkg/api"

	"k8s.io/api/core/v1"
)

func createDeployment(pc api.PodConfig, cc api.ContainerConfig) *Deployment {
	replicas := int32(1)
	d := NewDeployment(api.DeploymentConfig{
		Name:      "test",
		Namespace: "testns",
		Replicas:  &replicas,
	})

	pod := NewPod(pc)
	c, _ := NewContainer(cc)
	pod.AddContainer(c)
	d.AddPod(pod)

	return d
}

func TestDeploymentAddPod(t *testing.T) {
	pc := api.PodConfig{
		Name:           "pod",
		Namespace:      "testns",
		ServiceAccount: "testsa",
	}

	cc := api.ContainerConfig{
		Name:  "testcontainer",
		Image: "test",
	}

	d := createDeployment(pc, cc)

	if strings.Compare(d.Spec.Template.Name, pc.Name) != 0 {
		t.Errorf("deployment name is wrong.  Got %s expected %s\n", d.Spec.Template.Name, pc.Name)
	}
	if strings.Compare(d.Spec.Template.Namespace, pc.Namespace) != 0 {
		t.Errorf("deployment namespace is wrong.  Got %s expected %s\n", d.Spec.Template.Namespace, pc.Namespace)
	}
	if strings.Compare(d.Spec.Template.Spec.Containers[0].Name, cc.Name) != 0 {
		t.Errorf("container name in deployment is wrong.  Got %s expected %s\n", d.Spec.Template.Spec.Containers[0].Name, cc.Name)
	}
	if strings.Compare(d.Spec.Template.Spec.Containers[0].Image, cc.Image) != 0 {
		t.Errorf("pod image in deployment is wrong.  Got %s expected %s\n", d.Spec.Template.Spec.Containers[0].Image, cc.Image)
	}
}

func TestDeploymentRemovePod(t *testing.T) {
	pc := api.PodConfig{
		Name: "removepod",
	}

	d := createDeployment(pc, api.ContainerConfig{})
	d.RemovePod(pc.Name)

	if !reflect.DeepEqual(v1.PodTemplateSpec{}, d.Spec.Template) {
		t.Errorf("failed to remove pod %s from deployment", pc.Name)
	}
}
