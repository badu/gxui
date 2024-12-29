// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/mixins/parts"
)

type ContainerNoControlOuter interface {
	gxui.Container
	PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) // was outer.PaintChilder
	Paint(canvas gxui.Canvas)                                  // was outer.Painter
	LayoutChildren()                                           // was outer.LayoutChildren
}

type ContainerOuter interface {
	gxui.Container
	gxui.Control
	PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) // was outer.PaintChilder
	Paint(canvas gxui.Canvas)                                  // was outer.Painter
	LayoutChildren()                                           // was outer.LayoutChildren
}

type Container struct {
	parts.Attachable
	parts.Container
	parts.DrawPaint
	parts.InputEventHandler
	parts.Layoutable
	parts.Paddable
	parts.PaintChildren
	parts.Parentable
	parts.Visible
}

func (c *Container) Init(outer ContainerOuter, theme gxui.Theme) {
	c.Attachable.Init(outer)
	c.Container.Init(outer)
	c.DrawPaint.Init(outer, theme)
	c.InputEventHandler.Init(outer)
	c.Layoutable.Init(outer, theme)
	c.Paddable.Init(outer)
	c.PaintChildren.Init(outer)
	c.Parentable.Init(outer)
	c.Visible.Init(outer)
}
