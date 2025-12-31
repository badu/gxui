// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

type WindowDriver interface {
	// CreateWindowedViewport creates a new windowed Viewport with the specified width and height in device independent pixels.
	CreateWindowedViewport(width, height int, name string) Viewport

	// CreateFullscreenViewport creates a new fullscreen Viewport with the specified width and height in device independent pixels.
	// If width or height is 0, then the viewport adopts the current screen resolution.
	CreateFullscreenViewport(width, height int, name string) Viewport
	CreateCanvas(size math.Size) Canvas
	// Call queues f to be run on the UI go-routine, returning before f may have been called.
	// Call returns false if the driver has been terminated, in which case f may not be called.
	Call(callback func()) bool
}

type WindowImpl struct {
	PaintChildrenPart
	AttachablePart
	ContainerPart
	PaddablePart
	BackgroundBorderPainter
	driver                WindowDriver
	parent                *WindowImpl
	viewport              Viewport
	onClose               Event // Raised by viewport
	onResize              Event // Raised by viewport
	onMouseMove           Event // Raised by viewport
	onMouseEnter          Event // Raised by viewport
	onMouseExit           Event // Raised by viewport
	onMouseDown           Event // Raised by viewport
	onMouseUp             Event // Raised by viewport
	onMouseScroll         Event // Raised by viewport
	onKeyDown             Event // Raised by viewport
	onKeyUp               Event // Raised by viewport
	onKeyRepeat           Event // Raised by viewport
	onKeyStroke           Event // Raised by viewport
	onClick               Event // Raised by MouseController
	onDoubleClick         Event // Raised by MouseController
	mouseController       *MouseController
	keyboardController    *KeyboardController
	focusController       *FocusController
	viewportSubscriptions []EventSubscription
	windowedSize          math.Size
	layoutPending         bool
	drawPending           bool
	updatePending         bool
}

func (w *WindowImpl) Init(
	window *WindowImpl,
	driver WindowDriver,
	width, height int,
	title string,
) {
	w.BackgroundBorderPainter.Init(window)
	w.ContainerPart.Init(window)
	w.PaddablePart.Init(window)
	w.PaintChildrenPart.Init(window)

	w.parent = window
	w.driver = driver

	w.onResize = NewListener(func() {})
	w.onMouseMove = NewListener(func(MouseEvent) {})
	w.onMouseEnter = NewListener(func(MouseEvent) {})
	w.onMouseExit = NewListener(func(MouseEvent) {})
	w.onMouseDown = NewListener(func(MouseEvent) {})
	w.onMouseUp = NewListener(func(MouseEvent) {})
	w.onMouseScroll = NewListener(func(MouseEvent) {})
	w.onKeyDown = NewListener(func(KeyboardEvent) {})
	w.onKeyUp = NewListener(func(KeyboardEvent) {})
	w.onKeyRepeat = NewListener(func(KeyboardEvent) {})
	w.onKeyStroke = NewListener(func(KeyStrokeEvent) {})

	w.onClick = NewListener(func(MouseEvent) {})
	w.onDoubleClick = NewListener(func(MouseEvent) {})

	w.focusController = CreateFocusController(window)
	w.mouseController = CreateMouseController(window, w.focusController)
	w.keyboardController = CreateKeyboardController(window)

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
	offset := w.Padding().TopLeft()
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
		w.onClose = NewListener(func() {})
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
