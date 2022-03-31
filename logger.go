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
	"os"
	"strings"

	"github.com/xgfone/go-log/writer"
)

// Logger is the structured logger based on the key-value.
type Logger struct {
	*Output

	name    string
	level   int
	depth   int
	sampler Sampler
	fmtLvl  func(int) string

	// Key-Value Context
	hooks []Hook
	ctxs  []interface{}
	ctx   []byte
}

// New creates a new root logger, which encodes the log message as JSON
// and output the encoded log to os.Stderr.
func New(name string) Logger {
	return Logger{
		name:   name,
		level:  LvlDebug,
		Output: NewOutput(writer.SafeWriter(os.Stderr), nil),
	}
}

// Clone clones itself and returns a new one.
func (l Logger) Clone() Logger {
	return Logger{
		Output: l.Output,

		name:    l.name,
		level:   l.level,
		depth:   l.depth,
		sampler: l.sampler,

		hooks: append([]Hook{}, l.hooks...),
		ctxs:  append([]interface{}{}, l.ctxs...),
		ctx:   append([]byte{}, l.ctx...),
	}
}

// Name returns the name of the current logger.
func (l Logger) Name() string { return l.name }

// Depth returns the stack depth of the current logger.
func (l Logger) Depth() int { return l.depth }

// GetLevel returns the level of the current logger.
func (l Logger) GetLevel() int { return l.level }

// FormatLevel formats the level to string.
func (l Logger) FormatLevel(level int) string {
	if l.fmtLvl != nil {
		return l.fmtLvl(level)
	}
	return formatLevel(level)
}

// WithFormatLevel returns a new logger with the customized level formatter.
//
// If format is nil, use FormatLevel instead.
func (l Logger) WithFormatLevel(format func(level int) string) Logger {
	l.fmtLvl = format
	return l
}

// WithName returns a new logger with the name.
//
//   - If name is empty, the name of the new logger is equal to l.name.
//   - If both name and l.name are not empty, it is equal to l.Name()+"."+name.
//
func (l Logger) WithName(name string) Logger {
	if len(name) == 0 {
		name = l.name
	} else if len(l.name) > 0 {
		name = strings.Join([]string{l.name, name}, ".")
	}

	l = l.Clone()
	l.name = name
	return l
}

// WithDepth returns a new logger with the depth of the stack out of the logger.
func (l Logger) WithDepth(depth int) Logger {
	l = l.Clone()
	l.depth = depth
	return l
}

// WithLevel returns a new logger with the level.
func (l Logger) WithLevel(level int) Logger {
	checkLevel(level)
	l = l.Clone()
	l.level = level
	return l
}

// Contexts returns the key-value contexts.
func (l Logger) Contexts() (kvs []interface{}) { return l.ctxs }

// WithContext returns a new logger that appends the key-value context.
func (l Logger) WithContext(key string, value interface{}) Logger {
	l = l.Clone()
	l.ctx = l.Output.encoder.Encode(l.ctx, key, value)
	l.ctxs = append(l.ctxs, key, value)
	return l
}

// WithContexts returns a new logger that appends a set of the key-value contexts.
func (l Logger) WithContexts(kvs ...interface{}) Logger {
	l = l.Clone()
	l.appendContexts(kvs...)
	return l
}

func (l *Logger) appendContexts(kvs ...interface{}) {
	_len := len(kvs)
	if _len%2 != 0 {
		panic("the length of the key-value log contexts is not even")
	}

	for i := 0; i < _len; i += 2 {
		l.ctx = l.Output.encoder.Encode(l.ctx, kvs[i].(string), kvs[i+1])
	}
	l.ctxs = append(l.ctxs, kvs...)
}

// Write implements the interface io.Writer.
func (l Logger) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[:n-1]
	}
	l.getEmitter(l.level, 1).Printf(string(p))
	return
}
