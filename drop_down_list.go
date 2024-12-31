// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

type DropDownList interface {
	Focusable
	Container
	SetBubbleOverlay(BubbleOverlay)
	BubbleOverlay() BubbleOverlay
	Adapter() ListAdapter
	SetAdapter(ListAdapter)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	Selected() AdapterItem
	Select(AdapterItem)
	OnSelectionChanged(callback func(AdapterItem)) EventSubscription
	OnShowList(callback func()) EventSubscription
	OnHideList(callback func()) EventSubscription
}

type DropDownListImpl struct {
	ContainerBase
	BackgroundBorderPainter
	FocusablePart
	parent      BaseContainerParent
	styles      *StyleDefs
	list        List
	listShowing bool
	itemSize    math.Size
	overlay     BubbleOverlay
	selected    *Child
	onShowList  Event
	onHideList  Event
}

func (l *DropDownListImpl) Init(parent BaseContainerParent, driver Driver, styles *StyleDefs) {
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

func (l *DropDownListImpl) LayoutChildren() {
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

func (l *DropDownListImpl) DesiredSize(min, max math.Size) math.Size {
	if l.selected != nil {
		return l.selected.Control.DesiredSize(min, max).Expand(l.parent.Padding()).Clamp(min, max)
	} else {
		return l.itemSize.Expand(l.parent.Padding()).Clamp(min, max)
	}
}

func (l *DropDownListImpl) DataReplaced() {
	adapter := l.list.Adapter()
	itemSize := adapter.Size(l.styles)
	l.itemSize = itemSize
	l.parent.ReLayout()
}

func (l *DropDownListImpl) ListShowing() bool {
	return l.listShowing
}

func (l *DropDownListImpl) ShowList() bool {
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

func (l *DropDownListImpl) HideList() {
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

func (l *DropDownListImpl) List() List {
	return l.list
}

// InputEventHandlerPart override
func (l *DropDownListImpl) Click(ev MouseEvent) bool {
	l.InputEventHandlerPart.Click(ev)
	if l.ListShowing() {
		l.HideList()
	} else {
		l.ShowList()
	}
	return true
}

// gxui.DropDownList compliance
func (l *DropDownListImpl) SetBubbleOverlay(overlay BubbleOverlay) {
	l.overlay = overlay
}

func (l *DropDownListImpl) BubbleOverlay() BubbleOverlay {
	return l.overlay
}

func (l *DropDownListImpl) Adapter() ListAdapter {
	return l.list.Adapter()
}

func (l *DropDownListImpl) SetAdapter(adapter ListAdapter) {
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

func (l *DropDownListImpl) Selected() AdapterItem {
	return l.list.Selected()
}

func (l *DropDownListImpl) Select(item AdapterItem) {
	if l.list.Selected() == item {
		return
	}
	l.list.Select(item)
	l.LayoutChildren()
}

func (l *DropDownListImpl) OnSelectionChanged(callback func(AdapterItem)) EventSubscription {
	return l.list.OnSelectionChanged(callback)
}

func (l *DropDownListImpl) OnShowList(callback func()) EventSubscription {
	if l.onShowList == nil {
		l.onShowList = CreateEvent(callback)
	}

	return l.onShowList.Listen(callback)
}

func (l *DropDownListImpl) OnHideList(callback func()) EventSubscription {
	if l.onHideList == nil {
		l.onHideList = CreateEvent(callback)
	}

	return l.onHideList.Listen(callback)
}

// InputEventHandlerPart overrides
func (l *DropDownListImpl) KeyPress(event KeyboardEvent) (consume bool) {
	if event.Key == KeySpace || event.Key == KeyEnter {
		mouseEvent := MouseEvent{
			Button: MouseButtonLeft,
		}
		return l.Click(mouseEvent)
	}

	return l.InputEventHandlerPart.KeyPress(event)
}

// parts.ContainerPart overrides
func (l *DropDownListImpl) Paint(canvas Canvas) {
	rect := l.parent.Size().Rect()
	l.PaintBackground(canvas, rect)
	l.ContainerBase.Paint(canvas)
	l.PaintBorder(canvas, rect)
}
