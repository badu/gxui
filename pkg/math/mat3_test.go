// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"testing"

	"github.com/badu/gxui/test_helper"
)

func TestMat3InvertIdent(t *testing.T) {
	m := Mat3Ident.Invert()
	test_helper.AssertEquals(t, Mat3Ident, m)
}

func TestCreateMat3PositionToBarycentric(t *testing.T) {
	a := Vec2{+0.0, -1.0}
	b := Vec2{-1.0, 1.0}
	c := Vec2{+1.0, 1.0}
	m := CreateMat3PositionToBarycentric(a, b, c)
	test_helper.AssertEquals(t, Vec3{1.0, 0.0, 1.0}, a.Vec3(1).MulM(m))
	test_helper.AssertEquals(t, Vec3{0.0, 1.0, 1.0}, b.Vec3(1).MulM(m))
	test_helper.AssertEquals(t, Vec3{0.0, 0.0, 1.0}, c.Vec3(1).MulM(m))
}
