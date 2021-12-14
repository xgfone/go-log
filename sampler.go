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
	"strings"
	"sync"
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

// GetSampler returns the sampler.
//
// If no sampler is set, return nil.
func (e *Engine) GetSampler() Sampler { return e.sampler }

// SetSampler resets the sampler and returns itself, which is not thread-safe.
// For thread-safe, SwitchSampler may be used.
//
// If the sampler is nil, it will cancel the sampler.
func (e *Engine) SetSampler(sampler Sampler) *Engine {
	e.sampler = sampler
	return e
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

/// ----------------------------------------------------------------------- ///

var _ Sampler = &SimpleSampler{}

// SimpleSampler is a simple sampler.
//
// For the name, it supports not only the exact match but also the prefix match
// like "prefix1.prefix2.*".
type SimpleSampler struct {
	lock  sync.RWMutex
	names map[string]int
	value atomic.Value
	level int64
}

// NewSimpleSampler returns a new SimpleSampler with the default threshold level.
func NewSimpleSampler(defaultThresholdLevel int) *SimpleSampler {
	checkLevel(defaultThresholdLevel)
	s := &SimpleSampler{
		level: int64(defaultThresholdLevel),
		names: make(map[string]int),
	}
	s.value.Store(map[string]int{})
	return s
}

// Sample implements the interface Sampler.
func (s *SimpleSampler) Sample(name string, level int) bool {
	names := s.value.Load().(map[string]int)
	if len(names) > 0 {
		for lname, minLevel := range names {
			if nlen := len(lname); nlen > 0 && lname[nlen-1] == '*' {
				if strings.HasPrefix(name, lname[:nlen-1]) {
					return allowLevel(level, minLevel)
				}
			} else if lname == name {
				return allowLevel(level, minLevel)
			}
		}
	}

	return allowLevel(level, s.GetDefaultLevel())
}

// GetDefaultLevel returns the default threshold level.
func (s *SimpleSampler) GetDefaultLevel() (level int) {
	return int(atomic.LoadInt64(&s.level))
}

// SetDefaultLevel resets the default threshold level.
func (s *SimpleSampler) SetDefaultLevel(level int) {
	checkLevel(level)
	atomic.StoreInt64(&s.level, int64(level))
}

// GetNamedLevels returns all the named levels.
func (s *SimpleSampler) GetNamedLevels() map[string]int {
	s.lock.RLock()
	names := make(map[string]int, len(s.names))
	for name, level := range s.names {
		names[name] = level
	}
	s.lock.RUnlock()
	return names
}

// ResetNamedLevels resets the named levels.
//
// If no named levels is set, allow all the logs to be sampled.
//
// Notice: for the invalid levels, they are ignored.
func (s *SimpleSampler) ResetNamedLevels(names map[string]int) {
	s.lock.Lock()
	s.names = make(map[string]int, len(names))
	for name, level := range names {
		if LevelIsValid(level) {
			s.names[name] = level
		}
	}
	s.updateNames()
	s.lock.Unlock()
}

// AddNamedLevel adds the named level.
func (s *SimpleSampler) AddNamedLevel(name string, level int) {
	checkLevel(level)
	s.lock.Lock()
	if _, ok := s.names[name]; !ok {
		s.names[name] = level
		s.updateNames()
	}
	s.lock.Unlock()
}

// DelName deletes the named level by the name.
func (s *SimpleSampler) DelName(name string) {
	s.lock.Lock()
	if _, ok := s.names[name]; ok {
		delete(s.names, name)
		s.updateNames()
	}
	s.lock.Unlock()
}

func (s *SimpleSampler) updateNames() {
	names := make(map[string]int, len(s.names))
	for name, level := range s.names {
		names[name] = level
	}
	s.value.Store(names)
}
