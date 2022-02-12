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

import "log"

var stdlogDepth = 2

// StdLogger is equal to DefaultLogger.StdLogger(prefix, level).
func StdLogger(prefix string, level int) *log.Logger {
	return log.New(DefaultLogger.WithLevel(level).WithDepth(stdlogDepth), prefix, 0)
}

// StdLogger returns a new log.Logger based on the current logger engine
// with the prefix and level.
func (l Logger) StdLogger(prefix string, level int) *log.Logger {
	return log.New(l.WithLevel(level).WithDepth(stdlogDepth), prefix, 0)
}
