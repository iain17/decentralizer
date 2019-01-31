// Copyright 2018 github.com/ucirello
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

// Based on upspin.io/errors

// +build !debug

package errors

import (
	"fmt"
	"io"
	"testing"
)

func TestMarshal(t *testing.T) {
	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := E(Op("Get"), Invalid, "network unreachable")

	// Nested error.
	e2 := E(Op("Read"), Other, e1)

	b := MarshalError(e2)
	e3 := UnmarshalError(b)

	in := e2.(*Error)
	out := e3.(*Error)
	// Compare elementwise.
	if in.Op != out.Op {
		t.Errorf("expected Op %q; got %q", in.Op, out.Op)
	}
	if in.Kind != out.Kind {
		t.Errorf("expected kind %d; got %d", in.Kind, out.Kind)
	}
	// Note that error will have lost type information, so just check its Error string.
	if in.Err.Error() != out.Err.Error() {
		t.Errorf("expected Err %q; got %q", in.Err, out.Err)
	}
}

func TestSeparator(t *testing.T) {
	defer func(prev string) {
		Separator = prev
	}(Separator)
	Separator = ":: "

	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := E(Op("Get"), Invalid, "network unreachable")

	// Nested error.
	e2 := E(Op("Read"), Other, e1)

	want := "Read: invalid operation:: Get: network unreachable"
	if errorAsString(e2) != want {
		t.Errorf("expected %q; got %q", want, e2)
	}
}

func TestDoesNotChangePreviousError(t *testing.T) {
	err := E(Permission)
	err2 := E(Op("I will NOT modify err"), err)

	expected := "I will NOT modify err: permission denied"
	if errorAsString(err2) != expected {
		t.Fatalf("Expected %q, got %q", expected, err2)
	}
	kind := err.(*Error).Kind
	if kind != Permission {
		t.Fatalf("Expected kind %v, got %v", Permission, kind)
	}
}

func TestNoArgs(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = E()
}

type matchTest struct {
	err1, err2 error
	matched    bool
}

const (
	op  = Op("Op")
	op1 = Op("Op1")
	op2 = Op("Op2")
)

var matchTests = []matchTest{
	// Errors not of type *Error fail outright.
	{nil, nil, false},
	{io.EOF, io.EOF, false},
	{E(io.EOF), io.EOF, false},
	{io.EOF, E(io.EOF), false},
	// Success. We can drop fields from the first argument and still match.
	{E(io.EOF), E(io.EOF), true},
	// Failure.
	{E(io.EOF), E(io.ErrClosedPipe), false},
	{E(op1), E(op2), false},
	{E(Invalid), E(Permission), false},
}

func TestMatch(t *testing.T) {
	for _, test := range matchTests {
		matched := Match(test.err1, test.err2)
		if matched != test.matched {
			t.Errorf("Match(%q, %q)=%t; want %t", test.err1, test.err2, matched, test.matched)
		}
	}
}

type kindTest struct {
	err  error
	kind Kind
	want bool
}

var kindTests = []kindTest{
	// Non-Error errors.
	{nil, NotExist, false},
	{Str("not an *Error"), NotExist, false},

	// Basic comparisons.
	{E(NotExist), NotExist, true},
	{E(Exist), NotExist, false},
	{E("no kind"), NotExist, false},
	{E("no kind"), Other, false},

	// Nested *Error values.
	{E("Nesting", E(NotExist)), NotExist, true},
	{E("Nesting", E(Exist)), NotExist, false},
	{E("Nesting", E("no kind")), NotExist, false},
	{E("Nesting", E("no kind")), Other, false},
}

func TestKind(t *testing.T) {
	for _, test := range kindTests {
		got := Is(test.kind, test.err)
		if got != test.want {
			t.Errorf("Is(%q, %q)=%t; want %t", test.kind, test.err, got, test.want)
		}
	}
}

// errorAsString returns the string form of the provided error value.
// If the given string is an *Error, the stack information is removed
// before the value is stringified.
func errorAsString(err error) string {
	if e, ok := err.(*Error); ok {
		e2 := *e
		e2.stack = stack{}
		return e2.Error()
	}
	return err.Error()
}

func TestRootCause(t *testing.T) {
	rootErr := fmt.Errorf("root error cause")

	err1 := E(rootErr, "something bad happened")
	err2 := E(err1, "further up in stack")

	t.Log(err2)
	t.Log(err1)
	t.Log(rootErr)

	if RootCause(err2) != rootErr {
		t.Error("failed to detect error cause for err2:", err2)
	}
	if RootCause(err1) != rootErr {
		t.Error("failed to detect error cause for err1:", err1)
	}
	if RootCause(rootErr) != rootErr {
		t.Error("failed to detect error cause for the root cause itself:", rootErr)
	}

	var nilErr error
	if RootCause(nilErr) != nilErr {
		t.Error("failed to ignore nil error:")
	}
}
