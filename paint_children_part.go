// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type PaintChildrenParent interface {
	Container
	PaintChild(canvas Canvas, child *Child, idx int) // was outer.PaintChilder
	Size() math.Size                                 // was outer.Sized
	SetSize(newSize math.Size)                       // was outer.Sized
}

type PaintChildrenPart struct {
	parent PaintChildrenParent
}

func (p *PaintChildrenPart) Init(parent PaintChildrenParent) {
	p.parent = parent
}

func (p *PaintChildrenPart) Paint(canvas Canvas) {
	for i, v := range p.parent.Children() {
		if v.Control.IsVisible() {
			canvas.Push()
			canvas.AddClip(v.Control.Size().Rect().Offset(v.Offset))
			p.parent.PaintChild(canvas, v, i)
			canvas.Pop()
		}
	}
}

func (p *PaintChildrenPart) PaintChild(canvas Canvas, child *Child, idx int) {
	if childCanvas := child.Control.Draw(); childCanvas != nil {
		canvas.DrawCanvas(childCanvas, child.Offset)
	}
}
