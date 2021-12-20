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
	"sync/atomic"
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

// Sampler returns the sampler.
//
// If no sampler is set, return nil.
func (l Logger) Sampler() Sampler { return l.sampler }

// WithSampler returns a new logger with the sampler.
//
// If the sampler is nil, it will cancel the sampler.
func (l Logger) WithSampler(sampler Sampler) Logger {
	l = l.Clone()
	l.sampler = sampler
	return l
}

/// ----------------------------------------------------------------------- ///

var _ Sampler = &SwitchSampler{}

type samplerWrapper struct{ Sampler }

// SwitchSampler is a sampler to switch the proxy sampler.
type SwitchSampler struct {
	sampler atomic.Value
}

// NewSwitchSampler returns a new SwitchSampler with the wrapped sampler.
func NewSwitchSampler(sampler Sampler) *SwitchSampler {
	s := &SwitchSampler{}
	s.Set(sampler)
	return s
}

// Get returns the sampler.
func (s *SwitchSampler) Get() Sampler {
	return s.sampler.Load().(samplerWrapper).Sampler
}

// Set resets the sampler.
func (s *SwitchSampler) Set(sampler Sampler) {
	if sampler == nil {
		panic("SwitchSampler: sampler is nil")
	}
	s.sampler.Store(samplerWrapper{sampler})
}

// Sample implements the interface Sampler.
func (s *SwitchSampler) Sample(name string, level int) bool {
	return s.Get().Sample(name, level)
}
