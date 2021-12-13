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
	"errors"
	"fmt"
	"math"
	"os"
	"sync/atomic"
)

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
