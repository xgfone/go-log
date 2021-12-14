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
	"strings"
	"testing"
)

func TestLoggerStackDepth(t *testing.T) {
	buf := bytes.NewBufferString("")
	enc := NewJSONEncoder()
	enc.TimeKey = ""

	logger := New("test")
	logger.SetWriter(buf)
	logger.SetEncoder(enc)
	logger.AddHooks(Caller("caller"))

	logger.Info().Print("msg0")
	logger.Level(LvlInfo).Print("msg1")
	logger.Level(LvlInfo).Printf("msg2")
	logger.Level(LvlInfo).Kv("key1", "value1").Print("msg3")
	logger.Level(LvlInfo).Kv("key2", "value2").Printf("msg4")
	logger.Print("msg5")
	logger.Printf("msg6")
	logger.Kv("key3", "value3").Print("msg7")
	logger.Kv("key4", "value4").Printf("msg8")

	expects := []string{
		`{"lvl":"info","logger":"test","caller":"hook_test.go:33","msg":"msg0"}`,
		`{"lvl":"info","logger":"test","caller":"hook_test.go:34","msg":"msg1"}`,
		`{"lvl":"info","logger":"test","caller":"hook_test.go:35","msg":"msg2"}`,
		`{"lvl":"info","logger":"test","caller":"hook_test.go:36","key1":"value1","msg":"msg3"}`,
		`{"lvl":"info","logger":"test","caller":"hook_test.go:37","key2":"value2","msg":"msg4"}`,
		`{"lvl":"debug","logger":"test","caller":"hook_test.go:38","msg":"msg5"}`,
		`{"lvl":"debug","logger":"test","caller":"hook_test.go:39","msg":"msg6"}`,
		`{"lvl":"debug","logger":"test","caller":"hook_test.go:40","key3":"value3","msg":"msg7"}`,
		`{"lvl":"debug","logger":"test","caller":"hook_test.go:41","key4":"value4","msg":"msg8"}`,
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
