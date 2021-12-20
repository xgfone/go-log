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

// Encoder is used to encode the log record.
type Encoder interface {
	// Start starts to encode the log record into the buffer dst.
	Start(dst []byte, loggerName string, level int) []byte

	// Encode encodes the key-value with the stack depth into the buffer dst.
	Encode(dst []byte, key string, value interface{}) []byte

	// End ends to encode the log record with the message into the buffer dst.
	End(dst []byte, msg string) []byte
}
