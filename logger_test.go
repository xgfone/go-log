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
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("test").WithCtx(F("ctxkey", "ctxvalue"))

	encoder := logger.Encoder.(*JSONEncoder)
	encoder.SetWriter(StreamWriter(buf))
	encoder.TimeKey = ""

	logger.Info(`"test json encoder"`,
		F("nil", nil),
		F("bool", true),
		F("int", 10),
		F("int8", 11),
		F("int16", 12),
		F("int32", 13),
		F("int64", 14),
		F("uint", 15),
		F("uint8", 16),
		F("uint16", 17),
		F("uint32", 18),
		F("uint64", 19),
		F("float32", 20),
		F("float64", 21),
		F("string", "22"),
		F("error", errors.New("23")),
		F("duration", time.Duration(24)*time.Second),
		F("time", time.Date(2021, time.May, 25, 22, 52, 26, 0, time.UTC)),
		F("[]interface{}", []interface{}{"26", "27"}),
		F("[]string", []string{"28", "29"}),
		F("[]uint", []uint{30, 31}),
		F("[]uint64", []uint64{32, 33}),
		F("[]int", []int{34, 35}),
		F("[]int64", []int64{36, 37}),
		F("map[string]interface{}", map[string]interface{}{"a": "38", "b": "39"}),
		F("map[string]string", map[string]string{"c": "40", "d": "41"}),
	)

	type encoderT struct {
		Msg      string      `json:"msg"`
		Lvl      string      `json:"lvl"`
		Logger   string      `json:"logger"`
		Stack    string      `json:"stack"`
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
		Lvl:      "INFO",
		Logger:   "test",
		Stack:    "[github.com/xgfone/go-log/logger_test.go:34]",
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

func TestLoggerSF(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := New("test")
	logger.Ctxs = []Field{F("ctxkey", "ctxvalue")}

	encoder := logger.Encoder.(*JSONEncoder)
	encoder.SetWriter(StreamWriter(buf))
	encoder.TimeKey = ""

	type encoderT struct {
		Msg    string `json:"msg"`
		Lvl    string `json:"lvl"`
		Logger string `json:"logger"`
		Ctx    string `json:"ctxkey"`

		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}

	var result encoderT
	expect := encoderT{
		Msg:    "test json encoder",
		Lvl:    "INFO",
		Logger: "test",
		Ctx:    "ctxvalue",
		Key1:   "value1",
		Key2:   "value2",
	}

	logger.Infos("test json encoder", "key1", "value1", "key2", "value2")
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("err=%s, json=%s", err, buf.String())
	} else if !reflect.DeepEqual(result, expect) {
		t.Errorf("expect '%+v', but got '%+v'", expect, result)
	}

	buf.Reset()
	result.Key1 = expect.Key1
	result.Key2 = expect.Key2
	logger.Infof("test %s %s", "json", "encoder")
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Errorf("err=%s, json=%s", err, buf.String())
	} else if !reflect.DeepEqual(result, expect) {
		t.Error(buf.String())
		t.Errorf("expect '%+v', but got '%+v'", expect, result)
	}
}
