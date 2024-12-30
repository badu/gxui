// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type SplitterLayout interface {
	Control
	Container
	ChildWeight(Control) float32
	SetChildWeight(Control, float32)
	Orientation() Orientation
	SetOrientation(Orientation)
}

type SplitterLayoutParent interface {
	BaseContainerParent
	CreateSplitterBar() Control
}

type SplitterLayoutImpl struct {
	ContainerBase
	parent        SplitterLayoutParent
	driver        Driver
	styles        *StyleDefs
	orientation   Orientation
	splitterWidth int
	weights       map[Control]float32
}

func (l *SplitterLayoutImpl) Init(parent SplitterLayoutParent, driver Driver, styles *StyleDefs) {
	l.ContainerBase.Init(parent, driver)
	l.parent = parent
	l.driver = driver
	l.styles = styles
	l.weights = make(map[Control]float32)
	l.splitterWidth = 4
	l.SetMouseEventTarget(true)
}

func (l *SplitterLayoutImpl) LayoutChildren() {
	size := l.parent.Size().Contract(l.Padding())
	offset := l.Padding().LT()

	children := l.parent.Children()

	splitterCount := len(children) / 2

	splitterWidth := l.splitterWidth
	if l.orientation.Horizontal() {
		size.W -= splitterWidth * splitterCount
	} else {
		size.H -= splitterWidth * splitterCount
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
				childWidth := int(float32(size.W) * frac)
				childRect = math.CreateRect(trackedDist+childMargin.L, childMargin.T, trackedDist+childWidth-childMargin.R, size.H-childMargin.B)
				trackedDist += childWidth
			} else {
				childHeight := int(float32(size.H) * frac)
				childRect = math.CreateRect(childMargin.L, trackedDist+childMargin.T, size.W-childMargin.R, trackedDist+childHeight-childMargin.B)
				trackedDist += childHeight
			}
		} else {
			if l.orientation.Horizontal() {
				childRect = math.CreateRect(trackedDist, 0, trackedDist+splitterWidth, size.H)
			} else {
				childRect = math.CreateRect(0, trackedDist, size.W, trackedDist+splitterWidth)
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
	b.Init(b, l.driver, l.styles)
	b.OnSplitterDragged(func(wndPnt math.Point) { l.SplitterDragged(b, wndPnt) })
	return b
}

func (l *SplitterLayoutImpl) SplitterDragged(splitter Control, wndPnt math.Point) {
	o := l.orientation
	p := WindowToChild(wndPnt, l.parent)
	children := l.ContainerBase.Children()
	splitterIndex := children.IndexOf(splitter)
	childA, childB := children[splitterIndex-1], children[splitterIndex+1]
	boundsA, boundsB := childA.Bounds(), childB.Bounds()

	min, max := o.Major(boundsA.Min.XY()), o.Major(boundsB.Max.XY())
	frac := math.RampSat(float32(o.Major(p.XY())), float32(min), float32(max))

	netWeight := l.weights[childA.Control] + l.weights[childB.Control]
	l.weights[childA.Control] = netWeight * frac
	l.weights[childB.Control] = netWeight * (1.0 - frac)
	l.LayoutChildren()
}

// base.ContainerBase overrides
func (l *SplitterLayoutImpl) AddChildAt(index int, control Control) *Child {
	l.weights[control] = 1.0
	if len(l.ContainerBase.Children()) > 0 {
		l.ContainerBase.AddChildAt(index, l.parent.CreateSplitterBar())
		index++
	}
	return l.ContainerBase.AddChildAt(index, control)
}

func (l *SplitterLayoutImpl) RemoveChildAt(index int) {
	children := l.ContainerBase.Children()
	if len(children) > 1 {
		l.ContainerBase.RemoveChildAt(index + 1)
	}
	delete(l.weights, children[index].Control)
	l.ContainerBase.RemoveChildAt(index)
}
