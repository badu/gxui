// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"github.com/badu/gxui/math"
)

type List interface {
	Focusable
	Parent
	Adapter() ListAdapter
	SetAdapter(ListAdapter)
	SetOrientation(Orientation)
	Orientation() Orientation
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
	ScrollTo(AdapterItem)
	IsItemVisible(AdapterItem) bool
	ItemControl(AdapterItem) Control
	Selected() AdapterItem
	Select(AdapterItem) bool
	OnSelectionChanged(func(AdapterItem)) EventSubscription
	OnItemClicked(func(MouseEvent, AdapterItem)) EventSubscription
}

type ListParent interface {
	BaseContainerParent
	ContainsItem(AdapterItem) bool
	PaintBackground(c Canvas, r math.Rect)
	PaintMouseOverBackground(c Canvas, r math.Rect)
	PaintSelection(c Canvas, r math.Rect)
	PaintBorder(c Canvas, r math.Rect)
}

// ListAdapter is an interface used to visualize a flat set of items.
// Users of the ListAdapter should presume the data is unchanged until the
// OnDataChanged or OnDataReplaced events are fired.
type ListAdapter interface {
	// Count returns the total number of items.
	Count() int

	// ItemAt returns the AdapterItem for the item at index i. It is important
	// for the Adapter to return consistent AdapterItems for the same data, so
	// that selections can be persisted, or re-ordering animations can be played
	// when the dataset changes.
	// The AdapterItem returned must be equality-unique across all indices.
	ItemAt(index int) AdapterItem

	// ItemIndex returns the index of item, or -1 if the adapter does not contain
	// item.
	ItemIndex(item AdapterItem) int

	// Create returns a Control visualizing the item at the specified index.
	Create(driver Driver, styles *StyleDefs, index int) Control

	// Size returns the size that each of the item's controls will be displayed
	// at for the given theme.
	Size(styles *StyleDefs) math.Size

	// OnDataChanged registers f to be called when there is a partial change in
	// the items of the adapter. Scroll positions and selections should be
	// preserved if possible.
	// If recreateControls is true then each of the visible controls should be
	// recreated by re-calling Create().
	OnDataChanged(f func(recreateControls bool)) EventSubscription

	// OnDataReplaced registers f to be called when there is a complete
	// replacement of items in the adapter.
	OnDataReplaced(f func()) EventSubscription
}

type itemDetails struct {
	child               *Child
	index               int
	mark                int
	onClickSubscription EventSubscription
}

type ListImpl struct {
	ContainerBase
	BackgroundBorderPainter
	FocusablePart
	parent                   ListParent
	driver                   Driver
	styles                   *StyleDefs
	adapter                  ListAdapter
	scrollBar                ScrollBar
	scrollBarChild           *Child
	scrollBarEnabled         bool
	selectedItem             AdapterItem
	onSelectionChanged       Event
	details                  map[AdapterItem]itemDetails
	orientation              Orientation
	scrollOffset             int
	itemSize                 math.Size
	itemCount                int // Count number of items in the adapter
	layoutMark               int
	mousePosition            math.Point
	itemMouseOver            *Child
	onItemClicked            Event
	dataChangedSubscription  EventSubscription
	dataReplacedSubscription EventSubscription
}

func (l *ListImpl) Init(parent ListParent, driver Driver, styles *StyleDefs) {
	l.parent = parent
	l.ContainerBase.Init(parent, driver)
	l.BackgroundBorderPainter.Init(parent)
	l.FocusablePart.Init()
	l.driver = driver
	l.styles = styles
	l.scrollBar = CreateScrollBar(driver, styles)
	l.scrollBarChild = l.AddChild(l.scrollBar)
	l.scrollBarEnabled = true
	l.scrollBar.OnScroll(func(from, to int) { l.SetScrollOffset(from) })

	l.SetOrientation(Vertical)
	l.SetBackgroundBrush(TransparentBrush)
	l.SetMouseEventTarget(true)

	l.details = make(map[AdapterItem]itemDetails)
}

func (l *ListImpl) UpdateItemMouseOver() {
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

func (l *ListImpl) LayoutChildren() {
	if l.adapter == nil {
		l.parent.RemoveAll()
		return
	}

	if !l.RelayoutSuspended() {
		// Disable relayout on AddChild / RemoveChild as we're performing layout here.
		l.SetRelayoutSuspended(true)
		defer l.SetRelayoutSuspended(false)
	}

	size := l.parent.Size().Contract(l.Padding())
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
				panic(fmt.Errorf("adapter for control '%s' returned duplicate item (%v) for indices %v and %v", Path(l.parent), item, details.index, idx))
			}
		} else {
			control := l.adapter.Create(l.driver, l.styles, idx)
			details.onClickSubscription = control.OnClick(
				func(ev MouseEvent) {
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
			detail.onClickSubscription.Forget()
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

func (l *ListImpl) SetSize(size math.Size) {
	l.LayoutablePart.SetSize(size)
	// Ensure scroll offset is still valid
	l.SetScrollOffset(l.scrollOffset)
}

func (l *ListImpl) DesiredSize(min, max math.Size) math.Size {
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

	return size.Expand(l.parent.Padding()).Clamp(min, max)
}

func (l *ListImpl) ScrollBarEnabled(bool) bool {
	return l.scrollBarEnabled
}

func (l *ListImpl) SetScrollBarEnabled(enabled bool) {
	if l.scrollBarEnabled != enabled {
		l.scrollBarEnabled = enabled
		l.Relayout()
	}
}

func (l *ListImpl) SetScrollOffset(scrollOffset int) {
	if l.adapter == nil {
		return
	}

	size := l.parent.Size().Contract(l.parent.Padding())

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

func (l *ListImpl) MajorAxisItemSize() int {
	return l.orientation.Major(l.itemSize.WH())
}

func (l *ListImpl) VisibleItemRange(includePartiallyVisible bool) (int, int) {
	if l.itemCount == 0 {
		return 0, 0
	}

	size := l.parent.Size()
	padding := l.parent.Padding()
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

func (l *ListImpl) SizeChanged() {
	l.itemSize = l.adapter.Size(l.styles)
	l.scrollBar.SetScrollLimit(l.itemCount * l.MajorAxisItemSize())
	l.SetScrollOffset(l.scrollOffset)
	l.parent.Relayout()
}

func (l *ListImpl) DataChanged(recreateControls bool) {
	if recreateControls {
		for item, details := range l.details {
			details.onClickSubscription.Forget()
			l.RemoveChild(details.child.Control)
			delete(l.details, item)
		}
	}

	l.itemCount = l.adapter.Count()
	l.SizeChanged()
}

func (l *ListImpl) DataReplaced() {
	l.selectedItem = nil
	l.DataChanged(true)
}

func (l *ListImpl) Paint(canvas Canvas) {
	rect := l.parent.Size().Rect()
	l.parent.PaintBackground(canvas, rect)
	l.ContainerBase.Paint(canvas)
	l.parent.PaintBorder(canvas, rect)
}

func (l *ListImpl) PaintSelection(canvas Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, WhitePen, TransparentBrush)
}

func (l *ListImpl) PaintMouseOverBackground(canvas Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, TransparentPen, CreateBrush(Gray90))
}

func (l *ListImpl) SelectPrevious() {
	if l.selectedItem != nil {
		selectedIndex := l.adapter.ItemIndex(l.selectedItem)
		l.Select(l.adapter.ItemAt(math.Mod(selectedIndex-1, l.itemCount)))
	} else {
		l.Select(l.adapter.ItemAt(0))
	}
}

func (l *ListImpl) SelectNext() {
	if l.selectedItem != nil {
		selectedIndex := l.adapter.ItemIndex(l.selectedItem)
		l.Select(l.adapter.ItemAt(math.Mod(selectedIndex+1, l.itemCount)))
	} else {
		l.Select(l.adapter.ItemAt(0))
	}
}

func (l *ListImpl) ContainsItem(item AdapterItem) bool {
	return l.adapter != nil && l.adapter.ItemIndex(item) >= 0
}

func (l *ListImpl) RemoveAll() {
	for _, details := range l.details {
		details.onClickSubscription.Forget()
		l.parent.RemoveChild(details.child.Control)
	}
	l.details = make(map[AdapterItem]itemDetails)
}

// PaintChildrenPart overrides
func (l *ListImpl) PaintChild(canvas Canvas, child *Child, idx int) {
	if child == l.itemMouseOver {
		b := child.Bounds().Expand(child.Control.Margin())
		l.parent.PaintMouseOverBackground(canvas, b)
	}
	l.PaintChildrenPart.PaintChild(canvas, child, idx)
	if selected, found := l.details[l.selectedItem]; found {
		if child == selected.child {
			b := child.Bounds().Expand(child.Control.Margin())
			l.parent.PaintSelection(canvas, b)
		}
	}
}

// InputEventHandlerPart override
func (l *ListImpl) MouseMove(event MouseEvent) {
	l.InputEventHandlerPart.MouseMove(event)
	l.mousePosition = event.Point
	l.UpdateItemMouseOver()
}

func (l *ListImpl) MouseExit(event MouseEvent) {
	l.InputEventHandlerPart.MouseExit(event)
	l.itemMouseOver = nil
}

func (l *ListImpl) MouseScroll(event MouseEvent) bool {
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

func (l *ListImpl) KeyPress(event KeyboardEvent) bool {
	if l.itemCount > 0 {
		if l.orientation.Horizontal() {
			switch event.Key {
			case KeyLeft:
				l.SelectPrevious()
				return true
			case KeyRight:
				l.SelectNext()
				return true
			case KeyPageUp:
				l.SetScrollOffset(l.scrollOffset - l.Size().W)
				return true
			case KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Size().W)
				return true
			}
		} else {
			switch event.Key {
			case KeyUp:
				l.SelectPrevious()
				return true
			case KeyDown:
				l.SelectNext()
				return true
			case KeyPageUp:
				l.SetScrollOffset(l.scrollOffset - l.Size().H)
				return true
			case KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Size().H)
				return true
			}
		}
	}
	return l.ContainerBase.KeyPress(event)
}

// gxui.List compliance
func (l *ListImpl) Adapter() ListAdapter {
	return l.adapter
}

func (l *ListImpl) SetAdapter(adapter ListAdapter) {
	if l.adapter != adapter {
		if l.adapter != nil {
			l.dataChangedSubscription.Forget()
			l.dataReplacedSubscription.Forget()
		}
		l.adapter = adapter
		if l.adapter != nil {
			l.dataChangedSubscription = l.adapter.OnDataChanged(l.DataChanged)
			l.dataReplacedSubscription = l.adapter.OnDataReplaced(l.DataReplaced)
		}
		l.DataReplaced()
	}
}

func (l *ListImpl) Orientation() Orientation {
	return l.orientation
}

func (l *ListImpl) SetOrientation(o Orientation) {
	l.scrollBar.SetOrientation(o)
	if l.orientation != o {
		l.orientation = o
		l.Relayout()
	}
}

func (l *ListImpl) ScrollTo(item AdapterItem) {
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

func (l *ListImpl) IsItemVisible(item AdapterItem) bool {
	_, found := l.details[item]
	return found
}

func (l *ListImpl) ItemControl(item AdapterItem) Control {
	if item, found := l.details[item]; found {
		return item.child.Control
	}
	return nil
}

func (l *ListImpl) ItemClicked(ev MouseEvent, item AdapterItem) {
	if l.onItemClicked != nil {
		l.onItemClicked.Fire(ev, item)
	}
	l.Select(item)
}

func (l *ListImpl) OnItemClicked(f func(MouseEvent, AdapterItem)) EventSubscription {
	if l.onItemClicked == nil {
		l.onItemClicked = CreateEvent(f)
	}
	return l.onItemClicked.Listen(f)
}

func (l *ListImpl) Selected() AdapterItem {
	return l.selectedItem
}

func (l *ListImpl) Select(item AdapterItem) bool {
	if l.selectedItem != item {
		if !l.parent.ContainsItem(item) {
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

func (l *ListImpl) OnSelectionChanged(f func(AdapterItem)) EventSubscription {
	if l.onItemClicked == nil {
		l.onSelectionChanged = CreateEvent(f)
	}
	return l.onSelectionChanged.Listen(f)
}
