// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/badu/gxui"
)

type ContainerBaseNoControlOuter interface {
	gxui.Container
	PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) // was outer.PaintChilder
	Paint(canvas gxui.Canvas)                                  // was outer.Painter
	LayoutChildren()                                           // was outer.LayoutChildren
}

type ContainerBaseOuter interface {
	gxui.Container
	gxui.Control
	PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) // was outer.PaintChilder
	Paint(canvas gxui.Canvas)                                  // was outer.Painter
	LayoutChildren()                                           // was outer.LayoutChildren
}

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

func (c *ContainerBase) Init(outer ContainerBaseOuter, theme gxui.Theme) {
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
