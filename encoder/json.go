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

package encoder

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode/utf8"
)

const hex = "0123456789abcdef"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

// Now is used to get the current time.
var Now = time.Now

// FormatTime is used to format the time to the buffer dst.
var FormatTime = func(buf []byte, t time.Time) []byte {
	buf = append(buf, '"')
	buf = t.AppendFormat(buf, time.RFC3339Nano)
	return append(buf, '"')
}

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
	// If true, append a newline when emit the log record.
	//
	// Default: true
	Newline bool

	// TimeKey is the key name of the time when to emit the log record if not empty.
	//
	// Default: "t"
	TimeKey string

	// TimeFormatFunc is used to format time.Time.
	//
	// Default: FormatTime
	TimeFormatFunc func(dst []byte, t time.Time) []byte

	// LevelKey is the key name of the level if not empty.
	//
	// Default: "lvl"
	LevelKey string

	// LevelFormatFunc is used to format the level.
	LevelFormatFunc func(level int) string

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
func NewJSONEncoder(formatLevel func(level int) string) *JSONEncoder {
	return &JSONEncoder{
		Newline:   true,
		TimeKey:   "t",
		LevelKey:  "lvl",
		LoggerKey: "logger",
		MsgKey:    "msg",

		LevelFormatFunc: formatLevel,
	}
}

// Start implements the interface Encoder.
func (enc *JSONEncoder) Start(buf []byte, name string, level int) []byte {
	// JSON Start
	buf = append(buf, '{')

	// Time
	if len(enc.TimeKey) > 0 {
		buf = AppendJSONString(buf, enc.TimeKey)
		buf = append(buf, ':')
		buf = enc.appendTime(buf, Now())
		buf = append(buf, ',')
	}

	// Level
	if len(enc.LevelKey) > 0 && enc.LevelFormatFunc != nil {
		buf = AppendJSONString(buf, enc.LevelKey)
		buf = append(buf, ':')
		buf = AppendJSONString(buf, enc.LevelFormatFunc(level))
		buf = append(buf, ',')
	}

	// Logger
	if len(enc.LoggerKey) > 0 && len(name) > 0 {
		buf = AppendJSONString(buf, enc.LoggerKey)
		buf = append(buf, ':')
		buf = AppendJSONString(buf, name)
		buf = append(buf, ',')
	}

	return buf
}

// Encode implements the interface Encoder, which supports not only the basic
// or builtin types, but also other interfaces as follow:
//
//   - error
//   - fmt.Stringer
//   - json.Marshaler
//   - interface{ WriteJSON(io.Writer) }
//   - interface{ EncodeJSON(dst []byte) []byte }
//
func (enc *JSONEncoder) Encode(buf []byte, key string, value interface{}) []byte {
	buf = AppendJSONString(buf, key)
	buf = append(buf, ':')
	buf = enc.appendAny(buf, value)
	buf = append(buf, ',')
	return buf
}

// End implements the interface Encoder.
func (enc *JSONEncoder) End(buf []byte, msg string) []byte {
	// Msg
	buf = AppendJSONString(buf, enc.MsgKey)
	buf = append(buf, ':')
	buf = AppendJSONString(buf, msg)

	// JSON End
	buf = append(buf, '}')

	// Newline
	if enc.Newline {
		buf = append(buf, '\n')
	}

	return buf
}

func (enc *JSONEncoder) appendTime(buf []byte, t time.Time) []byte {
	if enc.TimeFormatFunc == nil {
		buf = FormatTime(buf, t)
	} else {
		buf = enc.TimeFormatFunc(buf, t)
	}
	return buf
}

type bufferWriter struct{ buf []byte }

func (w *bufferWriter) Write(p []byte) (n int, err error) {
	if n = len(p); n > 0 {
		w.buf = append(w.buf, p...)
	}
	return
}

func (enc *JSONEncoder) appendAny(buf []byte, any interface{}) []byte {
	switch v := any.(type) {
	case time.Duration:
		buf = append(buf, '"')
		buf = append(buf, v.String()...)
		buf = append(buf, '"')

	case time.Time:
		buf = enc.appendTime(buf, v)

	case nil:
		buf = append(buf, `null`...)

	case bool:
		if v {
			buf = append(buf, `true`...)
		} else {
			buf = append(buf, `false`...)
		}

	case int:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int8:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int16:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int32:
		buf = strconv.AppendInt(buf, int64(v), 10)

	case int64:
		buf = strconv.AppendInt(buf, v, 10)

	case uint:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint8:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint16:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint32:
		buf = strconv.AppendUint(buf, uint64(v), 10)

	case uint64:
		buf = strconv.AppendUint(buf, v, 10)

	case float32:
		buf = strconv.AppendFloat(buf, float64(v), 'f', -1, 32)

	case float64:
		buf = strconv.AppendFloat(buf, v, 'f', -1, 64)

	case string:
		buf = AppendJSONString(buf, v)

	case interface{ EncodeJSON([]byte) []byte }:
		buf = v.EncodeJSON(buf)

	case interface{ WriteJSON(io.Writer) }:
		w := &bufferWriter{buf: buf}
		v.WriteJSON(w)
		buf = w.buf

	case error:
		buf = AppendJSONString(buf, v.Error())

	case fmt.Stringer:
		buf = AppendJSONString(buf, v.String())

	case json.Marshaler:
		if data, err := v.MarshalJSON(); err != nil {
			buf = AppendJSONString(buf, fmt.Sprintf("JSONEncoderError: %s", err.Error()))
		} else {
			buf = append(buf, data...)
		}

	case []interface{}:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = enc.appendAny(buf, _v)
		}
		buf = append(buf, ']')

	case []string:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = AppendJSONString(buf, _v)
		}
		buf = append(buf, ']')

	case []uint:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = strconv.AppendUint(buf, uint64(_v), 10)
		}
		buf = append(buf, ']')

	case []int:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = strconv.AppendInt(buf, int64(_v), 10)
		}
		buf = append(buf, ']')

	case []int64:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = strconv.AppendInt(buf, _v, 10)
		}
		buf = append(buf, ']')

	case []uint64:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = strconv.AppendUint(buf, _v, 10)
		}
		buf = append(buf, ']')

	case map[string]interface{}:
		buf = append(buf, '{')
		var i int
		for key, value := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			i++

			buf = AppendJSONString(buf, key)
			buf = append(buf, ':')
			buf = enc.appendAny(buf, value)
		}
		buf = append(buf, '}')

	case map[string]string:
		buf = append(buf, '{')
		var i int
		for key, value := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			i++

			buf = AppendJSONString(buf, key)
			buf = append(buf, ':')
			buf = AppendJSONString(buf, value)
		}
		buf = append(buf, '}')

	default:
		if data, err := json.Marshal(v); err != nil {
			buf = AppendJSONString(buf, fmt.Sprintf("JSONEncoderError: %s", err.Error()))
		} else {
			buf = append(buf, data...)
		}
	}

	return buf
}

// AppendJSONString appends the string s as JSON into dst, then returns dst.
func AppendJSONString(dst []byte, s string) []byte {
	dst = append(dst, '"')

	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			dst = appendStringComplex(dst, s, i)
			return append(dst, '"')
		}
	}

	dst = append(dst, s...)
	return append(dst, '"')
}

// appendStringComplex is used by appendString to take over an in progress JSON
// string encoding that encountered a character that needs to be encoded.
func appendStringComplex(dst []byte, s string, i int) []byte {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a remplacement character code
				// in place of the invalid sequence.
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}

		if noEscapeTable[b] {
			i++
			continue
		}

		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}

		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}

	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
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
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendInt(dst, int64(value), 10)
	dst = append(dst, ',')
	return dst
}

// EncodeInt64 implements the interface Int64Encoder.
func (enc *JSONEncoder) EncodeInt64(dst []byte, key string, value int64) []byte {
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendInt(dst, value, 10)
	dst = append(dst, ',')
	return dst
}

// EncodeUint implements the interface UintEncoder.
func (enc *JSONEncoder) EncodeUint(dst []byte, key string, value uint) []byte {
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendUint(dst, uint64(value), 10)
	dst = append(dst, ',')
	return dst
}

// EncodeUint64 implements the interface Uint64Encoder.
func (enc *JSONEncoder) EncodeUint64(dst []byte, key string, value uint64) []byte {
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendUint(dst, value, 10)
	dst = append(dst, ',')
	return dst
}

// EncodeFloat64 implements the interface Float64Encoder.
func (enc *JSONEncoder) EncodeFloat64(dst []byte, key string, value float64) []byte {
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = strconv.AppendFloat(dst, value, 'f', -1, 64)
	dst = append(dst, ',')
	return dst
}

// EncodeBool implements the interface BoolEncoder.
func (enc *JSONEncoder) EncodeBool(dst []byte, key string, value bool) []byte {
	dst = AppendJSONString(dst, key)
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
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = AppendJSONString(dst, value)
	dst = append(dst, ',')
	return dst
}

// EncodeTime implements the interface TimeEncoder.
func (enc *JSONEncoder) EncodeTime(dst []byte, key string, value time.Time) []byte {
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = enc.appendTime(dst, value)
	dst = append(dst, ',')
	return dst
}

// EncodeDuration implements the interface DurationEncoder.
func (enc *JSONEncoder) EncodeDuration(dst []byte, key string, value time.Duration) []byte {
	dst = AppendJSONString(dst, key)
	dst = append(dst, ':')
	dst = append(dst, '"')
	dst = append(dst, value.String()...)
	dst = append(dst, '"')
	dst = append(dst, ',')
	return dst
}

// EncodeStringSlice implements the interface StringSliceEncoder.
func (enc *JSONEncoder) EncodeStringSlice(dst []byte, key string, value []string) []byte {
	dst = append(dst, '[')
	for i, v := range value {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = AppendJSONString(dst, v)
	}
	dst = append(dst, ']')
	return dst
}
