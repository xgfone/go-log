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

// Package writer provides some io writers.
package writer

import (
	"io"
	"sync"
)

// Discard is the writer to discard all the written data.
//
// For Go1.16+, it is equal to io.Discard. Or, it's an internal implementation.
var Discard io.Writer

// Close closes the writer if it has implemented the interface io.Closer.
func Close(writer io.Writer) (err error) {
	switch w := writer.(type) {
	case io.Closer:
		return w.Close()

	case WrappedWriter:
		return Close(w)

	default:
		return nil
	}
}

/// ----------------------------------------------------------------------- ///

// WrappedWriter is a writer which wraps and returns the inner writer.
type WrappedWriter interface {
	UnwrapWriter() io.Writer
	io.Writer
}

// UnwrapWriter recursively unwraps the wrapped innest writer from
// the given writer if it has implemented the interface WrappedWriter.
// Or return the original writer.
func UnwrapWriter(writer io.Writer) io.Writer {
	for {
		if w, ok := writer.(WrappedWriter); ok {
			writer = w.UnwrapWriter()
		} else {
			break
		}
	}
	return writer
}

/// ----------------------------------------------------------------------- ///

type safeWriter struct {
	lock sync.Mutex
	io.Writer
}

func (w *safeWriter) Close() (err error) {
	w.lock.Lock()
	err = Close(w.Writer)
	w.lock.Unlock()
	return
}

func (w *safeWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	n, err = w.Writer.Write(p)
	w.lock.Unlock()
	return
}

func (w *safeWriter) UnwrapWriter() io.Writer { return w.Writer }

// SafeWriter is guaranteed that only a single writing operation can proceed
// at a time.
//
// It's necessary for thread-safe concurrent writes.
func SafeWriter(w io.Writer) io.WriteCloser { return &safeWriter{Writer: w} }
