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

	"github.com/xgfone/go-log/writer"
)

// Output is used to handle the log output.
type Output struct {
	encoder Encoder
	writer  LevelWriter
}

// NewOutput returns a new log output.
//
// If the encoder is nil, use NewJSONEncoder() by default.
func NewOutput(writer io.Writer, encoder Encoder) *Output {
	if writer == nil {
		panic("writer is nil")
	}
	if encoder == nil {
		encoder = NewJSONEncoder()
	}
	return &Output{encoder: encoder, writer: ToLevelWriter(writer)}
}

// GetWriter returns the log writer.
func (o *Output) GetWriter() io.Writer { return o.writer }

// SetWriter resets the log writer to w.
func (o *Output) SetWriter(w io.Writer) { o.writer = ToLevelWriter(w) }

// GetEncoder returns the log encoder.
func (o *Output) GetEncoder() Encoder { return o.encoder }

// SetEncoder resets the log encoder to enc.
func (o *Output) SetEncoder(enc Encoder) { o.encoder = enc }

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
func (lw lvlWriter) Close() (err error)                      { return writer.Close(lw.Writer) }

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
