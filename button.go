// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

type ButtonType int

const (
	PushButton ButtonType = iota
	ToggleButton
)

type Button interface {
	LinearLayout
	Text() string
	SetText(string)
	Type() ButtonType
	SetType(ButtonType)
	IsChecked() bool
	SetChecked(bool)
}

type ButtonParent interface {
	BaseContainerParent
	IsChecked() bool
	SetChecked(bool)
}

type ButtonImpl struct {
	LinearLayoutImpl
	FocusablePart
	parent     ButtonParent
	app        App
	label      Label
	buttonType ButtonType
	checked    bool
}

func (b *ButtonImpl) Init(parent ButtonParent, app App) {
	b.LinearLayoutImpl.Init(parent, app)
	b.FocusablePart.Init()

	b.buttonType = PushButton
	b.app = app
	b.parent = parent
}

func (b *ButtonImpl) Label() Label {
	return b.label
}

func (b *ButtonImpl) Text() string {
	if b.label == nil {
		return ""
	}

	return b.label.Text()
}

func (b *ButtonImpl) SetText(text string) {
	if b.Text() == text {
		return
	}

	if text == "" {
		if b.label != nil {
			b.RemoveChild(b.label)
			b.label = nil
		}
	} else {
		if b.label == nil {
			b.label = b.app.CreateLabel()
			b.label.SetMargin(math.ZeroSpacing)
			b.AddChild(b.label)
		}
		b.label.SetText(text)
	}
}

func (b *ButtonImpl) Type() ButtonType {
	return b.buttonType
}

func (b *ButtonImpl) SetType(buttonType ButtonType) {
	if buttonType != b.buttonType {
		b.buttonType = buttonType
		b.parent.Redraw()
	}
}

func (b *ButtonImpl) IsChecked() bool {
	return b.checked
}

func (b *ButtonImpl) SetChecked(checked bool) {
	if checked != b.checked {
		b.checked = checked
		b.parent.Redraw()
	}
}

// InputEventHandlerPart override
func (b *ButtonImpl) Click(event MouseEvent) (consume bool) {
	if event.Button == MouseButtonLeft {
		if b.buttonType == ToggleButton {
			b.parent.SetChecked(!b.parent.IsChecked())
		}
		b.LinearLayoutImpl.Click(event)
		return true
	}
	return b.LinearLayoutImpl.Click(event)
}

func (b *ButtonImpl) KeyPress(event KeyboardEvent) (consume bool) {
	consume = b.LinearLayoutImpl.KeyPress(event)
	if event.Key == KeySpace || event.Key == KeyEnter {
		me := MouseEvent{
			Button: MouseButtonLeft,
		}
		return b.Click(me)
	}
	return
}
