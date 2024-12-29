// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parts

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type LayoutChildren interface {
	LayoutChildren()
}

type LayoutableOuter interface {
	Parent() gxui.Parent // was outer.Parenter
	Redraw()             // was outer.Redrawer
}

type Layoutable struct {
	outer             LayoutableOuter
	driver            gxui.Driver
	margin            math.Spacing
	size              math.Size
	relayoutRequested bool
	inLayoutChildren  bool // True when calling LayoutChildren
}

func (l *Layoutable) Init(outer LayoutableOuter, theme gxui.Theme) {
	l.outer = outer
	l.driver = theme.Driver()
}

func (l *Layoutable) SetMargin(margin math.Spacing) {
	l.margin = margin
	if p := l.outer.Parent(); p != nil {
		p.Relayout()
	}
}

func (l *Layoutable) Margin() math.Spacing {
	return l.margin
}

func (l *Layoutable) Size() math.Size {
	return l.size
}

func (l *Layoutable) SetSize(newSize math.Size) {
	if newSize.W < 0 {
		panic(fmt.Errorf("SetSize() called with a negative width. Size: %v", newSize))
	}
	if newSize.H < 0 {
		panic(fmt.Errorf("SetSize() called with a negative height. Size: %v", newSize))
	}

	sizeChanged := l.size != newSize
	l.size = newSize
	if l.relayoutRequested || sizeChanged {
		l.relayoutRequested = false
		l.inLayoutChildren = true

		impl, ok := l.outer.(LayoutChildren)
		if ok {
			impl.LayoutChildren()
		}

		l.inLayoutChildren = false
		l.outer.Redraw()
	}
}

func (l *Layoutable) Relayout() {
	l.driver.AssertUIGoroutine()
	if l.inLayoutChildren {
		panic("cannot call Relayout() while in LayoutChildren")
	}

	if !l.relayoutRequested {
		if p := l.outer.Parent(); p != nil {
			l.relayoutRequested = true
			p.Relayout()
		}
	}
}
