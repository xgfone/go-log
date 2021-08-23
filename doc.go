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

// Package log provides a simple, flexible, extensible, powerful and
// structured logging tool based on the level, which has done the better balance
// between the flexibility and the performance.
//
// Features
//
//   - The better performance.
//   - Lazy evaluation of expensive operations.
//   - Support the level inherited from the parent logger.
//   - Simple, Flexible, Extensible, Powerful and Structured.
//   - Avoid to allocate the memory on heap as far as possible.
//   - Child loggers which inherit and add their own private context.
//   - Built-in support for logging to files, syslog, etc. See `Writer`.
//
// Example
//
//   package main
//
//   import "github.com/xgfone/go-log"
//
//   func main() {
//       log.DefalutLogger.Level = log.LvlWarn
//
//       // Emit the log with the fields.
//       log.Info("log msg", log.F("key1", "value1"), log.F("key2", "value2"))
//       log.Error("log msg", log.F("key1", "value1"), log.F("key2", "value2"))
//
//       // Emit the log with key-values
//       log.Infos("log msg", "key1", "value1", "key2", "value2")
//       log.Errors("log msg", "key1", "value1", "key2", "value2")
//
//       // Emit the log with the formatter.
//       log.Infof("log %s", "msg")
//       log.Errorf("log %s", "msg")
//
//       // Output:
//       // {"t":"2021-05-28T22:07:07.394835+08:00","lvl":"ERROR","stack":"[main.go:10]","key1":"value1","key2":"value2","msg":"log msg"}
//       // {"t":"2021-05-28T22:07:07.395066+08:00","lvl":"ERROR","stack":"[main.go:14]","key1":"value1","key2":"value2","msg":"log msg"}
//       // {"t":"2021-05-28T22:07:07.3951+08:00","lvl":"ERROR","stack":"[main.go:18]","msg":"log msg"}
//   }
package log
