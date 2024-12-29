// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/mixins/base"
)

type ButtonOuter interface {
	base.ContainerBaseOuter
	IsChecked() bool
	SetChecked(bool)
}

type Button struct {
	LinearLayout
	base.FocusablePart
	outer      ButtonOuter
	theme      gxui.Theme
	label      gxui.Label
	buttonType gxui.ButtonType
	checked    bool
}

func (b *Button) Init(outer ButtonOuter, theme gxui.Theme) {
	b.LinearLayout.Init(outer, theme)
	b.FocusablePart.Init(outer)

	b.buttonType = gxui.PushButton
	b.theme = theme
	b.outer = outer
}

func (b *Button) Label() gxui.Label {
	return b.label
}

func (b *Button) Text() string {
	if b.label != nil {
		return b.label.Text()
	} else {
		return ""
	}
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
	} else {
		if b.label == nil {
			b.label = b.theme.CreateLabel()
			b.label.SetMargin(math.ZeroSpacing)
			b.AddChild(b.label)
		}
		b.label.SetText(text)
	}
}

func (b *Button) Type() gxui.ButtonType {
	return b.buttonType
}

func (b *Button) SetType(buttonType gxui.ButtonType) {
	if buttonType != b.buttonType {
		b.buttonType = buttonType
		b.outer.Redraw()
	}
}

func (b *Button) IsChecked() bool {
	return b.checked
}

func (b *Button) SetChecked(checked bool) {
	if checked != b.checked {
		b.checked = checked
		b.outer.Redraw()
	}
}

// InputEventHandlerPart override
func (b *Button) Click(event gxui.MouseEvent) (consume bool) {
	if event.Button == gxui.MouseButtonLeft {
		if b.buttonType == gxui.ToggleButton {
			b.outer.SetChecked(!b.outer.IsChecked())
		}
		b.LinearLayout.Click(event)
		return true
	}
	return b.LinearLayout.Click(event)
}

func (b *Button) KeyPress(event gxui.KeyboardEvent) (consume bool) {
	consume = b.LinearLayout.KeyPress(event)
	if event.Key == gxui.KeySpace || event.Key == gxui.KeyEnter {
		me := gxui.MouseEvent{
			Button: gxui.MouseButtonLeft,
		}
		return b.Click(me)
	}
	return
}
