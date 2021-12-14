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
	LvlDisable = int(0)
	LvlTrace   = int(10)
	LvlDebug   = int(20)
	LvlInfo    = int(40)
	LvlWarn    = int(60)
	LvlError   = int(80)
	LvlAlert   = int(100)
	LvlPanic   = int(120)
	LvlFatal   = int(127)
)

// LevelIsValid reports whether the level is valid.
func LevelIsValid(level int) bool {
	return LvlFatal >= level && level >= LvlDisable
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
	if level >= LvlDisable {
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
	case LvlDisable:
		return "disable"
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
	default:
		checkLevel(level)
		return fmt.Sprintf("Level(%d)", level)
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
//
func ParseLevel(level string, defaultLevel ...int) int {
	switch strings.ToLower(level) {
	case "trace", "T":
		return LvlTrace
	case "debug", "D":
		return LvlDebug
	case "info", "I":
		return LvlInfo
	case "warn", "W":
		return LvlWarn
	case "error", "E":
		return LvlError
	case "alert", "A":
		return LvlAlert
	case "panic", "P":
		return LvlPanic
	case "fatal", "F":
		return LvlFatal
	default:
		if len(defaultLevel) == 0 {
			panic(fmt.Errorf("unknown level '%s'", level))
		}

		checkLevel(defaultLevel[0])
		return defaultLevel[0]
	}
}
