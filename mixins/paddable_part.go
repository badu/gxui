// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui/math"
)

type PaddableOuter interface {
	LayoutChildren() // was outer.LayoutChildren
	Redraw()         // was outer.Redrawer
}

type PaddablePart struct {
	outer   PaddableOuter
	padding math.Spacing
}

func (p *PaddablePart) Init(outer PaddableOuter) {
	p.outer = outer
}

func (p *PaddablePart) SetPadding(m math.Spacing) {
	p.padding = m
	p.outer.LayoutChildren()
	p.outer.Redraw()
}

func (p *PaddablePart) Padding() math.Spacing {
	return p.padding
}
