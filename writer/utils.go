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
	"fmt"
	"strconv"
)

// ParseSize parses the size string. The size maybe have a unit suffix,
// such as "123", "123M, 123G". Valid size units are "b", "B", "k", "K",
// "m", "M", "g", "G", "t", "T", "p", "P", "e", "E". The lower units are 1000x,
// and the upper units are 1024x.
//
// Notice: "" will be considered as 0.
func ParseSize(s string) (size int64, err error) {
	if s == "" {
		return
	}

	var base int64
	switch _len := len(s) - 1; s[_len] {
	case 'b', 'B':
		s = s[:_len]
	case 'k':
		base = 1000
		s = s[:_len]
	case 'K':
		base = 1024
		s = s[:_len]
	case 'm':
		base = 1000000 // 1000**2
		s = s[:_len]
	case 'M':
		base = 1048576 // 1024**2
		s = s[:_len]
	case 'g':
		base = 1000000000 // 1000**3
		s = s[:_len]
	case 'G':
		base = 1073741824 // 1024**3
		s = s[:_len]
	case 't':
		base = 1000000000000 // 1000**4
		s = s[:_len]
	case 'T':
		base = 1099511627776 // 1024**4
		s = s[:_len]
	case 'p':
		base = 1000000000000000 // 1000**5
		s = s[:_len]
	case 'P':
		base = 1125899906842624 // 1024**5
		s = s[:_len]
	case 'e':
		base = 1000000000000000000 // 1000**6
		s = s[:_len]
	case 'E':
		base = 1152921504606846976 // 1024**6
		s = s[:_len]
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
	default:
		return 0, fmt.Errorf("unknown size string '%s'", s)
	}

	if size, err = strconv.ParseInt(s, 10, 64); err == nil && base > 1 {
		size *= base
	}

	return
}
