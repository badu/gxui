// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
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
	LinearLayoutImpl
	FocusablePart
	parent     ButtonParent
	driver     Driver
	styles     *StyleDefs
	label      *Label
	buttonType ButtonType
	checked    bool
}

func (b *Button) Init(parent ButtonParent, driver Driver, styles *StyleDefs) {
	b.LinearLayoutImpl.Init(parent, driver)
	b.FocusablePart.Init()

	b.buttonType = PushButton
	b.driver = driver
	b.styles = styles
	b.parent = parent
}

func (b *Button) Label() *Label {
	return b.label
}

func (b *Button) Text() string {
	if b.label == nil {
		return ""
	}

	return b.label.Text()
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
		b.label = CreateLabel(b.driver, b.styles)
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
