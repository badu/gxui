// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type VisibleOuter interface {
	Parent() gxui.Parent // was outer.Parenter
	Redraw()             // was outer.Redrawer
}

type VisiblePart struct {
	outer   VisibleOuter
	visible bool
}

func (v *VisiblePart) Init(outer VisibleOuter) {
	v.outer = outer
	v.visible = true
}

func (v *VisiblePart) IsVisible() bool {
	return v.visible
}

func (v *VisiblePart) SetVisible(visible bool) {
	if v.visible != visible {
		v.visible = visible
		if p := v.outer.Parent(); p != nil {
			p.Redraw()
		}
	}
}
