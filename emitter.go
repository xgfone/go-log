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
	"fmt"
	"sync"

	"github.com/xgfone/go-atexit"
	"github.com/xgfone/go-log/writer"
)

// DefaultBufferCap is the default capacity of the buffer to encode the log.
var DefaultBufferCap = 256

var emitterPool = sync.Pool{New: func() interface{} {
	return &Emitter{buffer: make([]byte, 0, DefaultBufferCap)}
}}

// Emitter is used to emit the log message.
type Emitter struct {
	writer  writer.LevelWriter
	encoder encoderProxy
	buffer  []byte
	level   int
}

// Enabled reports whether the log emitter is enabled.
func (e *Emitter) Enabled() bool { return e != nil }

// Err is equal to e.Kv("err", err).
func (e *Emitter) Err(err error) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.Encode(e.buffer, "err", err)
	return e
}

// Kv appends a key-value context into the log message and returns the emitter itself.
func (e *Emitter) Kv(key string, value interface{}) *Emitter {
	if e == nil {
		return nil
	}

	e.buffer = e.encoder.Encode(e.buffer, key, value)
	return e
}

// Kvs appends a set of the key-value contexts into the log message,
// and returns the emitter itself.
func (e *Emitter) Kvs(kvs ...interface{}) *Emitter {
	if e == nil {
		return nil
	}

	_len := len(kvs)
	if _len%2 != 0 {
		panic("the length of the key-value log contexts is not even")
	}

	for i := 0; i < _len; i += 2 {
		e.buffer = e.encoder.Encode(e.buffer, kvs[i].(string), kvs[i+1])
	}
	return e
}

// Print emits the log message to the underlying writer.
func (e *Emitter) Print(args ...interface{}) {
	if e == nil {
		return
	}

	e.emit(fmt.Sprint(args...))
}

// Printf emits the log message to the underlying writer.
func (e *Emitter) Printf(msg string, args ...interface{}) {
	if e == nil {
		return
	}

	if len(args) == 0 {
		e.emit(msg)
	} else {
		e.emit(fmt.Sprintf(msg, args...))
	}
}

func (e *Emitter) emit(msg string) {
	level := e.level
	e.buffer = e.encoder.End(e.buffer, msg)
	e.writer.WriteLevel(level, e.buffer)
	e.buffer = e.buffer[:0]
	emitterPool.Put(e)

	if level == LvlFatal {
		atexit.Exit(1)
	} else if level >= LvlPanic {
		panic(msg)
	}
}

func newEmitter(logger Logger, level int, depth int) *Emitter {
	if logger.isDisabled(level) {
		return nil
	}

	l := emitterPool.Get().(*Emitter)
	l.encoder = logger.Output.encoder
	l.writer = logger.Output.writer
	l.level = level

	l.buffer = l.encoder.Start(l.buffer, logger.name, logger.FormatLevel(level))
	l.buffer = append(l.buffer, logger.ctx...)
	for i, _len := 0, len(logger.hooks); i < _len; i++ {
		logger.hooks[i].Run(l, logger.name, level, depth+2)
	}

	return l
}

func (l Logger) getEmitter(level, depth int) *Emitter {
	return newEmitter(l, level, l.depth+depth)
}

// Log is convenient function to emit a log, which is equal to
// l.Level(level, 0).Kvs(keysAndValues...).Printf(msg).
func (l Logger) Log(level, depth int, msg string, keysAndValues ...interface{}) {
	l.Level(level, depth+1).Kvs(keysAndValues...).Printf(msg)
}

// Level returns an emitter with the level and the stack depth to emit the log.
func (l Logger) Level(level, depth int) *Emitter {
	checkLevel(level)
	return newEmitter(l, level, l.depth+depth)
}

// Trace is equal to l.Level(LvlTrace, 0).Kvs(kvs...).
func (l Logger) Trace(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlTrace, l.depth).Kvs(kvs...)
}

// Debug is equal to l.Level(LvlDebug, 0).Kvs(kvs...).
func (l Logger) Debug(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlDebug, l.depth).Kvs(kvs...)
}

// Info is equal to l.Level(LvlInfo, 0).Kvs(kvs...).
func (l Logger) Info(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlInfo, l.depth).Kvs(kvs...)
}

// Warn is equal to l.Level(LvlWarn, 0).Kvs(kvs...).
func (l Logger) Warn(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlWarn, l.depth).Kvs(kvs...)
}

// Error is equal to l.Level(LvlError, 0).Kvs(kvs...).
func (l Logger) Error(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlError, l.depth).Kvs(kvs...)
}

// Alert is equal to l.Level(LvlAlert, 0).Kvs(kvs...).
func (l Logger) Alert(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlAlert, l.depth).Kvs(kvs...)
}

// Panic is equal to l.Level(LvlPanic, 0).Kvs(kvs...).
func (l Logger) Panic(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlPanic, l.depth).Kvs(kvs...)
}

// Fatal is equal to l.Level(LvlFatal, 0).Kvs(kvs...).
func (l Logger) Fatal(kvs ...interface{}) *Emitter {
	return newEmitter(l, LvlFatal, l.depth).Kvs(kvs...)
}
