// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import "github.com/chewxy/math32"

type Vec2 struct {
	X float32
	Y float32
}

func (v Vec2) SqrLen() float32 {
	return v.Dot(v)
}

func (v Vec2) Len() float32 {
	return math32.Sqrt(v.SqrLen())
}

func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{}
	}

	return v.MulS(1.0 / l)
}

func (v Vec2) Tangent() Vec2 {
	return Vec2{X: -v.Y, Y: v.X}
}

func (v Vec2) Vec3(z float32) Vec3 {
	return Vec3{X: v.X, Y: v.Y, Z: z}
}

func (v Vec2) Add(vec2 Vec2) Vec2 {
	return Vec2{X: v.X + vec2.X, Y: v.Y + vec2.Y}
}

func (v Vec2) Sub(vec2 Vec2) Vec2 {
	return Vec2{X: v.X - vec2.X, Y: v.Y - vec2.Y}
}

func (v Vec2) Div(vec2 Vec2) Vec2 {
	return Vec2{X: v.X / vec2.X, Y: v.Y / vec2.Y}
}

func (v Vec2) Dot(vec2 Vec2) float32 {
	return v.X*vec2.X + v.Y*vec2.Y
}

func (v Vec2) Cross(vec2 Vec2) float32 {
	return v.X*vec2.Y - v.Y*vec2.X
}

func (v Vec2) MulS(size float32) Vec2 {
	return Vec2{X: v.X * size, Y: v.Y * size}
}

func (v Vec2) DivS(size float32) Vec2 {
	return Vec2{X: v.X / size, Y: v.Y / size}
}
