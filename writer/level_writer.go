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

package writer

import "io"

// LevelWriter is a writer with the level.
type LevelWriter interface {
	WriteLevel(level int, data []byte) (n int, err error)
	io.Writer
}

// ToLevelWriter converts the io.Writer to LevelWriter.
func ToLevelWriter(writer io.Writer) LevelWriter {
	if lw, ok := writer.(LevelWriter); ok {
		return lw
	}
	return lvlWriter{Writer: writer}
}

type lvlWriter struct{ io.Writer }

func (lw lvlWriter) UnwrapWriter() io.Writer                 { return lw.Writer }
func (lw lvlWriter) WriteLevel(l int, p []byte) (int, error) { return lw.Write(p) }
func (lw lvlWriter) Close() (err error)                      { return Close(lw.Writer) }

/// ----------------------------------------------------------------------- ///

// LevelSplitWriter returns a writer to write the log into the different writer
// by the level.
func LevelSplitWriter(defaultWriter io.Writer, levelWriters map[int]io.Writer) LevelWriter {
	lws := make(map[int]LevelWriter, len(levelWriters))
	for level, lw := range levelWriters {
		lws[level] = ToLevelWriter(lw)
	}
	return lvlSplitWriter{dw: ToLevelWriter(defaultWriter), lws: lws}
}

type lvlSplitWriter struct {
	lws map[int]LevelWriter
	dw  LevelWriter
}

func (w lvlSplitWriter) Write(p []byte) (int, error) { return w.dw.Write(p) }

func (w lvlSplitWriter) WriteLevel(level int, p []byte) (int, error) {
	if lw, ok := w.lws[level]; ok {
		return lw.WriteLevel(level, p)
	}
	return w.dw.WriteLevel(level, p)
}

type werrors []error

func (es werrors) Errors() []error { return es }
func (es werrors) Error() string {
	buf := make([]byte, 0, 128)
	for i, _len := 0, len(es); i < _len; i++ {
		buf = append(buf, es[i].Error()...)
	}
	return string(buf)
}

func (w lvlSplitWriter) Close() (err error) {
	var errors werrors
	if err := Close(w.dw); err != nil {
		errors = append(errors, err)
	}
	for _, lw := range w.lws {
		if err := Close(lw); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
