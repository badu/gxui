// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

type ButtonType int

const (
	PushButton ButtonType = iota
	ToggleButton
)

type ButtonParent interface {
	BaseContainerParent
	IsChecked() bool
	SetChecked(bool)
}

type Button struct {
	FocusablePart
	LinearLayoutImpl
	parent        ButtonParent
	canvasCreator CanvasCreator
	styles        *StyleDefs
	label         *Label
	buttonType    ButtonType
	checked       bool
}

func (b *Button) Init(parent ButtonParent, canvasCreator CanvasCreator, styles *StyleDefs) {
	b.LinearLayoutImpl.Init(parent, canvasCreator)
	b.FocusablePart.Init()

	b.buttonType = PushButton
	b.canvasCreator = canvasCreator
	b.styles = styles
	b.parent = parent

	b.SetBackgroundBrush(styles.ButtonDefaultStyle.Brush)
	b.SetBorderPen(styles.ButtonDefaultStyle.Pen)
	// TODO : @Badu - use styles
	b.SetPadding(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	b.SetMargin(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})

}

func (b *Button) Label() *Label {
	return b.label
}

func (b *Button) Text() string {
	if b.label == nil {
		return ""
	}

	return b.label.Text
}

func (b *Button) SetText(text string) {
	if b.Text() == text {
		return
	}

	if text == "" {
		if b.label != nil {
			b.RemoveChild(b.label)
			b.label = nil
		}
		return
	}

	if b.label == nil {
		b.label = CreateLabel(b.canvasCreator, b.styles)
		b.label.SetMargin(math.ZeroSpacing)
		b.AddChild(b.label)
	}
	b.label.SetText(text)
}

func (b *Button) Type() ButtonType {
	return b.buttonType
}

func (b *Button) SetType(buttonType ButtonType) {
	if buttonType == b.buttonType {
		return
	}

	b.buttonType = buttonType
	b.parent.Redraw()
}

func (b *Button) IsChecked() bool {
	return b.checked
}

func (b *Button) SetChecked(checked bool) {
	if checked == b.checked {
		return
	}

	b.checked = checked
	b.parent.Redraw()
}

// InputEventHandlerPart override
func (b *Button) Click(event MouseEvent) bool {
	if event.Button == MouseButtonLeft {
		if b.buttonType == ToggleButton {
			b.parent.SetChecked(!b.parent.IsChecked())
		}
		b.LinearLayoutImpl.Click(event)
		return true
	}

	return b.LinearLayoutImpl.Click(event)
}

func (b *Button) KeyPress(event KeyboardEvent) bool {
	consume := b.LinearLayoutImpl.KeyPress(event)
	if event.Key == KeySpace || event.Key == KeyEnter {
		return b.Click(MouseEvent{Button: MouseButtonLeft})
	}
	return consume
}

// Button internal overrides
func (b *Button) Paint(canvas Canvas) {
	pen := b.BorderPen()
	brush := b.BackgroundBrush()
	fontColor := b.styles.ButtonDefaultStyle.FontColor

	switch {
	case b.IsMouseDown(MouseButtonLeft) && b.IsMouseOver():
		pen = b.styles.ButtonPressedStyle.Pen
		brush = b.styles.ButtonPressedStyle.Brush
		fontColor = b.styles.ButtonPressedStyle.FontColor
	case b.IsMouseOver():
		pen = b.styles.ButtonOverStyle.Pen
		brush = b.styles.ButtonOverStyle.Brush
		fontColor = b.styles.ButtonOverStyle.FontColor
	}

	if label := b.Label(); label != nil {
		label.SetColor(fontColor)
	}

	rect := b.Size().Rect()

	canvas.DrawRoundedRect(rect, 2, 2, 2, 2, TransparentPen, brush)

	b.PaintChildrenPart.Paint(canvas)

	canvas.DrawRoundedRect(rect, 2, 2, 2, 2, pen, TransparentBrush)

	if b.IsChecked() {
		pen = b.styles.HighlightStyle.Pen
		brush = b.styles.HighlightStyle.Brush
		canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, pen, brush)
	}

	if b.HasFocus() {
		pen = b.styles.FocusedStyle.Pen
		brush = b.styles.FocusedStyle.Brush
		canvas.DrawRoundedRect(rect.ContractI(int(pen.Width)), 3.0, 3.0, 3.0, 3.0, pen, brush)
	}
}
