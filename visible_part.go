// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

// VisibleParent is minimal
type VisibleParent interface {
	Parent() Parent
}

type VisiblePart struct {
	parent  VisibleParent
	visible bool
}

func (v *VisiblePart) Init(parent VisibleParent) {
	v.parent = parent
	v.visible = true
}

func (v *VisiblePart) IsVisible() bool {
	return v.visible
}

func (v *VisiblePart) SetVisible(visible bool) {
	if v.visible == visible {
		return
	}

	v.visible = visible
	if grandParent := v.parent.Parent(); grandParent != nil {
		grandParent.Redraw()
	}
}
