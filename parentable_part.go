// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type ParentablePart struct {
	parent Parent
}

func (p *ParentablePart) Parent() Parent {
	return p.parent
}

func (p *ParentablePart) SetParent(parent Parent) {
	p.parent = parent
}
