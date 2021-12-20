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

// Package encoder provides the encoder to encode the log message.
package encoder

import "time"

// IntEncoder is used to encode the key and the value typed int.
type IntEncoder interface {
	EncodeInt(dst []byte, key string, value int) []byte
}

// Int64Encoder is used to encode the key and the value typed int64.
type Int64Encoder interface {
	EncodeInt64(dst []byte, key string, value int64) []byte
}

// UintEncoder is used to encode the key and the value typed uint.
type UintEncoder interface {
	EncodeUint(dst []byte, key string, value uint) []byte
}

// Uint64Encoder is used to encode the key and the value typed uint64.
type Uint64Encoder interface {
	EncodeUint64(dst []byte, key string, value uint64) []byte
}

// Float64Encoder is used to encode the key and the value typed float64.
type Float64Encoder interface {
	EncodeFloat64(dst []byte, key string, value float64) []byte
}

// BoolEncoder is used to encode the value typed bool.
type BoolEncoder interface {
	EncodeBool(dst []byte, key string, value bool) []byte
}

// StringEncoder is used to encode the value typed string.
type StringEncoder interface {
	EncodeString(dst []byte, key string, value string) []byte
}

// TimeEncoder is used to encode the value typed time.Time.
type TimeEncoder interface {
	EncodeTime(dst []byte, key string, value time.Time) []byte
}

// DurationEncoder is used to encode the value typed time.Duration.
type DurationEncoder interface {
	EncodeDuration(dst []byte, key string, value time.Duration) []byte
}
