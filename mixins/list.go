// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type itemDetails struct {
	child               *gxui.Child
	index               int
	mark                int
	onClickSubscription gxui.EventSubscription
}

type List struct {
	ContainerBase
	BackgroundBorderPainter
	FocusablePart
	outer                    gxui.ListOuter
	theme                    gxui.Theme
	adapter                  gxui.ListAdapter
	scrollBar                gxui.ScrollBar
	scrollBarChild           *gxui.Child
	scrollBarEnabled         bool
	selectedItem             gxui.AdapterItem
	onSelectionChanged       gxui.Event
	details                  map[gxui.AdapterItem]itemDetails
	orientation              gxui.Orientation
	scrollOffset             int
	itemSize                 math.Size
	itemCount                int // Count number of items in the adapter
	layoutMark               int
	mousePosition            math.Point
	itemMouseOver            *gxui.Child
	onItemClicked            gxui.Event
	dataChangedSubscription  gxui.EventSubscription
	dataReplacedSubscription gxui.EventSubscription
}

func (l *List) Init(outer gxui.ListOuter, theme gxui.Theme) {
	l.outer = outer
	l.ContainerBase.Init(outer, theme)
	l.BackgroundBorderPainter.Init(outer)
	l.FocusablePart.Init(outer)

	l.theme = theme
	l.scrollBar = theme.CreateScrollBar()
	l.scrollBarChild = l.AddChild(l.scrollBar)
	l.scrollBarEnabled = true
	l.scrollBar.OnScroll(func(from, to int) { l.SetScrollOffset(from) })

	l.SetOrientation(gxui.Vertical)
	l.SetBackgroundBrush(gxui.TransparentBrush)
	l.SetMouseEventTarget(true)

	l.details = make(map[gxui.AdapterItem]itemDetails)
}

func (l *List) UpdateItemMouseOver() {
	if !l.IsMouseOver() {
		if l.itemMouseOver != nil {
			l.itemMouseOver = nil
			l.Redraw()
		}
		return
	}

	for _, detail := range l.details {
		if detail.child.Bounds().Contains(l.mousePosition) {
			if l.itemMouseOver != detail.child {
				l.itemMouseOver = detail.child
				l.Redraw()
				return
			}
		}
	}
}

func (l *List) LayoutChildren() {
	if l.adapter == nil {
		l.outer.RemoveAll()
		return
	}

	if !l.RelayoutSuspended() {
		// Disable relayout on AddChild / RemoveChild as we're performing layout here.
		l.SetRelayoutSuspended(true)
		defer l.SetRelayoutSuspended(false)
	}

	size := l.outer.Size().Contract(l.Padding())
	offset := l.Padding().LT()

	var itemSize math.Size
	if l.orientation.Horizontal() {
		itemSize = math.Size{W: l.itemSize.W, H: size.H}
	} else {
		itemSize = math.Size{W: size.W, H: l.itemSize.H}
	}

	startIndex, endIndex := l.VisibleItemRange(true)
	majorAxisItemSize := l.MajorAxisItemSize()

	displacement := startIndex*majorAxisItemSize - l.scrollOffset

	mark := l.layoutMark
	l.layoutMark++

	for idx := startIndex; idx < endIndex; idx++ {
		item := l.adapter.ItemAt(idx)

		details, found := l.details[item]
		if found {
			if details.mark == mark {
				panic(fmt.Errorf("adapter for control '%s' returned duplicate item (%v) for indices %v and %v", gxui.Path(l.outer), item, details.index, idx))
			}
		} else {
			control := l.adapter.Create(l.theme, idx)
			details.onClickSubscription = control.OnClick(
				func(ev gxui.MouseEvent) {
					l.ItemClicked(ev, item)
				},
			)
			details.child = l.AddChildAt(0, control)
		}

		details.mark = mark
		details.index = idx

		l.details[item] = details

		child := details.child
		childMargin := child.Control.Margin()
		childSize := itemSize.Contract(childMargin).Max(math.ZeroSize)

		if l.orientation.Horizontal() {
			child.Layout(math.CreateRect(displacement, childMargin.T, displacement+childSize.W, childMargin.T+childSize.H).Offset(offset))
		} else {
			child.Layout(math.CreateRect(childMargin.L, displacement, childMargin.L+childSize.W, displacement+childSize.H).Offset(offset))
		}

		displacement += majorAxisItemSize
	}

	// Reap unused items
	for item, detail := range l.details {
		if detail.mark != mark {
			detail.onClickSubscription.Unlisten()
			l.RemoveChild(detail.child.Control)
			delete(l.details, item)
		}
	}

	if l.scrollBarEnabled {
		ss := l.scrollBar.DesiredSize(math.ZeroSize, size)
		if l.Orientation().Horizontal() {
			l.scrollBarChild.Layout(math.CreateRect(0, size.H-ss.H, size.W, size.H).Canon().Offset(offset))
		} else {
			l.scrollBarChild.Layout(math.CreateRect(size.W-ss.W, 0, size.W, size.H).Canon().Offset(offset))
		}

		// Only show the scroll bar if needed
		entireContentVisible := startIndex == 0 && endIndex == l.itemCount
		l.scrollBar.SetVisible(!entireContentVisible)
	}

	l.UpdateItemMouseOver()
}

func (l *List) SetSize(size math.Size) {
	l.LayoutablePart.SetSize(size)
	// Ensure scroll offset is still valid
	l.SetScrollOffset(l.scrollOffset)
}

func (l *List) DesiredSize(min, max math.Size) math.Size {
	if l.adapter == nil {
		return min
	}

	count := math.Max(l.itemCount, 1)

	var size math.Size
	if l.orientation.Horizontal() {
		size = math.Size{W: l.itemSize.W * count, H: l.itemSize.H}
	} else {
		size = math.Size{W: l.itemSize.W, H: l.itemSize.H * count}
	}

	if l.scrollBarEnabled {
		if l.orientation.Horizontal() {
			size.H += l.scrollBar.DesiredSize(min, max).H
		} else {
			size.W += l.scrollBar.DesiredSize(min, max).W
		}
	}

	return size.Expand(l.outer.Padding()).Clamp(min, max)
}

func (l *List) ScrollBarEnabled(bool) bool {
	return l.scrollBarEnabled
}

func (l *List) SetScrollBarEnabled(enabled bool) {
	if l.scrollBarEnabled != enabled {
		l.scrollBarEnabled = enabled
		l.Relayout()
	}
}

func (l *List) SetScrollOffset(scrollOffset int) {
	if l.adapter == nil {
		return
	}

	size := l.outer.Size().Contract(l.outer.Padding())

	if l.orientation.Horizontal() {
		maxScroll := math.Max(l.itemSize.W*l.itemCount-size.W, 0)
		scrollOffset = math.Clamp(scrollOffset, 0, maxScroll)
		l.scrollBar.SetScrollPosition(scrollOffset, scrollOffset+size.W)
	} else {
		maxScroll := math.Max(l.itemSize.H*l.itemCount-size.H, 0)
		scrollOffset = math.Clamp(scrollOffset, 0, maxScroll)
		l.scrollBar.SetScrollPosition(scrollOffset, scrollOffset+size.H)
	}

	if l.scrollOffset != scrollOffset {
		l.scrollOffset = scrollOffset
		l.LayoutChildren()
	}
}

func (l *List) MajorAxisItemSize() int {
	return l.orientation.Major(l.itemSize.WH())
}

func (l *List) VisibleItemRange(includePartiallyVisible bool) (int, int) {
	if l.itemCount == 0 {
		return 0, 0
	}

	size := l.outer.Size()
	padding := l.outer.Padding()
	majorAxisItemSize := l.MajorAxisItemSize()
	if majorAxisItemSize == 0 {
		return 0, 0
	}

	startIndex := l.scrollOffset
	if !includePartiallyVisible {
		startIndex += majorAxisItemSize - 1
	}

	endIndex := 0
	if l.orientation.Horizontal() {
		endIndex = l.scrollOffset + size.W - padding.W()
	} else {
		endIndex = l.scrollOffset + size.H - padding.H()
	}

	if includePartiallyVisible {
		endIndex += majorAxisItemSize - 1
	}

	startIndex = math.Max(startIndex/majorAxisItemSize, 0)
	endIndex = math.Min(endIndex/majorAxisItemSize, l.itemCount)

	return startIndex, endIndex
}

func (l *List) SizeChanged() {
	l.itemSize = l.adapter.Size(l.theme)
	l.scrollBar.SetScrollLimit(l.itemCount * l.MajorAxisItemSize())
	l.SetScrollOffset(l.scrollOffset)
	l.outer.Relayout()
}

func (l *List) DataChanged(recreateControls bool) {
	if recreateControls {
		for item, details := range l.details {
			details.onClickSubscription.Unlisten()
			l.RemoveChild(details.child.Control)
			delete(l.details, item)
		}
	}

	l.itemCount = l.adapter.Count()
	l.SizeChanged()
}

func (l *List) DataReplaced() {
	l.selectedItem = nil
	l.DataChanged(true)
}

func (l *List) Paint(canvas gxui.Canvas) {
	rect := l.outer.Size().Rect()
	l.outer.PaintBackground(canvas, rect)
	l.ContainerBase.Paint(canvas)
	l.outer.PaintBorder(canvas, rect)
}

func (l *List) PaintSelection(canvas gxui.Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, gxui.WhitePen, gxui.TransparentBrush)
}

func (l *List) PaintMouseOverBackground(canvas gxui.Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, gxui.TransparentPen, gxui.CreateBrush(gxui.Gray90))
}

func (l *List) SelectPrevious() {
	if l.selectedItem != nil {
		selectedIndex := l.adapter.ItemIndex(l.selectedItem)
		l.Select(l.adapter.ItemAt(math.Mod(selectedIndex-1, l.itemCount)))
	} else {
		l.Select(l.adapter.ItemAt(0))
	}
}

func (l *List) SelectNext() {
	if l.selectedItem != nil {
		selectedIndex := l.adapter.ItemIndex(l.selectedItem)
		l.Select(l.adapter.ItemAt(math.Mod(selectedIndex+1, l.itemCount)))
	} else {
		l.Select(l.adapter.ItemAt(0))
	}
}

func (l *List) ContainsItem(item gxui.AdapterItem) bool {
	return l.adapter != nil && l.adapter.ItemIndex(item) >= 0
}

func (l *List) RemoveAll() {
	for _, details := range l.details {
		details.onClickSubscription.Unlisten()
		l.outer.RemoveChild(details.child.Control)
	}
	l.details = make(map[gxui.AdapterItem]itemDetails)
}

// PaintChildrenPart overrides
func (l *List) PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) {
	if child == l.itemMouseOver {
		b := child.Bounds().Expand(child.Control.Margin())
		l.outer.PaintMouseOverBackground(canvas, b)
	}
	l.PaintChildrenPart.PaintChild(canvas, child, idx)
	if selected, found := l.details[l.selectedItem]; found {
		if child == selected.child {
			b := child.Bounds().Expand(child.Control.Margin())
			l.outer.PaintSelection(canvas, b)
		}
	}
}

// InputEventHandlerPart override
func (l *List) MouseMove(event gxui.MouseEvent) {
	l.InputEventHandlerPart.MouseMove(event)
	l.mousePosition = event.Point
	l.UpdateItemMouseOver()
}

func (l *List) MouseExit(event gxui.MouseEvent) {
	l.InputEventHandlerPart.MouseExit(event)
	l.itemMouseOver = nil
}

func (l *List) MouseScroll(event gxui.MouseEvent) bool {
	if event.ScrollY == 0 {
		return l.InputEventHandlerPart.MouseScroll(event)
	}

	prevOffset := l.scrollOffset
	if l.orientation.Horizontal() {
		delta := event.ScrollY * l.itemSize.W / 8
		l.SetScrollOffset(l.scrollOffset - delta)
	} else {
		delta := event.ScrollY * l.itemSize.H / 8
		l.SetScrollOffset(l.scrollOffset - delta)
	}

	return prevOffset != l.scrollOffset
}

func (l *List) KeyPress(event gxui.KeyboardEvent) bool {
	if l.itemCount > 0 {
		if l.orientation.Horizontal() {
			switch event.Key {
			case gxui.KeyLeft:
				l.SelectPrevious()
				return true
			case gxui.KeyRight:
				l.SelectNext()
				return true
			case gxui.KeyPageUp:
				l.SetScrollOffset(l.scrollOffset - l.Size().W)
				return true
			case gxui.KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Size().W)
				return true
			}
		} else {
			switch event.Key {
			case gxui.KeyUp:
				l.SelectPrevious()
				return true
			case gxui.KeyDown:
				l.SelectNext()
				return true
			case gxui.KeyPageUp:
				l.SetScrollOffset(l.scrollOffset - l.Size().H)
				return true
			case gxui.KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Size().H)
				return true
			}
		}
	}
	return l.ContainerBase.KeyPress(event)
}

// gxui.List compliance
func (l *List) Adapter() gxui.ListAdapter {
	return l.adapter
}

func (l *List) SetAdapter(adapter gxui.ListAdapter) {
	if l.adapter != adapter {
		if l.adapter != nil {
			l.dataChangedSubscription.Unlisten()
			l.dataReplacedSubscription.Unlisten()
		}
		l.adapter = adapter
		if l.adapter != nil {
			l.dataChangedSubscription = l.adapter.OnDataChanged(l.DataChanged)
			l.dataReplacedSubscription = l.adapter.OnDataReplaced(l.DataReplaced)
		}
		l.DataReplaced()
	}
}

func (l *List) Orientation() gxui.Orientation {
	return l.orientation
}

func (l *List) SetOrientation(o gxui.Orientation) {
	l.scrollBar.SetOrientation(o)
	if l.orientation != o {
		l.orientation = o
		l.Relayout()
	}
}

func (l *List) ScrollTo(item gxui.AdapterItem) {
	idx := l.adapter.ItemIndex(item)
	startIndex, endIndex := l.VisibleItemRange(false)
	if idx < startIndex {
		if l.Orientation().Horizontal() {
			l.SetScrollOffset(l.itemSize.W * idx)
		} else {
			l.SetScrollOffset(l.itemSize.H * idx)
		}
	} else if idx >= endIndex {
		count := endIndex - startIndex
		if l.Orientation().Horizontal() {
			l.SetScrollOffset(l.itemSize.W * (idx - count + 1))
		} else {
			l.SetScrollOffset(l.itemSize.H * (idx - count + 1))
		}
	}
}

func (l *List) IsItemVisible(item gxui.AdapterItem) bool {
	_, found := l.details[item]
	return found
}

func (l *List) ItemControl(item gxui.AdapterItem) gxui.Control {
	if item, found := l.details[item]; found {
		return item.child.Control
	}
	return nil
}

func (l *List) ItemClicked(ev gxui.MouseEvent, item gxui.AdapterItem) {
	if l.onItemClicked != nil {
		l.onItemClicked.Fire(ev, item)
	}
	l.Select(item)
}

func (l *List) OnItemClicked(f func(gxui.MouseEvent, gxui.AdapterItem)) gxui.EventSubscription {
	if l.onItemClicked == nil {
		l.onItemClicked = gxui.CreateEvent(f)
	}
	return l.onItemClicked.Listen(f)
}

func (l *List) Selected() gxui.AdapterItem {
	return l.selectedItem
}

func (l *List) Select(item gxui.AdapterItem) bool {
	if l.selectedItem != item {
		if !l.outer.ContainsItem(item) {
			return false
		}
		l.selectedItem = item
		if l.onSelectionChanged != nil {
			l.onSelectionChanged.Fire(item)
		}
		l.Redraw()
	}
	l.ScrollTo(item)
	return true
}

func (l *List) OnSelectionChanged(f func(gxui.AdapterItem)) gxui.EventSubscription {
	if l.onItemClicked == nil {
		l.onSelectionChanged = gxui.CreateEvent(f)
	}
	return l.onSelectionChanged.Listen(f)
}
