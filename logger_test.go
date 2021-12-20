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
	"fmt"
	"strings"
	"testing"

	"github.com/xgfone/go-atexit"
	jencoder "github.com/xgfone/go-log/encoder"
)

func newTestEncoder() Encoder {
	enc := jencoder.NewJSONEncoder(FormatLevel)
	enc.TimeKey = ""
	return enc
}

func testStrings(t *testing.T, prefix string, expects, results []string) {
	if len(expects) != len(results) {
		t.Errorf("%s: expect %d lines, but got %d",
			prefix, len(expects), len(results))
		return
	}

	for i, line := range expects {
		if results[i] != line {
			t.Errorf("%s: %d line: expect '%s', but got '%s'",
				prefix, i, line, results[i])
		}
	}
}

func TestLoggerEnabledLevel(t *testing.T) {
	logger := New("").WithLevel(LvlWarn)

	if logger.Enabled(LvlInfo) || !logger.Enabled(LvlError) {
		t.Error("fail")
	}

	SetGlobalLevel(LvlDebug)
	if logger.Enabled(LvlTrace) || !logger.Enabled(LvlInfo) {
		t.Error("fail")
	}

	SetGlobalLevel(LvlError)
	if logger.Enabled(LvlWarn) {
		t.Error("fail")
	}

	SetGlobalLevel(-1)
	if logger.Enabled(LvlDebug) || !logger.Enabled(LvlWarn) {
		t.Error("fail")
	}
}

func TestChildLoggerName(t *testing.T) {
	parent := New("parent")
	child1 := parent.WithName("child1")
	child2 := parent.WithName("child2")
	child3 := child2.WithName("child3")

	if name := child1.Name(); name != "parent.child1" {
		t.Errorf("expect the logger name '%s', but got '%s'", "parent.child1", name)
	}
	if name := child2.Name(); name != "parent.child2" {
		t.Errorf("expect the logger name '%s', but got '%s'", "parent.child2", name)
	}
	if name := child3.Name(); name != "parent.child2.child3" {
		t.Errorf("expect the logger name '%s', but got '%s'", "parent.child2.child3", name)
	}
}

func TestLevelPanicAndFatal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("").WithWriter(buf).WithEncoder(newTestEncoder())
	atexit.ExitFunc = func(code int) { fmt.Fprintf(buf, "exit with %d\n", code) }

	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintf(buf, "panic: %v\n", err)
			}
		}()
		logger.Panic().Printf("msg1")
	}()
	logger.Fatal().Printf("msg2")

	expects := []string{
		`{"lvl":"panic","msg":"msg1"}`,
		`panic: msg1`,
		`{"lvl":"fatal","msg":"msg2"}`,
		`exit with 1`,
		``,
	}
	testStrings(t, "panic_fatal", expects, strings.Split(buf.String(), "\n"))
}

func TestLevelDisable(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("").WithWriter(buf).WithEncoder(newTestEncoder()).WithLevel(LvlDisable)

	logger.Trace().Print("msg1")
	logger.Debug().Print("msg2")
	logger.Info().Print("msg3")
	logger.Warn().Print("msg4")
	logger.Error().Print("msg5")
	logger.Alert().Print("msg6")
	logger.Panic().Print("msg7")
	logger.Fatal().Print("msg8")
	if s := buf.String(); s != "" {
		t.Errorf("unexpected logs '%s'", s)
	}

	logger = logger.WithLevel(LvlTrace)
	global := GetGlobalLevel()
	defer SetGlobalLevel(global)
	SetGlobalLevel(LvlDisable)

	buf.Reset()
	logger.Trace().Print("msg1")
	logger.Debug().Print("msg2")
	logger.Info().Print("msg3")
	logger.Warn().Print("msg4")
	logger.Error().Print("msg5")
	logger.Alert().Print("msg6")
	logger.Panic().Print("msg7")
	logger.Fatal().Print("msg8")
	if s := buf.String(); s != "" {
		t.Errorf("unexpected logs '%s'", s)
	}
}

func TestSampler(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	sample := func(name string, lvl int) bool { return lvl > LvlWarn }
	logger := New("test").WithWriter(buf).
		WithEncoder(newTestEncoder()).
		WithSampler(SamplerFunc(sample))

	logger.Info().Print("msg1")
	logger.Error().Print("msg2")

	expect := `{"lvl":"error","logger":"test","msg":"msg2"}` + "\n"
	if result := buf.String(); result != expect {
		t.Errorf("expect '%s', but got '%s'", expect, result)
	}

	buf.Reset()
	GlobalDisableSampling(true)
	logger.Info().Print("msg1")
	logger.Error().Print("msg2")

	expects := []string{
		`{"lvl":"info","logger":"test","msg":"msg1"}`,
		`{"lvl":"error","logger":"test","msg":"msg2"}`,
		``,
	}
	testStrings(t, "sampler", expects, strings.Split(buf.String(), "\n"))
}
