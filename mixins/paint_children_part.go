// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type PaintChildrenOuter interface {
	gxui.Container
	PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) // was outer.PaintChilder
	Size() math.Size                                           // was outer.Sized
	SetSize(newSize math.Size)                                 // was outer.Sized
}

type PaintChildrenPart struct {
	outer PaintChildrenOuter
}

func (p *PaintChildrenPart) Init(outer PaintChildrenOuter) {
	p.outer = outer
}

func (p *PaintChildrenPart) Paint(canvas gxui.Canvas) {
	for i, v := range p.outer.Children() {
		if v.Control.IsVisible() {
			canvas.Push()
			canvas.AddClip(v.Control.Size().Rect().Offset(v.Offset))
			p.outer.PaintChild(canvas, v, i)
			canvas.Pop()
		}
	}
}

func (p *PaintChildrenPart) PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) {
	if childCanvas := child.Control.Draw(); childCanvas != nil {
		canvas.DrawCanvas(childCanvas, child.Offset)
	}
}
