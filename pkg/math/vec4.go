// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

type Vec4 struct {
	X, Y, Z, W float32
}

func (v Vec4) String() string {
	return fmt.Sprintf("(%.5f, %.5f, %.5f, %.5f)", v.X, v.Y, v.Z, v.W)
}

func (v Vec4) SqrLen() float32 {
	return v.Dot(v)
}

func (v Vec4) Len() float32 {
	return Sqrtf(v.SqrLen())
}

func (v Vec4) Normalize() Vec4 {
	l := v.Len()
	if l == 0 {
		return Vec4{}
	} else {
		return v.MulS(1.0 / v.Len())
	}
}

func (v Vec4) Neg() Vec4 {
	return Vec4{X: -v.X, Y: -v.Y, Z: -v.Z, W: -v.W}
}

func (v Vec4) XY() Vec2 {
	return Vec2{X: v.X, Y: v.Y}
}

func (v Vec4) Add(vec4 Vec4) Vec4 {
	return Vec4{X: v.X + vec4.X, Y: v.Y + vec4.Y, Z: v.Z + vec4.Z, W: v.W + vec4.W}
}

func (v Vec4) Sub(vec4 Vec4) Vec4 {
	return Vec4{X: v.X - vec4.X, Y: v.Y - vec4.Y, Z: v.Z - vec4.Z, W: v.W - vec4.W}
}

func (v Vec4) Mul(vec4 Vec4) Vec4 {
	return Vec4{X: v.X * vec4.X, Y: v.Y * vec4.Y, Z: v.Z * vec4.Z, W: v.W * vec4.W}
}

func (v Vec4) Div(vec4 Vec4) Vec4 {
	return Vec4{X: v.X / vec4.X, Y: v.Y / vec4.Y, Z: v.Z / vec4.Z, W: v.W / vec4.W}
}

func (v Vec4) Dot(vec4 Vec4) float32 {
	return v.X*vec4.X + v.Y*vec4.Y + v.Z*vec4.Z + v.W*vec4.W
}

func (v Vec4) MulS(size float32) Vec4 {
	return Vec4{X: v.X * size, Y: v.Y * size, Z: v.Z * size, W: v.W * size}
}

func (v Vec4) DivS(size float32) Vec4 {
	return Vec4{X: v.X / size, Y: v.Y / size, Z: v.Z / size, W: v.W / size}
}
