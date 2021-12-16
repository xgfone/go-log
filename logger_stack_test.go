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
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestLoggerStack(t *testing.T) {
	buf := bytes.NewBufferString("")
	logger := New("test").WithWriter(buf).
		WithEncoder(newTestEncoder()).
		WithHooks(Caller("caller"))

	logger.Info().Print("msg1")
	logger.Level(LvlInfo, 0).Printf("msg2")
	logger.Level(LvlInfo, 0).Kv("k1", "v1").Print("msg3")
	logger.Level(LvlInfo, 0).Kvs("k2", "v2").Printf("msg4")
	logger.WithLevel(LvlInfo).StdLog("").Printf("msg5")

	const prefix = `{"lvl":"info","logger":"test","caller":"logger_stack_test.go:`
	expects := []string{
		prefix + `30:TestLoggerStack","msg":"msg1"}`,
		prefix + `31:TestLoggerStack","msg":"msg2"}`,
		prefix + `32:TestLoggerStack","k1":"v1","msg":"msg3"}`,
		prefix + `33:TestLoggerStack","k2":"v2","msg":"msg4"}`,
		prefix + `34:TestLoggerStack","msg":"msg5"}`,
		``,
	}
	testStrings(t, "logger_stack", expects, strings.Split(buf.String(), "\n"))
}

func TestGlobalStack(t *testing.T) {
	buf := bytes.NewBufferString("")
	DefaultLogger.SetWriter(buf)
	DefaultLogger.Output.encoder.(*JSONEncoder).TimeKey = ""

	Info().Printf("msg1")
	Level(LvlInfo, 0).Print("msg2")
	IfErr(errors.New("error"), "msg3")
	Ef(errors.New("error"), "msg4")
	StdLog("").Printf("msg5")

	expects := []string{
		`{"lvl":"info","caller":"logger_stack_test.go:53:TestGlobalStack","msg":"msg1"}`,
		`{"lvl":"info","caller":"logger_stack_test.go:54:TestGlobalStack","msg":"msg2"}`,
		`{"lvl":"error","caller":"logger_stack_test.go:55:TestGlobalStack","err":"error","msg":"msg3"}`,
		`{"lvl":"error","caller":"logger_stack_test.go:56:TestGlobalStack","err":"error","msg":"msg4"}`,
		`{"lvl":"debug","caller":"logger_stack_test.go:57:TestGlobalStack","msg":"msg5"}`,
		``,
	}
	testStrings(t, "global_stack", expects, strings.Split(buf.String(), "\n"))
}

func TestStdLog(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("").WithWriter(buf).
		WithEncoder(newTestEncoder()).
		WithHooks(Caller("caller"))

	stdlog1 := logger.StdLog("")
	stdlog1.Print("msg1")
	stdlog1.Println("msg2")

	stdlog2 := logger.StdLog("stdlog: ")
	stdlog2.Print("msg3")

	expects := []string{
		`{"lvl":"debug","caller":"logger_stack_test.go:77:TestStdLog","msg":"msg1"}`,
		`{"lvl":"debug","caller":"logger_stack_test.go:78:TestStdLog","msg":"msg2"}`,
		`{"lvl":"debug","caller":"logger_stack_test.go:81:TestStdLog","msg":"stdlog: msg3"}`,
		``,
	}
	testStrings(t, "stdlog", expects, strings.Split(buf.String(), "\n"))
}
