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
	logger := New("test").
		WithWriter(buf).
		WithEncoder(newTestEncoder()).
		WithHooks(Caller("caller"))

	logger.Info().Print("msg1")
	logger.Level(LvlInfo, 0).Printf("msg2")
	logger.Level(LvlInfo, 0).Kv("k1", "v1").Print("msg3")
	logger.Log(LvlInfo, 0, "msg4", "k2", "v2")
	logger.WithLevel(LvlInfo).Write([]byte("msg5"))

	const prefix = `{"lvl":"info","logger":"test","caller":"logger_stack_test.go:`
	expects := []string{
		prefix + `TestLoggerStack:31","msg":"msg1"}`,
		prefix + `TestLoggerStack:32","msg":"msg2"}`,
		prefix + `TestLoggerStack:33","k1":"v1","msg":"msg3"}`,
		prefix + `TestLoggerStack:34","k2":"v2","msg":"msg4"}`,
		prefix + `TestLoggerStack:35","msg":"msg5"}`,
		``,
	}
	testStrings(t, "logger_stack", expects, strings.Split(buf.String(), "\n"))
}

func TestGlobalStack(t *testing.T) {
	buf := bytes.NewBufferString("")
	DefaultLogger.SetWriter(buf)
	DefaultLogger.Output.SetEncoder(newTestEncoder())

	Info().Printf("msg1")
	LevelLog(LvlInfo, 0).Print("msg2")
	IfErr(errors.New("error"), "msg3")
	Ef(errors.New("error"), "msg4")

	expects := []string{
		`{"lvl":"info","caller":"logger_stack_test.go:TestGlobalStack:54","msg":"msg1"}`,
		`{"lvl":"info","caller":"logger_stack_test.go:TestGlobalStack:55","msg":"msg2"}`,
		`{"lvl":"error","caller":"logger_stack_test.go:TestGlobalStack:56","err":"error","msg":"msg3"}`,
		`{"lvl":"error","caller":"logger_stack_test.go:TestGlobalStack:57","err":"error","msg":"msg4"}`,
		``,
	}
	testStrings(t, "global_stack", expects, strings.Split(buf.String(), "\n"))
}
