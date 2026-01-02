// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

type SplitterLayoutParent interface {
	Control

	CreateSplitterBar() Control

	Children() Children
	ReLayout()
	Redraw()

	AddChildAt(index int, child Control) *Child
	RemoveChildAt(index int)

	Attached() bool
	OnAttach(callback func()) EventSubscription
	OnDetach(callback func()) EventSubscription
	IsVisible() bool
	Size() math.Size
	Parent() Parent
	Paint(canvas Canvas)
	LayoutChildren()
	PaintChild(canvas Canvas, child *Child, idx int)
}

type SplitterLayoutImpl struct {
	InputEventHandlerPart
	PaintChildrenPart
	ParentablePart
	DrawPaintPart
	AttachablePart
	VisiblePart
	ContainerPart
	PaddablePart
	LayoutablePart
	parent        SplitterLayoutParent
	canvasCreator CanvasCreator
	styles        *StyleDefs
	weights       map[Control]float32
	orientation   Orientation
	splitterWidth int
}

func (l *SplitterLayoutImpl) Init(parent SplitterLayoutParent, driver Driver, styles *StyleDefs) {
	l.ContainerPart.Init(parent)
	l.DrawPaintPart.Init(parent, driver)
	l.InputEventHandlerPart.Init()
	l.LayoutablePart.Init(parent)
	l.PaddablePart.Init(parent)
	l.PaintChildrenPart.Init(parent)
	l.VisiblePart.Init(parent)

	l.parent = parent
	l.canvasCreator = driver
	l.styles = styles
	l.weights = make(map[Control]float32)
	l.splitterWidth = 4
	l.SetMouseEventTarget(true)
}

func (l *SplitterLayoutImpl) LayoutChildren() {
	size := l.parent.Size().Contract(l.Padding())
	offset := l.Padding().TopLeft()

	children := l.parent.Children()

	splitterCount := len(children) / 2

	splitterWidth := l.splitterWidth
	if l.orientation.Horizontal() {
		size.Width -= splitterWidth * splitterCount
	} else {
		size.Height -= splitterWidth * splitterCount
	}

	netWeight := float32(0.0)
	for i, c := range children {
		if isSplitter := (i & 1) == 1; !isSplitter {
			netWeight += l.weights[c.Control]
		}
	}

	trackedDist := 0
	for i, child := range children {
		var childRect math.Rect
		if isSplitter := (i & 1) == 1; !isSplitter {
			childMargin := child.Control.Margin()
			frac := l.weights[child.Control] / netWeight
			if l.orientation.Horizontal() {
				childWidth := int(float32(size.Width) * frac)
				childRect = math.CreateRect(trackedDist+childMargin.Left, childMargin.Top, trackedDist+childWidth-childMargin.Right, size.Height-childMargin.Bottom)
				trackedDist += childWidth
			} else {
				childHeight := int(float32(size.Height) * frac)
				childRect = math.CreateRect(childMargin.Left, trackedDist+childMargin.Top, size.Width-childMargin.Right, trackedDist+childHeight-childMargin.Bottom)
				trackedDist += childHeight
			}
		} else {
			if l.orientation.Horizontal() {
				childRect = math.CreateRect(trackedDist, 0, trackedDist+splitterWidth, size.Height)
			} else {
				childRect = math.CreateRect(0, trackedDist, size.Width, trackedDist+splitterWidth)
			}
			trackedDist += splitterWidth
		}
		child.Layout(childRect.Offset(offset).Canon())
	}
}

func (l *SplitterLayoutImpl) ChildWeight(child Control) float32 {
	return l.weights[child]
}

func (l *SplitterLayoutImpl) SetChildWeight(child Control, weight float32) {
	if l.weights[child] != weight {
		l.weights[child] = weight
		l.LayoutChildren()
	}
}

func (l *SplitterLayoutImpl) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *SplitterLayoutImpl) Orientation() Orientation {
	return l.orientation
}

func (l *SplitterLayoutImpl) SetOrientation(o Orientation) {
	if l.orientation != o {
		l.orientation = o
		l.LayoutChildren()
	}
}

func (l *SplitterLayoutImpl) CreateSplitterBar() Control {
	b := &SplitterBar{}
	b.Init(b, l.canvasCreator, l.styles)
	b.OnSplitterDragged(func(wndPnt math.Point) { l.SplitterDragged(b, wndPnt) })
	return b
}

func (l *SplitterLayoutImpl) SplitterDragged(splitter Control, windowPoint math.Point) {
	o := l.orientation
	point := WindowToChild(windowPoint, l.parent)
	children := l.ContainerPart.Children()
	splitterIndex := children.IndexOf(splitter)
	childA, childB := children[splitterIndex-1], children[splitterIndex+1]
	boundsA, boundsB := childA.Bounds(), childB.Bounds()

	minB, maxB := o.Major(boundsA.Min.XY()), o.Major(boundsB.Max.XY())
	frac := math.RampSat(float32(o.Major(point.XY())), float32(minB), float32(maxB))

	netWeight := l.weights[childA.Control] + l.weights[childB.Control]
	l.weights[childA.Control] = netWeight * frac
	l.weights[childB.Control] = netWeight * (1.0 - frac)
	l.LayoutChildren()
}

// base.ContainerBase overrides
func (l *SplitterLayoutImpl) AddChildAt(index int, control Control) *Child {
	l.weights[control] = 1.0
	if len(l.ContainerPart.Children()) > 0 {
		l.ContainerPart.AddChildAt(index, l.parent.CreateSplitterBar())
		index++
	}
	return l.ContainerPart.AddChildAt(index, control)
}

func (l *SplitterLayoutImpl) RemoveChildAt(index int) {
	children := l.ContainerPart.Children()
	if len(children) > 1 {
		l.ContainerPart.RemoveChildAt(index + 1)
	}
	delete(l.weights, children[index].Control)
	l.ContainerPart.RemoveChildAt(index)
}
