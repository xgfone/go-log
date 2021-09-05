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
	"reflect"
	"strings"
	"testing"
)

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

func TestLoggerInheritance(t *testing.T) {
	parent := New("parent")
	parent.Ctxs = nil
	child := parent.WithName("child")

	buf := bytes.NewBuffer(nil)
	encoder := parent.Encoder.(*JSONEncoder)
	encoder.SetWriter(StreamWriter(buf))
	encoder.TimeKey = ""

	parent.Info("parent info 1")
	child.Info("child info 1")

	parent.SetLevel(LvlWarn)
	parent.Info("parent info 2")
	child.Info("child info 2")
	if lvl := child.GetLevel(); lvl != LvlWarn {
		t.Errorf("child logger expect level '%s', but got '%s'\n", LvlWarn, lvl)
	}

	child.SetLevel(LvlInfo)
	parent.Info("parent info 3")
	child.Info("child info 3")
	if lvl := child.GetLevel(); lvl != LvlInfo {
		t.Errorf("child logger expect level '%s', but got '%s'\n", LvlInfo, lvl)
	}

	child.UnsetLevel()
	child.Info("child info 4")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expect %d lines, but got %d: %v\n", 3, len(lines), lines)
	} else {
		expects := []string{
			`{"lvl":"INFO","logger":"parent","msg":"parent info 1"}`,
			`{"lvl":"INFO","logger":"child","msg":"child info 1"}`,
			`{"lvl":"INFO","logger":"child","msg":"child info 3"}`,
		}
		for i := 0; i < 3; i++ {
			if expects[i] != lines[i] {
				t.Errorf("%d: expect '%s', but got '%s'\n", i, expects[i], lines[i])
			}
		}
	}
}
