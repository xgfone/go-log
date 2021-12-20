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
	"testing"
	"time"

	"github.com/xgfone/go-log/writer"
)

var (
	bMessage  = "message"
	bCtxKey   = "key"
	bCtxValue = "value"
)

func newBenchLogger() Logger {
	return New("").WithWriter(writer.Discard).WithEncoder(newTestEncoder())
}

func BenchmarkJSONEncoderDisabled(b *testing.B) {
	logger := newBenchLogger().WithLevel(LvlDisable)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf(bMessage)
		}
	})
}

func BenchmarkJSONEncoderEmpty(b *testing.B) {
	logger := newBenchLogger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf("")
		}
	})
}

func BenchmarkJSONEncoderInfo(b *testing.B) {
	logger := newBenchLogger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf(bMessage)
		}
	})
}

func BenchmarkJSONEncoderWith8Contexts(b *testing.B) {
	logger := newBenchLogger()
	for i := 0; i < 8; i++ {
		logger = logger.WithContext(bCtxKey, bCtxValue)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf(bMessage)
		}
	})
}

func BenchmarkJSONEncoderWith8KeyValues(b *testing.B) {
	logger := newBenchLogger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Kv(bCtxKey, bCtxValue).
				Printf(bMessage)
		}
	})
}

func BenchmarkJSONEncoderWithOptimized8KVs(b *testing.B) {
	logger := newBenchLogger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().
				Int(bCtxKey, 111).
				Int64(bCtxKey, 222).
				Uint(bCtxKey, 333).
				Uint64(bCtxKey, 444).
				Float64(bCtxKey, 555).
				Bool(bCtxKey, true).
				Str(bCtxKey, bCtxValue).
				Duration(bCtxKey, time.Second).
				Printf(bMessage)
		}
	})
}
