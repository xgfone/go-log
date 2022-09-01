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
	"strings"
	"testing"
)

type lwriter struct{ io.Writer }

func (lw lwriter) UnwrapWriter() io.Writer                     { return lw.Writer }
func (lw lwriter) WriteLevel(_ int, _ []byte) (n int, e error) { return }
func (lw lwriter) Flush() (e error)                            { return }

func TestClose(t *testing.T) {
	Close(lwriter{bytes.NewBuffer(nil)})
}

func TestLevelWriterFlush(t *testing.T) {
	lw := ToLevelWriter(lwriter{Discard})
	if _, ok := lw.(Flusher); !ok {
		t.Error("expect a Flusher writer, but not")
	}
}

func TestFlush(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	w := SafeWriter(BufferWriter(buf, 0))

	io.WriteString(buf, "------")
	io.WriteString(w, "111111")
	io.WriteString(buf, "++++++")
	Flush(w)
	io.WriteString(buf, "@@@@@@")
	io.WriteString(w, "222222")
	io.WriteString(buf, "######")
	Close(w)
	io.WriteString(buf, "$$$$$$")

	expect := strings.Join([]string{
		"------",
		"++++++",
		"111111",
		"@@@@@@",
		"######",
		"222222",
		"$$$$$$",
	}, "")
	if s := buf.String(); s != expect {
		t.Errorf("expect '%s', but got '%s'", expect, s)
	}
}
