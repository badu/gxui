// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type WindowParent interface {
	Container

	// Title returns the title of the window.
	// This is usually the text displayed at the top of the window.
	Title() string

	// SetTitle changes the title of the window.
	SetTitle(string)

	// Scale returns the display scaling for this window.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	Scale() float32

	// SetScale alters the display scaling for this window.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	SetScale(float32)

	// Position returns position of the window.
	Position() math.Point

	// SetPosition changes position of the window.
	SetPosition(math.Point)

	// Fullscreen returns true if the window is currently full-screen.
	Fullscreen() bool

	// SetFullscreen makes the window either full-screen or windowed.
	SetFullscreen(bool)

	// Show makes the window visible.
	Show()

	// Hide makes the window invisible.
	Hide()

	// Close destroys the window.
	// Once the window is closed, no further calls should be made to it.
	Close()

	// Focus returns the control currently with focus.
	Focus() Focusable

	// SetFocus gives the specified control Focus, returning true on success or
	// false if the control cannot be given focus.
	SetFocus(Control) bool

	// BackgroundBrush returns the brush used to draw the window background.
	BackgroundBrush() Brush

	// SetBackgroundBrush sets the brush used to draw the window background.
	SetBackgroundBrush(Brush)

	// BorderPen returns the pen used to draw the window border.
	BorderPen() Pen

	// SetBorderPen sets the pen used to draw the window border.
	SetBorderPen(Pen)

	Click(MouseEvent)
	DoubleClick(MouseEvent)
	KeyPress(KeyboardEvent)
	KeyStroke(KeyStrokeEvent)

	// Events
	OnClose(func()) EventSubscription
	OnResize(func()) EventSubscription
	OnClick(func(MouseEvent)) EventSubscription
	OnDoubleClick(func(MouseEvent)) EventSubscription
	OnMouseMove(func(MouseEvent)) EventSubscription
	OnMouseEnter(func(MouseEvent)) EventSubscription
	OnMouseExit(func(MouseEvent)) EventSubscription
	OnMouseDown(func(MouseEvent)) EventSubscription
	OnMouseUp(func(MouseEvent)) EventSubscription
	OnMouseScroll(func(MouseEvent)) EventSubscription
	OnKeyDown(func(KeyboardEvent)) EventSubscription
	OnKeyUp(func(KeyboardEvent)) EventSubscription
	OnKeyRepeat(func(KeyboardEvent)) EventSubscription
	OnKeyStroke(func(KeyStrokeEvent)) EventSubscription

	Attached() bool                                  // was outer.Attachable
	Attach()                                         // was outer.Attachable
	Detach()                                         // was outer.Attachable
	OnAttach(callback func()) EventSubscription      // was outer.Attachable
	OnDetach(callback func()) EventSubscription      // was outer.Attachable
	IsVisible() bool                                 // was outer.IsVisibler
	LayoutChildren()                                 // was outer.LayoutChildren
	PaintChild(canvas Canvas, child *Child, idx int) // was outer.PaintChilder
	Paint(canvas Canvas)                             // was outer.Painter
	Parent() Parent                                  // was outer.Parenter
	Size() math.Size                                 // was outer.Sized
	SetSize(newSize math.Size)                       // was outer.Sized
}

type WindowImpl struct {
	PaintChildrenPart
	driver             Driver
	parent             WindowParent
	viewport           Viewport
	onClose            Event // Raised by viewport
	onResize           Event // Raised by viewport
	onMouseMove        Event // Raised by viewport
	onMouseEnter       Event // Raised by viewport
	onMouseExit        Event // Raised by viewport
	onMouseDown        Event // Raised by viewport
	onMouseUp          Event // Raised by viewport
	onMouseScroll      Event // Raised by viewport
	onKeyDown          Event // Raised by viewport
	onKeyUp            Event // Raised by viewport
	onKeyRepeat        Event // Raised by viewport
	onKeyStroke        Event // Raised by viewport
	onClick            Event // Raised by MouseController
	onDoubleClick      Event // Raised by MouseController
	mouseController    *MouseController
	keyboardController *KeyboardController
	focusController    *FocusController
	AttachablePart
	viewportSubscriptions []EventSubscription
	ContainerPart
	PaddablePart
	BackgroundBorderPainter
	windowedSize  math.Size
	layoutPending bool
	drawPending   bool
	updatePending bool
}

func (w *WindowImpl) requestUpdate() {
	if !w.updatePending {
		w.updatePending = true
		w.driver.Call(w.update)
	}
}

func (w *WindowImpl) update() {
	if !w.Attached() {
		// WindowImpl was detached between requestUpdate() and update()
		w.updatePending = false
		w.layoutPending = false
		w.drawPending = false
		return
	}

	w.updatePending = false
	if w.layoutPending {
		w.layoutPending = false
		w.drawPending = true
		w.parent.LayoutChildren()
	}

	if w.drawPending {
		w.drawPending = false
		w.Draw()
	}
}

func (w *WindowImpl) Init(parent WindowParent, driver Driver, width, height int, title string) {
	w.BackgroundBorderPainter.Init(parent)
	w.ContainerPart.Init(parent)
	w.PaddablePart.Init(parent)
	w.PaintChildrenPart.Init(parent)

	w.parent = parent
	w.driver = driver

	w.onResize = CreateEvent(func() {})
	w.onMouseMove = CreateEvent(func(MouseEvent) {})
	w.onMouseEnter = CreateEvent(func(MouseEvent) {})
	w.onMouseExit = CreateEvent(func(MouseEvent) {})
	w.onMouseDown = CreateEvent(func(MouseEvent) {})
	w.onMouseUp = CreateEvent(func(MouseEvent) {})
	w.onMouseScroll = CreateEvent(func(MouseEvent) {})
	w.onKeyDown = CreateEvent(func(KeyboardEvent) {})
	w.onKeyUp = CreateEvent(func(KeyboardEvent) {})
	w.onKeyRepeat = CreateEvent(func(KeyboardEvent) {})
	w.onKeyStroke = CreateEvent(func(KeyStrokeEvent) {})

	w.onClick = CreateEvent(func(MouseEvent) {})
	w.onDoubleClick = CreateEvent(func(MouseEvent) {})

	switch window := parent.(type) {
	case *WindowImpl:
		w.focusController = CreateFocusController(window)
		w.mouseController = CreateMouseController(window, w.focusController)
		w.keyboardController = CreateKeyboardController(window)
	default:
		panic("parent should be *WindowImpl")
	}

	w.onResize.Listen(
		func() {
			w.parent.LayoutChildren()
			w.Draw()
		},
	)

	w.SetBorderPen(TransparentPen)

	w.setViewport(driver.CreateWindowedViewport(width, height, title))

	// TODO : @Badu - maybe this is not a good idea (window should show upon demand, since we might have loading to do)
	// WindowImpl starts shown
	w.Attach()
}

func (w *WindowImpl) Draw() Canvas {
	// TODO : the DrawPaintPart has similar functionality, except setting the canvas to the viewport - embed DrawPaintPart in Window
	if size := w.viewport.SizeDips(); size != math.ZeroSize {
		canvas := w.driver.CreateCanvas(size)
		w.parent.Paint(canvas)
		canvas.Complete()
		w.viewport.SetCanvas(canvas)
		return canvas
	} else {
		return nil
	}
}

func (w *WindowImpl) Paint(canvas Canvas) {
	w.PaintBackground(canvas, canvas.Size().Rect())
	w.PaintChildrenPart.Paint(canvas)
	w.PaintBorder(canvas, canvas.Size().Rect())
}

func (w *WindowImpl) LayoutChildren() {
	size := w.Size().Contract(w.Padding()).Max(math.ZeroSize)
	offset := w.Padding().LT()
	for _, child := range w.parent.Children() {
		child.Layout(child.Control.DesiredSize(math.ZeroSize, size).Rect().Offset(offset))
	}
}

func (w *WindowImpl) Size() math.Size {
	return w.viewport.SizeDips()
}

func (w *WindowImpl) SetSize(size math.Size) {
	w.viewport.SetSizeDips(size)
}

func (w *WindowImpl) Parent() Parent {
	return nil
}

func (w *WindowImpl) Viewport() Viewport {
	return w.viewport
}

func (w *WindowImpl) Title() string {
	return w.viewport.Title()
}

func (w *WindowImpl) SetTitle(title string) {
	w.viewport.SetTitle(title)
}

func (w *WindowImpl) Scale() float32 {
	return w.viewport.Scale()
}

func (w *WindowImpl) SetScale(scale float32) {
	w.viewport.SetScale(scale)
}

func (w *WindowImpl) Position() math.Point {
	return w.viewport.Position()
}

func (w *WindowImpl) SetPosition(point math.Point) {
	w.viewport.SetPosition(point)
}

func (w *WindowImpl) Fullscreen() bool {
	return w.viewport.Fullscreen()
}

func (w *WindowImpl) SetFullscreen(fullscreen bool) {
	title := w.viewport.Title()
	if fullscreen != w.Fullscreen() {
		old := w.viewport
		if fullscreen {
			w.windowedSize = old.SizeDips()
			w.setViewport(w.driver.CreateFullscreenViewport(0, 0, title))
		} else {
			width, height := w.windowedSize.WH()
			w.setViewport(w.driver.CreateWindowedViewport(width, height, title))
		}
		old.Close()
	}
}

func (w *WindowImpl) Show() {
	w.Attach()
	w.viewport.Show()
}

func (w *WindowImpl) Hide() {
	w.Detach()
	w.viewport.Hide()
}

func (w *WindowImpl) Close() {
	w.Detach()
	w.viewport.Close()
}

func (w *WindowImpl) Focus() Focusable {
	return w.focusController.Focus()
}

func (w *WindowImpl) SetFocus(control Control) bool {
	focusController := w.focusController
	if control == nil {
		focusController.SetFocus(nil)
		return true
	}

	if target := focusController.Focusable(control); target != nil {
		focusController.SetFocus(target)
		return true
	}

	return false
}

func (w *WindowImpl) IsVisible() bool {
	return true
}

func (w *WindowImpl) OnClose(callback func()) EventSubscription {
	if w.onClose == nil {
		w.onClose = CreateEvent(func() {})
	}

	return w.onClose.Listen(callback)
}

func (w *WindowImpl) OnResize(callback func()) EventSubscription {
	return w.onResize.Listen(callback)
}

func (w *WindowImpl) OnClick(callback func(MouseEvent)) EventSubscription {
	return w.onClick.Listen(callback)
}

func (w *WindowImpl) OnDoubleClick(callback func(MouseEvent)) EventSubscription {
	return w.onDoubleClick.Listen(callback)
}

func (w *WindowImpl) OnMouseMove(callback func(MouseEvent)) EventSubscription {
	return w.onMouseMove.Listen(
		func(ev MouseEvent) {
			ev.Window = w
			ev.WindowPoint = ev.Point
			callback(ev)
		},
	)
}

func (w *WindowImpl) OnMouseEnter(callback func(MouseEvent)) EventSubscription {
	return w.onMouseEnter.Listen(
		func(ev MouseEvent) {
			ev.Window = w
			ev.WindowPoint = ev.Point
			callback(ev)
		},
	)
}

func (w *WindowImpl) OnMouseExit(callback func(MouseEvent)) EventSubscription {
	return w.onMouseExit.Listen(
		func(ev MouseEvent) {
			ev.Window = w
			ev.WindowPoint = ev.Point
			callback(ev)
		},
	)
}

func (w *WindowImpl) OnMouseDown(callback func(MouseEvent)) EventSubscription {
	return w.onMouseDown.Listen(
		func(ev MouseEvent) {
			ev.Window = w
			ev.WindowPoint = ev.Point
			callback(ev)
		},
	)
}

func (w *WindowImpl) OnMouseUp(callback func(MouseEvent)) EventSubscription {
	return w.onMouseUp.Listen(
		func(ev MouseEvent) {
			ev.Window = w
			ev.WindowPoint = ev.Point
			callback(ev)
		},
	)
}

func (w *WindowImpl) OnMouseScroll(callback func(MouseEvent)) EventSubscription {
	return w.onMouseScroll.Listen(
		func(ev MouseEvent) {
			ev.Window = w
			ev.WindowPoint = ev.Point
			callback(ev)
		},
	)
}

func (w *WindowImpl) OnKeyDown(callback func(KeyboardEvent)) EventSubscription {
	return w.onKeyDown.Listen(callback)
}

func (w *WindowImpl) OnKeyUp(callback func(KeyboardEvent)) EventSubscription {
	return w.onKeyUp.Listen(callback)
}

func (w *WindowImpl) OnKeyRepeat(callback func(KeyboardEvent)) EventSubscription {
	return w.onKeyRepeat.Listen(callback)
}

func (w *WindowImpl) OnKeyStroke(callback func(KeyStrokeEvent)) EventSubscription {
	return w.onKeyStroke.Listen(callback)
}

func (w *WindowImpl) ReLayout() {
	w.layoutPending = true
	w.requestUpdate()
}

func (w *WindowImpl) Redraw() {
	w.drawPending = true
	w.requestUpdate()
}

func (w *WindowImpl) Click(event MouseEvent) {
	w.onClick.Emit(event)
}

func (w *WindowImpl) DoubleClick(event MouseEvent) {
	w.onDoubleClick.Emit(event)
}

func (w *WindowImpl) KeyPress(event KeyboardEvent) {
	if event.Key == KeyTab {
		if event.Modifier&ModShift != 0 {
			w.focusController.FocusPrev()
		} else {
			w.focusController.FocusNext()
		}
	}
}
func (w *WindowImpl) KeyStroke(event KeyStrokeEvent) {}

func (w *WindowImpl) setViewport(viewport Viewport) {
	for _, subscription := range w.viewportSubscriptions {
		subscription.Forget()
	}

	w.viewport = viewport

	w.viewportSubscriptions = []EventSubscription{
		viewport.OnClose(
			func() {
				if w.onClose != nil {
					w.onClose.Emit()
				}
			},
		),
		viewport.OnResize(func() { w.onResize.Emit() }),
		viewport.OnMouseMove(func(ev MouseEvent) { w.onMouseMove.Emit(ev) }),
		viewport.OnMouseEnter(func(ev MouseEvent) { w.onMouseEnter.Emit(ev) }),
		viewport.OnMouseExit(func(ev MouseEvent) { w.onMouseExit.Emit(ev) }),
		viewport.OnMouseDown(func(ev MouseEvent) { w.onMouseDown.Emit(ev) }),
		viewport.OnMouseUp(func(ev MouseEvent) { w.onMouseUp.Emit(ev) }),
		viewport.OnMouseScroll(func(ev MouseEvent) { w.onMouseScroll.Emit(ev) }),
		viewport.OnKeyDown(func(ev KeyboardEvent) { w.onKeyDown.Emit(ev) }),
		viewport.OnKeyUp(func(ev KeyboardEvent) { w.onKeyUp.Emit(ev) }),
		viewport.OnKeyRepeat(func(ev KeyboardEvent) { w.onKeyRepeat.Emit(ev) }),
		viewport.OnKeyStroke(func(ev KeyStrokeEvent) { w.onKeyStroke.Emit(ev) }),
	}

	w.ReLayout()
}
