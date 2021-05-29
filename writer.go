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
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

// WriterLevel is an interface to write the data with the level.
type WriterLevel interface {
	WriteLevel(level Level, data []byte) (n int, err error)
}

// Writer is the interface to write the log to the underlying storage.
type Writer interface {
	WriterLevel
	io.Closer
}

// WriterFunc is a function writer.
type WriterFunc func(level Level, data []byte) (n int, err error)

// WriteLevel implements the interface Writer.
func (w WriterFunc) WriteLevel(l Level, p []byte) (int, error) { return w(l, p) }

// Close implements the interface io.Closer.
func (w WriterFunc) Close() error { return nil }

// IOWriter is the writer implementing the interface io.Writer.
type IOWriter struct {
	Level Level
	Writer
}

// NewIOWriter returns a new IOWriter.
func NewIOWriter(w Writer, l Level) IOWriter { return IOWriter{Level: l, Writer: w} }

// Write implements the interface io.Writer.
func (w IOWriter) Write(p []byte) (int, error) { return w.WriteLevel(w.Level, p) }

type streamWriter struct {
	writer WriterLevel
	closer io.Closer
	io.Writer
}

func (w streamWriter) WriteLevel(l Level, p []byte) (int, error) {
	if w.writer != nil {
		return w.writer.WriteLevel(l, p)
	}
	return w.Writer.Write(p)
}

func (w streamWriter) Close() (err error) {
	if w.closer != nil {
		err = w.closer.Close()
	}
	return
}

// StreamWriter converts io.Writer to Writer.
func StreamWriter(w io.Writer) Writer {
	sw := streamWriter{Writer: w}

	if lw, ok := w.(WriterLevel); ok {
		sw.writer = lw
	}

	if closer, ok := w.(io.Closer); ok {
		sw.closer = closer
	}

	return sw
}

// DiscardWriter discards all the data.
func DiscardWriter() Writer {
	return WriterFunc(func(l Level, p []byte) (int, error) { return len(p), nil })
}

type levelWriter struct {
	Level Level
	Writer
}

func (w levelWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	if l >= w.Level {
		n, err = w.WriteLevel(l, p)
	} else {
		n = len(p)
	}
	return
}

// LevelWriter filters the logs whose level is less than lvl.
func LevelWriter(level Level, writer Writer) Writer {
	return levelWriter{Level: level, Writer: writer}
}

type safeWriter struct {
	lock sync.Mutex
	Writer
}

func (w *safeWriter) Close() error {
	w.lock.Lock()
	err := w.Writer.Close()
	w.lock.Unlock()
	return err
}

func (w *safeWriter) WriteLevel(l Level, p []byte) (n int, err error) {
	w.lock.Lock()
	n, err = w.Writer.WriteLevel(l, p)
	w.lock.Unlock()
	return
}

// SafeWriter is guaranteed that only a single writing operation can proceed
// at a time.
//
// It's necessary for thread-safe concurrent writes.
func SafeWriter(writer Writer) Writer { return &safeWriter{Writer: writer} }

type splitWriter struct {
	writers map[Level]Writer
	twriter Writer //  LvlTrace
	dwriter Writer //  LvlDebug
	iwriter Writer //  LvlInfo
	wwriter Writer //  LvlWarn
	ewriter Writer //  LvlError
	fwriter Writer //  LvlFatal
}

func (w splitWriter) WriteLevel(level Level, data []byte) (n int, err error) {
	var writer Writer
	switch level {
	case LvlTrace:
		writer = w.twriter
	case LvlDebug:
		writer = w.dwriter
	case LvlInfo:
		writer = w.iwriter
	case LvlWarn:
		writer = w.wwriter
	case LvlError:
		writer = w.ewriter
	case LvlFatal:
		writer = w.fwriter
	}

	if writer == nil {
		writer = w.writers[level]
	}

	if writer != nil {
		n, err = writer.WriteLevel(level, data)
	} else {
		n = len(data)
	}

	return
}

func (w splitWriter) Close() error {
	for _, writer := range w.writers {
		writer.Close()
	}
	return nil
}

// SplitWriter returns a level-separated writer, which will write the log record
// into the separated writer.
func SplitWriter(writers map[Level]Writer) Writer {
	return splitWriter{
		writers: writers,
		twriter: writers[LvlTrace],
		dwriter: writers[LvlDebug],
		iwriter: writers[LvlInfo],
		wwriter: writers[LvlWarn],
		ewriter: writers[LvlError],
		fwriter: writers[LvlFatal],
	}
}

// FileWriter returns a writer based the file, which uses NewSizedRotatingFile
// to generate the file writer. If filename is "", however, it will return
// an os.Stdout writer instead.
//
// filesize is parsed by ParseSize to get the size of the log file.
// If it is "", it is "100M" by default.
//
// filenum is the number of the log file. If it is 0 or negative,
// it will be reset to 100.
//
// Notice: if the directory in where filename is does not exist, it will be
// created automatically.
func FileWriter(filename, filesize string, filenum int) Writer {
	var w io.WriteCloser = os.Stdout
	if filename != "" {
		if filesize == "" {
			filesize = "100M"
		}

		size, err := ParseSize(filesize)
		if err != nil {
			panic(err)
		} else if filenum <= 0 {
			filenum = 100
		}

		if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			panic(err)
		}

		w = NewSizedRotatingFile(filename, int(size), filenum)
	}

	return StreamWriter(w)
}

// NewSizedRotatingFile returns a new SizedRotatingFile, which is not thread-safe.
//
// Default:
//   fileperm: 0644
//   filesize: 100 * 1024 * 1024
//   filenum:  0
//
func NewSizedRotatingFile(filename string, filesize, filenum int,
	fileperm ...os.FileMode) *SizedRotatingFile {
	var filemode os.FileMode = 0644
	if len(fileperm) > 0 && fileperm[0] > 0 {
		filemode = fileperm[0]
	}

	if filenum <= 0 {
		filesize = int(math.MaxInt32)
	} else if filesize <= 0 {
		filesize = 100 * 1024 * 1024
	}

	return &SizedRotatingFile{
		filename:    filename,
		filemode:    filemode,
		maxSize:     filesize,
		backupCount: filenum,
	}
}

// SizedRotatingFile is a file rotating logging writer based on the size.
type SizedRotatingFile struct {
	file        *os.File
	filemode    os.FileMode
	filename    string
	maxSize     int
	backupCount int
	nbytes      int
	closed      int32
}

// Close implements io.Closer.
func (f *SizedRotatingFile) Close() (err error) {
	if atomic.CompareAndSwapInt32(&f.closed, 0, 1) {
		err = f.close()
	}
	return
}

// Flush flushes the data to the underlying disk.
func (f *SizedRotatingFile) Flush() (err error) {
	if f.file != nil {
		err = f.file.Sync()
	}
	return
}

// Write implements io.Writer.
func (f *SizedRotatingFile) Write(data []byte) (n int, err error) {
	if atomic.LoadInt32(&f.closed) == 1 {
		return 0, errors.New("the file has been closed")
	}

	if f.file == nil {
		if err = f.open(); err != nil {
			return
		}
	}

	if f.nbytes+len(data) > f.maxSize {
		if err = f.doRollover(); err != nil {
			return
		}
	}

	if n, err = f.file.Write(data); err != nil {
		return
	}

	f.nbytes += n
	return
}

func (f *SizedRotatingFile) open() (err error) {
	file, err := os.OpenFile(f.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, f.filemode)
	if err != nil {
		return
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return
	}

	f.nbytes = int(info.Size())
	f.file = file
	return
}

func (f *SizedRotatingFile) close() (err error) {
	if f.file != nil {
		err = f.file.Close()
		f.file = nil
	}
	return
}

func (f *SizedRotatingFile) doRollover() (err error) {
	if f.backupCount > 0 {
		if err = f.close(); err != nil {
			return fmt.Errorf("failed to close the rotating file '%s': %s", f.filename, err)
		}

		if !fileIsExist(f.filename) {
			return nil
		} else if n, err := fileSize(f.filename); err != nil {
			return fmt.Errorf("failed to get the size of the rotating file '%s': %s",
				f.filename, err)
		} else if n == 0 {
			return nil
		}

		for _, i := range ranges(f.backupCount-1, 0, -1) {
			sfn := fmt.Sprintf("%s.%d", f.filename, i)
			dfn := fmt.Sprintf("%s.%d", f.filename, i+1)
			if fileIsExist(sfn) {
				if fileIsExist(dfn) {
					os.Remove(dfn)
				}
				if err = os.Rename(sfn, dfn); err != nil {
					return fmt.Errorf("failed to rename the rotating file '%s' to '%s': %s",
						sfn, dfn, err)
				}
			}
		}

		dfn := f.filename + ".1"
		if fileIsExist(dfn) {
			if err = os.Remove(dfn); err != nil {
				return fmt.Errorf("failed to remove the rotating file '%s': %s", dfn, err)
			}
		}
		if fileIsExist(f.filename) {
			if err = os.Rename(f.filename, dfn); err != nil {
				return fmt.Errorf("failed to rename the rotating file '%s' to '%s': %s",
					f.filename, dfn, err)
			}
		}

		err = f.open()
	}

	return
}

func fileIsExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// fileSize returns the size of the file as how many bytes.
func fileSize(fp string) (int64, error) {
	f, e := os.Stat(fp)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

func ranges(start, stop, step int) (r []int) {
	if step > 0 {
		for start < stop {
			r = append(r, start)
			start += step
		}
		return
	} else if step < 0 {
		for start > stop {
			r = append(r, start)
			start += step
		}
		return
	}

	panic(fmt.Errorf("step must not be 0"))
}
