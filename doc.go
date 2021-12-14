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
//       log.Log(log.LvlError, 1).Kvs(kvs...).Kv("err", err).Printf(msg)
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
//       log.SetWriter(writer).SetLevel(log.ParseLevel(loglevel))
//       defer writer.Close()
//
//       // Emit the log.
//       log.Print("msg1")
//       log.Printf("msg%d", 2)
//       log.Kv("key1", "value1").Print("msg3")
//       log.Debug().Kv("key2", "value2").Print("msg4") // no log output.
//       log.Info().Kv("key3", "value3").Print("msg5")
//       log.Log(log.LvlInfo, 0).Kv("key4", "value4").Printf("msg6")
//       logError(nil, "msg7", "key5", "value5", "key6", 666, "key7", "value7")
//       logError(errors.New("error"), "msg8", "key8", 888, "key9", "value9")
//
//       // For Clild Logger
//       child1Logger := log.WithName("child1")
//       child2Logger := child1Logger.New("child2")
//       child1Logger.Kv("ckey1", "cvalue1").Print("msg9")
//       child2Logger.Printf("msg10")
//
//       // $ go run main.go
//       // {"t":"2021-12-12T11:41:11.2844234+08:00","lvl":"info","caller":"main.go:32:main","msg":"msg1"}
//       // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:33:main","msg":"msg2"}
//       // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:34:main","key1":"value1","msg":"msg3"}
//       // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:36:main","key3":"value3","msg":"msg5"}
//       // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"info","caller":"main.go:37:main","key4":"value4","msg":"msg6"}
//       // {"t":"2021-12-12T11:41:11.2918549+08:00","lvl":"error","caller":"main.go:39:main","key8":888,"key9":"value9","err":"error","msg":"msg8"}
//       // {"t":"2021-12-12T12:22:15.2466635+08:00","lvl":"info","logger":"child1","caller":"main.go:44:main","ckey1":"cvalue1","msg":"msg9"}
//       // {"t":"2021-12-12T12:22:15.2466635+08:00","lvl":"info","logger":"child1.child2","caller":"main.go:45:main","msg":"msg10"}
//   }
package log
