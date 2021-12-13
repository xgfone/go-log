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

// FormatTime is used to format the time to the buffer dst.
var FormatTime = func(buf []byte, t time.Time) []byte {
	buf = append(buf, '"')
	buf = t.AppendFormat(buf, time.RFC3339Nano)
	return append(buf, '"')
}

// Now is used to get the current time.
var Now = time.Now

// Encoder is used to encode the log record.
type Encoder interface {
	// Start starts to encode the log record into the buffer dst.
	Start(dst []byte, loggerName string, level int) []byte

	// Encode encodes the key-value with the stack depth into the buffer dst.
	Encode(dst []byte, key string, value interface{}) []byte

	// End ends to encode the log record with the message into the buffer dst.
	End(dst []byte, msg string) []byte
}

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

	// LoggerKey is the key name of the logger name.
	//
	// Default: "logger"
	LoggerKey string

	// LevelKey is the key name of the level if not empty.
	//
	// Default: "lvl"
	LevelKey string

	// MsgKey is the key name of the message.
	//
	// Default: "msg"
	MsgKey string
}

// NewJSONEncoder returns a new JSONEncoder.
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
func (enc *JSONEncoder) Start(buf []byte, name string, level int) []byte {
	// JSON Start
	buf = append(buf, '{')

	// Time
	if len(enc.TimeKey) > 0 {
		buf = enc.appendString(buf, enc.TimeKey)
		buf = append(buf, ':')
		buf = FormatTime(buf, Now())
		buf = append(buf, ',')
	}

	// Level
	if len(enc.LevelKey) > 0 {
		buf = enc.appendString(buf, enc.LevelKey)
		buf = append(buf, ':')
		buf = enc.appendString(buf, FormatLevel(level))
		buf = append(buf, ',')
	}

	// Logger
	if len(enc.LoggerKey) > 0 && len(name) > 0 {
		buf = enc.appendString(buf, enc.LoggerKey)
		buf = append(buf, ':')
		buf = enc.appendString(buf, name)
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
	buf = enc.appendString(buf, key)
	buf = append(buf, ':')
	buf = enc.appendAny(buf, value)
	buf = append(buf, ',')
	return buf
}

// End implements the interface Encoder.
func (enc *JSONEncoder) End(buf []byte, msg string) []byte {
	// Msg
	buf = enc.appendString(buf, enc.MsgKey)
	buf = append(buf, ':')
	buf = enc.appendString(buf, msg)

	// JSON End
	buf = append(buf, '}')

	// Newline
	if enc.Newline {
		buf = append(buf, '\n')
	}

	return buf
}

func (enc *JSONEncoder) appendString(buf []byte, s string) []byte {
	buf = append(buf, '"')

	// Loop through each character in the string.
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			buf = appendStringComplex(buf, s, i)
			return append(buf, '"')
		}
	}

	buf = append(buf, s...)
	return append(buf, '"')
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
		buf = FormatTime(buf, v)

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
		buf = enc.appendString(buf, v)

	case interface{ EncodeJSON([]byte) []byte }:
		buf = v.EncodeJSON(buf)

	case interface{ WriteJSON(io.Writer) }:
		w := &bufferWriter{buf: buf}
		v.WriteJSON(w)
		buf = w.buf

	case error:
		buf = enc.appendString(buf, v.Error())

	case fmt.Stringer:
		buf = enc.appendString(buf, v.String())

	case json.Marshaler:
		if data, err := v.MarshalJSON(); err != nil {
			buf = enc.appendString(buf, fmt.Sprintf("JSONEncoderError: %s", err.Error()))
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
			buf = enc.appendString(buf, _v)
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

			buf = enc.appendString(buf, key)
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

			buf = enc.appendString(buf, key)
			buf = append(buf, ':')
			buf = enc.appendString(buf, value)
		}
		buf = append(buf, '}')

	default:
		if data, err := json.Marshal(v); err != nil {
			buf = enc.appendString(buf, fmt.Sprintf("JSONEncoderError: %s", err.Error()))
		} else {
			buf = append(buf, data...)
		}
	}

	return buf
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
