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

package logf

import glog "github.com/xgfone/go-log"

// Logger is an logger interface based on the format.
type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Alertf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// NewLogger returns a new Logger based on the format.
func NewLogger(log *glog.Engine, depth int) Logger {
	if log == nil {
		log = glog.DefaultLogger
	}
	return logger{Engine: log, depth: fixDepth(depth)}
}

type logger struct {
	*glog.Engine
	depth int
}

var fixDepth = func(depth int) int { return depth + 1 }

func (l logger) Tracef(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlTrace, l.depth).Printf(format, args...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlDebug, l.depth).Printf(format, args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlInfo, l.depth).Printf(format, args...)
}

func (l logger) Warnf(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlWarn, l.depth).Printf(format, args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlError, l.depth).Printf(format, args...)
}

func (l logger) Alertf(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlAlert, l.depth).Printf(format, args...)
}

func (l logger) Panicf(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlPanic, l.depth).Printf(format, args...)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.Engine.Logger(glog.LvlFatal, l.depth).Printf(format, args...)
}
