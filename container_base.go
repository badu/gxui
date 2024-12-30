// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

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

func (c *ContainerBase) Init(outer ParentBaseContainer, theme App) {
	c.AttachablePart.Init()
	c.ContainerPart.Init(outer)
	c.DrawPaintPart.Init(outer, theme)
	c.InputEventHandlerPart.Init()
	c.LayoutablePart.Init(outer, theme)
	c.PaddablePart.Init(outer)
	c.PaintChildrenPart.Init(outer)
	c.ParentablePart.Init()
	c.VisiblePart.Init(outer)
}
