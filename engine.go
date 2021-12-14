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
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/xgfone/go-log/writer"
)

var globalSampling = int32(1)

func globalSamplingIsEnabled() bool {
	return atomic.LoadInt32(&globalSampling) == 1
}

// GlobalDisableSampling is used to disable all the samplings globally.
func GlobalDisableSampling(disable bool) {
	if disable {
		atomic.StoreInt32(&globalSampling, 0)
	} else {
		atomic.StoreInt32(&globalSampling, 1)
	}
}

// CallerFormatFunc is used to format the line and line of the caller.
var CallerFormatFunc = func(file string, line int) string {
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Caller returns a callback function that returns the caller "file:line".
func Caller(key string) Hook {
	return HookFunc(func(logger Logger, name string, level, depth int) {
		if _, file, line, ok := runtime.Caller(depth + 1); ok {
			logger.Kv(key, CallerFormatFunc(file, line))
		}
	})
}

// Hook is used to add the dynamic value into the log record.
type Hook interface {
	Run(logger Logger, loggerName string, level int, depth int)
}

// HookFunc is a function hook.
type HookFunc func(logger Logger, name string, level int, depth int)

// Run implements the interface Hook.
func (f HookFunc) Run(logger Logger, name string, level int, depth int) {
	f(logger, name, level, depth+1)
}

// Sampler is used to sample the log message.
type Sampler interface {
	// Sample reports whether the log message should be sampled.
	// If the log message should be sampled, return true. Or, return false,
	// that's, the log message will be discarded.
	Sample(loggerName string, level int) bool
}

// SamplerFunc is a function sampler.
type SamplerFunc func(loggerName string, level int) bool

// Sample implements the interface Sampler.
func (f SamplerFunc) Sample(name string, lvl int) bool { return f(name, lvl) }

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
		name:   e.name,
		level:  e.level,
		depth:  e.depth,
		Output: e.Output,

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

// SetSampler resets the sampler and returns itself.
//
// If the sampler is nil, it will cancel the sampler.
func (e *Engine) SetSampler(sampler Sampler) *Engine {
	e.sampler = sampler
	return e
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

// AddHooks appends the hooks and returns itself.
func (e *Engine) AddHooks(hooks ...Hook) *Engine {
	e.hooks = append(e.hooks, hooks...)
	return e
}

// ResetHooks resets the hooks and returns itself.
func (e *Engine) ResetHooks(hooks ...Hook) *Engine {
	e.hooks = append([]Hook{}, hooks...)
	return e
}

// ResetCtxs resets the key-value contexts and returns itself.
func (e *Engine) ResetCtxs() *Engine {
	e.ctx, e.origs = nil, nil
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
	e.ctx = nil
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
