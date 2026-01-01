// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type PaintChildrenParent interface {
	// Container
	Children() Children
	PaintChild(canvas Canvas, child *Child, idx int)
}

type PaintChildrenPart struct {
	parent PaintChildrenParent
}

func (p *PaintChildrenPart) Init(parent PaintChildrenParent) {
	p.parent = parent
}

func (p *PaintChildrenPart) Paint(canvas Canvas) {
	for i, v := range p.parent.Children() {
		if !v.Control.IsVisible() {
			continue
		}

		canvas.Push()
		canvas.AddClip(v.Control.Size().Rect().Offset(v.Offset))
		p.parent.PaintChild(canvas, v, i)
		canvas.Pop()
	}
}

func (p *PaintChildrenPart) PaintChild(canvas Canvas, child *Child, idx int) {
	childCanvas := child.Control.Draw()

	if childCanvas != nil {
		canvas.DrawCanvas(childCanvas, child.Offset)
	}
}
