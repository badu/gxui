// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"

	"github.com/badu/gxui/math"
)

type LayoutChildren interface {
	LayoutChildren()
}

type LayoutableParent interface {
	Parent() Parent // was outer.Parenter
	Redraw()        // was outer.Redrawer
}

type LayoutablePart struct {
	parent           LayoutableParent
	margin           math.Spacing
	size             math.Size
	reLayoutNeeded   bool
	inLayoutChildren bool // True when calling LayoutChildren
}

func (l *LayoutablePart) Init(parent LayoutableParent) {
	l.parent = parent
}

func (l *LayoutablePart) SetMargin(margin math.Spacing) {
	l.margin = margin
	if p := l.parent.Parent(); p != nil {
		p.ReLayout()
	}
}

func (l *LayoutablePart) Margin() math.Spacing {
	return l.margin
}

func (l *LayoutablePart) Size() math.Size {
	return l.size
}

func (l *LayoutablePart) SetSize(newSize math.Size) {
	if newSize.W < 0 {
		panic(fmt.Errorf("SetSize() called with a negative width. Size: %v", newSize))
	}
	if newSize.H < 0 {
		panic(fmt.Errorf("SetSize() called with a negative height. Size: %v", newSize))
	}

	sizeChanged := l.size != newSize
	l.size = newSize
	if l.reLayoutNeeded || sizeChanged {
		l.reLayoutNeeded = false
		l.inLayoutChildren = true

		impl, ok := l.parent.(LayoutChildren)
		if ok {
			impl.LayoutChildren()
		}

		l.inLayoutChildren = false
		l.parent.Redraw()
	}
}

func (l *LayoutablePart) ReLayout() {
	if l.inLayoutChildren {
		panic("cannot call ReLayout() while in LayoutChildren")
	}

	if !l.reLayoutNeeded {
		if p := l.parent.Parent(); p != nil {
			l.reLayoutNeeded = true
			p.ReLayout()
		}
	}
}
