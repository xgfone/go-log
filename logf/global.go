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

import "github.com/xgfone/go-log"

// Tracef is equal to log.DefaultLogger.Trace().Printf(msg, args...).
func Tracef(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlTrace, 1).Printf(msg, args...)
}

// Debugf is equal to log.DefaultLogger.Debug().Printf(msg, args...).
func Debugf(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlDebug, 1).Printf(msg, args...)
}

// Infof is equal to log.DefaultLogger.Info().Printf(msg, args...).
func Infof(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlInfo, 1).Printf(msg, args...)
}

// Warnf is equal to log.DefaultLogger.Warn().Printf(msg, args...).
func Warnf(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlWarn, 1).Printf(msg, args...)
}

// Errorf is equal to log.DefaultLogger.Error().Printf(msg, args...).
func Errorf(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlError, 1).Printf(msg, args...)
}

// Alertf is equal to log.DefaultLogger.Alert().Printf(msg, args...).
func Alertf(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlAlert, 1).Printf(msg, args...)
}

// Panicf is equal to log.DefaultLogger.Panic().Printf(msg, args...).
func Panicf(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlPanic, 1).Printf(msg, args...)
}

// Fatalf is equal to log.DefaultLogger.Fatal().Printf(msg, args...).
func Fatalf(msg string, args ...interface{}) {
	log.DefaultLogger.Level(log.LvlFatal, 1).Printf(msg, args...)
}
