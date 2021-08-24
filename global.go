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

// DefalutLogger is the default global logger.
var DefalutLogger = New("root")

// GetLevel is equal to DefalutLogger.GetLevel() to return the level
// of the default logger.
func GetLevel() Level { return DefalutLogger.GetLevel() }

// SetLevel is equal to DefalutLogger.SetLevel(level) to reset the level
// of the default logger.
func SetLevel(level Level) { DefalutLogger.SetLevel(level) }

// WithCtx is equal to DefalutLogger.WithCtx(fields...).
func WithCtx(fields ...Field) *Logger { return DefalutLogger.WithCtx(fields...) }

// WithName is equal to DefalutLogger.WithName(name).
func WithName(name string) *Logger { return DefalutLogger.WithName(name) }

// WithLevel is equal to DefalutLogger.WithLevel(level).
func WithLevel(level Level) *Logger { return DefalutLogger.WithLevel(level) }

// WithEncoder is equal to DefalutLogger.WithEncoder(enc).
func WithEncoder(enc Encoder) *Logger { return DefalutLogger.WithEncoder(enc) }

// WithDepth is equal to DefalutLogger.WithDepth(depth).
func WithDepth(depth int) *Logger { return DefalutLogger.WithDepth(depth) }

// Trace is equal to DefalutLogger.Trace(msg, fields...).
func Trace(msg string, fields ...Field) {
	DefalutLogger.Log(LvlTrace, 1, msg, nil, fields)
}

// Debug is equal to DefalutLogger.Debug(msg, fields...).
func Debug(msg string, fields ...Field) {
	DefalutLogger.Log(LvlDebug, 1, msg, nil, fields)
}

// Info is equal to DefalutLogger.Info(msg, fields...).
func Info(msg string, fields ...Field) {
	DefalutLogger.Log(LvlInfo, 1, msg, nil, fields)
}

// Warn is equal to DefalutLogger.Warn(msg, fields...).
func Warn(msg string, fields ...Field) {
	DefalutLogger.Log(LvlWarn, 1, msg, nil, fields)
}

// Error is equal to DefalutLogger.Error(msg, fields...).
func Error(msg string, fields ...Field) {
	DefalutLogger.Log(LvlError, 1, msg, nil, fields)
}

// Fatal is equal to DefalutLogger.Fatal(msg, fields...).
func Fatal(msg string, fields ...Field) {
	DefalutLogger.Log(LvlFatal, 1, msg, nil, fields)
}

// Traces is equal to DefalutLogger.Traces.
func Traces(msg string, keyAndValues ...interface{}) {
	DefalutLogger.logs(LvlTrace, 1, msg, keyAndValues)
}

// Debugs is equal to DefalutLogger.Debugs.
func Debugs(msg string, keyAndValues ...interface{}) {
	DefalutLogger.logs(LvlDebug, 1, msg, keyAndValues)
}

// Infos is equal to DefalutLogger.Infos.
func Infos(msg string, keyAndValues ...interface{}) {
	DefalutLogger.logs(LvlInfo, 1, msg, keyAndValues)
}

// Warns is equal to DefalutLogger.Warns.
func Warns(msg string, keyAndValues ...interface{}) {
	DefalutLogger.logs(LvlWarn, 1, msg, keyAndValues)
}

// Errors equal to DefalutLogger.Errors.
func Errors(msg string, keyAndValues ...interface{}) {
	DefalutLogger.logs(LvlError, 1, msg, keyAndValues)
}

// Fatals is equal to DefalutLogger.Fatals.
func Fatals(msg string, keyAndValues ...interface{}) {
	DefalutLogger.logs(LvlFatal, 1, msg, keyAndValues)
}

// Tracef is equal to DefalutLogger.Tracef(msg, args...).
func Tracef(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlTrace, 1, msg, args, nil)
}

// Debugf is equal to DefalutLogger.Debugf(msg, args...).
func Debugf(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlDebug, 1, msg, args, nil)
}

// Infof is equal to DefalutLogger.Infof(msg, args...).
func Infof(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlInfo, 1, msg, args, nil)
}

// Warnf is equal to DefalutLogger.Warnf(msg, args...).
func Warnf(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlWarn, 1, msg, args, nil)
}

// Errorf is equal to DefalutLogger.Errorf(msg, args...).
func Errorf(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlError, 1, msg, args, nil)
}

// Fatalf is equal to DefalutLogger.Fatalf(msg, args...).
func Fatalf(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlFatal, 1, msg, args, nil)
}

// Printf is equal to DefalutLogger.Infof(msg, args...).
func Printf(msg string, args ...interface{}) {
	DefalutLogger.Log(LvlInfo, 1, msg, args, nil)
}

// Ef is equal to DefalutLogger.Error(fmt.Sprintf(msg, args), E(err)).
func Ef(err error, format string, args ...interface{}) {
	DefalutLogger.Log(LvlError, 1, format, args, []Field{E(err)})
}
