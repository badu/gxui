// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Point struct {
	X, Y int
}

func NewPoint(X, Y int) Point {
	return Point{X: X, Y: Y}
}

func (p Point) Add(point Point) Point {
	return Point{X: p.X + point.X, Y: p.Y + point.Y}
}

func (p Point) AddX(x int) Point {
	return Point{X: p.X + x, Y: p.Y}
}

func (p Point) AddY(y int) Point {
	return Point{X: p.X, Y: p.Y + y}
}

func (p Point) Sub(point Point) Point {
	return Point{X: p.X - point.X, Y: p.Y - point.Y}
}

func (p Point) Neg() Point {
	return Point{X: -p.X, Y: -p.Y}
}

func (p Point) SqrLen() int {
	return p.Dot(p)
}

func (p Point) Len() float32 {
	return Sqrtf(float32(p.SqrLen()))
}

func (p Point) Dot(o Point) int {
	return p.X*o.X + p.Y*o.Y
}

func (p Point) XY() (x, y int) {
	return p.X, p.Y
}

func (p Point) Vec2() Vec2 {
	return Vec2{X: float32(p.X), Y: float32(p.Y)}
}

func (p Point) Vec3(z float32) Vec3 {
	return Vec3{X: float32(p.X), Y: float32(p.Y), Z: z}
}

func (p Point) Scale(s Vec2) Point {
	return Point{X: int(float32(p.X) * s.X), Y: int(float32(p.Y) * s.Y)}
}

func (p Point) ScaleS(s float32) Point {
	return Point{X: int(float32(p.X) * s), Y: int(float32(p.Y) * s)}
}

func (p Point) ScaleX(s float32) Point {
	return Point{X: int(float32(p.X) * s), Y: p.Y}
}

func (p Point) ScaleY(s float32) Point {
	return Point{X: p.X, Y: int(float32(p.Y) * s)}
}

func (p Point) Size() Size {
	return Size{Width: p.X, Height: p.Y}
}

func (p Point) Min(point Point) Point {
	return Point{X: min(p.X, point.X), Y: min(p.Y, point.Y)}
}

func (p Point) Max(point Point) Point {
	return Point{X: max(p.X, point.X), Y: max(p.Y, point.Y)}
}

func (p Point) Clamp(min, max Point) Point {
	return p.Min(max).Max(min)
}

func (p Point) Remap(from, to Rect) Point {
	return p.Sub(from.Min).
		ScaleX(float32(to.Width()) / float32(from.Width())).
		ScaleY(float32(to.Height()) / float32(from.Height())).
		Add(to.Min)
}
