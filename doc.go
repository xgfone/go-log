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
// structured logger based on the level, which has done the better balance
// between the flexibility and the performance. It collects the log message
// with the key-value contexts, encodes them into the buffer, then writes
// the encoded log from the buffer into the underlying writer.
//
// Features
//
//   - Support `Go1.7+`.
//   - Compatible with the stdlib `log.Printf`.
//   - The better performance:
//       - Lazy evaluation of expensive operations.
//       - Avoid to allocate the memory on heap as far as possible.
//       - Encode in real time or pre-encode the key-value contexts into the buffer cache.
//   - Simple, Flexible, Extensible, Powerful and Structured.
//   - Support to customize the log encoder and writer.
//   - Provide the simple and easy-used api interface.
//
// Example
//
//   package main
//
//   import (
//       "errors"
//       "flag"
//
//       "github.com/xgfone/go-log"
//   )
//
//   var logfile string
//   var loglevel string
//
//   func logError(err error, msg string, kvs ...interface{}) {
//       if err == nil {
//           return
//       }
//       log.Level(log.LvlError, 1).Kvs(kvs...).Kv("err", err).Printf(msg)
//   }
//
//   func main() {
//       // Parse the CLI options.
//       flag.StringVar(&logfile, "logfile", "", "The log file path, default to stderr.")
//       flag.StringVar(&loglevel, "loglevel", "info", "The log level, such as debug, info, etc.")
//       flag.Parse()
//
//       // Configure the logger.
//       writer := log.FileWriter(logfile, "100M", 100)
//       defer writer.Close()
//       log.SetWriter(writer)
//       log.SetLevel(log.ParseLevel(loglevel))
//
//       // Emit the log.
//       log.Debug().Kv("key1", "value1").Print("msg1") // no log output.
//       log.Info().Kv("key2", "value2").Print("msg2")
//       log.Level(log.LvlInfo, 0).Kv("key3", "value3").Printf("msg3")
//       logError(nil, "msg4", "key4", "value4", "key5", 555, "key6", "value6")
//       logError(errors.New("error"), "msg7", "key8", 888, "key9", "value9")
//
//       // For Clild Logger
//       child1Logger := log.WithName("child1")
//       child2Logger := child1Logger.WithName("child2")
//       child1Logger.Info().Kv("ckey1", "cvalue1").Print("msg8")
//       child2Logger.Info().Kv("ckey2", "cvalue2").Printf("msg9")
//
//       // $ go run main.go
//       // {"t":"2021-12-17T00:04:44.8609884+08:00","lvl":"info","caller":"main.go:34:main","key2":"value2","msg":"msg2"}
//       // {"t":"2021-12-17T00:04:44.8660577+08:00","lvl":"info","caller":"main.go:35:main","key3":"value3","msg":"msg3"}
//       // {"t":"2021-12-17T00:04:44.8671207+08:00","lvl":"error","caller":"main.go:37:main","key8":888,"key9":"value9","err":"error","msg":"msg7"}
//       // {"t":"2021-12-17T00:04:44.8671207+08:00","lvl":"info","logger":"child1","caller":"main.go:42:main","ckey1":"cvalue1","msg":"msg8"}
//       // {"t":"2021-12-17T00:04:44.8678731+08:00","lvl":"info","logger":"child1.child2","caller":"main.go:43:main","ckey2":"cvalue2","msg":"msg9"}
//   }
package log
