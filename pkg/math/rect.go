// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Rect struct {
	Min, Max Point
}

func CreateRect(minX, minY, maxX, maxY int) Rect {
	return Rect{Min: Point{X: minX, Y: minY}, Max: Point{X: maxX, Y: maxY}}
}

func (r Rect) Middle() Point {
	return Point{
		X: (r.Min.X + r.Max.X) / 2,
		Y: (r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rect) Width() int {
	return r.Max.X - r.Min.X
}

func (r Rect) Height() int {
	return r.Max.Y - r.Min.Y
}

func (r Rect) TopLeft() Point {
	return r.Min
}

func (r Rect) TopCenter() Point {
	return Point{X: (r.Min.X + r.Max.X) / 2, Y: r.Min.Y}
}

func (r Rect) TopRight() Point {
	return Point{X: r.Max.X, Y: r.Min.Y}
}

func (r Rect) BottomLeft() Point {
	return Point{X: r.Min.X, Y: r.Max.Y}
}

func (r Rect) BottomCenter() Point {
	return Point{X: (r.Min.X + r.Max.X) / 2, Y: r.Max.Y}
}

func (r Rect) BottomRight() Point {
	return r.Max
}

func (r Rect) MiddleLeft() Point {
	return Point{X: r.Min.X, Y: (r.Min.Y + r.Max.Y) / 2}
}

func (r Rect) MiddleRight() Point {
	return Point{X: r.Max.X, Y: (r.Min.Y + r.Max.Y) / 2}
}

func (r Rect) Size() Size {
	return Size{Width: r.Max.X - r.Min.X, Height: r.Max.Y - r.Min.Y}
}

func (r Rect) ScaleAt(p Point, s Vec2) Rect {
	return Rect{
		Min: p.Add(r.Min.Sub(p).Scale(s)),
		Max: p.Add(r.Max.Sub(p).Scale(s)),
	}
}

func (r Rect) ScaleS(s float32) Rect {
	return Rect{Min: r.Min.ScaleS(s), Max: r.Max.ScaleS(s)}
}

func (r Rect) Offset(p Point) Rect {
	return Rect{Min: r.Min.Add(p), Max: r.Max.Add(p)}
}

func (r Rect) OffsetX(x int) Rect {
	return r.Offset(Point{X: x})
}

func (r Rect) OffsetY(y int) Rect {
	return r.Offset(Point{Y: y})
}

func (r Rect) ClampXY(x, y int) (int, int) {
	return Clamp(x, r.Min.X, r.Max.X), Clamp(y, r.Min.Y, r.Max.Y)
}

func (r Rect) Lerp(vec2 Vec2) Point {
	return r.Min.Add(r.Size().Scale(vec2).Point())
}

func (r Rect) Frac(point Point) Vec2 {
	return point.Sub(r.Min).Vec2().Div(r.Size().Vec2())
}

func (r Rect) Remap(from, to Rect) Rect {
	return Rect{Min: r.Min.Remap(from, to), Max: r.Max.Remap(from, to)}
}

func (r Rect) Expand(spacing Spacing) Rect {
	return Rect{
		Min: Point{r.Min.X - spacing.Left, r.Min.Y - spacing.Top},
		Max: Point{r.Max.X + spacing.Right, r.Max.Y + spacing.Bottom},
	}.Canon()
}

func (r Rect) ExpandI(size int) Rect {
	return Rect{
		Min: Point{r.Min.X - size, r.Min.Y - size},
		Max: Point{r.Max.X + size, r.Max.Y + size},
	}.Canon()
}

func (r Rect) Contract(spacing Spacing) Rect {
	return Rect{
		Min: Point{r.Min.X + spacing.Left, r.Min.Y + spacing.Top},
		Max: Point{r.Max.X - spacing.Right, r.Max.Y - spacing.Bottom},
	}.Canon()
}

func (r Rect) ContractI(size int) Rect {
	return Rect{
		Min: Point{r.Min.X + size, r.Min.Y + size},
		Max: Point{r.Max.X - size, r.Max.Y - size},
	}.Canon()
}

func (r Rect) Union(rect Rect) Rect {
	return Rect{Min: r.Min.Min(rect.Min), Max: r.Max.Max(rect.Max)}
}

func (r Rect) Intersect(rect Rect) Rect {
	return Rect{
		Min: r.Min.Max(rect.Min),
		Max: r.Max.Min(rect.Max),
	}.Canon()
}

func (r Rect) Constrain(rect Rect) Rect {
	overflowMin := rect.Min.Sub(r.Min).Max(ZeroPoint)
	overflowMax := rect.Max.Sub(r.Max).Min(ZeroPoint)
	return Rect{
		Min: r.Min.Add(overflowMax).Max(rect.Min),
		Max: r.Max.Add(overflowMin).Min(rect.Max),
	}
}

func (r Rect) Canon() Rect {
	return Rect{
		Min: r.Min.Min(r.Max),
		Max: r.Min.Max(r.Max),
	}
}

func (r Rect) Contains(point Point) bool {
	return r.Min.X <= point.X && r.Min.Y <= point.Y && r.Max.X > point.X && r.Max.Y > point.Y
}
