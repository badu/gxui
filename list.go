// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"

	"github.com/badu/gxui/pkg/math"
)

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

	// ItemAt returns the AdapterItem for the item at index i.
	// It is important for the Adapter to return consistent AdapterItems for the same data, so that selections can be persisted, or re-ordering animations can be played when the dataset changes.
	// The AdapterItem returned must be equality-unique across all indices.
	ItemAt(index int) AdapterItem

	// ItemIndex returns the index of item, or -1 if the adapter does not contain item.
	ItemIndex(item AdapterItem) int

	// Create returns a Control visualizing the item at the specified index.
	Create(driver Driver, styles *StyleDefs, index int) Control

	// Size returns the size that each of the item's controls will be displayed at for the given theme.
	Size(styles *StyleDefs) math.Size

	// OnDataChanged registers f to be called when there is a partial change in the items of the adapter.
	// Scroll positions and selections should be preserved if possible.
	// If recreateControls is true then each of the visible controls should be recreated by re-calling Create().
	OnDataChanged(callback func(recreateControls bool)) EventSubscription

	// OnDataReplaced registers f to be called when there is a complete replacement of items in the adapter.
	OnDataReplaced(callback func()) EventSubscription
}

type itemDetails struct {
	onClickSubscription EventSubscription
	child               *Child
	index               int
	mark                int
}

type ListImpl struct {
	ContainerBase
	FocusablePart
	BackgroundBorderPainter
	parent                   ListParent
	driver                   Driver
	adapter                  ListAdapter
	onSelectionChanged       Event
	onItemClicked            Event
	dataChangedSubscription  EventSubscription
	dataReplacedSubscription EventSubscription
	selectedItem             AdapterItem
	scrollBar                *ScrollBarImpl
	styles                   *StyleDefs
	scrollBarChild           *Child
	itemMouseOver            *Child
	details                  map[AdapterItem]itemDetails
	itemSize                 math.Size
	mousePosition            math.Point
	orientation              Orientation
	scrollOffset             int
	itemCount                int // Count number of items in the adapter
	layoutMark               int
	hiddenItemCount          int
	scrollBarEnabled         bool
}

func (l *ListImpl) Init(parent ListParent, driver Driver, styles *StyleDefs) {
	l.parent = parent
	l.driver = driver
	l.styles = styles

	l.ContainerBase.Init(parent, driver)
	l.BackgroundBorderPainter.Init(parent)
	l.FocusablePart.Init()

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
		// Disable reLayout on AddChild / RemoveChild as we're performing layout here.
		l.SetRelayoutSuspended(true)
		defer l.SetRelayoutSuspended(false)
	}

	size := l.parent.Size().Contract(l.Padding())
	offset := l.Padding().TopLeft()

	var itemSize math.Size
	if l.orientation.Horizontal() {
		itemSize = math.Size{Width: l.itemSize.Width, Height: size.Height}
	} else {
		itemSize = math.Size{Width: size.Width, Height: l.itemSize.Height}
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
			child.Layout(math.CreateRect(displacement, childMargin.Top, displacement+childSize.Width, childMargin.Top+childSize.Height).Offset(offset))
		} else {
			child.Layout(math.CreateRect(childMargin.Left, displacement, childMargin.Left+childSize.Width, displacement+childSize.Height).Offset(offset))
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
		scrollSize := l.scrollBar.DesiredSize(math.ZeroSize, size)
		if l.Orientation().Horizontal() {
			l.scrollBarChild.Layout(math.CreateRect(0, size.Height-scrollSize.Height, size.Width, size.Height).Canon().Offset(offset))
		} else {
			l.scrollBarChild.Layout(math.CreateRect(size.Width-scrollSize.Width, 0, size.Width, size.Height).Canon().Offset(offset))
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

func (l *ListImpl) DesiredSize(minSize, maxSize math.Size) math.Size {
	if l.adapter == nil {
		return minSize
	}

	count := max(l.itemCount, 1)

	var size math.Size
	if l.orientation.Horizontal() {
		size = math.Size{Width: l.itemSize.Width * count, Height: l.itemSize.Height}
	} else {
		size = math.Size{Width: l.itemSize.Width, Height: l.itemSize.Height * count}
	}

	if l.scrollBarEnabled {
		if l.orientation.Horizontal() {
			size.Height += l.scrollBar.DesiredSize(minSize, maxSize).Height
		} else {
			size.Width += l.scrollBar.DesiredSize(minSize, maxSize).Width
		}
	}

	return size.Expand(l.parent.Padding()).Clamp(minSize, maxSize)
}

func (l *ListImpl) ScrollBarEnabled(bool) bool {
	return l.scrollBarEnabled
}

func (l *ListImpl) SetScrollBarEnabled(enabled bool) {
	if l.scrollBarEnabled == enabled {
		return
	}

	l.scrollBarEnabled = enabled
	l.ReLayout()
}

func (l *ListImpl) SetScrollOffset(scrollOffset int) {
	if l.adapter == nil {
		return
	}

	size := l.parent.Size().Contract(l.parent.Padding())

	if l.orientation.Horizontal() {
		maxScroll := max(l.itemSize.Width*l.itemCount-size.Width, 0)
		scrollOffset = math.Clamp(scrollOffset, 0, maxScroll)
		l.scrollBar.SetScrollPosition(scrollOffset, scrollOffset+size.Width)
	} else {
		maxScroll := max(l.itemSize.Height*l.itemCount-size.Height, 0)
		scrollOffset = math.Clamp(scrollOffset, 0, maxScroll)
		l.scrollBar.SetScrollPosition(scrollOffset, scrollOffset+size.Height)
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
		endIndex = l.scrollOffset + size.Width - padding.Width()
	} else {
		endIndex = l.scrollOffset + size.Height - padding.Height()
	}

	if includePartiallyVisible {
		endIndex += majorAxisItemSize - 1
	}

	startIndex = max(startIndex/majorAxisItemSize, 0)
	endIndex = min(endIndex/majorAxisItemSize, l.itemCount)

	return startIndex, endIndex
}

func (l *ListImpl) SizeChanged() {
	l.itemSize = l.adapter.Size(l.styles)
	l.scrollBar.SetScrollLimit(l.itemCount * l.MajorAxisItemSize())
	l.SetScrollOffset(l.scrollOffset)
	l.parent.ReLayout()
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
	if l.HasFocus() {
		rect := l.Size().Rect().ContractI(1)
		canvas.DrawRoundedRect(rect, 3.0, 3.0, 3.0, 3.0, l.styles.FocusedStyle.Pen, l.styles.FocusedStyle.Brush)
	}
}

func (l *ListImpl) PaintSelection(canvas Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, l.styles.HighlightStyle.Pen, l.styles.HighlightStyle.Brush)
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
		bounds := child.Bounds().Expand(child.Control.Margin())
		l.parent.PaintMouseOverBackground(canvas, bounds)
	}
	l.PaintChildrenPart.PaintChild(canvas, child, idx)
	if selected, found := l.details[l.selectedItem]; found {
		if child == selected.child {
			bounds := child.Bounds().Expand(child.Control.Margin())
			l.parent.PaintSelection(canvas, bounds)
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
		delta := event.ScrollY * l.itemSize.Width / 8
		l.SetScrollOffset(l.scrollOffset - delta)
	} else {
		delta := event.ScrollY * l.itemSize.Height / 8
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
				l.SetScrollOffset(l.scrollOffset - l.Size().Width)
				return true

			case KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Size().Width)
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
				l.SetScrollOffset(l.scrollOffset - l.Size().Height)
				return true

			case KeyPageDown:
				l.SetScrollOffset(l.scrollOffset + l.Size().Height)
				return true

			}
		}
	}
	return l.ContainerBase.KeyPress(event)
}

func (l *ListImpl) Adapter() ListAdapter {
	return l.adapter
}

func (l *ListImpl) SetAdapter(adapter ListAdapter) {
	if l.adapter == adapter {
		return
	}

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

func (l *ListImpl) Orientation() Orientation {
	return l.orientation
}

func (l *ListImpl) SetOrientation(o Orientation) {
	l.scrollBar.SetOrientation(o)
	if l.orientation == o {
		return
	}

	l.orientation = o
	l.ReLayout()
}

func (l *ListImpl) ScrollTo(item AdapterItem) {
	idx := l.adapter.ItemIndex(item)
	startIndex, endIndex := l.VisibleItemRange(false)
	if idx < startIndex {
		if l.Orientation().Horizontal() {
			l.SetScrollOffset(l.itemSize.Width * idx)
		} else {
			l.SetScrollOffset(l.itemSize.Height * idx)
		}
	} else if idx >= endIndex {
		count := endIndex - startIndex
		if l.Orientation().Horizontal() {
			l.SetScrollOffset(l.itemSize.Width * (idx - count + 1))
		} else {
			l.SetScrollOffset(l.itemSize.Height * (idx - count + 1))
		}
	}
}

func (l *ListImpl) IsItemVisible(item AdapterItem) bool {
	_, found := l.details[item]
	return found
}

func (l *ListImpl) ItemControl(item AdapterItem) Control {
	if control, found := l.details[item]; found {
		return control.child.Control
	}
	return nil
}

func (l *ListImpl) ItemClicked(event MouseEvent, item AdapterItem) {
	if l.onItemClicked != nil {
		l.onItemClicked.Emit(event, item)
	}
	l.Select(item)
}

func (l *ListImpl) OnItemClicked(callback func(event MouseEvent, item AdapterItem)) EventSubscription {
	if l.onItemClicked == nil {
		l.onItemClicked = CreateEvent(callback)
	}

	return l.onItemClicked.Listen(callback)
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
			l.onSelectionChanged.Emit(item)
		}

		l.Redraw()
	}

	l.ScrollTo(item)
	return true
}

func (l *ListImpl) OnSelectionChanged(callback func(item AdapterItem)) EventSubscription {
	if l.onItemClicked == nil {
		l.onSelectionChanged = CreateEvent(callback)
	}

	return l.onSelectionChanged.Listen(callback)
}

func (l *ListImpl) ChangeHiddenCount(value int) {
	l.hiddenItemCount += value
}
