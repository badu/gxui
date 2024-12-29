// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type WindowOuter interface {
	gxui.Window
	Attached() bool                                            // was outer.Attachable
	Attach()                                                   // was outer.Attachable
	Detach()                                                   // was outer.Attachable
	OnAttach(callback func()) gxui.EventSubscription           // was outer.Attachable
	OnDetach(callback func()) gxui.EventSubscription           // was outer.Attachable
	IsVisible() bool                                           // was outer.IsVisibler
	LayoutChildren()                                           // was outer.LayoutChildren
	PaintChild(canvas gxui.Canvas, child *gxui.Child, idx int) // was outer.PaintChilder
	Paint(canvas gxui.Canvas)                                  // was outer.Painter
	Parent() gxui.Parent                                       // was outer.Parenter
	Size() math.Size                                           // was outer.Sized
	SetSize(newSize math.Size)                                 // was outer.Sized
}

type Window struct {
	AttachablePart
	BackgroundBorderPainter
	ContainerPart
	PaddablePart
	PaintChildrenPart
	driver                gxui.Driver
	outer                 WindowOuter
	viewport              gxui.Viewport
	windowedSize          math.Size
	mouseController       *gxui.MouseController
	keyboardController    *gxui.KeyboardController
	focusController       *gxui.FocusController
	layoutPending         bool
	drawPending           bool
	updatePending         bool
	onClose               gxui.Event // Raised by viewport
	onResize              gxui.Event // Raised by viewport
	onMouseMove           gxui.Event // Raised by viewport
	onMouseEnter          gxui.Event // Raised by viewport
	onMouseExit           gxui.Event // Raised by viewport
	onMouseDown           gxui.Event // Raised by viewport
	onMouseUp             gxui.Event // Raised by viewport
	onMouseScroll         gxui.Event // Raised by viewport
	onKeyDown             gxui.Event // Raised by viewport
	onKeyUp               gxui.Event // Raised by viewport
	onKeyRepeat           gxui.Event // Raised by viewport
	onKeyStroke           gxui.Event // Raised by viewport
	onClick               gxui.Event // Raised by MouseController
	onDoubleClick         gxui.Event // Raised by MouseController
	viewportSubscriptions []gxui.EventSubscription
}

func (w *Window) requestUpdate() {
	if !w.updatePending {
		w.updatePending = true
		w.driver.Call(w.update)
	}
}

func (w *Window) update() {
	if !w.Attached() {
		// Window was detached between requestUpdate() and update()
		w.updatePending = false
		w.layoutPending = false
		w.drawPending = false
		return
	}
	w.updatePending = false
	if w.layoutPending {
		w.layoutPending = false
		w.drawPending = true
		w.outer.LayoutChildren()
	}
	if w.drawPending {
		w.drawPending = false
		w.Draw()
	}
}

func (w *Window) Init(outer WindowOuter, driver gxui.Driver, width, height int, title string) {
	w.AttachablePart.Init(outer)
	w.BackgroundBorderPainter.Init(outer)
	w.ContainerPart.Init(outer)
	w.PaddablePart.Init(outer)
	w.PaintChildrenPart.Init(outer)
	w.outer = outer
	w.driver = driver

	w.onClose = gxui.CreateEvent(func() {})
	w.onResize = gxui.CreateEvent(func() {})
	w.onMouseMove = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseEnter = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseExit = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseDown = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseUp = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onMouseScroll = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onKeyDown = gxui.CreateEvent(func(gxui.KeyboardEvent) {})
	w.onKeyUp = gxui.CreateEvent(func(gxui.KeyboardEvent) {})
	w.onKeyRepeat = gxui.CreateEvent(func(gxui.KeyboardEvent) {})
	w.onKeyStroke = gxui.CreateEvent(func(gxui.KeyStrokeEvent) {})

	w.onClick = gxui.CreateEvent(func(gxui.MouseEvent) {})
	w.onDoubleClick = gxui.CreateEvent(func(gxui.MouseEvent) {})

	w.focusController = gxui.CreateFocusController(outer)
	w.mouseController = gxui.CreateMouseController(outer, w.focusController)
	w.keyboardController = gxui.CreateKeyboardController(outer)

	w.onResize.Listen(func() {
		w.outer.LayoutChildren()
		w.Draw()
	})

	w.SetBorderPen(gxui.TransparentPen)

	w.setViewport(driver.CreateWindowedViewport(width, height, title))

	// Window starts shown
	w.Attach()
}

func (w *Window) Draw() gxui.Canvas {
	if s := w.viewport.SizeDips(); s != math.ZeroSize {
		c := w.driver.CreateCanvas(s)
		w.outer.Paint(c)
		c.Complete()
		w.viewport.SetCanvas(c)
		return c
	} else {
		return nil
	}
}

func (w *Window) Paint(canvas gxui.Canvas) {
	w.PaintBackground(canvas, canvas.Size().Rect())
	w.PaintChildrenPart.Paint(canvas)
	w.PaintBorder(canvas, canvas.Size().Rect())
}

func (w *Window) LayoutChildren() {
	size := w.Size().Contract(w.Padding()).Max(math.ZeroSize)
	offset := w.Padding().LT()
	for _, child := range w.outer.Children() {
		child.Layout(child.Control.DesiredSize(math.ZeroSize, size).Rect().Offset(offset))
	}
}

func (w *Window) Size() math.Size {
	return w.viewport.SizeDips()
}

func (w *Window) SetSize(size math.Size) {
	w.viewport.SetSizeDips(size)
}

func (w *Window) Parent() gxui.Parent {
	return nil
}

func (w *Window) Viewport() gxui.Viewport {
	return w.viewport
}

func (w *Window) Title() string {
	return w.viewport.Title()
}

func (w *Window) SetTitle(title string) {
	w.viewport.SetTitle(title)
}

func (w *Window) Scale() float32 {
	return w.viewport.Scale()
}

func (w *Window) SetScale(scale float32) {
	w.viewport.SetScale(scale)
}

func (w *Window) Position() math.Point {
	return w.viewport.Position()
}

func (w *Window) SetPosition(point math.Point) {
	w.viewport.SetPosition(point)
}

func (w *Window) Fullscreen() bool {
	return w.viewport.Fullscreen()
}

func (w *Window) SetFullscreen(fullscreen bool) {
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

func (w *Window) Show() {
	w.Attach()
	w.viewport.Show()
}

func (w *Window) Hide() {
	w.Detach()
	w.viewport.Hide()
}

func (w *Window) Close() {
	w.Detach()
	w.viewport.Close()
}

func (w *Window) Focus() gxui.Focusable {
	return w.focusController.Focus()
}

func (w *Window) SetFocus(control gxui.Control) bool {
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

func (w *Window) IsVisible() bool {
	return true
}

func (w *Window) OnClose(callback func()) gxui.EventSubscription {
	return w.onClose.Listen(callback)
}

func (w *Window) OnResize(callback func()) gxui.EventSubscription {
	return w.onResize.Listen(callback)
}

func (w *Window) OnClick(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onClick.Listen(callback)
}

func (w *Window) OnDoubleClick(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onDoubleClick.Listen(callback)
}

func (w *Window) OnMouseMove(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseMove.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		callback(ev)
	})
}

func (w *Window) OnMouseEnter(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseEnter.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		callback(ev)
	})
}

func (w *Window) OnMouseExit(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseExit.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		callback(ev)
	})
}

func (w *Window) OnMouseDown(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseDown.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		callback(ev)
	})
}

func (w *Window) OnMouseUp(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseUp.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		callback(ev)
	})
}

func (w *Window) OnMouseScroll(callback func(gxui.MouseEvent)) gxui.EventSubscription {
	return w.onMouseScroll.Listen(func(ev gxui.MouseEvent) {
		ev.Window = w
		ev.WindowPoint = ev.Point
		callback(ev)
	})
}

func (w *Window) OnKeyDown(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.onKeyDown.Listen(callback)
}

func (w *Window) OnKeyUp(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.onKeyUp.Listen(callback)
}

func (w *Window) OnKeyRepeat(callback func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return w.onKeyRepeat.Listen(callback)
}

func (w *Window) OnKeyStroke(callback func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return w.onKeyStroke.Listen(callback)
}

func (w *Window) Relayout() {
	w.layoutPending = true
	w.requestUpdate()
}

func (w *Window) Redraw() {
	w.drawPending = true
	w.requestUpdate()
}

func (w *Window) Click(event gxui.MouseEvent) {
	w.onClick.Fire(event)
}

func (w *Window) DoubleClick(event gxui.MouseEvent) {
	w.onDoubleClick.Fire(event)
}

func (w *Window) KeyPress(event gxui.KeyboardEvent) {
	if event.Key == gxui.KeyTab {
		if event.Modifier&gxui.ModShift != 0 {
			w.focusController.FocusPrev()
		} else {
			w.focusController.FocusNext()
		}
	}
}
func (w *Window) KeyStroke(event gxui.KeyStrokeEvent) {}

func (w *Window) setViewport(viewport gxui.Viewport) {
	for _, subscription := range w.viewportSubscriptions {
		subscription.Unlisten()
	}

	w.viewport = viewport

	w.viewportSubscriptions = []gxui.EventSubscription{
		viewport.OnClose(func() { w.onClose.Fire() }),
		viewport.OnResize(func() { w.onResize.Fire() }),
		viewport.OnMouseMove(func(ev gxui.MouseEvent) { w.onMouseMove.Fire(ev) }),
		viewport.OnMouseEnter(func(ev gxui.MouseEvent) { w.onMouseEnter.Fire(ev) }),
		viewport.OnMouseExit(func(ev gxui.MouseEvent) { w.onMouseExit.Fire(ev) }),
		viewport.OnMouseDown(func(ev gxui.MouseEvent) { w.onMouseDown.Fire(ev) }),
		viewport.OnMouseUp(func(ev gxui.MouseEvent) { w.onMouseUp.Fire(ev) }),
		viewport.OnMouseScroll(func(ev gxui.MouseEvent) { w.onMouseScroll.Fire(ev) }),
		viewport.OnKeyDown(func(ev gxui.KeyboardEvent) { w.onKeyDown.Fire(ev) }),
		viewport.OnKeyUp(func(ev gxui.KeyboardEvent) { w.onKeyUp.Fire(ev) }),
		viewport.OnKeyRepeat(func(ev gxui.KeyboardEvent) { w.onKeyRepeat.Fire(ev) }),
		viewport.OnKeyStroke(func(ev gxui.KeyStrokeEvent) { w.onKeyStroke.Fire(ev) }),
	}

	w.Relayout()
}
