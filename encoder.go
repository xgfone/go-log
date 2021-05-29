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
	"sync"
	"time"
)

var bufpool = sync.Pool{
	New: func() interface{} { return newBuilder(256) },
}

// Record represents a log record.
type Record struct {
	Name  string    // The logger name, which may be empty
	Time  time.Time // The start time when to emit the log
	Depth int       // The stack depth of the caller

	Lvl    Level   // The log level
	Msg    string  // The log message
	Ctxs   []Field // The SHARED key-value contexts. DON'T MODIFY IT!
	Fields []Field // The key-value pairs
}

// Encoder is used to encode the log record and to write it into the writer.
type Encoder interface {
	// Writer returns the writer.
	Writer() Writer

	// SetWriter resets the writer.
	SetWriter(Writer)

	// Encode encodes the log record and writes it into the writer.
	Encode(Record)
}

type nothingEncoder struct{ w Writer }

func (e *nothingEncoder) Writer() Writer     { return e.w }
func (e *nothingEncoder) SetWriter(w Writer) { e.w = w }
func (e *nothingEncoder) Encode(Record)      {}

// NothingEncoder encodes nothing.
func NothingEncoder() Encoder { return &nothingEncoder{} }

// JSONEncoder is an encoder to encode the key-values log as json.
type JSONEncoder struct {
	// If true, append a newline when emit the log record.
	//
	// Default: false
	Newline bool

	// MsgKey is the key name of the message.
	//
	// Default: "msg"
	MsgKey string

	// TimeKey is the key name of the time when to emit the log record if not empty.
	TimeKey string

	// TimeFmt is the layout of time.Time.
	TimeFmt string

	// LevelKey is the key name of the level if not empty.
	LevelKey string

	// LoggerKey is the key name of the logger name.
	LoggerKey string

	w Writer
}

// NewJSONEncoder returns a new JSONEncoder with the writer,
// which supports not only the built-in types but also the interfaces as follow:
//   fmt.Stringer
//   json.Marshaler
//   interface { WriteJSON(io.Writer) }
//
// Default:
//   Newline: true,
//   TimeKey: "t"
//   TimeFmt: time.RFC3339Nano
//   LevelKey: "lvl"
//   LoggerKey: "logger"
//
func NewJSONEncoder(w Writer) *JSONEncoder {
	return &JSONEncoder{
		LoggerKey: "logger",
		LevelKey:  "lvl",
		TimeKey:   "t",
		TimeFmt:   time.RFC3339Nano,
		Newline:   true,
		w:         w,
	}
}

// Writer implements the interface Encoder.
func (e *JSONEncoder) Writer() Writer { return e.w }

// SetWriter implements the interface Encoder.
func (e *JSONEncoder) SetWriter(w Writer) { e.w = w }

// Encode implements the interface Encoder.
func (e *JSONEncoder) Encode(r Record) {
	r.Depth++
	buf := bufpool.Get().(*builder)

	// JSON Start
	buf.AppendByte('{')

	// Time
	if e.TimeKey != "" {
		buf.AppendKeyAsJSON(e.TimeKey)
		buf.AppendByte(':')
		if r.Time.IsZero() {
			buf.AppendTimeAsJSON(time.Now(), e.TimeFmt)
		} else {
			buf.AppendTimeAsJSON(r.Time, e.TimeFmt)
		}
		buf.AppendByte(',')
	}

	// Level
	if e.LevelKey != "" {
		buf.AppendKeyAsJSON(e.LevelKey)
		buf.AppendByte(':')
		buf.AppendKeyAsJSON(r.Lvl.String())
		buf.AppendByte(',')
	}

	// Logger
	if r.Name != "" && e.LoggerKey != "" {
		buf.AppendKeyAsJSON(e.LoggerKey)
		buf.AppendByte(':')
		buf.AppendKeyAsJSON(r.Name)
		buf.AppendByte(',')
	}

	// Ctxs and Fields
	jsonEncodeFields(buf, r.Ctxs, r.Depth, e.TimeFmt)
	jsonEncodeFields(buf, r.Fields, r.Depth, e.TimeFmt)

	// Msg
	if e.MsgKey == "" {
		buf.AppendString(`"msg"`)
	} else {
		buf.AppendKeyAsJSON(e.MsgKey)
	}
	buf.AppendByte(':')
	buf.AppendStringAsJSON(r.Msg)

	// JSON End
	buf.AppendByte('}')

	if e.Newline {
		buf.AppendByte('\n')
	}

	buf.WriteLevel(e.w, r.Lvl)
	buf.Reset()
	bufpool.Put(buf)
}

func jsonEncodeFields(buf *builder, fields []Field, depth int, timeFmt string) {
	depth++
	for _, field := range fields {
		buf.AppendKeyAsJSON(field.Key())
		buf.AppendByte(':')
		if s, ok := field.(StackField); ok {
			buf.AppendAnyAsJSON(s.Stack(depth), timeFmt)
		} else {
			buf.AppendAnyAsJSON(field.Value(), timeFmt)
		}
		buf.AppendByte(',')
	}
}
