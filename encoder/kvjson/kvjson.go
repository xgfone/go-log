// Copyright 2022 xgfone
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

// Package kvjson is used to encode the json key-value pair.
package kvjson

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// FormatTime is used to format the time to the buffer dst.
var FormatTime = func(buf []byte, t time.Time) []byte {
	buf = append(buf, '"')
	buf = t.AppendFormat(buf, time.RFC3339Nano)
	return append(buf, '"')
}

// JSON is used to encode the key-value json.
type JSON struct {
	// TimeFormatFunc is used to format time.Time.
	//
	// Default: FormatTime
	TimeFormatFunc func(dst []byte, t time.Time) []byte
}

// EncodeStart writes the json start character "{".
func (j JSON) EncodeStart(buf []byte) []byte {
	return append(buf, '{')
}

// EncodeEnd writes the json end character "}".
func (j JSON) EncodeEnd(buf []byte) []byte {
	return append(buf, '}')
}

// EncodeKV encodes the key-value pair, which supports not only the basic
// or builtin types, but also other interfaces as follow:
//
//   - error
//   - fmt.Stringer
//   - json.Marshaler
//   - interface{ WriteJSON(io.Writer) }
//   - interface{ EncodeJSON(dst []byte) []byte }
func (j JSON) EncodeKV(buf []byte, key string, value interface{}) []byte {
	buf = AppendJSONString(buf, key)
	buf = append(buf, ':')
	buf = j.appendAny(buf, value)
	buf = append(buf, ',')
	return buf
}

// Encode is the alias of EncodeKV.
func (j JSON) Encode(buf []byte, key string, value interface{}) []byte {
	return j.EncodeKV(buf, key, value)
}

// EncodeTime encodes the time.
func (j JSON) EncodeTime(buf []byte, t time.Time) []byte {
	return j.appendTime(buf, t)
}

// EncodeAny encodes the any value as JSON.
func (j JSON) EncodeAny(buf []byte, v interface{}) []byte {
	return j.appendAny(buf, v)
}

type bufferWriter struct{ buf []byte }

func (w *bufferWriter) Write(p []byte) (n int, err error) {
	if n = len(p); n > 0 {
		w.buf = append(w.buf, p...)
	}
	return
}

func (j JSON) appendTime(buf []byte, t time.Time) []byte {
	if j.TimeFormatFunc == nil {
		buf = FormatTime(buf, t)
	} else {
		buf = j.TimeFormatFunc(buf, t)
	}
	return buf
}

func (j JSON) appendAny(buf []byte, any interface{}) []byte {
	switch v := any.(type) {
	case time.Duration:
		buf = append(buf, '"')
		buf = append(buf, v.String()...)
		buf = append(buf, '"')

	case time.Time:
		buf = j.appendTime(buf, v)

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

	case json.RawMessage:
		if v == nil {
			buf = append(buf, "null"...)
		} else {
			buf = append(buf, v...)
		}

	case json.Marshaler:
		if data, err := v.MarshalJSON(); err != nil {
			buf = AppendJSONString(buf, fmt.Sprintf("JSONError: %s", err.Error()))
		} else {
			buf = append(buf, data...)
		}

	case []interface{}:
		buf = append(buf, '[')
		for i, _v := range v {
			if i > 0 {
				buf = append(buf, ',')
			}
			buf = j.appendAny(buf, _v)
		}
		buf = append(buf, ']')

	case []byte:
		if v == nil {
			buf = append(buf, "null"...)
		} else {
			buf = append(buf, v...)
		}

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
			buf = j.appendAny(buf, value)
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
			buf = AppendJSONString(buf, fmt.Sprintf("JSONError: %s", err.Error()))
		} else {
			buf = append(buf, data...)
		}
	}

	return buf
}
