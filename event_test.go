// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/test_helper"
	"testing"
)

func TestEventNoArgs(t *testing.T) {
	e := CreateEvent(func() {})

	fired := false
	e.Listen(func() { fired = true })
	test_helper.AssertEquals(t, false, fired)

	e.Emit()
	test_helper.AssertEquals(t, true, fired)
}

type TestEvent func(i int, s string, b bool)

func TestEventExactArgs(t *testing.T) {
	e := CreateEvent(func(int, string, bool, int, int, bool) {})

	fired := false
	e.Listen(func(i1 int, s string, b1 bool, i2, i3 int, b2 bool) {
		test_helper.AssertEquals(t, 1, i1)
		test_helper.AssertEquals(t, "hello", s)
		test_helper.AssertEquals(t, false, b1)
		test_helper.AssertEquals(t, 2, i2)
		test_helper.AssertEquals(t, 3, i3)
		test_helper.AssertEquals(t, true, b2)
		fired = true
	})
	test_helper.AssertEquals(t, false, fired)

	e.Emit(1, "hello", false, 2, 3, true)
	test_helper.AssertEquals(t, true, fired)
}

func TestEventNilArgs(t *testing.T) {
	e := CreateEvent(func(chan int, func(), interface{}, map[int]int, *int, []int) {})

	fired := false
	e.Listen(func(c chan int, f func(), i interface{}, m map[int]int, p *int, s []int) {
		test_helper.AssertEquals(t, true, nil == c)
		test_helper.AssertEquals(t, true, nil == f)
		test_helper.AssertEquals(t, true, nil == i)
		test_helper.AssertEquals(t, true, nil == m)
		test_helper.AssertEquals(t, true, nil == p)
		test_helper.AssertEquals(t, true, nil == s)
		fired = true
	})
	test_helper.AssertEquals(t, false, fired)

	e.Emit(nil, nil, nil, nil, nil, nil)
	test_helper.AssertEquals(t, true, fired)
}

func TestEventMixedVariadic(t *testing.T) {
	e := CreateEvent(func(int, int, ...int) {})

	fired := false
	e.Listen(func(a, b int, cde ...int) {
		test_helper.AssertEquals(t, 3, len(cde))

		test_helper.AssertEquals(t, 0, a)
		test_helper.AssertEquals(t, 1, b)
		test_helper.AssertEquals(t, 2, cde[0])
		test_helper.AssertEquals(t, 3, cde[1])
		test_helper.AssertEquals(t, 4, cde[2])
		fired = true
	})
	e.Emit(0, 1, 2, 3, 4)
	test_helper.AssertEquals(t, true, fired)
}

func TestEventSingleVariadic(t *testing.T) {
	e := CreateEvent(func(...int) {})

	fired := false
	e.Listen(func(va ...int) {
		test_helper.AssertEquals(t, 3, len(va))

		test_helper.AssertEquals(t, 2, va[0])
		test_helper.AssertEquals(t, 3, va[1])
		test_helper.AssertEquals(t, 4, va[2])
		fired = true
	})
	e.Emit(2, 3, 4)
	test_helper.AssertEquals(t, true, fired)
}

func TestEventEmptyVariadic(t *testing.T) {
	e := CreateEvent(func(...int) {})

	fired := false
	e.Listen(func(va ...int) {
		test_helper.AssertEquals(t, 0, len(va))
		fired = true
	})
	e.Emit()
	test_helper.AssertEquals(t, true, fired)
}

func TestEventChaining(t *testing.T) {
	e1 := CreateEvent(func(int, string, bool, int, int, bool) {})
	e2 := CreateEvent(func(int, string, bool, int, int, bool) {})

	e1.Listen(e2)

	fired := false
	e2.Listen(func(i1 int, s string, b1 bool, i2, i3 int, b2 bool) {
		test_helper.AssertEquals(t, 1, i1)
		test_helper.AssertEquals(t, "hello", s)
		test_helper.AssertEquals(t, false, b1)
		test_helper.AssertEquals(t, 2, i2)
		test_helper.AssertEquals(t, 3, i3)
		test_helper.AssertEquals(t, true, b2)
		fired = true
	})
	test_helper.AssertEquals(t, false, fired)

	e1.Emit(1, "hello", false, 2, 3, true)
	test_helper.AssertEquals(t, true, fired)
}

func TestEventUnlisten(t *testing.T) {
	eI := CreateEvent(func() {})
	eJ := CreateEvent(func() {})
	eK := CreateEvent(func() {})

	e := CreateEvent(func() {})
	e.Listen(eI)
	e.Listen(eJ)
	e.Listen(eK)

	i, j, k := 0, 0, 0
	subI := eI.Listen(func() { i++ })
	subJ := eJ.Listen(func() { j++ })
	subK := eK.Listen(func() { k++ })

	test_helper.AssertEquals(t, 0, i)
	test_helper.AssertEquals(t, 0, j)
	test_helper.AssertEquals(t, 0, k)

	e.Emit()
	test_helper.AssertEquals(t, 1, i)
	test_helper.AssertEquals(t, 1, j)
	test_helper.AssertEquals(t, 1, k)

	subJ.Forget()

	e.Emit()
	test_helper.AssertEquals(t, 2, i)
	test_helper.AssertEquals(t, 1, j)
	test_helper.AssertEquals(t, 2, k)

	subK.Forget()

	e.Emit()
	test_helper.AssertEquals(t, 3, i)
	test_helper.AssertEquals(t, 1, j)
	test_helper.AssertEquals(t, 2, k)

	subI.Forget()

	e.Emit()
	test_helper.AssertEquals(t, 3, i)
	test_helper.AssertEquals(t, 1, j)
	test_helper.AssertEquals(t, 2, k)
}

// TODO: Add tests for early signature mismatch failures
