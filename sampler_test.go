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

import "os"

func ExampleSimpleSampler() {
	// For example test
	GlobalDisableSampling(false)
	encoder := NewJSONEncoder()
	encoder.TimeKey = ""

	sampler := NewSimpleSampler(LvlInfo)
	sampler.ResetNamedLevels(map[string]int{"root": LvlError})
	sampler.AddNamedLevel("root.child1.*", LvlWarn)

	logger := New("root").WithSampler(sampler)
	logger.SetWriter(os.Stdout)
	logger.SetEncoder(encoder)

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

func ExampleSwitchSampler() {
	// For example test
	GlobalDisableSampling(false)
	encoder := NewJSONEncoder()
	encoder.TimeKey = ""

	sampler1 := NewSimpleSampler(LvlInfo)
	sampler1.ResetNamedLevels(map[string]int{"root": LvlWarn})

	switchSampler := NewSwitchSampler(sampler1)
	logger := New("root").WithSampler(switchSampler)
	logger.SetWriter(os.Stdout)
	logger.SetEncoder(encoder)

	logger.Debug().Print("msg1")
	logger.Info().Print("msg2")
	logger.Warn().Print("msg3")
	logger.Error().Print("msg4")

	sampler2 := NewSimpleSampler(LvlInfo)
	sampler2.ResetNamedLevels(map[string]int{"root": LvlError})
	switchSampler.Set(sampler2)

	logger.Debug().Print("msg5")
	logger.Info().Print("msg6")
	logger.Warn().Print("msg7")
	logger.Error().Print("msg8")

	// Output:
	// {"lvl":"warn","logger":"root","msg":"msg3"}
	// {"lvl":"error","logger":"root","msg":"msg4"}
	// {"lvl":"error","logger":"root","msg":"msg8"}
}
