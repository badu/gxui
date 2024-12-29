// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type ParentablePart struct {
	parent gxui.Parent
}

func (p *ParentablePart) Init() {}

func (p *ParentablePart) Parent() gxui.Parent {
	return p.parent
}

func (p *ParentablePart) SetParent(parent gxui.Parent) {
	p.parent = parent
}
