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
	"io"
	"log"
	"os"
	"strings"

	"github.com/xgfone/go-log/writer"
)

type kvctx struct {
	Key   string
	Value interface{}
}

// Engine is a logger engine.
type Engine struct {
	Output  *Output
	sampler Sampler

	name  string
	level int
	depth int

	// Key-Value Context
	hooks []Hook
	origs []kvctx
	ctx   []byte
}

// New creates a new root logger engine, which has no parent engine,
// to encode the log record as JSON and output the log to os.Stderr.
func New(name string) *Engine {
	return &Engine{
		name:   name,
		level:  LvlDebug,
		Output: NewOutput(writer.SafeWriter(os.Stderr), nil),
	}
}

// Name returns the name of the current log engine.
func (e *Engine) Name() string { return e.name }

// SetName resets the name of the log engine.
func (e *Engine) SetName(name string) *Engine {
	e.name = name
	return e
}

// SetDepth sets the depth of the stack out of the log engine and returns itself.
func (e *Engine) SetDepth(depth int) *Engine {
	e.depth = depth
	return e
}

// SetLevel resets the level.
//
// Notice: the level must not be negative.
func (e *Engine) SetLevel(level int) *Engine {
	checkLevel(level)
	e.level = level
	return e
}

// GetLevel returns the level of the logger engine.
func (e *Engine) GetLevel() int { return e.level }

// Enabled reports whether the level of the engine is enabled compared to
// the global level.
//
// Notice: if the global level is not set, it is equal to be enabled.
func (e *Engine) Enabled() bool {
	if e.level == LvlDisable {
		return false
	}

	global := GetGlobalLevel()
	if global < LvlTrace {
		return true
	}

	return e.level >= global
}

// Enable reports whether the given level is enabled.
func (e *Engine) Enable(level int) bool {
	checkLevel(level)
	return !e.isDisabled(level)
}

func (e *Engine) isDisabled(level int) bool {
	if level == LvlDisable {
		return true
	}

	global := GetGlobalLevel()
	if global < LvlTrace {
		return e.disabled(level, e.level)
	}
	return e.disabled(level, global)
}

func (e *Engine) disabled(logLevel, minThresholdLevel int) bool {
	if logLevel < minThresholdLevel {
		return true
	}

	if e.sampler != nil && globalSamplingIsEnabled() {
		return !e.sampler.Sample(e.name, logLevel)
	}

	return false
}

// Clone clones itself and returns a new one.
func (e *Engine) Clone() *Engine {
	return &Engine{
		name:    e.name,
		level:   e.level,
		depth:   e.depth,
		Output:  e.Output,
		sampler: e.sampler,

		hooks: append([]Hook{}, e.hooks...),
		origs: append([]kvctx{}, e.origs...),
		ctx:   append([]byte{}, e.ctx...),
	}
}

// New is the same as Clone, but the name of the new engine is equal to
// e.Name()+"."+name.
func (e *Engine) New(name string) *Engine {
	ee := e.Clone()
	if len(e.name) == 0 {
		ee.name = name
	} else if len(name) > 0 {
		ee.name = strings.Join([]string{ee.name, name}, ".")
	}
	return ee
}

// SetEncoder is the convenient function to set the encoder of the output to enc,
// which is equal to e.Output.SetEncoder(enc), but reencodes all the key-value
// contexts.
func (e *Engine) SetEncoder(enc Encoder) *Engine {
	e.Output.SetEncoder(enc)
	e.EncodeCtxs()
	return e
}

// SetWriter is the convenient function to set the writer of the output to w.
func (e *Engine) SetWriter(w io.Writer) *Engine {
	e.Output.SetWriter(w)
	return e
}

// ResetCtxs resets the key-value contexts to kvs and returns itself.
func (e *Engine) ResetCtxs(kvs ...interface{}) *Engine {
	e.ctx, e.origs = e.ctx[:0], e.origs[:0]
	return e.AppendCtxs(kvs...)
}

// AppendCtxs appends a set of the key-value contexts and returns itself.
func (e *Engine) AppendCtxs(kvs ...interface{}) *Engine {
	_len := len(kvs)
	if _len%2 != 0 {
		panic("the length of the key-value log contexts is not even")
	}
	for i := 0; i < _len; i += 2 {
		e.AppendCtx(kvs[i].(string), kvs[i+1])
	}
	return e
}

// AppendCtx appends a key-value context and returns itself.
func (e *Engine) AppendCtx(key string, value interface{}) *Engine {
	e.ctx = e.Output.encoder.Encode(e.ctx, key, value)
	e.origs = append(e.origs, kvctx{Key: key, Value: value})
	return e
}

// EncodeCtxs re-pre-encodes all the key-value contexts.
func (e *Engine) EncodeCtxs() {
	e.ctx = e.ctx[:0]
	for _, kv := range e.origs {
		e.ctx = e.Output.encoder.Encode(e.ctx, kv.Key, kv.Value)
	}
}

// Write implements the interface io.Writer.
func (e *Engine) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[:n-1]
	}
	e.getLogger(e.level, 1).Printf(string(p))
	return
}

// StdLog returns a new log.Logger based on the current logger engine.
func (e *Engine) StdLog(prefix string) *log.Logger {
	return log.New(e.Clone().SetDepth(2), prefix, 0)
}
