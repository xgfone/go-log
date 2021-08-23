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
	"errors"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/xgfone/go-atexit"
)

var fixDepth = func(depth int) int { return depth }

// Logger is a logger implementation.
type Logger struct {
	Name     string
	Ctxs     []Field
	Depth    int
	Encoder  Encoder
	ExitCode int

	logger atomic.Value
	level  uint32
}

// New creates a new Logger to encode the log record as JSON
// and output the log to os.Stdout.
func New(name string) *Logger {
	ctxs := []Field{CallerStack("stack", true)}
	encoder := NewJSONEncoder(SafeWriter(StreamWriter(os.Stdout)))
	return &Logger{
		Name:     name,
		Ctxs:     ctxs,
		level:    uint32(LvlDebug),
		Encoder:  encoder,
		ExitCode: 1,
	}
}

// NewSimpleLogger returns a new simple logger.
func NewSimpleLogger(name, level, filepath, filesize string, filenum int) *Logger {
	log := New(name)
	log.level = uint32(NameToLevel(level))
	if filepath != "" {
		log.Encoder.SetWriter(SafeWriter(FileWriter(filepath, filesize, filenum)))
	}
	return log
}

// GetLevel returns the level. It will return the level of the parent logger,
// however, if the logger is a child and does not set the level.
func (l *Logger) GetLevel() Level {
	if v := l.logger.Load(); v != nil {
		if logger := v.(*Logger); logger != nil {
			return logger.GetLevel()
		}
	}
	return Level(atomic.LoadUint32(&l.level))
}

// SetLevel resets the level.
func (l *Logger) SetLevel(level Level) {
	atomic.StoreUint32(&l.level, uint32(level))
	l.logger.Store((*Logger)(nil))
}

// SetParent sets the parent logger to inherit its level.
//
// If logger is nil, clear the parent logger.
func (l *Logger) SetParent(logger *Logger) { l.logger.Store(logger) }

// StdLog converts the Logger to the std log.
func (l *Logger) StdLog(prefix string, flags ...int) *log.Logger {
	flag := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	if len(flags) > 0 {
		flag = flags[0]
	}
	return log.New(NewIOWriter(l.Encoder.Writer(), l.GetLevel()), prefix, flag)
}

// New clones itself as the parent and returns a new one as the child.
func (l *Logger) New() *Logger {
	var ctxs []Field
	if len(l.Ctxs) != 0 {
		ctxs = append([]Field{}, l.Ctxs...)
	}

	logger := &Logger{
		Ctxs:    ctxs,
		Name:    l.Name,
		Depth:   l.Depth,
		Encoder: l.Encoder,
	}
	logger.logger.Store(l)
	return logger
}

// WithName returns a new Logger with the new name.
func (l *Logger) WithName(name string) *Logger {
	ll := l.New()
	ll.Name = name
	return ll
}

// WithLevel returns a new Logger with the new level.
func (l *Logger) WithLevel(level Level) *Logger {
	ll := l.New()
	ll.SetLevel(level)
	return ll
}

// WithEncoder returns a new Logger with the new encoder.
func (l *Logger) WithEncoder(e Encoder) *Logger {
	ll := l.New()
	ll.Encoder = e
	return ll
}

// WithDepth returns a new Logger, which will increase the depth.
func (l *Logger) WithDepth(depth int) *Logger {
	ll := l.New()
	ll.Depth += depth
	return ll
}

// WithCtx returns a new Logger with the new context fields.
func (l *Logger) WithCtx(ctxs ...Field) *Logger {
	ll := l.New()
	ll.Ctxs = append(ll.Ctxs, ctxs...)
	return ll
}

// Log emits the logs with the level and the depth.
//
// If lvl is equal to LvlFatal, the program exits with ExitCode.
func (l *Logger) Log(lvl Level, depth int, msgfmt string, msgargs []interface{},
	fields []Field) {
	if lvl < l.GetLevel() {
		return
	}

	if len(msgargs) != 0 {
		msgfmt = fmt.Sprintf(msgfmt, msgargs...)
	}

	r := Record{
		Name:   l.Name,
		Depth:  l.Depth + 1 + fixDepth(depth),
		Lvl:    lvl,
		Msg:    msgfmt,
		Ctxs:   l.Ctxs,
		Fields: fields,
	}
	l.Encoder.Encode(r)

	if lvl == LvlFatal {
		atexit.Exit(l.ExitCode)
	}
}

// Trace is equal to Log(LvlTrace, 1, msg, nil, fields).
func (l *Logger) Trace(msg string, fields ...Field) { l.Log(LvlTrace, 1, msg, nil, fields) }

// Debug is equal to Log(LvlDebug, 1, msg, nil, fields).
func (l *Logger) Debug(msg string, fields ...Field) { l.Log(LvlDebug, 1, msg, nil, fields) }

// Info is equal to Log(LvlInfo, 1, msg, nil, fields).
func (l *Logger) Info(msg string, fields ...Field) { l.Log(LvlInfo, 1, msg, nil, fields) }

// Warn is equal to Log(LvlWarn, 1, msg, nil, fields).
func (l *Logger) Warn(msg string, fields ...Field) { l.Log(LvlWarn, 1, msg, nil, fields) }

// Error is equal to Log(LvlError, 1, msg, nil, fields).
func (l *Logger) Error(msg string, fields ...Field) { l.Log(LvlError, 1, msg, nil, fields) }

// Fatal is equal to Log(LvlFatal, 1, msg, nil, fields).
func (l *Logger) Fatal(msg string, fields ...Field) { l.Log(LvlFatal, 1, msg, nil, fields) }

// Tracef is equal to Log(LvlTrace, 1, msg, args, nil).
func (l *Logger) Tracef(msg string, args ...interface{}) { l.Log(LvlTrace, 1, msg, args, nil) }

// Debugf is equal to Log(LvlDebug, 1, msg, args, nil).
func (l *Logger) Debugf(msg string, args ...interface{}) { l.Log(LvlDebug, 1, msg, args, nil) }

// Infof is equal to Log(LvlInfo, 1, msg, args, nil).
func (l *Logger) Infof(msg string, args ...interface{}) { l.Log(LvlInfo, 1, msg, args, nil) }

// Warnf is equal to Log(LvlWarn, 1, msg, args, nil).
func (l *Logger) Warnf(msg string, args ...interface{}) { l.Log(LvlWarn, 1, msg, args, nil) }

// Errorf is equal to Log(LvlError, 1, msg, args, nil).
func (l *Logger) Errorf(msg string, args ...interface{}) { l.Log(LvlError, 1, msg, args, nil) }

// Fatalf is equal to Log(LvlFatal, 1, msg, args, nil).
func (l *Logger) Fatalf(msg string, args ...interface{}) { l.Log(LvlFatal, 1, msg, args, nil) }

// Printf is equal to Infof(msg, args...).
func (l *Logger) Printf(msg string, args ...interface{}) { l.Log(LvlInfo, 1, msg, args, nil) }

func (l *Logger) logs(lvl Level, depth int, msg string, keyAndValues []interface{}) {
	_len := len(keyAndValues)
	if _len%2 != 0 {
		panic(errors.New("Logger: the number of keyAndValues is not even"))
	}

	_len /= 2
	fields := make([]Field, _len)
	for i := 0; i < _len; i++ {
		j := i * 2
		fields[i] = F(keyAndValues[j].(string), keyAndValues[j+1])
	}

	l.Log(lvl, depth+1, msg, nil, fields)
}

// Traces is the same as Trace, but convert keyAndValues to []Field.
func (l *Logger) Traces(msg string, keyAndValues ...interface{}) {
	l.logs(LvlTrace, 1, msg, keyAndValues)
}

// Debugs is the same as Debug, but convert keyAndValues to []Field.
func (l *Logger) Debugs(msg string, keyAndValues ...interface{}) {
	l.logs(LvlDebug, 1, msg, keyAndValues)
}

// Infos is the same as Info, but convert keyAndValues to []Field.
func (l *Logger) Infos(msg string, keyAndValues ...interface{}) {
	l.logs(LvlInfo, 1, msg, keyAndValues)
}

// Warns is the same as Warn, but convert keyAndValues to []Field.
func (l *Logger) Warns(msg string, keyAndValues ...interface{}) {
	l.logs(LvlWarn, 1, msg, keyAndValues)
}

// Errors is the same as Error, but convert keyAndValues to []Field.
func (l *Logger) Errors(msg string, keyAndValues ...interface{}) {
	l.logs(LvlError, 1, msg, keyAndValues)
}

// Fatals is the same as Fatal, but convert keyAndValues to []Field.
func (l *Logger) Fatals(msg string, keyAndValues ...interface{}) {
	l.logs(LvlFatal, 1, msg, keyAndValues)
}
