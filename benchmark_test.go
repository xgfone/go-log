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

	"github.com/xgfone/go-log/writer"
)

type nothingEncoder struct{}

func (enc nothingEncoder) Start(b []byte, n string, l int) []byte          { return b }
func (enc nothingEncoder) Encode(b []byte, k string, v interface{}) []byte { return b }
func (enc nothingEncoder) End(b []byte, m string) []byte                   { return b }

func newTestJSONEncoder() Encoder {
	enc := NewJSONEncoder()
	enc.TimeKey = ""
	return enc
}

func BenchmarkLevelDisabled(b *testing.B) {
	logger := New("").SetWriter(writer.Discard).SetEncoder(newTestJSONEncoder())
	logger.SetLevel(LvlError)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf("message")
		}
	})
}

func BenchmarkNothingEncoder(b *testing.B) {
	logger := New("").SetWriter(writer.Discard).SetEncoder(nothingEncoder{})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf("message")
		}
	})
}

func BenchmarkJSONEncoderWithoutCtxsAndKVs(b *testing.B) {
	logger := New("").SetWriter(writer.Discard).SetEncoder(newTestJSONEncoder())

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf("message")
		}
	})
}

func BenchmarkJSONEncoderWith8Contexts(b *testing.B) {
	logger := New("").SetWriter(writer.Discard).SetEncoder(newTestJSONEncoder())
	logger = logger.AppendCtx("k1", "v1").AppendCtx("k2", "v2").
		AppendCtx("k3", "v3").AppendCtx("k4", "v4").AppendCtx("k5", "v5").
		AppendCtx("k6", "v6").AppendCtx("k7", "v7").AppendCtx("k8", "v8")

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Printf("message")
		}
	})
}

func BenchmarkJSONEncoderWith8KeyValues(b *testing.B) {
	logger := New("").SetWriter(writer.Discard).SetEncoder(newTestJSONEncoder())

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Kv("k1", "v1").Kv("k2", "v2").Kv("k3", "v3").
				Kv("k4", "v4").Kv("k5", "v5").Kv("k6", "v6").Kv("k7", "v7").
				Kv("k8", "v8").Printf("message")
		}
	})
}
