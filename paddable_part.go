// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type PaddableParent interface {
	LayoutChildren() // was outer.LayoutChildren
	Redraw()         // was outer.Redrawer
}

type PaddablePart struct {
	parent  PaddableParent
	padding math.Spacing
}

func (p *PaddablePart) Init(parent PaddableParent) {
	p.parent = parent
}

func (p *PaddablePart) SetPadding(m math.Spacing) {
	p.padding = m
	p.parent.LayoutChildren()
	p.parent.Redraw()
}

func (p *PaddablePart) Padding() math.Spacing {
	return p.padding
}
