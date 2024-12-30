// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type BubbleOverlay interface {
	Control
	Show(control Control, target math.Point)
	Hide()
}

type BubbleOverlayImpl struct {
	ContainerBase
	outer       ContainerBaseOuter
	targetPoint math.Point
	arrowLength int
	arrowWidth  int
	brush       Brush
	pen         Pen
}

func (o *BubbleOverlayImpl) Init(outer ContainerBaseOuter, theme Theme) {
	o.ContainerBase.Init(outer, theme)
	o.outer = outer
	o.arrowLength = 20
	o.arrowWidth = 15
}

func (o *BubbleOverlayImpl) LayoutChildren() {
	for _, child := range o.outer.Children() {
		bounds := o.outer.Size().Rect().Contract(o.outer.Padding())
		arrowPadding := math.CreateSpacing(o.arrowLength)
		cm := child.Control.Margin()
		cs := child.Control.DesiredSize(math.ZeroSize, bounds.Size().Contract(cm).Max(math.ZeroSize))
		cr := cs.Expand(arrowPadding).EdgeAlignedFit(bounds, o.targetPoint).Contract(arrowPadding)
		child.Layout(cr)
	}
}

func (o *BubbleOverlayImpl) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (o *BubbleOverlayImpl) Show(control Control, point math.Point) {
	o.Hide()
	o.outer.AddChild(control)
	o.targetPoint = point
}

func (o *BubbleOverlayImpl) Hide() {
	o.outer.RemoveAll()
}

func (o *BubbleOverlayImpl) Brush() Brush {
	return o.brush
}

func (o *BubbleOverlayImpl) SetBrush(brush Brush) {
	if o.brush != brush {
		o.brush = brush
		o.Redraw()
	}
}

func (o *BubbleOverlayImpl) Pen() Pen {
	return o.pen
}

func (o *BubbleOverlayImpl) SetPen(pen Pen) {
	if o.pen != pen {
		o.pen = pen
		o.Redraw()
	}
}

func (o *BubbleOverlayImpl) Paint(canvas Canvas) {
	if !o.IsVisible() {
		return
	}

	for _, child := range o.outer.Children() {
		expandedBounds := child.Bounds().Expand(o.outer.Padding())
		targetPoint := o.targetPoint
		halfWidth := o.arrowWidth / 2
		var polygon Polygon

		switch {
		case targetPoint.X < expandedBounds.Min.X:
			/*
			    A-----------------B
			    G                 |
			 F                    |
			    E                 |
			    D-----------------C
			*/
			polygon = Polygon{
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
			polygon = Polygon{
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
			polygon = Polygon{
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
			polygon = Polygon{
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
