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

import "testing"

func BenchmarkNothingEncoder(b *testing.B) {
	logger := New("").WithEncoder(NothingEncoder())
	logger.Ctxs = nil

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("message")
		}
	})
}

func BenchmarkJSONEncoderWithoutCtxField(b *testing.B) {
	logger := New("").WithEncoder(NewJSONEncoder(DiscardWriter()))
	logger.Ctxs = nil

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("message")
		}
	})
}

func BenchmarkJSONEncoderWith10CtxFields(b *testing.B) {
	logger := New("").WithEncoder(NewJSONEncoder(DiscardWriter()))
	logger.Ctxs = nil
	logger = logger.WithCtx(F("k1", "v1"), F("k2", "v2"), F("k3", "v3"),
		F("k4", "v4"), F("k5", "v5"), F("k6", "v6"), F("k7", "v7"),
		F("k8", "v8"), F("k9", "v9"), F("k10", "v10"))

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("message")
		}
	})
}
