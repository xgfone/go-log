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
func NewLogger(logger glog.Logger, depth int) Logger {
	return loggerf{Logger: logger, depth: fixDepth(depth)}
}

type loggerf struct {
	glog.Logger
	depth int
}

var fixDepth = func(depth int) int { return depth + 1 }

func (l loggerf) Tracef(format string, args ...interface{}) {
	l.Level(glog.LvlTrace, l.depth).Printf(format, args...)
}

func (l loggerf) Debugf(format string, args ...interface{}) {
	l.Level(glog.LvlDebug, l.depth).Printf(format, args...)
}

func (l loggerf) Infof(format string, args ...interface{}) {
	l.Level(glog.LvlInfo, l.depth).Printf(format, args...)
}

func (l loggerf) Warnf(format string, args ...interface{}) {
	l.Level(glog.LvlWarn, l.depth).Printf(format, args...)
}

func (l loggerf) Errorf(format string, args ...interface{}) {
	l.Level(glog.LvlError, l.depth).Printf(format, args...)
}

func (l loggerf) Alertf(format string, args ...interface{}) {
	l.Level(glog.LvlAlert, l.depth).Printf(format, args...)
}

func (l loggerf) Panicf(format string, args ...interface{}) {
	l.Level(glog.LvlPanic, l.depth).Printf(format, args...)
}

func (l loggerf) Fatalf(format string, args ...interface{}) {
	l.Level(glog.LvlFatal, l.depth).Printf(format, args...)
}
