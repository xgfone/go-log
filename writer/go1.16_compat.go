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

//go:build !go1.16
// +build !go1.16

package writer

import (
	"os"
	"path/filepath"
	"strings"
)

func init() { Discard = discard{} }

type discard struct{}

func (discard) Write(p []byte) (int, error)       { return len(p), nil }
func (discard) WriteString(s string) (int, error) { return len(s), nil }

func listdir(dir, prefix string) (files map[string]int64) {
	files = make(map[string]int64)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if name := info.Name(); strings.HasPrefix(name, prefix) {
			files[name] = info.Size()
		}
		return nil
	})
	return
}
