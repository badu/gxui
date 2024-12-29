// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/badu/gxui"
)

type ParentableOuter interface{}

type ParentablePart struct {
	outer  ParentableOuter
	parent gxui.Parent
}

func (p *ParentablePart) Init(outer ParentableOuter) {
	p.outer = outer
}

func (p *ParentablePart) Parent() gxui.Parent {
	return p.parent
}

func (p *ParentablePart) SetParent(parent gxui.Parent) {
	p.parent = parent
}
