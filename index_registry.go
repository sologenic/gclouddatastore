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
)

var indexedFieldRegistry sync.Map // reflect.Type (struct) -> map[string]struct{}

// RegisterIndexedFields marks property name paths on typ as indexed when saving
// structs. typ must be a struct type or pointer-to-struct. Paths use the same
// dotted names as saved properties (including flatten prefixes, e.g. "I.X").
// Tags on a field still apply; use this to opt fields into indexing without
// datastore struct tags.
func RegisterIndexedFields(typ reflect.Type, paths ...string) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic("datastore: RegisterIndexedFields requires a struct type")
	}
	m := make(map[string]struct{})
	if v, ok := indexedFieldRegistry.Load(typ); ok {
		for k := range v.(map[string]struct{}) {
			m[k] = struct{}{}
		}
	}
	for _, p := range paths {
		m[p] = struct{}{}
	}
	indexedFieldRegistry.Store(typ, m)
}

func isRegisteredIndexed(typ reflect.Type, name string) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v, ok := indexedFieldRegistry.Load(typ)
	if !ok {
		return false
	}
	_, ok = v.(map[string]struct{})[name]
	return ok
}
