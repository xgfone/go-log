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

// DefaultLogger is the default global logger.
var DefaultLogger = New("root")

// GetLevel is equal to DefaultLogger.GetLevel() to return the level
// of the default logger.
func GetLevel() Level { return DefaultLogger.GetLevel() }

// SetLevel is equal to DefaultLogger.SetLevel(level) to reset the level
// of the default logger.
func SetLevel(level Level) { DefaultLogger.SetLevel(level) }

// WithCtx is equal to DefaultLogger.WithCtx(fields...).
func WithCtx(fields ...Field) *Logger { return DefaultLogger.WithCtx(fields...) }

// WithName is equal to DefaultLogger.WithName(name).
func WithName(name string) *Logger { return DefaultLogger.WithName(name) }

// WithLevel is equal to DefaultLogger.WithLevel(level).
func WithLevel(level Level) *Logger { return DefaultLogger.WithLevel(level) }

// WithEncoder is equal to DefaultLogger.WithEncoder(enc).
func WithEncoder(enc Encoder) *Logger { return DefaultLogger.WithEncoder(enc) }

// WithDepth is equal to DefaultLogger.WithDepth(depth).
func WithDepth(depth int) *Logger { return DefaultLogger.WithDepth(depth) }

// Trace is equal to DefaultLogger.Trace(msg, fields...).
func Trace(msg string, fields ...Field) {
	DefaultLogger.Log(LvlTrace, 1, msg, nil, fields)
}

// Debug is equal to DefaultLogger.Debug(msg, fields...).
func Debug(msg string, fields ...Field) {
	DefaultLogger.Log(LvlDebug, 1, msg, nil, fields)
}

// Info is equal to DefaultLogger.Info(msg, fields...).
func Info(msg string, fields ...Field) {
	DefaultLogger.Log(LvlInfo, 1, msg, nil, fields)
}

// Warn is equal to DefaultLogger.Warn(msg, fields...).
func Warn(msg string, fields ...Field) {
	DefaultLogger.Log(LvlWarn, 1, msg, nil, fields)
}

// Error is equal to DefaultLogger.Error(msg, fields...).
func Error(msg string, fields ...Field) {
	DefaultLogger.Log(LvlError, 1, msg, nil, fields)
}

// Fatal is equal to DefaultLogger.Fatal(msg, fields...).
func Fatal(msg string, fields ...Field) {
	DefaultLogger.Log(LvlFatal, 1, msg, nil, fields)
}

// Traces is equal to DefaultLogger.Traces.
func Traces(msg string, keyAndValues ...interface{}) {
	DefaultLogger.logs(LvlTrace, 1, msg, keyAndValues)
}

// Debugs is equal to DefaultLogger.Debugs.
func Debugs(msg string, keyAndValues ...interface{}) {
	DefaultLogger.logs(LvlDebug, 1, msg, keyAndValues)
}

// Infos is equal to DefaultLogger.Infos.
func Infos(msg string, keyAndValues ...interface{}) {
	DefaultLogger.logs(LvlInfo, 1, msg, keyAndValues)
}

// Warns is equal to DefaultLogger.Warns.
func Warns(msg string, keyAndValues ...interface{}) {
	DefaultLogger.logs(LvlWarn, 1, msg, keyAndValues)
}

// Errors equal to DefaultLogger.Errors.
func Errors(msg string, keyAndValues ...interface{}) {
	DefaultLogger.logs(LvlError, 1, msg, keyAndValues)
}

// Fatals is equal to DefaultLogger.Fatals.
func Fatals(msg string, keyAndValues ...interface{}) {
	DefaultLogger.logs(LvlFatal, 1, msg, keyAndValues)
}

// Tracef is equal to DefaultLogger.Tracef(msg, args...).
func Tracef(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlTrace, 1, msg, args, nil)
}

// Debugf is equal to DefaultLogger.Debugf(msg, args...).
func Debugf(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlDebug, 1, msg, args, nil)
}

// Infof is equal to DefaultLogger.Infof(msg, args...).
func Infof(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlInfo, 1, msg, args, nil)
}

// Warnf is equal to DefaultLogger.Warnf(msg, args...).
func Warnf(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlWarn, 1, msg, args, nil)
}

// Errorf is equal to DefaultLogger.Errorf(msg, args...).
func Errorf(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlError, 1, msg, args, nil)
}

// Fatalf is equal to DefaultLogger.Fatalf(msg, args...).
func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlFatal, 1, msg, args, nil)
}

// Printf is equal to DefaultLogger.Infof(msg, args...).
func Printf(msg string, args ...interface{}) {
	DefaultLogger.Log(LvlInfo, 1, msg, args, nil)
}

// Ef is equal to DefaultLogger.Error(fmt.Sprintf(msg, args), E(err)).
func Ef(err error, format string, args ...interface{}) {
	DefaultLogger.Log(LvlError, 1, format, args, []Field{E(err)})
}

// IfErr logs the message and fields with the ERROR level
// only if err is not equal to nil.
func IfErr(err error, msg string, fields ...Field) {
	if err != nil {
		if len(fields) == 0 {
			fields = []Field{E(err)}
		} else {
			fields = append(fields, E(err))
		}

		DefaultLogger.Log(LvlError, 1, msg, nil, fields)
	}
}
