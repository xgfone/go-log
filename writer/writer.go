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

// Package writer provides some log writers.
package writer

import (
	"bufio"
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
		return Close(w.UnwrapWriter())

	default:
		return nil
	}
}

// Flush flushes the writer if it has implemented the interface Flusher.
func Flush(writer io.Writer) (err error) {
	switch w := writer.(type) {
	case Flusher:
		return w.Flush()

	case WrappedWriter:
		return Flush(w.UnwrapWriter())

	default:
		return nil
	}
}

// Flusher is used to flush the data in the writer to the underlying storage media.
type Flusher interface {
	Flush() error
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
	writer LevelWriter
	lock   sync.Mutex
}

func (w *safeWriter) Close() (err error) {
	w.lock.Lock()
	err = Close(w.writer)
	w.lock.Unlock()
	return
}

func (w *safeWriter) Flush() (err error) {
	w.lock.Lock()
	err = Flush(w.writer)
	w.lock.Unlock()
	return
}

func (w *safeWriter) Write(p []byte) (n int, err error) {
	w.lock.Lock()
	n, err = w.writer.Write(p)
	w.lock.Unlock()
	return
}

func (w *safeWriter) WriteLevel(level int, p []byte) (n int, err error) {
	w.lock.Lock()
	n, err = w.writer.WriteLevel(level, p)
	w.lock.Unlock()
	return
}

func (w *safeWriter) UnwrapWriter() io.Writer { return w.writer }

// SafeWriter is guaranteed that only a single writing operation can proceed
// at a time, which implements the interface LevelWriter, WrappedWriter and
// Flusher.
//
// It's necessary for thread-safe concurrent writes.
func SafeWriter(writer io.Writer) io.WriteCloser {
	if writer == nil {
		panic("SafeWriter: the wrapped writer is nil")
	}
	return &safeWriter{writer: ToLevelWriter(writer)}
}

/// ----------------------------------------------------------------------- ///

type bufWriter struct {
	*bufio.Writer
	w io.Writer
}

func (w bufWriter) UnwrapWriter() io.Writer     { return w.w }
func (w bufWriter) Write(p []byte) (int, error) { return w.Writer.Write(p) }
func (w bufWriter) Close() error                { w.Flush(); return Close(w.w) }

// BufferWriter returns a buffer writer, which implements the interfaces
// Flusher and WrappedWriter, writes the data into the buffer and flushes
// all the datas into the wrapped writer when the buffer is full.
//
// If bufSize is equal to or less than 0, it is 4096 by default.
func BufferWriter(writer io.Writer, bufSize int) io.WriteCloser {
	if writer == nil {
		panic("BufferWriter: the wrapped writer is nil")
	}
	return bufWriter{w: writer, Writer: bufio.NewWriterSize(writer, bufSize)}
}
