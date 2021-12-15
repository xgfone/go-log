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

//go:build go1.16
// +build go1.16

package writer

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSizedRotatingFile(t *testing.T) {
	const filename = "test_file_writer.log"
	size, err := ParseSize("15")
	if err != nil {
		t.Fatal(err)
	}

	data := []byte("0123456789")
	logfiles := make(map[string]int64)

	file := SafeWriter(NewSizedRotatingFile(filename, int(size), 3))
	defer func() {
		Close(file)
		for name := range logfiles {
			os.Remove(name)
		}
	}()

	for i := 0; i < 10; i++ {
		n, err := file.Write(data)
		if err != nil {
			t.Error(err)
		} else if _len := len(data); n != _len {
			t.Errorf("expect write %d bytes, but only wrote %d bytes", _len, n)
		}
	}

	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if name := info.Name(); strings.HasPrefix(name, filename) {
			logfiles[name] = info.Size()
		}
		return nil
	})

	if len(logfiles) != 4 {
		t.Errorf("expect %d log files, but got %d", 4, len(logfiles))
	} else {
		for name, size := range logfiles {
			switch name {
			case filename, filename + ".1", filename + ".2", filename + ".3":
			default:
				t.Errorf("unexpeced log filename '%s'", name)
			}

			if size != 10 {
				t.Errorf("expect log file size %d, gut got %d", 10, size)
			}
		}
	}
}