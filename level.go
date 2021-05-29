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
)

// Predefine some levels.
const (
	LvlTrace = Level(0)
	LvlDebug = Level(50)
	LvlInfo  = Level(100)
	LvlWarn  = Level(150)
	LvlError = Level(200)
	LvlFatal = Level(255)
)

// Level is the level of the log.
type Level uint8

func (l Level) String() string {
	switch l {
	case LvlTrace:
		return "TRACE"
	case LvlDebug:
		return "DEBUG"
	case LvlInfo:
		return "INFO"
	case LvlWarn:
		return "WARN"
	case LvlError:
		return "ERROR"
	case LvlFatal:
		return "FATAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", l)
	}
}

// NameToLevel returns a Level by the level name.
//
// Support the level name, which is case insensitive:
//   TRACE
//   DEBUG
//   INFO
//   WARN
//   ERROR
//   FATAL
func NameToLevel(level string, defaultLevel ...Level) Level {
	switch strings.ToUpper(level) {
	case "TRACE":
		return LvlTrace
	case "DEBUG":
		return LvlDebug
	case "INFO":
		return LvlInfo
	case "WARN":
		return LvlWarn
	case "ERROR":
		return LvlError
	case "FATAL":
		return LvlFatal
	default:
		if len(defaultLevel) > 0 {
			return defaultLevel[0]
		}
		panic(fmt.Errorf("unknown level '%s'", level))
	}
}
