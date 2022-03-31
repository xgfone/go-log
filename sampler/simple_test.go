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

package sampler

import (
	"os"

	"github.com/xgfone/go-log"
	"github.com/xgfone/go-log/encoder"
)

func ExampleSimpleSampler() {
	enc := encoder.NewJSONEncoder()
	enc.TimeKey = "" // Disable the time for the test example

	sampler := NewSimpleSampler(log.LvlInfo)
	sampler.ResetNamedLevels(map[string]int{"root": log.LvlError})
	sampler.AddNamedLevel("root.child1.*", log.LvlWarn)

	logger := log.New("root").WithSampler(sampler)
	logger.SetWriter(os.Stdout)
	logger.SetEncoder(enc)

	logger.Debug().Print("msg11")
	logger.Info().Print("msg12")
	logger.Warn().Print("msg13")
	logger.Error().Print("msg14")

	clogger := logger.WithName("child1")
	clogger.Debug().Print("msg21")
	clogger.Info().Print("msg22")
	clogger.Warn().Print("msg23")
	clogger.Error().Print("msg24")

	cclogger := clogger.WithName("child2")
	cclogger.Debug().Print("msg31")
	cclogger.Info().Print("msg32")
	cclogger.Warn().Print("msg33")
	cclogger.Error().Print("msg34")

	// Output:
	// {"lvl":"error","logger":"root","msg":"msg14"}
	// {"lvl":"info","logger":"root.child1","msg":"msg22"}
	// {"lvl":"warn","logger":"root.child1","msg":"msg23"}
	// {"lvl":"error","logger":"root.child1","msg":"msg24"}
	// {"lvl":"warn","logger":"root.child1.child2","msg":"msg33"}
	// {"lvl":"error","logger":"root.child1.child2","msg":"msg34"}
}
