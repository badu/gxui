// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/mixins/base"
	"github.com/badu/gxui/mixins/parts"
)

type DropDownListOuter interface {
	base.ContainerOuter
}

type DropDownList struct {
	base.Container
	parts.BackgroundBorderPainter
	parts.Focusable
	outer       DropDownListOuter
	theme       gxui.Theme
	list        gxui.List
	listShowing bool
	itemSize    math.Size
	overlay     gxui.BubbleOverlay
	selected    *gxui.Child
	onShowList  gxui.Event
	onHideList  gxui.Event
}

func (l *DropDownList) Init(outer DropDownListOuter, theme gxui.Theme) {
	l.outer = outer
	l.Container.Init(outer, theme)
	l.BackgroundBorderPainter.Init(outer)
	l.Focusable.Init(outer)

	l.theme = theme
	l.list = theme.CreateList()
	l.list.OnSelectionChanged(
		func(item gxui.AdapterItem) {
			l.outer.RemoveAll()
			adapter := l.list.Adapter()
			if item != nil && adapter != nil {
				l.selected = l.AddChild(adapter.Create(l.theme, adapter.ItemIndex(item)))
			} else {
				l.selected = nil
			}
			l.Relayout()
		},
	)

	l.list.OnItemClicked(
		func(gxui.MouseEvent, gxui.AdapterItem) {
			l.HideList()
		},
	)
	l.list.OnKeyPress(
		func(ev gxui.KeyboardEvent) {
			switch ev.Key {
			case gxui.KeyEnter, gxui.KeyEscape:
				l.HideList()
			}
		},
	)

	l.list.OnLostFocus(l.HideList)
	l.OnDetach(l.HideList)
	l.SetMouseEventTarget(true)
}

func (l *DropDownList) LayoutChildren() {
	if !l.RelayoutSuspended() {
		// Disable relayout on AddChild / RemoveChild as we're performing layout here.
		l.SetRelayoutSuspended(true)
		defer l.SetRelayoutSuspended(false)
	}

	if l.selected != nil {
		size := l.outer.Size().Contract(l.Padding()).Max(math.ZeroSize)
		offset := l.Padding().LT()
		l.selected.Layout(size.Rect().Offset(offset))
	}
}

func (l *DropDownList) DesiredSize(min, max math.Size) math.Size {
	if l.selected != nil {
		return l.selected.Control.DesiredSize(min, max).Expand(l.outer.Padding()).Clamp(min, max)
	} else {
		return l.itemSize.Expand(l.outer.Padding()).Clamp(min, max)
	}
}

func (l *DropDownList) DataReplaced() {
	adapter := l.list.Adapter()
	itemSize := adapter.Size(l.theme)
	l.itemSize = itemSize
	l.outer.Relayout()
}

func (l *DropDownList) ListShowing() bool {
	return l.listShowing
}

func (l *DropDownList) ShowList() bool {
	if l.listShowing || l.overlay == nil {
		return false
	}

	l.listShowing = true
	size := l.Size()
	at := math.Point{X: size.W / 2, Y: size.H}
	l.overlay.Show(l.list, gxui.TransformCoordinate(at, l.outer, l.overlay))

	gxui.SetFocus(l.list)

	if l.onShowList != nil {
		l.onShowList.Fire()
	}

	return true
}

func (l *DropDownList) HideList() {
	if l.listShowing {
		l.listShowing = false
		l.overlay.Hide()

		if l.Attached() {
			gxui.SetFocus(l)
		}

		if l.onHideList != nil {
			l.onHideList.Fire()
		}
	}
}

func (l *DropDownList) List() gxui.List {
	return l.list
}

// InputEventHandler override
func (l *DropDownList) Click(ev gxui.MouseEvent) bool {
	l.InputEventHandler.Click(ev)
	if l.ListShowing() {
		l.HideList()
	} else {
		l.ShowList()
	}
	return true
}

// gxui.DropDownList compliance
func (l *DropDownList) SetBubbleOverlay(overlay gxui.BubbleOverlay) {
	l.overlay = overlay
}

func (l *DropDownList) BubbleOverlay() gxui.BubbleOverlay {
	return l.overlay
}

func (l *DropDownList) Adapter() gxui.ListAdapter {
	return l.list.Adapter()
}

func (l *DropDownList) SetAdapter(adapter gxui.ListAdapter) {
	if l.list.Adapter() != adapter {
		l.list.SetAdapter(adapter)
		if adapter != nil {
			adapter.OnDataChanged(func(bool) { l.DataReplaced() })
			adapter.OnDataReplaced(l.DataReplaced)
		}
		// TODO: Unlisten
		l.DataReplaced()
	}
}

func (l *DropDownList) Selected() gxui.AdapterItem {
	return l.list.Selected()
}

func (l *DropDownList) Select(item gxui.AdapterItem) {
	if l.list.Selected() != item {
		l.list.Select(item)
		l.LayoutChildren()
	}
}

func (l *DropDownList) OnSelectionChanged(callback func(gxui.AdapterItem)) gxui.EventSubscription {
	return l.list.OnSelectionChanged(callback)
}

func (l *DropDownList) OnShowList(callback func()) gxui.EventSubscription {
	if l.onShowList == nil {
		l.onShowList = gxui.CreateEvent(callback)
	}
	return l.onShowList.Listen(callback)
}

func (l *DropDownList) OnHideList(callback func()) gxui.EventSubscription {
	if l.onHideList == nil {
		l.onHideList = gxui.CreateEvent(callback)
	}
	return l.onHideList.Listen(callback)
}

// InputEventHandler overrides
func (l *DropDownList) KeyPress(event gxui.KeyboardEvent) (consume bool) {
	if event.Key == gxui.KeySpace || event.Key == gxui.KeyEnter {
		mouseEvent := gxui.MouseEvent{
			Button: gxui.MouseButtonLeft,
		}
		return l.Click(mouseEvent)
	}
	return l.InputEventHandler.KeyPress(event)
}

// parts.Container overrides
func (l *DropDownList) Paint(canvas gxui.Canvas) {
	rect := l.outer.Size().Rect()
	l.PaintBackground(canvas, rect)
	l.Container.Paint(canvas)
	l.PaintBorder(canvas, rect)
}
