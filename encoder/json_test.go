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

package encoder

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestJSONEncoder(t *testing.T) {
	var buf []byte
	enc := NewJSONEncoder()
	enc.TimeKey = ""

	buf = enc.Start(buf, "", "")
	buf = enc.EncodeInt(buf, "k1", 111)
	buf = enc.EncodeInt64(buf, "k2", 222)
	buf = enc.EncodeUint(buf, "k3", 333)
	buf = enc.EncodeUint64(buf, "k4", 444)
	buf = enc.EncodeFloat64(buf, "k5", 555)
	buf = enc.EncodeBool(buf, "k6", true)
	buf = enc.EncodeString(buf, "k7", "test")
	buf = enc.EncodeTime(buf, "k8", time.Date(2021, time.December, 20, 22, 40, 0, 0, time.UTC))
	buf = enc.EncodeDuration(buf, "k9", time.Second*10)

	buf = enc.Encode(buf, "nil", nil)
	buf = enc.Encode(buf, "bool", true)
	buf = enc.Encode(buf, "int", 10)
	buf = enc.Encode(buf, "int8", 11)
	buf = enc.Encode(buf, "int16", 12)
	buf = enc.Encode(buf, "int32", 13)
	buf = enc.Encode(buf, "int64", 14)
	buf = enc.Encode(buf, "uint", 15)
	buf = enc.Encode(buf, "uint8", 16)
	buf = enc.Encode(buf, "uint16", 17)
	buf = enc.Encode(buf, "uint32", 18)
	buf = enc.Encode(buf, "uint64", 19)
	buf = enc.Encode(buf, "float32", 20)
	buf = enc.Encode(buf, "float64", 21)
	buf = enc.Encode(buf, "string", "22")
	buf = enc.Encode(buf, "error", errors.New("23"))
	buf = enc.Encode(buf, "duration", time.Duration(24)*time.Second)
	buf = enc.Encode(buf, "time", time.Date(2021, time.May, 25, 22, 52, 26, 0, time.UTC))
	buf = enc.Encode(buf, "[]interface{}", []interface{}{"26", "27"})
	buf = enc.EncodeStringSlice(buf, "[]string", []string{"28", "29"})
	buf = enc.Encode(buf, "[]uint", []uint{30, 31})
	buf = enc.Encode(buf, "[]uint64", []uint64{32, 33})
	buf = enc.Encode(buf, "[]int", []int{34, 35})
	buf = enc.Encode(buf, "[]int64", []int64{36, 37})
	buf = enc.Encode(buf, "map[string]interface{}", map[string]interface{}{"a": "38", "b": "39"})
	buf = enc.Encode(buf, "map[string]string", map[string]string{"c": "40", "d": "41"})
	buf = enc.Encode(buf, "bytes1", []byte(`"42"`))
	buf = enc.Encode(buf, "bytes2", []byte(`[43, 44]`))
	buf = enc.End(buf, `"test json encoder"`)

	type encoderT struct {
		Msg string `json:"msg"`

		K1 int     `json:"k1"`
		K2 int64   `json:"k2"`
		K3 uint    `json:"k3"`
		K4 uint64  `json:"k4"`
		K5 float64 `json:"k5"`
		K6 bool    `json:"k6"`
		K7 string  `json:"k7"`
		K8 string  `json:"k8"`
		K9 string  `json:"k9"`

		Nil      interface{} `json:"nil"`
		Bool     bool        `json:"bool"`
		Int      int         `json:"int"`
		Int8     int8        `json:"int8"`
		Int16    int16       `json:"int16"`
		Int32    int32       `json:"int32"`
		Int64    int64       `json:"int64"`
		Uint     uint        `json:"uint"`
		Uint8    uint8       `json:"uint8"`
		Uint16   uint16      `json:"uint16"`
		Uint32   uint32      `json:"uint32"`
		Uint64   uint64      `json:"uint64"`
		Float32  float32     `json:"float32"`
		Float64  float64     `json:"float64"`
		String   string      `json:"string"`
		Error    string      `json:"error"`
		Duration string      `json:"duration"`
		Time     string      `json:"time"`

		Interfaces []interface{} `json:"[]interface{}"`
		Strings    []string      `json:"[]string"`
		Uints      []uint        `json:"[]uint"`
		Uint64s    []uint64      `json:"[]uint64"`
		Ints       []int         `json:"[]int"`
		Int64s     []int64       `json:"[]int64"`

		MapString    map[string]string      `json:"map[string]string"`
		MapInterface map[string]interface{} `json:"map[string]interface{}"`

		Bytes1 string `json:"bytes1"`
		Bytes2 []int  `json:"bytes2"`
	}

	expect := encoderT{
		Msg: `"test json encoder"`,

		K1: 111,
		K2: 222,
		K3: 333,
		K4: 444,
		K5: 555,
		K6: true,
		K7: "test",
		K8: "2021-12-20T22:40:00Z",
		K9: "10s",

		Nil:      nil,
		Bool:     true,
		Int:      10,
		Int8:     11,
		Int16:    12,
		Int32:    13,
		Int64:    14,
		Uint:     15,
		Uint8:    16,
		Uint16:   17,
		Uint32:   18,
		Uint64:   19,
		Float32:  20,
		Float64:  21,
		String:   "22",
		Error:    "23",
		Duration: "24s",
		Time:     "2021-05-25T22:52:26Z",

		Interfaces: []interface{}{"26", "27"},
		Strings:    []string{"28", "29"},
		Uints:      []uint{30, 31},
		Uint64s:    []uint64{32, 33},
		Ints:       []int{34, 35},
		Int64s:     []int64{36, 37},

		MapInterface: map[string]interface{}{"a": "38", "b": "39"},
		MapString:    map[string]string{"c": "40", "d": "41"},

		Bytes1: "42",
		Bytes2: []int{43, 44},
	}

	var result encoderT
	if err := json.Unmarshal(buf, &result); err != nil {
		t.Error(err, string(buf))
	} else if !reflect.DeepEqual(result, expect) {
		t.Errorf("expect '%+v', but got '%+v'", expect, result)
	}
}
