// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) String() string {
	return fmt.Sprintf("(%.5f, %.5f, %.5f)", v.X, v.Y, v.Z)
}

func (v Vec3) SqrLen() float32 {
	return v.Dot(v)
}

func (v Vec3) Len() float32 {
	return Sqrtf(v.SqrLen())
}

func (v Vec3) Normalize() Vec3 {
	l := v.Len()
	if l == 0 {
		return Vec3{0, 0, 0}
	} else {
		return v.MulS(1.0 / v.Len())
	}
}

func (v Vec3) Neg() Vec3 {
	return Vec3{X: -v.X, Y: -v.Y, Z: -v.Z}
}

func (v Vec3) XY() Vec2 {
	return Vec2{X: v.X, Y: v.Y}
}

func (v Vec3) Add(vec3 Vec3) Vec3 {
	return Vec3{X: v.X + vec3.X, Y: v.Y + vec3.Y, Z: v.Z + vec3.Z}
}

func (v Vec3) Sub(vec3 Vec3) Vec3 {
	return Vec3{X: v.X - vec3.X, Y: v.Y - vec3.Y, Z: v.Z - vec3.Z}
}

func (v Vec3) Mul(vec3 Vec3) Vec3 {
	return Vec3{X: v.X * vec3.X, Y: v.Y * vec3.Y, Z: v.Z * vec3.Z}
}

func (v Vec3) Div(vec3 Vec3) Vec3 {
	return Vec3{X: v.X / vec3.X, Y: v.Y / vec3.Y, Z: v.Z / vec3.Z}
}

func (v Vec3) Dot(vec3 Vec3) float32 {
	return v.X*vec3.X + v.Y*vec3.Y + v.Z*vec3.Z
}

func (v Vec3) Cross(vec3 Vec3) Vec3 {
	return Vec3{
		X: v.Y*vec3.Z - v.Z*vec3.Y,
		Y: v.Z*vec3.X - v.X*vec3.Z,
		Z: v.X*vec3.Y - v.Y*vec3.X,
	}
}

//	╭          ╮
//	│ M₀ M₁ M₂ │
//
// [V₀, V₁, V₂] ⨯ │ M₃ M₄ M₅ │ = [R₀, R₁, R₂]
//
//	│ M₆ M₇ M₈ │
//	╰          ╯
//
// R₀ = V₀ • M₀ + V₁ • M₃ + V₂ • M₆
// R₁ = V₀ • M₁ + V₁ • M₄ + V₂ • M₇
// R₂ = V₀ • M₂ + V₁ • M₅ + V₂ • M₈
func (v Vec3) MulM(mat3 Mat3) Vec3 {
	a := mat3.Row(0).MulS(v.X)
	b := mat3.Row(1).MulS(v.Y)
	c := mat3.Row(2).MulS(v.Z)
	return a.Add(b).Add(c)
}

func (v Vec3) MulS(size float32) Vec3 {
	return Vec3{v.X * size, v.Y * size, v.Z * size}
}

func (v Vec3) DivS(size float32) Vec3 {
	return Vec3{v.X / size, v.Y / size, v.Z / size}
}
