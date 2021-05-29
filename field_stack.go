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

import (
	"fmt"

	"github.com/go-stack/stack"
)

// StackField is used to get the stack of the caller.
type StackField interface {
	Stack(depth int) interface{}
	Field
}

type stackField struct {
	key   string
	stack func(depth int) interface{}
}

func (f stackField) Key() string                 { return f.key }
func (f stackField) Value() interface{}          { panic("StackField.Value(): cannot be called") }
func (f stackField) Stack(depth int) interface{} { return f.stack(depth + 1) }

// StackFieldFunc returns a new StackField.
func StackFieldFunc(key string, getStack func(depth int) interface{}) StackField {
	return stackField{key: key, stack: getStack}
}

// Caller returns a StackField that returns the caller "file:line".
//
// If fullPath is true, the file is the full path but removing the GOPATH prefix.
func Caller(key string, fullPath ...bool) StackField {
	format := "%v"
	if len(fullPath) > 0 && fullPath[0] {
		format = "%+v"
	}

	return StackFieldFunc(key, func(depth int) interface{} {
		return fmt.Sprintf(format, stack.Caller(depth+1))
	})
}

// CallerStack returns a StackField returning the caller stack without runtime.
//
// If fullPath is true, the file is the full path but removing the GOPATH prefix.
func CallerStack(key string, fullPath ...bool) StackField {
	format := "%v"
	if len(fullPath) > 0 && fullPath[0] {
		format = "%+v"
	}

	return StackFieldFunc(key, func(depth int) interface{} {
		s := stack.Trace().TrimBelow(stack.Caller(depth + 1)).TrimRuntime()
		if len(s) > 0 {
			return fmt.Sprintf(format, s)
		}
		return ""
	})
}
