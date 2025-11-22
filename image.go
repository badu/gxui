// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type ScalingMode int

const (
	Scaling1to1 ScalingMode = iota
	ScalingExpandGreedy
	ScalingExplicitSize
)

type AspectMode int

const (
	AspectStretch = iota
	AspectCorrectLetterbox
	AspectCorrectCrop
)

type Image struct {
	ControlBase
	BackgroundBorderPainter
	parent       ControlBaseParent
	texture      Texture
	canvas       Canvas
	explicitSize math.Size
	scalingMode  ScalingMode
	aspectMode   AspectMode
}

func (i *Image) calculateDrawRect() math.Rect {
	rect := i.parent.Size().Rect()
	texW, texH := i.texture.Size().WH()
	aspectSrc := float32(texH) / float32(texW)
	aspectDst := float32(rect.H()) / float32(rect.W())
	switch i.aspectMode {
	case AspectCorrectLetterbox, AspectCorrectCrop:
		if (aspectDst < aspectSrc) != (i.aspectMode == AspectCorrectLetterbox) {
			contract := rect.H() - int(float32(rect.W())*aspectSrc)
			rect = rect.Contract(math.Spacing{T: contract / 2, B: contract / 2})
		} else {
			contract := rect.W() - int(float32(rect.H())/aspectSrc)
			rect = rect.Contract(math.Spacing{L: contract / 2, R: contract / 2})
		}
	default:
		//
	}
	return rect
}

func (i *Image) Init(parent ControlBaseParent, driver Driver) {
	i.parent = parent
	i.ControlBase.Init(parent, driver)
	i.BackgroundBorderPainter.Init(parent)
	i.SetBorderPen(TransparentPen)
	i.SetBackgroundBrush(TransparentBrush)
}

func (i *Image) Texture() Texture {
	return i.texture
}

func (i *Image) SetTexture(texture Texture) {
	if i.texture == texture {
		return
	}

	i.texture = texture
	i.canvas = nil
	i.parent.ReLayout()
}

func (i *Image) Canvas() Canvas {
	return i.canvas
}

func (i *Image) SetCanvas(canvas Canvas) {
	if !canvas.IsComplete() {
		panic("SetCanvas() called with an incomplete canvas")
	}

	if i.canvas == canvas {
		return
	}

	i.canvas = canvas
	i.texture = nil
	i.parent.ReLayout()
}

func (i *Image) ScalingMode() ScalingMode {
	return i.scalingMode
}

func (i *Image) SetScalingMode(mode ScalingMode) {
	if i.scalingMode == mode {
		return
	}

	i.scalingMode = mode
	i.parent.ReLayout()
}

func (i *Image) AspectMode() AspectMode {
	return i.aspectMode
}

func (i *Image) SetAspectMode(mode AspectMode) {
	if i.aspectMode == mode {
		return
	}

	i.aspectMode = mode
	i.parent.Redraw()
}

func (i *Image) SetExplicitSize(explicitSize math.Size) {
	if i.explicitSize != explicitSize {
		i.explicitSize = explicitSize
		i.parent.ReLayout()
	}
	i.SetScalingMode(ScalingExplicitSize)
}

func (i *Image) PixelAt(point math.Point) (math.Point, bool) {
	rect := i.calculateDrawRect()
	if tex := i.Texture(); tex != nil {
		size := tex.SizePixels()
		point = point.Sub(rect.Min).
			ScaleX(float32(size.W) / float32(rect.W())).
			ScaleY(float32(size.H) / float32(rect.H()))
		if size.Rect().Contains(point) {
			return point, true
		}
	}
	return math.Point{X: -1, Y: -1}, false
}

func (i *Image) DesiredSize(min, max math.Size) math.Size {
	size := max
	switch i.scalingMode {
	case ScalingExplicitSize:
		size = i.explicitSize
	case Scaling1to1:
		switch {
		case i.texture != nil:
			size = i.texture.Size()
		case i.canvas != nil:
			size = i.canvas.Size()
		}
	}
	return size.Expand(math.CreateSpacing(int(i.BorderPen().Width))).Clamp(min, max)
}

func (i *Image) Paint(canvas Canvas) {
	rect := i.parent.Size().Rect()
	i.PaintBackground(canvas, rect)
	switch {
	case i.texture != nil:
		canvas.DrawTexture(i.texture, i.calculateDrawRect())
	case i.canvas != nil:
		canvas.DrawCanvas(i.canvas, math.ZeroPoint)
	}
	i.PaintBorder(canvas, rect)
}
