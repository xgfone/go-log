// Copyright 2021 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

// Field represents a key-value pair.
type Field interface {
	Key() string
	Value() interface{}
}

type field struct {
	key   string
	value interface{}
}

func (f field) Key() string { return f.key }
func (f field) Value() interface{} {
	switch v := f.value.(type) {
	case func() interface{}:
		return v()
	case func() string:
		return v()
	default:
		return v
	}
}

// E is equal to F("err", err).
func E(err error) Field { return field{key: "err", value: err} }

// F returns a new Field. If value is "func() interface{}" or "func() string",
// it will be evaluated when the log is emitted, that's, it is lazy.
func F(key string, value interface{}) Field { return field{key: key, value: value} }

// FieldFunc returns a higher-order function to create the field with the key.
func FieldFunc(key string) func(value interface{}) Field {
	return func(value interface{}) Field {
		return F(key, value)
	}
}
