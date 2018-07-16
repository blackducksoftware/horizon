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
	"bytes"
	"strings"
	"testing"

	"github.com/blackducksoftware/horizon/pkg/api"
)

func TestAddData(t *testing.T) {
	s := NewSecret(api.SecretConfig{
		Name:      "name",
		Namespace: "ns",
		Type:      api.SecretTypeOpaque,
	})
	data := map[string][]byte{"key1": []byte("data")}
	s.AddData(data)
	for k, v := range data {
		if objv, ok := s.obj.Data[k]; !ok {
			t.Errorf("%s key missing", k)
		} else if bytes.Compare(objv, v) != 0 {
			t.Errorf("expected %v got %v", v, objv)
		}
	}
	delete(data, "key1")
	if !(len(s.obj.Data) > 0) {
		t.Errorf("expected %d got %d", len(data), len(s.obj.Data))
	}
}

func TestAddStringData(t *testing.T) {
	s := NewSecret(api.SecretConfig{
		Name:      "name",
		Namespace: "ns",
		Type:      api.SecretTypeOpaque,
	})
	data := map[string]string{"key1": "data"}
	s.AddStringData(data)
	for k, v := range data {
		if objv, ok := s.obj.StringData[k]; !ok {
			t.Errorf("%s key missing", k)
		} else if strings.Compare(objv, v) != 0 {
			t.Errorf("expected %v got %v", v, objv)
		}
	}
	delete(data, "key1")
	if !(len(s.obj.StringData) > 0) {
		t.Errorf("expected %d got %d", len(data), len(s.obj.Data))
	}
}
