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

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestMultiLevelWriter(t *testing.T) {
	infolog := "info log data content"
	errorlog := "error log data content"

	const errfilename = "test.error.log"
	const logfilename = "test.log"
	const logfilesize = "100M"
	const logfilenum = 100

	defer os.Remove(errfilename)
	defer os.Remove(logfilename)

	defaultLogWriter := NewSizedRotatingFile(logfilename, 1024*1024, logfilenum)
	levelWriters := map[int]io.Writer{
		80: NewSizedRotatingFile(errfilename, 1024*1024, logfilenum),
	}

	mwriter := LevelSplitWriter(defaultLogWriter, levelWriters)
	mwriter.WriteLevel(40, []byte(infolog))
	mwriter.WriteLevel(80, []byte(errorlog))
	Close(mwriter)

	if logdata, err := ioutil.ReadFile(logfilename); err != nil {
		t.Error(err)
	} else if s := string(logdata); s != infolog {
		t.Errorf("expect log '%s', but got '%s'", infolog, s)
	}

	if logdata, err := ioutil.ReadFile(errfilename); err != nil {
		t.Error(err)
	} else if s := string(logdata); s != errorlog {
		t.Errorf("expect log '%s', but got '%s'", infolog, s)
	}
}

func TestClose(t *testing.T) {
	Close(lwriter{bytes.NewBuffer(nil)})
}

type lwriter struct{ io.Writer }

func (lw lwriter) UnwrapWriter() io.Writer { return lw.Writer }
