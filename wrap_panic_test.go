// Copyright 2022 xgfone
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

//go:build go1.12
// +build go1.12

package log

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func handleBusiness() {
	func(s string) {
		panic(s)
	}("test")
}

func testHandleBusiness() {
	defer WrapPanic()
	handleBusiness()
}

func TestWrapPanic(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	DefaultLogger = New("").WithWriter(buf).WithHooks(Caller("caller"))

	testHandleBusiness()

	results := strings.Split(buf.String(), "\n")
	if len(results) != 2 {
		t.Errorf("expect %d line logs, but got %d", 2, len(results))
		return
	}

	var r struct {
		Caller string   `json:"caller"`
		Stacks []string `json:"stacks"`
	}

	if err := json.Unmarshal([]byte(results[0]), &r); err != nil {
		t.Fatal(err)
	}

	if caller := "wrap_panic_test.go:func1:29"; r.Caller != caller {
		t.Errorf("expect caller '%s', but got '%s'", caller, r.Caller)
	}

	expects := []string{
		"github.com/xgfone/go-log/wrap_panic_test.go:func1:29",
		"github.com/xgfone/go-log/wrap_panic_test.go:handleBusiness:30",
		"github.com/xgfone/go-log/wrap_panic_test.go:testHandleBusiness:35",
		"github.com/xgfone/go-log/wrap_panic_test.go:TestWrapPanic:42",
	}

	for i, stack := range r.Stacks { // Remove the testing.go.
		if strings.HasPrefix(stack, "testing/testing.go:") {
			r.Stacks = r.Stacks[:i]
			break
		}
	}

	if len(expects) != len(r.Stacks) {
		t.Errorf("expect %d stacks, but got %d", len(expects), len(r.Stacks))
		t.Error(r.Stacks)
	} else {
		for i, line := range expects {
			if r.Stacks[i] != line {
				t.Errorf("%d: expect stack '%s', but got '%s'", i, line, r.Stacks[i])
			}
		}
	}
}

func TestGetCallStack(t *testing.T) {
	stacks := GetCallStack(1)
	for i, stack := range stacks {
		if strings.HasPrefix(stack, "testing/") {
			stacks = stacks[:i]
			break
		}
	}

	expects := []string{
		"github.com/xgfone/go-log/hook.go:GetCallStack:69",
		"github.com/xgfone/go-log/wrap_panic_test.go:TestGetCallStack:90",
	}

	if len(expects) != len(stacks) {
		t.Fatalf("expect %d line, but got %d: %v", len(expects), len(stacks), stacks)
	}

	for i, line := range expects {
		if line != stacks[i] {
			t.Errorf("%d: expect '%s', but got '%s'", i, line, stacks[i])
		}
	}
}
