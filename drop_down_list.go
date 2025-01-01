// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

type DropDownList struct {
	ContainerBase
	BackgroundBorderPainter
	FocusablePart
	parent      BaseContainerParent
	styles      *StyleDefs
	list        *ListImpl
	listShowing bool
	itemSize    math.Size
	overlay     *BubbleOverlay
	selected    *Child
	onShowList  Event
	onHideList  Event
}

func (l *DropDownList) Init(parent BaseContainerParent, driver Driver, styles *StyleDefs) {
	l.parent = parent
	l.styles = styles
	l.ContainerBase.Init(parent, driver)
	l.BackgroundBorderPainter.Init(parent)
	l.FocusablePart.Init()

	l.list = CreateList(driver, styles)
	l.list.OnSelectionChanged(
		func(item AdapterItem) {
			l.parent.RemoveAll()
			adapter := l.list.Adapter()
			if item != nil && adapter != nil {
				l.selected = l.AddChild(adapter.Create(driver, styles, adapter.ItemIndex(item)))
			} else {
				l.selected = nil
			}
			l.ReLayout()
		},
	)

	l.list.OnItemClicked(
		func(event MouseEvent, item AdapterItem) {
			l.HideList()
		},
	)

	l.list.OnKeyPress(
		func(event KeyboardEvent) {
			switch event.Key {
			case KeyEnter, KeyEscape:
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
		size := l.parent.Size().Contract(l.Padding()).Max(math.ZeroSize)
		offset := l.Padding().LT()
		l.selected.Layout(size.Rect().Offset(offset))
	}
}

func (l *DropDownList) DesiredSize(min, max math.Size) math.Size {
	if l.selected != nil {
		return l.selected.Control.DesiredSize(min, max).Expand(l.parent.Padding()).Clamp(min, max)
	} else {
		return l.itemSize.Expand(l.parent.Padding()).Clamp(min, max)
	}
}

func (l *DropDownList) DataReplaced() {
	adapter := l.list.Adapter()
	itemSize := adapter.Size(l.styles)
	l.itemSize = itemSize
	l.parent.ReLayout()
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
	l.overlay.Show(l.list, TransformCoordinate(at, l.parent, l.overlay))

	SetFocus(l.list)

	if l.onShowList != nil {
		l.onShowList.Emit()
	}

	return true
}

func (l *DropDownList) HideList() {
	if !l.listShowing {
		return
	}

	l.listShowing = false
	l.overlay.Hide()

	if l.Attached() {
		SetFocus(l)
	}

	if l.onHideList != nil {
		l.onHideList.Emit()
	}
}

func (l *DropDownList) List() *ListImpl {
	return l.list
}

// InputEventHandlerPart override
func (l *DropDownList) Click(ev MouseEvent) bool {
	l.InputEventHandlerPart.Click(ev)
	if l.ListShowing() {
		l.HideList()
	} else {
		l.ShowList()
	}
	return true
}

func (l *DropDownList) SetBubbleOverlay(overlay *BubbleOverlay) {
	l.overlay = overlay
}

func (l *DropDownList) BubbleOverlay() *BubbleOverlay {
	return l.overlay
}

func (l *DropDownList) Adapter() ListAdapter {
	return l.list.Adapter()
}

func (l *DropDownList) SetAdapter(adapter ListAdapter) {
	if l.list.Adapter() == adapter {
		return
	}
	l.list.SetAdapter(adapter)
	if adapter != nil {
		adapter.OnDataChanged(func(bool) { l.DataReplaced() })
		adapter.OnDataReplaced(l.DataReplaced)
	}
	// TODO: Event.Forget()
	l.DataReplaced()
}

func (l *DropDownList) Selected() AdapterItem {
	return l.list.Selected()
}

func (l *DropDownList) Select(item AdapterItem) {
	if l.list.Selected() == item {
		return
	}
	l.list.Select(item)
	l.LayoutChildren()
}

func (l *DropDownList) OnSelectionChanged(callback func(AdapterItem)) EventSubscription {
	return l.list.OnSelectionChanged(callback)
}

func (l *DropDownList) OnShowList(callback func()) EventSubscription {
	if l.onShowList == nil {
		l.onShowList = CreateEvent(callback)
	}

	return l.onShowList.Listen(callback)
}

func (l *DropDownList) OnHideList(callback func()) EventSubscription {
	if l.onHideList == nil {
		l.onHideList = CreateEvent(callback)
	}

	return l.onHideList.Listen(callback)
}

// InputEventHandlerPart overrides
func (l *DropDownList) KeyPress(event KeyboardEvent) (consume bool) {
	if event.Key == KeySpace || event.Key == KeyEnter {
		mouseEvent := MouseEvent{
			Button: MouseButtonLeft,
		}
		return l.Click(mouseEvent)
	}

	return l.InputEventHandlerPart.KeyPress(event)
}

// parts.ContainerPart overrides
func (l *DropDownList) Paint(canvas Canvas) {
	rect := l.parent.Size().Rect()
	l.PaintBackground(canvas, rect)
	l.ContainerBase.Paint(canvas)
	l.PaintBorder(canvas, rect)
}
