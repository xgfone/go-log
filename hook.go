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
	"path/filepath"
	"runtime"
)

// Hook is used to add the dynamic value into the log record.
type Hook interface {
	Run(e *Emitter, loggerName string, level int, depth int)
}

// HookFunc is a function hook.
type HookFunc func(e *Emitter, name string, level int, depth int)

// Run implements the interface Hook.
func (f HookFunc) Run(e *Emitter, name string, level int, depth int) {
	f(e, name, level, depth+1)
}

// Hooks returns the hooks.
func (l Logger) Hooks() []Hook { return l.hooks }

// WithHooks returns a new logger with the hooks.
func (l Logger) WithHooks(hooks ...Hook) Logger {
	l = l.Clone()
	l.hooks = append([]Hook{}, hooks...)
	return l
}

// CallerFormatFunc is used to format the file, name and line of the caller.
var CallerFormatFunc = func(file, name string, line int) string {
	name = filepath.Ext(name)
	if len(name) > 0 && name[0] == '.' {
		name = name[1:]
	}
	return fmt.Sprintf("%s:%d:%s", filepath.Base(file), line, name)
}

// Caller returns a callback function that returns the caller "file:line".
func Caller(key string) Hook {
	return HookFunc(func(e *Emitter, name string, level, depth int) {
		if pc, file, line, ok := runtime.Caller(depth + 1); ok {
			f := runtime.FuncForPC(pc)
			e.Kv(key, CallerFormatFunc(file, f.Name(), line))
		}
	})
}
