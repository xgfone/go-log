// Copyright 2021~2022 xgfone
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

package encoder

import (
	"strconv"
	"time"

	"github.com/xgfone/go-log/encoder/kvjson"
)

// Now is used to get the current time.
var Now = time.Now

var (
	_ IntEncoder         = &JSONEncoder{}
	_ Int64Encoder       = &JSONEncoder{}
	_ UintEncoder        = &JSONEncoder{}
	_ Uint64Encoder      = &JSONEncoder{}
	_ Float64Encoder     = &JSONEncoder{}
	_ BoolEncoder        = &JSONEncoder{}
	_ StringEncoder      = &JSONEncoder{}
	_ TimeEncoder        = &JSONEncoder{}
	_ DurationEncoder    = &JSONEncoder{}
	_ StringSliceEncoder = &JSONEncoder{}
)

// JSONEncoder is a log encoder to encode the log record as JSON.
type JSONEncoder struct {
	kvjson.JSON

	// If true, append a newline when emit the log record.
	//
	// Default: true
	Newline bool

	// TimeKey is the key name of the time when to emit the log record if not empty.
	//
	// Default: "t"
	TimeKey string

	// LevelKey is the key name of the level if not empty.
	//
	// Default: "lvl"
	LevelKey string

	// LoggerKey is the key name of the logger name.
	//
	// Default: "logger"
	LoggerKey string

	// MsgKey is the key name of the message.
	//
	// Default: "msg"
	MsgKey string
}

// NewJSONEncoder returns a new JSONEncoder.
//
// If formatLevel is nil, disable to format the level.
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		Newline:   true,
		TimeKey:   "t",
		LevelKey:  "lvl",
		LoggerKey: "logger",
		MsgKey:    "msg",
	}
}

// Start implements the interface Encoder.
func (enc *JSONEncoder) Start(buf []byte, name, level string) []byte {
	// JSON Start
	buf = append(buf, '{')

	// Time
	if len(enc.TimeKey) > 0 {
		buf = kvjson.AppendJSONString(buf, enc.TimeKey)
		buf = append(buf, ':')
		buf = enc.JSON.EncodeTime(buf, Now())
		buf = append(buf, ',')
	}

	// Level
	if len(enc.LevelKey) > 0 {
		buf = kvjson.AppendJSONString(buf, enc.LevelKey)
		buf = append(buf, ':')
		buf = kvjson.AppendJSONString(buf, level)
		buf = append(buf, ',')
	}

	// Logger
	if len(enc.LoggerKey) > 0 && len(name) > 0 {
		buf = kvjson.AppendJSONString(buf, enc.LoggerKey)
		buf = append(buf, ':')
		buf = kvjson.AppendJSONString(buf, name)
		buf = append(buf, ',')
	}

	return buf
}

// Encode implements the interface Encoder by using kvjson.JSON.EncodeKV.
func (enc *JSONEncoder) Encode(buf []byte, key string, value interface{}) []byte {
	return enc.JSON.EncodeKV(buf, key, value)
}

// End implements the interface Encoder.
func (enc *JSONEncoder) End(buf []byte, msg string) []byte {
	// Msg
	buf = kvjson.AppendJSONString(buf, enc.MsgKey)
	buf = append(buf, ':')
	buf = kvjson.AppendJSONString(buf, msg)

	// JSON End
	buf = append(buf, '}')

	// Newline
	if enc.Newline {
		buf = append(buf, '\n')
	}

	return buf
}

/// ----------------------------------------------------------------------- ///

var (
	_ IntEncoder      = &JSONEncoder{}
	_ Int64Encoder    = &JSONEncoder{}
	_ UintEncoder     = &JSONEncoder{}
	_ Uint64Encoder   = &JSONEncoder{}
	_ Float64Encoder  = &JSONEncoder{}
	_ BoolEncoder     = &JSONEncoder{}
	_ StringEncoder   = &JSONEncoder{}
	_ TimeEncoder     = &JSONEncoder{}
	_ DurationEncoder = &JSONEncoder{}
)

// EncodeInt implements the interface IntEncoder.
func (enc *JSONEncoder) EncodeInt(dst []byte, key string, value int) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendInt(dst, int64(value), 10)
	dst = append(dst, ',')
	return dst
}

// EncodeInt64 implements the interface Int64Encoder.
func (enc *JSONEncoder) EncodeInt64(dst []byte, key string, value int64) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendInt(dst, value, 10)
	dst = append(dst, ',')
	return dst
}

// EncodeUint implements the interface UintEncoder.
func (enc *JSONEncoder) EncodeUint(dst []byte, key string, value uint) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendUint(dst, uint64(value), 10)
	dst = append(dst, ',')
	return dst
}

// EncodeUint64 implements the interface Uint64Encoder.
func (enc *JSONEncoder) EncodeUint64(dst []byte, key string, value uint64) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendUint(dst, value, 10)
	dst = append(dst, ',')
	return dst
}

// EncodeFloat64 implements the interface Float64Encoder.
func (enc *JSONEncoder) EncodeFloat64(dst []byte, key string, value float64) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendFloat(dst, value, 'f', -1, 64)
	dst = append(dst, ',')
	return dst
}

// EncodeBool implements the interface BoolEncoder.
func (enc *JSONEncoder) EncodeBool(dst []byte, key string, value bool) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	if value {
		dst = append(dst, `true`...)
	} else {
		dst = append(dst, `false`...)
	}
	dst = append(dst, ',')
	return dst
}

// EncodeString implements the interface StringEncoder.
func (enc *JSONEncoder) EncodeString(dst []byte, key string, value string) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = kvjson.AppendJSONString(dst, value)
	dst = append(dst, ',')
	return dst
}

// EncodeTime implements the interface TimeEncoder.
func (enc *JSONEncoder) EncodeTime(dst []byte, key string, value time.Time) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = enc.JSON.EncodeTime(dst, value)
	dst = append(dst, ',')
	return dst
}

// EncodeDuration implements the interface DurationEncoder.
func (enc *JSONEncoder) EncodeDuration(dst []byte, key string, value time.Duration) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = append(dst, '"')
	dst = append(dst, value.String()...)
	dst = append(dst, '"')
	dst = append(dst, ',')
	return dst
}

// EncodeStringSlice implements the interface StringSliceEncoder.
func (enc *JSONEncoder) EncodeStringSlice(dst []byte, key string, value []string) []byte {
	dst = kvjson.AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = append(dst, '[')
	for i, v := range value {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = kvjson.AppendJSONString(dst, v)
	}
	dst = append(dst, ']')
	dst = append(dst, ',')
	return dst
}
