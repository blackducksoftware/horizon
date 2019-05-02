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
	"strings"
	"testing"

	"github.com/blackducksoftware/horizon/pkg/api"

	"k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FakeDeployment struct {
	*v1.Deployment
	LabelSelectorFuncs
}

func TestCreateLabelSelector(t *testing.T) {
	testcases := []struct {
		name                string
		labels              map[string]string
		expressions         []api.ExpressionRequirementConfig
		expectedExpressions []metav1.LabelSelectorRequirement
	}{
		{
			name:   "full config",
			labels: map[string]string{"key1": "value1", "key2": "value2"},
			expressions: []api.ExpressionRequirementConfig{
				{
					Key:    "expKey",
					Op:     api.ExpressionRequirementOpIn,
					Values: []string{"expValue1", "expValue2"},
				},
				{
					Key:    "expKey2",
					Op:     api.ExpressionRequirementOpExists,
					Values: []string{"expValue3"},
				},
			},
			expectedExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "expKey",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"expValue1", "expValue2"},
				},
				{
					Key:      "expKey2",
					Operator: metav1.LabelSelectorOpExists,
					Values:   []string{"expValue3"},
				},
			},
		},
		{
			name:                "labels only",
			labels:              map[string]string{"key1": "value1", "key2": "value2"},
			expressions:         []api.ExpressionRequirementConfig{},
			expectedExpressions: []metav1.LabelSelectorRequirement{},
		},
		{
			name:   "expressions only",
			labels: map[string]string{},
			expressions: []api.ExpressionRequirementConfig{
				{
					Key:    "expKey",
					Op:     api.ExpressionRequirementOpIn,
					Values: []string{"expValue1", "expValue2"},
				},
				{
					Key:    "expKey2",
					Op:     api.ExpressionRequirementOpExists,
					Values: []string{"expValue3"},
				},
			},
			expectedExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "expKey",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"expValue1", "expValue2"},
				},
				{
					Key:      "expKey2",
					Operator: metav1.LabelSelectorOpExists,
					Values:   []string{"expValue3"},
				},
			},
		},
		{
			name:                "empty",
			labels:              map[string]string{},
			expressions:         []api.ExpressionRequirementConfig{},
			expectedExpressions: []metav1.LabelSelectorRequirement{},
		},
	}

	for _, tc := range testcases {
		s := createLabelSelector(api.SelectorConfig{Labels: tc.labels, Expressions: tc.expressions})

		if !reflect.DeepEqual(s.MatchLabels, tc.labels) {
			t.Errorf("%s: incorrect labels.  Expected %+v, got %+v", tc.name, tc.labels, s.MatchLabels)
		}

		if !reflect.DeepEqual(s.MatchExpressions, tc.expectedExpressions) {
			t.Errorf("%s: incorrect expressions.  Expected %+v, got %+v", tc.name, tc.expectedExpressions, s.MatchExpressions)
		}
	}
}

func TestCreateLabelSelectorRequirement(t *testing.T) {
	testcases := []struct {
		name       string
		key        string
		values     []string
		op         api.ExpressionRequirementOp
		expectedOp metav1.LabelSelectorOperator
	}{
		{
			name:       "OpIn",
			key:        "key1",
			values:     []string{"val1", "val2"},
			op:         api.ExpressionRequirementOpIn,
			expectedOp: metav1.LabelSelectorOpIn,
		},
		{
			name:       "OpNotIn",
			key:        "key1",
			values:     []string{"val1", "val2"},
			op:         api.ExpressionRequirementOpNotIn,
			expectedOp: metav1.LabelSelectorOpNotIn,
		},
		{
			name:       "OpExists",
			key:        "key1",
			values:     []string{"val1", "val2"},
			op:         api.ExpressionRequirementOpExists,
			expectedOp: metav1.LabelSelectorOpExists,
		},
		{
			name:       "OpDoesNotExist",
			key:        "key1",
			values:     []string{"val1", "val2"},
			op:         api.ExpressionRequirementOpDoesNotExist,
			expectedOp: metav1.LabelSelectorOpDoesNotExist,
		},
		{
			name:   "empty and unset operator",
			key:    "",
			values: []string{},
		},
	}

	for _, tc := range testcases {
		req := createLabelSelectorRequirement(api.ExpressionRequirementConfig{
			Key:    tc.key,
			Op:     tc.op,
			Values: tc.values,
		})

		if !strings.EqualFold(req.Key, tc.key) {
			t.Errorf("%s: incorrect key.  Expected %s, got %s", tc.name, tc.key, req.Key)
		}

		if req.Operator != tc.expectedOp {
			t.Errorf("%s: incorrect operator.  Expected %v got %v", tc.name, tc.expectedOp, req.Operator)
		}

		if !reflect.DeepEqual(req.Values, tc.values) {
			t.Errorf("%s: incorrect values. Got %+v expected %+v", tc.name, tc.values, req.Values)
		}
	}
}

func TestAddMatchLabelsSelectors(t *testing.T) {
	testcases := []struct {
		name      string
		new       map[string]string
		selectors *metav1.LabelSelector
		expected  *metav1.LabelSelector
	}{
		{
			name:      "no existing match label selectors",
			new:       map[string]string{"key1": "value1"},
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{}},
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "value1"}},
		},
		{
			name:      "exiting match label selectors, no overlap",
			new:       map[string]string{"key1": "value1", "key2": "value2"},
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{"existingKey": "existingValue"}},
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{"existingKey": "existingValue", "key1": "value1", "key2": "value2"}},
		},
		{
			name:      "existing match label selectors with overlap",
			new:       map[string]string{"key1": "value1", "key2": "value2"},
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "existing"}},
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "value1", "key2": "value2"}},
		},
		{
			name:      "nil selector",
			new:       map[string]string{"key1": "value1", "key2": "value2"},
			selectors: nil,
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "value1", "key2": "value2"}},
		},
	}

	for _, tc := range testcases {
		d := v1.Deployment{
			Spec: v1.DeploymentSpec{
				Selector: tc.selectors,
			},
		}
		obj := FakeDeployment{&d, LabelSelectorFuncs{&d}}

		obj.AddMatchLabelsSelectors(tc.new)

		if !reflect.DeepEqual(tc.expected, d.Spec.Selector) {
			t.Errorf("%s: wrong match label selectors.  Expected %+v, got %+v", tc.name, tc.expected, d.Spec.Selector)
		}
	}
}

func TestRemoveMatchLabelsSelectors(t *testing.T) {
	testcases := []struct {
		name      string
		selectors *metav1.LabelSelector
		remove    []string
		expected  *metav1.LabelSelector
	}{
		{
			name:      "no existing match label selectors",
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{}},
			remove:    []string{"invalid"},
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{}},
		},
		{
			name:      "exiting match label selectors",
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "value1", "key2": "value2"}},
			remove:    []string{"key1"},
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{"key2": "value2"}},
		},
		{
			name:      "match label selectors to remove doesn't exist",
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "existing"}},
			remove:    []string{"invalid"},
			expected:  &metav1.LabelSelector{MatchLabels: map[string]string{"key1": "existing"}},
		},
		{
			name:      "nil selector",
			selectors: nil,
			remove:    []string{"key1"},
			expected:  nil,
		},
	}

	for _, tc := range testcases {
		d := v1.Deployment{
			Spec: v1.DeploymentSpec{
				Selector: tc.selectors,
			},
		}
		obj := FakeDeployment{&d, LabelSelectorFuncs{&d}}
		obj.RemoveMatchLabelsSelectors(tc.remove)

		if !reflect.DeepEqual(tc.expected, d.Spec.Selector) {
			t.Errorf("%s: wrong match label selectors.  Expected %+v, got %+v", tc.name, tc.expected, d.Spec.Selector)
		}
	}
}

func TestAddMatchExpressionsSelector(t *testing.T) {
	testcases := []struct {
		name      string
		new       api.ExpressionRequirementConfig
		selectors *metav1.LabelSelector
		expected  *metav1.LabelSelector
	}{
		{
			name: "no existing match expressions selectors",
			new: api.ExpressionRequirementConfig{
				Key:    "key1",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			selectors: &metav1.LabelSelector{},
			expected: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value1", "value2"},
				},
			}},
		},
		{
			name: "existing match expressions selectors",
			new: api.ExpressionRequirementConfig{
				Key:    "key1",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			selectors: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key2",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value3", "value4"},
				},
			}},
			expected: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key2",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value3", "value4"},
				},
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value1", "value2"},
				},
			}},
		},
		{
			name: "nil selector",
			new: api.ExpressionRequirementConfig{
				Key:    "key1",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			selectors: nil,
			expected: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value1", "value2"},
				},
			}},
		},
	}

	for _, tc := range testcases {
		d := v1.Deployment{
			Spec: v1.DeploymentSpec{
				Selector: tc.selectors,
			},
		}
		obj := FakeDeployment{&d, LabelSelectorFuncs{&d}}

		obj.AddMatchExpressionsSelector(tc.new)

		if !reflect.DeepEqual(tc.expected, d.Spec.Selector) {
			t.Errorf("%s: wrong match expressions selectors.  Expected %+v, got %+v", tc.name, tc.expected, d.Spec.Selector)
		}
	}
}

func TestRemoveMatchExpressionsSelector(t *testing.T) {
	testcases := []struct {
		name      string
		selectors *metav1.LabelSelector
		remove    api.ExpressionRequirementConfig
		expected  *metav1.LabelSelector
	}{
		{
			name:      "no existing match expressions selectors",
			selectors: &metav1.LabelSelector{MatchLabels: map[string]string{}},
			remove: api.ExpressionRequirementConfig{
				Key:    "key1",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			expected: &metav1.LabelSelector{MatchLabels: map[string]string{}},
		},
		{
			name: "exiting match expressions selectors",
			selectors: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key2",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value3", "value4"},
				},
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value1", "value2"},
				},
			}},
			remove: api.ExpressionRequirementConfig{
				Key:    "key1",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			expected: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key2",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value3", "value4"},
				},
			}},
		},
		{
			name: "match expressions selectors to remove doesn't exist",
			selectors: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value1", "value2"},
				},
			}},
			remove: api.ExpressionRequirementConfig{
				Key:    "key2",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			expected: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      "key1",
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"value1", "value2"},
				},
			}},
		},
		{
			name:      "nil selector",
			selectors: nil,
			remove: api.ExpressionRequirementConfig{
				Key:    "key1",
				Op:     api.ExpressionRequirementOpIn,
				Values: []string{"value1", "value2"},
			},
			expected: nil,
		},
	}

	for _, tc := range testcases {
		d := v1.Deployment{
			Spec: v1.DeploymentSpec{
				Selector: tc.selectors,
			},
		}
		obj := FakeDeployment{&d, LabelSelectorFuncs{&d}}
		obj.RemoveMatchExpressionsSelector(tc.remove)

		if !reflect.DeepEqual(tc.expected, d.Spec.Selector) {
			t.Errorf("%s: wrong match expressions selectors.  Expected %+v, got %+v", tc.name, tc.expected, d.Spec.Selector)
		}
	}
}
