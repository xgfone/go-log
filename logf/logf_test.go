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
	"os"

	"github.com/xgfone/go-log"
	"github.com/xgfone/go-log/encoder"
)

func ExampleLogger() {
	enc := encoder.NewJSONEncoder()
	enc.TimeKey = "" // For test, we disable the log time
	_logger := log.New("").WithHooks(log.Caller("caller"))
	_logger.SetWriter(os.Stdout)
	_logger.SetEncoder(enc)

	logger := NewLogger(_logger, 0)
	logger.Tracef("%s msg", "trace")
	logger.Debugf("%s msg", "debug")
	logger.Infof("%s msg", "info")
	logger.Warnf("%s msg", "warn")
	logger.Errorf("%s msg", "error")
	logger.Alertf("%s msg", "alert")

	// Output:
	// {"lvl":"debug","caller":"logf_test.go:ExampleLogger:33","msg":"debug msg"}
	// {"lvl":"info","caller":"logf_test.go:ExampleLogger:34","msg":"info msg"}
	// {"lvl":"warn","caller":"logf_test.go:ExampleLogger:35","msg":"warn msg"}
	// {"lvl":"error","caller":"logf_test.go:ExampleLogger:36","msg":"error msg"}
	// {"lvl":"alert","caller":"logf_test.go:ExampleLogger:37","msg":"alert msg"}
}
