// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/mixins/base"
	"github.com/badu/gxui/mixins/parts"
)

type LinearLayoutOuter interface {
	base.ContainerOuter
}

type LinearLayout struct {
	base.Container
	parts.LinearLayout
	parts.BackgroundBorderPainter
}

func (l *LinearLayout) Init(outer LinearLayoutOuter, theme gxui.Theme) {
	l.Container.Init(outer, theme)
	l.LinearLayout.Init(outer)
	l.BackgroundBorderPainter.Init(outer)
	l.SetMouseEventTarget(true)
	l.SetBackgroundBrush(gxui.TransparentBrush)
	l.SetBorderPen(gxui.TransparentPen)
}

func (l *LinearLayout) Paint(canvas gxui.Canvas) {
	rect := l.Size().Rect()
	l.BackgroundBorderPainter.PaintBackground(canvas, rect)
	l.PaintChildren.Paint(canvas)
	l.BackgroundBorderPainter.PaintBorder(canvas, rect)
}
