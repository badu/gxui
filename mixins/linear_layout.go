// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type LinearLayout struct {
	ContainerBase
	LinearLayoutPart
	BackgroundBorderPainter
}

func (l *LinearLayout) Init(outer ContainerBaseOuter, theme gxui.Theme) {
	l.ContainerBase.Init(outer, theme)
	l.LinearLayoutPart.Init(outer)
	l.BackgroundBorderPainter.Init(outer)
	l.SetMouseEventTarget(true)
	l.SetBackgroundBrush(gxui.TransparentBrush)
	l.SetBorderPen(gxui.TransparentPen)
}

func (l *LinearLayout) Paint(canvas gxui.Canvas) {
	rect := l.Size().Rect()
	l.BackgroundBorderPainter.PaintBackground(canvas, rect)
	l.PaintChildrenPart.Paint(canvas)
	l.BackgroundBorderPainter.PaintBorder(canvas, rect)
}
