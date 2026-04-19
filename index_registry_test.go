// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datastore

import (
	"reflect"
	"sync"
	"testing"
)

func TestRegisterIndexedFields(t *testing.T) {
	t.Cleanup(func() {
		indexedFieldRegistry = sync.Map{}
	})

	type row struct {
		A int
		B string
	}

	RegisterIndexedFields(reflect.TypeOf(row{}), "A")

	props, err := SaveStruct(&row{A: 1, B: "x"})
	if err != nil {
		t.Fatal(err)
	}
	var gotA, gotB bool
	for _, p := range props {
		switch p.Name {
		case "A":
			gotA = p.Index
		case "B":
			gotB = p.Index
		}
	}
	if !gotA {
		t.Error("registered field A should be indexed")
	}
	if gotB {
		t.Error("field B should not be indexed without registration")
	}
}
