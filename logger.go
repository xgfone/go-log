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
)

// LevelLogger is the logger interface with the level to emit the log.
type LevelLogger interface {
	Level(level int) Logger
	Logger
}

// Logger is the logger interface to emit the log.
type Logger interface {
	// Enabled reports whether the logger is enabled.
	Enabled() bool

	// Kv and Kvs append the key-value contexts and return the logger itself.
	Kv(key string, value interface{}) Logger
	Kvs(kvs ...interface{}) Logger

	// Print and Printf log the message and end the logger.
	Printf(msg string, args ...interface{})
	Print(args ...interface{})
}

type logger struct {
	writer  LevelWriter
	encoder Encoder
	buffer  []byte
	level   int
}

func (l *logger) Enabled() bool { return l != nil }

func (l *logger) Kv(key string, value interface{}) Logger {
	if l != nil {
		l.buffer = l.encoder.Encode(l.buffer, key, value)
	}
	return l
}

func (l *logger) Kvs(kvs ...interface{}) Logger {
	if l != nil {
		_len := len(kvs)
		if _len%2 != 0 {
			panic("the length of the key-value log contexts is not even")
		}
		for i := 0; i < _len; i += 2 {
			l.Kv(kvs[i].(string), kvs[i+1])
		}
	}
	return l
}

func (l *logger) Print(args ...interface{}) {
	if l == nil {
		return
	}

	l.emit(fmt.Sprint(args...))
}

func (l *logger) Printf(msg string, args ...interface{}) {
	if l == nil {
		return
	}

	if len(args) == 0 {
		l.emit(msg)
	} else {
		l.emit(fmt.Sprintf(msg, args...))
	}
}

func (l *logger) emit(msg string) {
	l.buffer = l.encoder.End(l.buffer, msg)
	l.writer.WriteLevel(l.level, l.buffer)
	l.buffer = l.buffer[:0]
	loggerPool.Put(l)

	if l.level == LvlFatal {
		atexit.Exit(1)
	} else if l.level >= LvlPanic {
		panic(msg)
	}
}

// DefaultBufferCap is the default capacity of the buffer to encode the log.
var DefaultBufferCap = 256

var loggerPool = sync.Pool{New: func() interface{} {
	return &logger{buffer: make([]byte, 0, DefaultBufferCap)}
}}

func newLogger(engine *Engine, level int, depth int) *logger {
	if engine.isDisabled(level) {
		return nil
	}

	l := loggerPool.Get().(*logger)
	l.level = level
	l.writer = engine.Output.writer
	l.encoder = engine.Output.encoder

	l.buffer = l.encoder.Start(l.buffer, engine.name, l.level)
	l.buffer = append(l.buffer, engine.ctx...)
	for i, _len := 0, len(engine.hooks); i < _len; i++ {
		engine.hooks[i].Run(l, engine.name, level, depth+2)
	}

	return l
}

var _ LevelLogger = &Engine{}

// Kv implements the interface Logger.
func (e *Engine) Kv(key string, value interface{}) Logger {
	return newLogger(e, e.level, e.depth).Kv(key, value)
}

// Kvs implements the interface Logger.
func (e *Engine) Kvs(kvs ...interface{}) Logger {
	return newLogger(e, e.level, e.depth).Kvs(kvs...)
}

// Print implements the interface Logger.
func (e *Engine) Print(args ...interface{}) {
	newLogger(e, e.level, e.depth).Print(args...)
}

// Printf implements the interface Logger.
func (e *Engine) Printf(msg string, args ...interface{}) {
	newLogger(e, e.level, e.depth).Printf(msg, args...)
}

func (e *Engine) getLogger(level, depth int) Logger {
	return newLogger(e, level, e.depth+depth)
}

// Logger returns a logger with the level and the additional stack depth
// to emit the log.
func (e *Engine) Logger(level, depth int) Logger {
	checkLevel(level)
	return newLogger(e, level, e.depth+depth)
}

// Level implements the interface LevelLogger to emit the log based on the level,
// which is equal to e.Logger(level, 0).
func (e *Engine) Level(level int) Logger {
	checkLevel(level)
	return newLogger(e, level, e.depth)
}

// Trace is equal to e.Level(LvlTrace).
func (e *Engine) Trace() Logger { return newLogger(e, LvlTrace, e.depth) }

// Debug is equal to e.Level(LvlDebug).
func (e *Engine) Debug() Logger { return newLogger(e, LvlDebug, e.depth) }

// Info is equal to e.Level(LvlInfo).
func (e *Engine) Info() Logger { return newLogger(e, LvlInfo, e.depth) }

// Warn is equal to e.Level(LvlWarn).
func (e *Engine) Warn() Logger { return newLogger(e, LvlWarn, e.depth) }

// Error is equal to e.Level(LvlError).
func (e *Engine) Error() Logger { return newLogger(e, LvlError, e.depth) }

// Alert is equal to e.Level(LvlAlert).
func (e *Engine) Alert() Logger { return newLogger(e, LvlAlert, e.depth) }

// Panic is equal to e.Panic(LvlPanic).
func (e *Engine) Panic() Logger { return newLogger(e, LvlPanic, e.depth) }

// Fatal is equal to e.Level(LvlFatal).
func (e *Engine) Fatal() Logger { return newLogger(e, LvlFatal, e.depth) }
