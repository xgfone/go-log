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
	"time"

	"github.com/xgfone/go-log/encoder"
)

// Encoder is used to encode the log record.
type Encoder interface {
	// Start starts to encode the log record into the buffer dst.
	Start(dst []byte, loggerName string, level int) []byte

	// Encode encodes the key-value with the stack depth into the buffer dst.
	Encode(dst []byte, key string, value interface{}) []byte

	// End ends to encode the log record with the message into the buffer dst.
	End(dst []byte, msg string) []byte
}

type encoderProxy struct {
	encoder.IntEncoder
	encoder.Int64Encoder
	encoder.UintEncoder
	encoder.Uint64Encoder
	encoder.Float64Encoder
	encoder.BoolEncoder
	encoder.StringEncoder
	encoder.TimeEncoder
	encoder.DurationEncoder
	Encoder
}

func newEncoder(orig Encoder) (enc encoderProxy) {
	enc.Encoder = orig

	var ok bool
	if enc.IntEncoder, ok = orig.(encoder.IntEncoder); !ok {
		enc.IntEncoder = intEncoder{orig}
	}
	if enc.Int64Encoder, ok = orig.(encoder.Int64Encoder); !ok {
		enc.Int64Encoder = int64Encoder{orig}
	}
	if enc.UintEncoder, ok = orig.(encoder.UintEncoder); !ok {
		enc.UintEncoder = uintEncoder{orig}
	}
	if enc.Uint64Encoder, ok = orig.(encoder.Uint64Encoder); !ok {
		enc.Uint64Encoder = uint64Encoder{orig}
	}
	if enc.Float64Encoder, ok = orig.(encoder.Float64Encoder); !ok {
		enc.Float64Encoder = float64Encoder{orig}
	}
	if enc.BoolEncoder, ok = orig.(encoder.BoolEncoder); !ok {
		enc.BoolEncoder = boolEncoder{orig}
	}
	if enc.StringEncoder, ok = orig.(encoder.StringEncoder); !ok {
		enc.StringEncoder = strEncoder{orig}
	}
	if enc.TimeEncoder, ok = orig.(encoder.TimeEncoder); !ok {
		enc.TimeEncoder = timeEncoder{orig}
	}
	if enc.DurationEncoder, ok = orig.(encoder.DurationEncoder); !ok {
		enc.DurationEncoder = durationEncoder{orig}
	}
	return
}

/// ----------------------------------------------------------------------- ///

type intEncoder struct{ Encoder }

func (e intEncoder) EncodeInt(dst []byte, key string, value int) []byte {
	return e.Encode(dst, key, value)
}

// Int is equal to e.Kv(key, value), but optimized for the value typed int.
func (e *Emitter) Int(key string, value int) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeInt(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type int64Encoder struct{ Encoder }

func (e int64Encoder) EncodeInt64(dst []byte, key string, value int64) []byte {
	return e.Encode(dst, key, value)
}

// Int64 is equal to e.Kv(key, value), but optimized for the value typed int64.
func (e *Emitter) Int64(key string, value int64) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeInt64(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type uintEncoder struct{ Encoder }

func (e uintEncoder) EncodeUint(dst []byte, key string, value uint) []byte {
	return e.Encode(dst, key, value)
}

// Uint is equal to e.Kv(key, value), but optimized for the value typed uint.
func (e *Emitter) Uint(key string, value uint) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeUint(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type uint64Encoder struct{ Encoder }

func (e uint64Encoder) EncodeUint64(dst []byte, key string, value uint64) []byte {
	return e.Encode(dst, key, value)
}

// Uint64 is equal to e.Kv(key, value), but optimized for the value typed uint64.
func (e *Emitter) Uint64(key string, value uint64) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeUint64(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type float64Encoder struct{ Encoder }

func (e float64Encoder) EncodeFloat64(dst []byte, key string, value float64) []byte {
	return e.Encode(dst, key, value)
}

// Float64 is equal to e.Kv(key, value), but optimized for the value typed float64.
func (e *Emitter) Float64(key string, value float64) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeFloat64(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type boolEncoder struct{ Encoder }

func (e boolEncoder) EncodeBool(dst []byte, key string, value bool) []byte {
	return e.Encode(dst, key, value)
}

// Bool is equal to e.Kv(key, value), but optimized for the value typed bool.
func (e *Emitter) Bool(key string, value bool) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeBool(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type strEncoder struct{ Encoder }

func (e strEncoder) EncodeString(dst []byte, key string, value string) []byte {
	return e.Encode(dst, key, value)
}

// Str is equal to e.Kv(key, value), but optimized for the value typed string.
func (e *Emitter) Str(key string, value string) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeString(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type timeEncoder struct{ Encoder }

func (e timeEncoder) EncodeTime(dst []byte, key string, value time.Time) []byte {
	return e.Encode(dst, key, value)
}

// Time is equal to e.Kv(key, value), but optimized for the value typed time.Time.
func (e *Emitter) Time(key string, value time.Time) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeTime(e.buffer, key, value)
	return e
}

/// ----------------------------------------------------------------------- ///

type durationEncoder struct{ Encoder }

func (e durationEncoder) EncodeDuration(dst []byte, key string, value time.Duration) []byte {
	return e.Encode(dst, key, value)
}

// Duration is equal to e.Kv(key, value), but optimized for the value typed time.Duration.
func (e *Emitter) Duration(key string, value time.Duration) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.EncodeDuration(e.buffer, key, value)
	return e
}
