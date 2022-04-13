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

package log

import "io"

// DefaultLogger is the default global logger.
var DefaultLogger = New("").WithHooks(Caller("caller"))

// SetWriter is eqaul to DefaultLogger.SetWriter(w).
func SetWriter(w io.Writer) { DefaultLogger.SetWriter(w) }

// SetEncoder is eqaul to DefaultLogger.SetEncoder(enc).
func SetEncoder(enc Encoder) { DefaultLogger.SetEncoder(enc) }

// SetLevel is equal to DefaultLogger.SetLevel(level).
func SetLevel(level int) { DefaultLogger.SetLevel(level) }

// GetLevel is equal to DefaultLogger.GetLevel().
func GetLevel() int { return DefaultLogger.GetLevel() }

// Clone is equal to DefaultLogger.Clone().
func Clone() Logger { return DefaultLogger.Clone() }

// WithName is equal to DefaultLogger.New(name).
func WithName(name string) Logger { return DefaultLogger.WithName(name) }

// WithDepth is equal to DefaultLogger.WithDepth(depth).
func WithDepth(depth int) Logger { return DefaultLogger.WithDepth(depth) }

// WithLevel is equal to DefaultLogger.WithLevel(level).
func WithLevel(level int) Logger { return DefaultLogger.WithLevel(level) }

// WithHooks is equal to DefaultLogger.WithHooks(hooks...).
func WithHooks(hooks ...Hook) Logger { return DefaultLogger.WithHooks(hooks...) }

// WithContext is equal to DefaultLogger.WithContext(key, value).
func WithContext(key string, value interface{}) Logger {
	return DefaultLogger.WithContext(key, value)
}

// WithContexts is equal to DefaultLogger.WithContexts(kvs...).
func WithContexts(kvs ...interface{}) Logger {
	return DefaultLogger.WithContexts(kvs...)
}

// WithLevelFormat is equal to DefaultLogger.WithLevelFormat(format).
func WithLevelFormat(format func(level int) string) Logger {
	return DefaultLogger.WithLevelFormat(format)
}

// SetLevelFormat is equal to DefaultLogger.SetLevelFormat(format).
func SetLevelFormat(format func(level int) string) {
	DefaultLogger.SetLevelFormat(format)
}

// LevelLog is equal to DefaultLogger.Level(level, depth).
func LevelLog(level, depth int) *Emitter { return DefaultLogger.Level(level, depth+1) }

// Trace is equal to DefaultLogger.Trace().
func Trace() *Emitter { return DefaultLogger.getEmitter(LvlTrace, 1) }

// Debug is equal to DefaultLogger.Debug().
func Debug() *Emitter { return DefaultLogger.getEmitter(LvlDebug, 1) }

// Info is equal to DefaultLogger.Info().
func Info() *Emitter { return DefaultLogger.getEmitter(LvlInfo, 1) }

// Warn is equal to DefaultLogger.Warn().
func Warn() *Emitter { return DefaultLogger.getEmitter(LvlWarn, 1) }

// Error is equal to DefaultLogger.Error().
func Error() *Emitter { return DefaultLogger.getEmitter(LvlError, 1) }

// Alert is equal to DefaultLogger.Alert()).
func Alert() *Emitter { return DefaultLogger.getEmitter(LvlAlert, 1) }

// Panic is equal to DefaultLogger.Panic().
func Panic() *Emitter { return DefaultLogger.getEmitter(LvlPanic, 1) }

// Fatal is equal to DefaultLogger.Fatal().
func Fatal() *Emitter { return DefaultLogger.getEmitter(LvlFatal, 1) }

// Ef is equal to DefaultLogger.Error().Kv("err", err).Printf(format, args...).
func Ef(err error, format string, args ...interface{}) {
	DefaultLogger.getEmitter(LvlError, 1).Kv("err", err).Printf(format, args...)
}

// IfErr is the same as Err, but only if err is not equal to nil.
func IfErr(err error, msg string, keysAndValues ...interface{}) {
	if err != nil {
		DefaultLogger.getEmitter(LvlError, 1).Kvs(keysAndValues...).
			Kv("err", err).Printf(msg)
	}
}

// WrapPanic wraps and logs the panic, which should be called directly with defer,
// For example,
//   defer WrapPanic()
//   defer WrapPanic("key1", "value1")
//   defer WrapPanic("key1", "value1", "key2", "value2")
//   defer WrapPanic("key1", "value1", "key2", "value2", "key3", "value3")
func WrapPanic(keysAndValues ...interface{}) {
	if r := recover(); r != nil {
		DefaultLogger.getEmitter(LvlError, 2).Kvs(keysAndValues...).
			StrSlice("stacks", GetCallStack(4)).Kv("panic", r).
			Printf("wrap a panic")
	}
}
