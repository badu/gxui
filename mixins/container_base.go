// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type ContainerBase struct {
	AttachablePart
	ContainerPart
	DrawPaintPart
	InputEventHandlerPart
	LayoutablePart
	PaddablePart
	PaintChildrenPart
	ParentablePart
	VisiblePart
}

func (c *ContainerBase) Init(outer gxui.ContainerBaseOuter, theme gxui.Theme) {
	c.AttachablePart.Init(outer)
	c.ContainerPart.Init(outer)
	c.DrawPaintPart.Init(outer, theme)
	c.InputEventHandlerPart.Init(outer)
	c.LayoutablePart.Init(outer, theme)
	c.PaddablePart.Init(outer)
	c.PaintChildrenPart.Init(outer)
	c.ParentablePart.Init(outer)
	c.VisiblePart.Init(outer)
}
