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
	"strings"
	"sync/atomic"
)

// Predefine some levels.
const (
	LvlTrace   = int(0)
	LvlDebug   = int(20)
	LvlInfo    = int(40)
	LvlWarn    = int(60)
	LvlError   = int(80)
	LvlAlert   = int(100)
	LvlPanic   = int(120)
	LvlFatal   = int(126)
	LvlDisable = int(127)
)

// LevelIsValid reports whether the level is valid, that's, [LvlTrace, LvlDisable].
func LevelIsValid(level int) bool {
	return LvlTrace <= level && level <= LvlDisable
}

func checkLevel(level int) {
	if !LevelIsValid(level) {
		panic(fmt.Errorf("invalid level '%d'", level))
	}
}

var globalLevel = int64(-1)

// SetGlobalLevel sets the global level, which will act on all the logger.
//
// If the level is negative, it will unset the global level.
func SetGlobalLevel(level int) {
	if level >= LvlTrace {
		checkLevel(level)
	}
	atomic.StoreInt64(&globalLevel, int64(level))
}

// GetGlobalLevel returns the global level setting.
//
// Notice: if the returned level value is negative, it represents
// that no global level is set.
func GetGlobalLevel() int { return int(atomic.LoadInt64(&globalLevel)) }

// FormatLevel is used to format the level to string.
var FormatLevel func(level int) string = formatLevel

func formatLevel(level int) string {
	switch level {
	case LvlTrace:
		return "trace"
	case LvlDebug:
		return "debug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "error"
	case LvlAlert:
		return "alert"
	case LvlPanic:
		return "panic"
	case LvlFatal:
		return "fatal"
	case LvlDisable:
		return "disable"
	default:
		checkLevel(level)
		switch {
		case level < LvlDebug:
			return fmt.Sprintf("trace%d", level-LvlTrace)

		case level < LvlInfo:
			return fmt.Sprintf("debug%d", level-LvlDebug)

		case level < LvlWarn:
			return fmt.Sprintf("info%d", level-LvlInfo)

		case level < LvlError:
			return fmt.Sprintf("warn%d", level-LvlWarn)

		case level < LvlAlert:
			return fmt.Sprintf("error%d", level-LvlError)

		case level < LvlPanic:
			return fmt.Sprintf("alert%d", level-LvlAlert)

		case level < LvlFatal:
			return fmt.Sprintf("panic%d", level-LvlPanic)

		case level < LvlDisable:
			return fmt.Sprintf("fatal%d", level-LvlFatal)

		default:
			return fmt.Sprintf("Level(%d)", level)
		}
	}
}

// ParseLevel parses a string to the level.
//
// Support the level string as follow, which is case insensitive:
//
//   trace
//   debug
//   info
//   warn
//   error
//   alert
//   panic
//   fatal
//   disable
//
func ParseLevel(level string, defaultLevel ...int) int {
	switch strings.ToLower(level) {
	case "trace":
		return LvlTrace
	case "debug":
		return LvlDebug
	case "info":
		return LvlInfo
	case "warn":
		return LvlWarn
	case "error":
		return LvlError
	case "alert":
		return LvlAlert
	case "panic":
		return LvlPanic
	case "fatal":
		return LvlFatal
	case "disable":
		return LvlDisable
	default:
		if len(defaultLevel) == 0 {
			panic(fmt.Errorf("unknown level '%s'", level))
		}

		checkLevel(defaultLevel[0])
		return defaultLevel[0]
	}
}

// Enabled reports whether the given level is enabled.
func (l Logger) Enabled(level int) bool {
	checkLevel(level)
	return !l.isDisabled(level)
}

func (l Logger) isDisabled(level int) bool {
	if level == LvlDisable {
		return true
	}

	global := GetGlobalLevel()
	if global < LvlTrace {
		return l.disabled(level, l.level)
	}
	return l.disabled(level, global)
}

func (l Logger) disabled(logLevel, minThresholdLevel int) bool {
	if logLevel < minThresholdLevel {
		return true
	}

	if l.sampler != nil && globalSamplingIsEnabled() {
		return !l.sampler.Sample(l.name, logLevel)
	}

	return false
}
