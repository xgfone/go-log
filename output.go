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
	"os"
	"path/filepath"

	jencoder "github.com/xgfone/go-log/encoder"
	"github.com/xgfone/go-log/writer"
)

// Output is used to handle the log output.
type Output struct {
	encoder encoderProxy
	writer  writer.LevelWriter
}

// NewOutput returns a new log output.
//
// If the encoder is nil, use JSONEncoder in the sub-package encoder by default.
func NewOutput(w io.Writer, encoder Encoder) *Output {
	if w == nil {
		panic("writer is nil")
	}
	if encoder == nil {
		encoder = jencoder.NewJSONEncoder(FormatLevel)
	}
	return &Output{
		encoder: newEncoder(encoder),
		writer:  writer.ToLevelWriter(w),
	}
}

func (o *Output) clone() *Output {
	return &Output{encoder: o.encoder, writer: o.writer}
}

// Writer is the alias of GetWriter.
func (o *Output) Writer() io.Writer {
	return o.writer
}

// GetWriter returns the log writer.
func (o *Output) GetWriter() io.Writer {
	return o.writer
}

// GetEncoder returns the log encoder.
func (o *Output) GetEncoder() Encoder {
	return o.encoder.Encoder
}

// SetWriter resets the log writer to w.
func (o *Output) SetWriter(w io.Writer) {
	if w == nil {
		panic("Output: the log writer is nil")
	}
	o.writer = writer.ToLevelWriter(w)
}

// SetEncoder resets the log encoder to enc.
func (o *Output) SetEncoder(enc Encoder) {
	if enc == nil {
		panic("Output: the log encoder is nil")
	}
	o.encoder = newEncoder(enc)
}

// WithEncoder returns a new logger with the new output created the new encoder
// and the original writer, which will also re-encode all the key-value contexts.
func (l Logger) WithEncoder(encoder Encoder) Logger {
	l.Output = l.Output.clone()
	l.Output.SetEncoder(encoder)

	ctxs := l.ctxs
	l.ctx, l.ctxs = nil, nil
	return l.WithContexts(ctxs...)
}

// WithWriter returns a new logger with the writer.
func (l Logger) WithWriter(writer io.Writer) Logger {
	l = l.Clone()
	l.Output = l.Output.clone()
	l.Output.SetWriter(writer)
	return l
}

/// ----------------------------------------------------------------------- ///

// FileWriter returns a writer based the file, which uses NewSizedRotatingFile
// to generate the file writer. If filename is "", however, it will return
// an os.Stderr writer instead.
//
// filesize is parsed by ParseSize to get the size of the log file.
// If it is "", it is "100M" by default.
//
// filenum is the number of the log file. If it is 0 or negative,
// it will be reset to 100.
//
// Notice: if the directory in where filename is does not exist, it will be
// created automatically.
func FileWriter(filename, filesize string, filenum int) io.WriteCloser {
	if filename == "" {
		return os.Stderr
	}

	if filesize == "" {
		filesize = "100M"
	}

	size, err := writer.ParseSize(filesize)
	if err != nil {
		panic(err)
	} else if filenum <= 0 {
		filenum = 100
	}

	if err = os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		panic(err)
	}

	return writer.NewSizedRotatingFile(filename, int(size), filenum)
}
