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
)

func ExampleLogger() {
	// For test, we disable the log time
	log.DefaultLogger.Output.GetEncoder().(*log.JSONEncoder).TimeKey = ""
	log.DefaultLogger.SetWriter(os.Stdout)

	logger := NewLogger(nil, 0)
	logger.Tracef("%s msg", "trace")
	logger.Debugf("%s msg", "debug")
	logger.Infof("%s msg", "info")
	logger.Warnf("%s msg", "warn")
	logger.Errorf("%s msg", "error")
	logger.Alertf("%s msg", "alert")

	// Output:
	// {"lvl":"debug","caller":"logf_test.go:30","msg":"debug msg"}
	// {"lvl":"info","caller":"logf_test.go:31","msg":"info msg"}
	// {"lvl":"warn","caller":"logf_test.go:32","msg":"warn msg"}
	// {"lvl":"error","caller":"logf_test.go:33","msg":"error msg"}
	// {"lvl":"alert","caller":"logf_test.go:34","msg":"alert msg"}
}
