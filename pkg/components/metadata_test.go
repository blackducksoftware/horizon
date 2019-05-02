/*
Copyright (C) 2019 Synopsys, Inc.

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
	"testing"

	"k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FakeNamespace struct {
	*v1.Namespace
	MetadataFuncs
}

func TestAddAnnotations(t *testing.T) {
	testcases := []struct {
		Name     string
		Existing map[string]string
		New      map[string]string
		Expected map[string]string
	}{
		{
			Name:     "no existing annotations",
			Existing: map[string]string{},
			New:      map[string]string{"key1": "value1", "key2": "value2"},
			Expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			Name:     "exiting annotations, no overlap",
			Existing: map[string]string{"existingKey": "existingValue"},
			New:      map[string]string{"key1": "value1", "key2": "value2"},
			Expected: map[string]string{"existingKey": "existingValue", "key1": "value1", "key2": "value2"},
		},
		{
			Name:     "existing annotations with overlap",
			Existing: map[string]string{"key1": "existing"},
			New:      map[string]string{"key1": "value1", "key2": "value2"},
			Expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tc := range testcases {
		n := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: tc.Existing,
			},
		}
		c := FakeNamespace{&n, MetadataFuncs{&n}}
		c.AddAnnotations(tc.New)

		if !reflect.DeepEqual(tc.Expected, c.Annotations) {
			t.Errorf("%s: wrong annotations.  Expected %+v, got %+v", tc.Name, tc.Expected, c.Annotations)
		}
	}
}

func TestRemoveAnnotations(t *testing.T) {
	testcases := []struct {
		Name     string
		Existing map[string]string
		Remove   []string
		Expected map[string]string
	}{
		{
			Name:     "no existing annotations",
			Existing: map[string]string{},
			Remove:   []string{"invalid"},
			Expected: map[string]string{},
		},
		{
			Name:     "exiting annotations",
			Existing: map[string]string{"key1": "value1", "key2": "value2"},
			Remove:   []string{"key1"},
			Expected: map[string]string{"key2": "value2"},
		},
		{
			Name:     "annotation to remove doesn't exist",
			Existing: map[string]string{"key1": "existing"},
			Remove:   []string{"invalid"},
			Expected: map[string]string{"key1": "existing"},
		},
	}

	for _, tc := range testcases {
		n := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: tc.Existing,
			},
		}
		c := FakeNamespace{&n, MetadataFuncs{&n}}
		c.RemoveAnnotations(tc.Remove)

		if !reflect.DeepEqual(tc.Expected, c.Annotations) {
			t.Errorf("%s: wrong annotations.  Expected %+v, got %+v", tc.Name, tc.Expected, c.Annotations)
		}
	}
}

func TestAddLabels(t *testing.T) {
	testcases := []struct {
		Name     string
		Existing map[string]string
		New      map[string]string
		Expected map[string]string
	}{
		{
			Name:     "no existing labels",
			Existing: map[string]string{},
			New:      map[string]string{"key1": "value1", "key2": "value2"},
			Expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			Name:     "exiting labels, no overlap",
			Existing: map[string]string{"existingKey": "existingValue"},
			New:      map[string]string{"key1": "value1", "key2": "value2"},
			Expected: map[string]string{"existingKey": "existingValue", "key1": "value1", "key2": "value2"},
		},
		{
			Name:     "existing labels with overlap",
			Existing: map[string]string{"key1": "existing"},
			New:      map[string]string{"key1": "value1", "key2": "value2"},
			Expected: map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tc := range testcases {
		n := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Labels: tc.Existing,
			},
		}
		c := FakeNamespace{&n, MetadataFuncs{&n}}
		c.AddLabels(tc.New)

		if !reflect.DeepEqual(tc.Expected, c.Labels) {
			t.Errorf("%s: wrong labels.  Expected %+v, got %+v", tc.Name, tc.Expected, c.Labels)
		}
	}
}

func TestRemoveLabels(t *testing.T) {
	testcases := []struct {
		Name     string
		Existing map[string]string
		Remove   []string
		Expected map[string]string
	}{
		{
			Name:     "no existing labels",
			Existing: map[string]string{},
			Remove:   []string{"invalid"},
			Expected: map[string]string{},
		},
		{
			Name:     "exiting labels",
			Existing: map[string]string{"key1": "value1", "key2": "value2"},
			Remove:   []string{"key1"},
			Expected: map[string]string{"key2": "value2"},
		},
		{
			Name:     "label to remove doesn't exist",
			Existing: map[string]string{"key1": "existing"},
			Remove:   []string{"invalid"},
			Expected: map[string]string{"key1": "existing"},
		},
	}

	for _, tc := range testcases {
		n := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Labels: tc.Existing,
			},
		}
		c := FakeNamespace{&n, MetadataFuncs{&n}}
		c.RemoveLabels(tc.Remove)

		if !reflect.DeepEqual(tc.Expected, c.Labels) {
			t.Errorf("%s: wrong labels.  Expected %+v, got %+v", tc.Name, tc.Expected, c.Labels)
		}
	}
}

func TestAddFinalizers(t *testing.T) {
	testcases := []struct {
		Name     string
		Existing []string
		New      []string
		Expected []string
	}{
		{
			Name:     "no existing finalizers",
			Existing: []string{},
			New:      []string{"final1", "final2"},
			Expected: []string{"final1", "final2"},
		},
		{
			Name:     "exiting finalizers, no overlap",
			Existing: []string{"existing"},
			New:      []string{"final1", "final2"},
			Expected: []string{"existing", "final1", "final2"},
		},
		{
			Name:     "existing finalizers with overlap",
			Existing: []string{"final1"},
			New:      []string{"final1", "final2"},
			Expected: []string{"final1", "final2"},
		},
	}

	for _, tc := range testcases {
		n := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Finalizers: tc.Existing,
			},
		}
		c := FakeNamespace{&n, MetadataFuncs{&n}}
		c.AddFinalizers(tc.New)

		if !reflect.DeepEqual(tc.Expected, c.Finalizers) {
			t.Errorf("%s: wrong finalizers.  Expected %+v, got %+v", tc.Name, tc.Expected, c.Finalizers)
		}
	}
}

func TestRemoveFinalizers(t *testing.T) {
	testcases := []struct {
		Name     string
		Existing []string
		Remove   string
		Expected []string
	}{
		{
			Name:     "no existing finalizers",
			Existing: []string{},
			Remove:   "invalid",
			Expected: []string{},
		},
		{
			Name:     "exiting finalizers",
			Existing: []string{"final1", "final2"},
			Remove:   "final1",
			Expected: []string{"final2"},
		},
		{
			Name:     "finalizer to remove doesn't exist",
			Existing: []string{"final1"},
			Remove:   "invalid",
			Expected: []string{"final1"},
		},
		{
			Name:     "remove last finalizer in the list",
			Existing: []string{"final1", "final2"},
			Remove:   "final2",
			Expected: []string{"final1"},
		},
	}

	for _, tc := range testcases {
		n := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Finalizers: tc.Existing,
			},
		}
		c := FakeNamespace{&n, MetadataFuncs{&n}}
		c.RemoveFinalizer(tc.Remove)

		if !reflect.DeepEqual(tc.Expected, c.Finalizers) {
			t.Errorf("%s: wrong finalizers.  Expected %+v, got %+v", tc.Name, tc.Expected, c.Finalizers)
		}
	}
}
