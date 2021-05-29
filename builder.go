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
	"strings"
	"time"
)

type builder struct {
	buf []byte
}

func newBuilder(n int) *builder { return &builder{buf: make([]byte, 0, n)} }

func (b *builder) Reset()                            { b.buf = b.buf[:0] }
func (b *builder) Write(p []byte) (int, error)       { b.buf = append(b.buf, p...); return len(p), nil }
func (b *builder) WriteString(s string) (int, error) { b.buf = append(b.buf, s...); return len(s), nil }

func (b *builder) AppendByte(p byte)     { b.buf = append(b.buf, p) }
func (b *builder) AppendString(s string) { b.buf = append(b.buf, s...) }
func (b *builder) AppendInt(i int64)     { b.buf = strconv.AppendInt(b.buf, i, 10) }
func (b *builder) AppendUint(i uint64)   { b.buf = strconv.AppendUint(b.buf, i, 10) }
func (b *builder) AppendFloat(f float64, bitSize int) {
	b.buf = strconv.AppendFloat(b.buf, f, 'f', -1, bitSize)
}

func (b *builder) AppendKeyAsJSON(key string) {
	b.buf = append(b.buf, '"')
	b.buf = append(b.buf, key...)
	b.buf = append(b.buf, '"')
}

func (b *builder) AppendTimeAsJSON(t time.Time, layout string) {
	b.buf = append(b.buf, '"')
	b.buf = t.AppendFormat(b.buf, layout)
	b.buf = append(b.buf, '"')
}

func (b *builder) AppendStringAsJSON(s string) {
	if s == "" {
		b.buf = append(b.buf, '"', '"')
	} else if strings.IndexByte(s, '"') > -1 {
		b.buf = strconv.AppendQuote(b.buf, s)
	} else {
		b.buf = append(b.buf, '"')
		b.buf = append(b.buf, s...)
		b.buf = append(b.buf, '"')
	}
}

func (b *builder) AppendAnyAsJSON(value interface{}, timeFmt string) {
	switch v := value.(type) {
	case interface{ WriteJSON(io.Writer) }:
		v.WriteJSON(b)
	case time.Duration:
		b.buf = append(b.buf, '"')
		b.buf = append(b.buf, v.String()...)
		b.buf = append(b.buf, '"')
	case time.Time:
		b.AppendTimeAsJSON(v, timeFmt)
	case fmt.Stringer:
		b.AppendStringAsJSON(v.String())
	case json.Marshaler:
		if data, err := v.MarshalJSON(); err != nil {
			b.AppendStringAsJSON(fmt.Sprintf("JSONEncoderError: %s", err.Error()))
		} else {
			b.buf = append(b.buf, data...)
		}
	case nil:
		b.AppendString(`null`)
	case bool:
		if v {
			b.AppendString(`true`)
		} else {
			b.AppendString(`false`)
		}
	case int:
		b.AppendInt(int64(v))
	case int8:
		b.AppendInt(int64(v))
	case int16:
		b.AppendInt(int64(v))
	case int32:
		b.AppendInt(int64(v))
	case int64:
		b.AppendInt(v)
	case uint:
		b.AppendUint(uint64(v))
	case uint8:
		b.AppendUint(uint64(v))
	case uint16:
		b.AppendUint(uint64(v))
	case uint32:
		b.AppendUint(uint64(v))
	case uint64:
		b.AppendUint(v)
	case float32:
		b.AppendFloat(float64(v), 32)
	case float64:
		b.AppendFloat(v, 64)
	case error:
		b.AppendStringAsJSON(v.Error())
	case string:
		b.AppendStringAsJSON(v)
	case []interface{}:
		b.AppendByte('[')
		for i, _v := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			b.AppendAnyAsJSON(_v, timeFmt)
		}
		b.AppendByte(']')
	case []string:
		b.AppendByte('[')
		for i, _v := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			b.AppendStringAsJSON(_v)
		}
		b.AppendByte(']')
	case []uint:
		b.AppendByte('[')
		for i, _v := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			b.AppendUint(uint64(_v))
		}
		b.AppendByte(']')
	case []int:
		b.AppendByte('[')
		for i, _v := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			b.AppendInt(int64(_v))
		}
		b.AppendByte(']')
	case []int64:
		b.AppendByte('[')
		for i, _v := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			b.AppendInt(_v)
		}
		b.AppendByte(']')
	case []uint64:
		b.AppendByte('[')
		for i, _v := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			b.AppendUint(_v)
		}
		b.AppendByte(']')
	case map[string]interface{}:
		b.AppendByte('{')
		var i int
		for key, value := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			i++

			b.AppendStringAsJSON(key)
			b.AppendByte(':')
			b.AppendAnyAsJSON(value, timeFmt)
		}
		b.AppendByte('}')
	case map[string]string:
		b.AppendByte('{')
		var i int
		for key, value := range v {
			if i > 0 {
				b.AppendByte(',')
			}
			i++

			b.AppendStringAsJSON(key)
			b.AppendByte(':')
			b.AppendStringAsJSON(value)
		}
		b.AppendByte('}')
	default:
		if data, err := json.Marshal(v); err != nil {
			b.AppendStringAsJSON(fmt.Sprintf("JSONEncoderError: %s", err.Error()))
		} else {
			b.buf = append(b.buf, data...)
		}
	}
}

func (b *builder) WriteLevel(w Writer, l Level) (n int, err error) {
	if n = len(b.buf); n > 0 {
		n, err = w.WriteLevel(l, b.buf)
	}
	return
}
