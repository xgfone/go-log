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
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGlobal(t *testing.T) {
	buf := bytes.NewBufferString("")
	DefaultLogger.SetWriter(buf)
	DefaultLogger.Output.encoder.(*JSONEncoder).TimeKey = ""

	Info().Printf("msg1")
	Infof("msg2")
	Printf("msg3")
	Print("msg4")
	Ef(errors.New("error"), "msg5")
	IfErr(errors.New("error"), "msg6", "k", "v")
	StdLog("stdlog: ").Print("msg7")

	expects := []string{
		`{"lvl":"info","caller":"logger_test.go:32","msg":"msg1"}`,
		`{"lvl":"info","caller":"logger_test.go:33","msg":"msg2"}`,
		`{"lvl":"debug","caller":"logger_test.go:34","msg":"msg3"}`,
		`{"lvl":"debug","caller":"logger_test.go:35","msg":"msg4"}`,
		`{"lvl":"error","caller":"logger_test.go:36","err":"error","msg":"msg5"}`,
		`{"lvl":"error","caller":"logger_test.go:37","k":"v","err":"error","msg":"msg6"}`,
		`{"lvl":"debug","caller":"logger_test.go:38","msg":"stdlog: msg7"}`,
		"",
	}
	if lines := strings.Split(buf.String(), "\n"); len(lines) != len(expects) {
		t.Errorf("expect %d lines, but got %d", len(expects), len(lines))
	} else {
		for i, line := range lines {
			if line != expects[i] {
				t.Errorf("%d: expect line '%s', but '%s'", i, expects[i], line)
			}
		}
	}
}

func TestStdLog(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("test").SetWriter(buf).AddHooks(Caller("caller"))
	logger.Output.encoder.(*JSONEncoder).TimeKey = ""

	stdlog1 := logger.StdLog("")
	stdlog1.Print("msg1")
	stdlog1.Println("msg2")

	stdlog2 := logger.StdLog("stdlog: ")
	stdlog2.Print("msg3")

	expects := []string{
		`{"lvl":"debug","logger":"test","caller":"logger_test.go:67","msg":"msg1"}`,
		`{"lvl":"debug","logger":"test","caller":"logger_test.go:68","msg":"msg2"}`,
		`{"lvl":"debug","logger":"test","caller":"logger_test.go:71","msg":"stdlog: msg3"}`,
		``,
	}
	if lines := strings.Split(buf.String(), "\n"); len(lines) != len(expects) {
		t.Errorf("expect %d lines, but got %d", len(expects), len(lines))
	} else {
		for i, line := range expects {
			if lines[i] != line {
				t.Errorf("%d: expect line '%s', but got '%s'", i, line, lines[i])
			}
		}
	}
}

func TestLoggerEnabledLevel(t *testing.T) {
	logger := New("")
	logger.SetLevel(LvlInfo)

	if !logger.Enabled() {
		t.Error("fail")
	}

	SetGlobalLevel(LvlDebug)
	if !logger.Enabled() {
		t.Error("fail")
	}

	SetGlobalLevel(LvlWarn)
	if logger.Enabled() {
		t.Error("fail")
	}

	SetGlobalLevel(-1)
	if logger.Enable(LvlDebug) || !logger.Enable(LvlWarn) {
		t.Error("fail")
	}
}

func TestChildLoggerName(t *testing.T) {
	parent := New("parent")
	child1 := parent.New("child1")
	child2 := parent.New("child2")
	child3 := child2.New("child3")

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

func TestLevelDisable(t *testing.T) {
	enc := NewJSONEncoder()
	enc.TimeKey = ""
	buf := bytes.NewBuffer(nil)
	logger := New("").SetWriter(buf).SetEncoder(enc)

	logger.SetLevel(LvlDisable)
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

	logger.SetLevel(LvlTrace)
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

func TestLoggerSampler(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	sample := func(name string, lvl int) bool { return lvl > LvlWarn }
	logger := New("test").SetSampler(SamplerFunc(sample)).SetWriter(buf)
	logger.Output.encoder.(*JSONEncoder).TimeKey = ""
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
	expect = `{"lvl":"info","logger":"test","msg":"msg1"}` + "\n" +
		`{"lvl":"error","logger":"test","msg":"msg2"}` + "\n"
	if result := buf.String(); result != expect {
		t.Errorf("expect '%s', but got '%s'", expect, result)
	}
}

func TestLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("test").AppendCtx("ctxkey", "ctxvalue")
	logger.Output.GetEncoder().(*JSONEncoder).TimeKey = ""
	logger.SetWriter(buf)

	logger.Info().
		Kv("nil", nil).
		Kv("bool", true).
		Kv("int", 10).
		Kv("int8", 11).
		Kv("int16", 12).
		Kv("int32", 13).
		Kv("int64", 14).
		Kv("uint", 15).
		Kv("uint8", 16).
		Kv("uint16", 17).
		Kv("uint32", 18).
		Kv("uint64", 19).
		Kv("float32", 20).
		Kv("float64", 21).
		Kv("string", "22").
		Kv("error", errors.New("23")).
		Kv("duration", time.Duration(24)*time.Second).
		Kv("time", time.Date(2021, time.May, 25, 22, 52, 26, 0, time.UTC)).
		Kv("[]interface{}", []interface{}{"26", "27"}).
		Kv("[]string", []string{"28", "29"}).
		Kv("[]uint", []uint{30, 31}).
		Kv("[]uint64", []uint64{32, 33}).
		Kv("[]int", []int{34, 35}).
		Kv("[]int64", []int64{36, 37}).
		Kv("map[string]interface{}", map[string]interface{}{"a": "38", "b": "39"}).
		Kv("map[string]string", map[string]string{"c": "40", "d": "41"}).
		Print(`"test json encoder"`)

	type encoderT struct {
		Msg      string      `json:"msg"`
		Lvl      string      `json:"lvl"`
		Logger   string      `json:"logger"`
		Ctx      string      `json:"ctxkey"`
		Nil      interface{} `json:"nil"`
		Bool     bool        `json:"bool"`
		Int      int         `json:"int"`
		Int8     int8        `json:"int8"`
		Int16    int16       `json:"int16"`
		Int32    int32       `json:"int32"`
		Int64    int64       `json:"int64"`
		Uint     uint        `json:"uint"`
		Uint8    uint8       `json:"uint8"`
		Uint16   uint16      `json:"uint16"`
		Uint32   uint32      `json:"uint32"`
		Uint64   uint64      `json:"uint64"`
		Float32  float32     `json:"float32"`
		Float64  float64     `json:"float64"`
		String   string      `json:"string"`
		Error    string      `json:"error"`
		Duration string      `json:"duration"`
		Time     string      `json:"time"`

		Interfaces []interface{} `json:"[]interface{}"`
		Strings    []string      `json:"[]string"`
		Uints      []uint        `json:"[]uint"`
		Uint64s    []uint64      `json:"[]uint64"`
		Ints       []int         `json:"[]int"`
		Int64s     []int64       `json:"[]int64"`

		MapString    map[string]string      `json:"map[string]string"`
		MapInterface map[string]interface{} `json:"map[string]interface{}"`
	}

	expect := encoderT{
		Msg:      `"test json encoder"`,
		Lvl:      "info",
		Logger:   "test",
		Ctx:      "ctxvalue",
		Nil:      nil,
		Bool:     true,
		Int:      10,
		Int8:     11,
		Int16:    12,
		Int32:    13,
		Int64:    14,
		Uint:     15,
		Uint8:    16,
		Uint16:   17,
		Uint32:   18,
		Uint64:   19,
		Float32:  20,
		Float64:  21,
		String:   "22",
		Error:    "23",
		Duration: "24s",
		Time:     "2021-05-25T22:52:26Z",

		Interfaces: []interface{}{"26", "27"},
		Strings:    []string{"28", "29"},
		Uints:      []uint{30, 31},
		Uint64s:    []uint64{32, 33},
		Ints:       []int{34, 35},
		Int64s:     []int64{36, 37},

		MapInterface: map[string]interface{}{"a": "38", "b": "39"},
		MapString:    map[string]string{"c": "40", "d": "41"},
	}

	var result encoderT
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(result, expect) {
		t.Errorf("expect '%+v', but got '%+v'", expect, result)
	}
}
