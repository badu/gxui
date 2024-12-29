// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type SplitterLayout struct {
	ContainerBase
	outer         gxui.SplitterLayoutOuter
	theme         gxui.Theme
	orientation   gxui.Orientation
	splitterWidth int
	weights       map[gxui.Control]float32
}

func (l *SplitterLayout) Init(outer gxui.SplitterLayoutOuter, theme gxui.Theme) {
	l.ContainerBase.Init(outer, theme)
	l.outer = outer
	l.theme = theme
	l.weights = make(map[gxui.Control]float32)
	l.splitterWidth = 4
	l.SetMouseEventTarget(true)
}

func (l *SplitterLayout) LayoutChildren() {
	size := l.outer.Size().Contract(l.Padding())
	offset := l.Padding().LT()

	children := l.outer.Children()

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

func (l *SplitterLayout) ChildWeight(child gxui.Control) float32 {
	return l.weights[child]
}

func (l *SplitterLayout) SetChildWeight(child gxui.Control, weight float32) {
	if l.weights[child] != weight {
		l.weights[child] = weight
		l.LayoutChildren()
	}
}

func (l *SplitterLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *SplitterLayout) Orientation() gxui.Orientation {
	return l.orientation
}

func (l *SplitterLayout) SetOrientation(o gxui.Orientation) {
	if l.orientation != o {
		l.orientation = o
		l.LayoutChildren()
	}
}

func (l *SplitterLayout) CreateSplitterBar() gxui.Control {
	b := &SplitterBar{}
	b.Init(b, l.theme)
	b.OnSplitterDragged(func(wndPnt math.Point) { l.SplitterDragged(b, wndPnt) })
	return b
}

func (l *SplitterLayout) SplitterDragged(splitter gxui.Control, wndPnt math.Point) {
	o := l.orientation
	p := gxui.WindowToChild(wndPnt, l.outer)
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
func (l *SplitterLayout) AddChildAt(index int, control gxui.Control) *gxui.Child {
	l.weights[control] = 1.0
	if len(l.ContainerBase.Children()) > 0 {
		l.ContainerBase.AddChildAt(index, l.outer.CreateSplitterBar())
		index++
	}
	return l.ContainerBase.AddChildAt(index, control)
}

func (l *SplitterLayout) RemoveChildAt(index int) {
	children := l.ContainerBase.Children()
	if len(children) > 1 {
		l.ContainerBase.RemoveChildAt(index + 1)
	}
	delete(l.weights, children[index].Control)
	l.ContainerBase.RemoveChildAt(index)
}
