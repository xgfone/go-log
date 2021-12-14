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
	"io"
	"log"
)

// DefaultLogger is the default global logger.
var DefaultLogger = New("").AddHooks(Caller("caller"))

// StdLog is equal to DefaultLogger.StdLog(prefix).
func StdLog(prefix string) *log.Logger {
	return log.New(DefaultLogger.Clone().SetDepth(2), prefix, 0)
}

// GetLevel is equal to DefaultLogger.GetLevel().
func GetLevel() int { return DefaultLogger.GetLevel() }

// SetDepth is equal to DefaultLogger.SetDepth(depth).
func SetDepth(depth int) *Engine { return DefaultLogger.SetDepth(depth) }

// SetLevel is equal to DefaultLogger.SetLevel(level).
func SetLevel(level int) *Engine { return DefaultLogger.SetLevel(level) }

// SetWriter is eqaul to DefaultLogger.SetWriter(w).
func SetWriter(w io.Writer) *Engine { return DefaultLogger.SetWriter(w) }

// SetEncoder is eqaul to DefaultLogger.SetEncoder(enc).
func SetEncoder(enc Encoder) *Engine { return DefaultLogger.SetEncoder(enc) }

// AddHooks is equal to DefaultLogger.AddHooks(hooks...).
func AddHooks(hooks ...Hook) *Engine { return DefaultLogger.AddHooks(hooks...) }

// ResetHooks is equal to DefaultLogger.ResetHooks(hooks...).
func ResetHooks(hooks ...Hook) *Engine { return DefaultLogger.ResetHooks(hooks...) }

// ResetCtxs is equal to DefaultLogger.ResetCtxs().
func ResetCtxs() *Engine { return DefaultLogger.ResetCtxs() }

// AppendCtx is equal to DefaultLogger.AppendCtx(key, value).
func AppendCtx(key string, value interface{}) *Engine {
	return DefaultLogger.AppendCtx(key, value)
}

// Clone is equal to DefaultLogger.Clone().
func Clone() *Engine { return DefaultLogger.Clone() }

// WithName is equal to DefaultLogger.New(name).
func WithName(name string) *Engine { return DefaultLogger.New(name) }

// Log is equal to DefaultLogger.Logger(level, depth).
func Log(level, depth int) Logger { return DefaultLogger.Logger(level, depth+1) }

// Trace is equal to DefaultLogger.Trace().
func Trace() Logger { return DefaultLogger.getLogger(LvlTrace, 1) }

// Debug is equal to DefaultLogger.Debug().
func Debug() Logger { return DefaultLogger.getLogger(LvlDebug, 1) }

// Info is equal to DefaultLogger.Info().
func Info() Logger { return DefaultLogger.getLogger(LvlInfo, 1) }

// Warn is equal to DefaultLogger.Warn().
func Warn() Logger { return DefaultLogger.getLogger(LvlWarn, 1) }

// Error is equal to DefaultLogger.Error().
func Error() Logger { return DefaultLogger.getLogger(LvlError, 1) }

// Panic is equal to DefaultLogger.Panic().
func Panic() Logger { return DefaultLogger.getLogger(LvlPanic, 1) }

// Fatal is equal to DefaultLogger.Fatal().
func Fatal() Logger { return DefaultLogger.getLogger(LvlFatal, 1) }

// Tracef is equal to DefaultLogger.Trace().Printf(msg, args...).
func Tracef(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlTrace, 1).Printf(msg, args...)
}

// Debugf is equal to DefaultLogger.Debug().Printf(msg, args...).
func Debugf(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlDebug, 1).Printf(msg, args...)
}

// Infof is equal to DefaultLogger.Info().Printf(msg, args...).
func Infof(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlInfo, 1).Printf(msg, args...)
}

// Warnf is equal to DefaultLogger.Warn().Printf(msg, args...).
func Warnf(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlWarn, 1).Printf(msg, args...)
}

// Errorf is equal to DefaultLogger.Error().Printf(msg, args...).
func Errorf(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlError, 1).Printf(msg, args...)
}

// Panicf is equal to DefaultLogger.Panic().Printf(msg, args...).
func Panicf(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlPanic, 1).Printf(msg, args...)
}

// Fatalf is equal to DefaultLogger.Fatal().Printf(msg, args...).
func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlFatal, 1).Printf(msg, args...)
}

// Kv is equal to DefaultLogger.Kv(key, value).
func Kv(key string, value interface{}) Logger {
	return DefaultLogger.getLogger(DefaultLogger.level, 1).Kv(key, value)
}

// Kvs is equal to DefaultLogger.Kvs(key, value).
func Kvs(kvs ...interface{}) Logger {
	return DefaultLogger.getLogger(DefaultLogger.level, 1).Kvs(kvs...)
}

// Printf is equal to DefaultLogger.Printf(msg, args...).
func Printf(msg string, args ...interface{}) {
	DefaultLogger.getLogger(DefaultLogger.level, 1).Printf(msg, args...)
}

// Print is equal to DefaultLogger.Print(msg, args...).
func Print(args ...interface{}) {
	DefaultLogger.getLogger(DefaultLogger.level, 1).Print(args...)
}

// Ef is equal to DefaultLogger.Error().Kv("err", err).Printf(msg, args...).
func Ef(err error, msg string, args ...interface{}) {
	DefaultLogger.getLogger(LvlError, 1).Kv("err", err).Printf(msg, args...)
}

// IfErr logs the message and key-values with the ERROR level
// only if err is not equal to nil.
func IfErr(err error, msg string, kvs ...interface{}) { ifErr(err, msg, kvs...) }

// WrapPanic wraps and logs the panic.
func WrapPanic(kvs ...interface{}) { ifErr(recover(), "panic", kvs...) }

func ifErr(err interface{}, msg string, kvs ...interface{}) {
	if err == nil {
		return
	}
	DefaultLogger.getLogger(LvlError, 2).Kvs(kvs...).Kv("err", err).Printf(msg)
}
