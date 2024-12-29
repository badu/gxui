// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type BubbleOverlay struct {
	ContainerBase
	outer       ContainerBaseOuter
	targetPoint math.Point
	arrowLength int
	arrowWidth  int
	brush       gxui.Brush
	pen         gxui.Pen
}

func (o *BubbleOverlay) Init(outer ContainerBaseOuter, theme gxui.Theme) {
	o.ContainerBase.Init(outer, theme)
	o.outer = outer
	o.arrowLength = 20
	o.arrowWidth = 15
}

func (o *BubbleOverlay) LayoutChildren() {
	for _, child := range o.outer.Children() {
		bounds := o.outer.Size().Rect().Contract(o.outer.Padding())
		arrowPadding := math.CreateSpacing(o.arrowLength)
		cm := child.Control.Margin()
		cs := child.Control.DesiredSize(math.ZeroSize, bounds.Size().Contract(cm).Max(math.ZeroSize))
		cr := cs.Expand(arrowPadding).EdgeAlignedFit(bounds, o.targetPoint).Contract(arrowPadding)
		child.Layout(cr)
	}
}

func (o *BubbleOverlay) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (o *BubbleOverlay) Show(control gxui.Control, point math.Point) {
	o.Hide()
	o.outer.AddChild(control)
	o.targetPoint = point
}

func (o *BubbleOverlay) Hide() {
	o.outer.RemoveAll()
}

func (o *BubbleOverlay) Brush() gxui.Brush {
	return o.brush
}

func (o *BubbleOverlay) SetBrush(brush gxui.Brush) {
	if o.brush != brush {
		o.brush = brush
		o.Redraw()
	}
}

func (o *BubbleOverlay) Pen() gxui.Pen {
	return o.pen
}

func (o *BubbleOverlay) SetPen(pen gxui.Pen) {
	if o.pen != pen {
		o.pen = pen
		o.Redraw()
	}
}

func (o *BubbleOverlay) Paint(canvas gxui.Canvas) {
	if !o.IsVisible() {
		return
	}

	for _, child := range o.outer.Children() {
		expandedBounds := child.Bounds().Expand(o.outer.Padding())
		targetPoint := o.targetPoint
		halfWidth := o.arrowWidth / 2
		var polygon gxui.Polygon

		switch {
		case targetPoint.X < expandedBounds.Min.X:
			/*
			    A-----------------B
			    G                 |
			 F                    |
			    E                 |
			    D-----------------C
			*/
			polygon = gxui.Polygon{
				/*A*/ {Position: expandedBounds.TL(), RoundedRadius: 5},
				/*B*/ {Position: expandedBounds.TR(), RoundedRadius: 5},
				/*C*/ {Position: expandedBounds.BR(), RoundedRadius: 5},
				/*D*/ {Position: expandedBounds.BL(), RoundedRadius: 5},
				/*E*/ {Position: math.Point{X: expandedBounds.Min.X, Y: math.Clamp(targetPoint.Y+halfWidth, expandedBounds.Min.Y+halfWidth, expandedBounds.Max.Y)}, RoundedRadius: 0},
				/*F*/ {Position: targetPoint, RoundedRadius: 0},
				/*G*/ {Position: math.Point{X: expandedBounds.Min.X, Y: math.Clamp(targetPoint.Y-halfWidth, expandedBounds.Min.Y, expandedBounds.Max.Y-halfWidth)}, RoundedRadius: 0},
			}
			// fmt.Printf("A: %+v\n", polygon)
		case targetPoint.X > expandedBounds.Max.X:
			/*
			   A-----------------B
			   |                 C
			   |                    D
			   |                 E
			   G-----------------F
			*/
			polygon = gxui.Polygon{
				/*A*/ {Position: expandedBounds.TL(), RoundedRadius: 5},
				/*B*/ {Position: expandedBounds.TR(), RoundedRadius: 5},
				/*C*/ {Position: math.Point{X: expandedBounds.Max.X, Y: math.Clamp(targetPoint.Y-halfWidth, expandedBounds.Min.Y, expandedBounds.Max.Y-halfWidth)}, RoundedRadius: 0},
				/*D*/ {Position: targetPoint, RoundedRadius: 0},
				/*E*/ {Position: math.Point{X: expandedBounds.Max.X, Y: math.Clamp(targetPoint.Y+halfWidth, expandedBounds.Min.Y+halfWidth, expandedBounds.Max.Y)}, RoundedRadius: 0},
				/*F*/ {Position: expandedBounds.BR(), RoundedRadius: 5},
				/*G*/ {Position: expandedBounds.BL(), RoundedRadius: 5},
			}
			// fmt.Printf("B: %+v\n", polygon)
		case targetPoint.Y < expandedBounds.Min.Y:
			/*
			                 C
			                / \
			   A-----------B   D-E
			   |                 |
			   |                 |
			   G-----------------F
			*/
			polygon = gxui.Polygon{
				/*A*/ {Position: expandedBounds.TL(), RoundedRadius: 5},
				/*B*/ {Position: math.Point{X: math.Clamp(targetPoint.X-halfWidth, expandedBounds.Min.X, expandedBounds.Max.X-halfWidth), Y: expandedBounds.Min.Y}, RoundedRadius: 0},
				/*C*/ {Position: targetPoint, RoundedRadius: 0},
				/*D*/ {Position: math.Point{X: math.Clamp(targetPoint.X+halfWidth, expandedBounds.Min.X+halfWidth, expandedBounds.Max.X), Y: expandedBounds.Min.Y}, RoundedRadius: 0},
				/*E*/ {Position: expandedBounds.TR(), RoundedRadius: 5},
				/*F*/ {Position: expandedBounds.BR(), RoundedRadius: 5},
				/*G*/ {Position: expandedBounds.BL(), RoundedRadius: 5},
			}
			// fmt.Printf("C: %+v\n", polygon)
		default:
			/*
			   A-----------------B
			   |                 |
			   |                 |
			   G-----------F   D-C
			                \ /
			                 E
			*/
			polygon = gxui.Polygon{
				/*A*/ {Position: expandedBounds.TL(), RoundedRadius: 5},
				/*B*/ {Position: expandedBounds.TR(), RoundedRadius: 5},
				/*C*/ {Position: expandedBounds.BR(), RoundedRadius: 5},
				/*D*/ {Position: math.Point{X: math.Clamp(targetPoint.X+halfWidth, expandedBounds.Min.X+halfWidth, expandedBounds.Max.X), Y: expandedBounds.Max.Y}, RoundedRadius: 0},
				/*E*/ {Position: targetPoint, RoundedRadius: 0},
				/*F*/ {Position: math.Point{X: math.Clamp(targetPoint.X-halfWidth, expandedBounds.Min.X, expandedBounds.Max.X-halfWidth), Y: expandedBounds.Max.Y}, RoundedRadius: 0},
				/*G*/ {Position: expandedBounds.BL(), RoundedRadius: 5},
			}
			// fmt.Printf("D: %+v\n", polygon)
		}
		canvas.DrawPolygon(polygon, o.pen, o.brush)
	}

	o.PaintChildrenPart.Paint(canvas)
}
