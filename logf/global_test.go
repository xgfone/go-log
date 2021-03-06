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

package logf

import (
	"bytes"
	"strings"
	"testing"

	"github.com/xgfone/go-log"
	"github.com/xgfone/go-log/encoder"
)

func TestGlobal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log.DefaultLogger.SetWriter(buf)
	enc := encoder.NewJSONEncoder()
	enc.TimeKey = ""
	log.DefaultLogger.Output.SetEncoder(enc)

	Tracef("msg%d", 1)
	Debugf("msg%d", 2)
	Infof("msg%d", 3)
	Warnf("msg%d", 4)
	Errorf("msg%d", 5)
	Alertf("msg%d", 6)

	expects := []string{
		`{"lvl":"debug","caller":"global_test.go:TestGlobal:34","msg":"msg2"}`,
		`{"lvl":"info","caller":"global_test.go:TestGlobal:35","msg":"msg3"}`,
		`{"lvl":"warn","caller":"global_test.go:TestGlobal:36","msg":"msg4"}`,
		`{"lvl":"error","caller":"global_test.go:TestGlobal:37","msg":"msg5"}`,
		`{"lvl":"alert","caller":"global_test.go:TestGlobal:38","msg":"msg6"}`,
		``,
	}
	lines := strings.Split(buf.String(), "\n")
	if len(expects) != len(lines) {
		t.Errorf("expect %d log lines, but got %d:", len(expects), len(lines))
	} else {
		for i, line := range expects {
			if lines[i] != line {
				t.Errorf("%d: expect log line '%s', but got '%s'", i, line, lines[i])
			}
		}
	}
}
