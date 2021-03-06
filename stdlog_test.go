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
	"log"
	"strings"
	"testing"
)

func TestStdLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	DefaultLogger.SetWriter(buf)
	DefaultLogger.Output.SetEncoder(newTestEncoder())
	logger := New("").WithWriter(buf).WithEncoder(newTestEncoder()).
		WithHooks(Caller("caller"))

	stdlog1 := logger.StdLogger("", LvlDebug)
	stdlog1.Print("msg1")
	stdlog1.Println("msg2")

	log.SetFlags(0)
	log.SetOutput(logger.WithDepth(stdlogDepth))
	log.Printf("msg3")

	StdLogger("", LvlDebug).Printf("msg4")

	expects := []string{
		`{"lvl":"debug","caller":"stdlog_test.go:TestStdLogger:32","msg":"msg1"}`,
		`{"lvl":"debug","caller":"stdlog_test.go:TestStdLogger:33","msg":"msg2"}`,
		`{"lvl":"debug","caller":"stdlog_test.go:TestStdLogger:37","msg":"msg3"}`,
		`{"lvl":"debug","caller":"stdlog_test.go:TestStdLogger:39","msg":"msg4"}`,
		``,
	}
	testStrings(t, "stdlog", expects, strings.Split(buf.String(), "\n"))
}
