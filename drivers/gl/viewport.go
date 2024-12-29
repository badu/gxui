// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/gl/platform"
	"github.com/badu/gxui/math"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

const viewportDebugEnabled = false

const clearColorR = 0.5
const clearColorG = 0.5
const clearColorB = 0.5

type ViewportImpl struct {
	sync.Mutex

	driver                  *DriverImpl
	context                 *context
	window                  *glfw.Window
	canvas                  *CanvasImpl
	fullscreen              bool
	scaling                 float32
	sizeDipsUnscaled        math.Size
	sizeDips                math.Size
	sizePixels              math.Size
	position                math.Point
	title                   string
	pendingMouseMoveEvent   *gxui.MouseEvent
	pendingMouseScrollEvent *gxui.MouseEvent
	scrollAccumX            float64
	scrollAccumY            float64
	destroyed               bool
	redrawCount             uint32

	// Broadcasts to application thread
	onClose       gxui.Event // ()
	onResize      gxui.Event // ()
	onMouseMove   gxui.Event // (gxui.MouseEvent)
	onMouseEnter  gxui.Event // (gxui.MouseEvent)
	onMouseExit   gxui.Event // (gxui.MouseEvent)
	onMouseDown   gxui.Event // (gxui.MouseEvent)
	onMouseUp     gxui.Event // (gxui.MouseEvent)
	onMouseScroll gxui.Event // (gxui.MouseEvent)
	onKeyDown     gxui.Event // (gxui.KeyboardEvent)
	onKeyUp       gxui.Event // (gxui.KeyboardEvent)
	onKeyRepeat   gxui.Event // (gxui.KeyboardEvent)
	onKeyStroke   gxui.Event // (gxui.KeyStrokeEvent)
	// Broadcasts to driver thread
	onDestroy gxui.Event
}

func NewViewport(driver *DriverImpl, width, height int, title string, fullscreen bool) *ViewportImpl {
	result := &ViewportImpl{fullscreen: fullscreen, scaling: 1, title: title}

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Samples, 4)

	var monitor *glfw.Monitor
	if fullscreen {
		monitor = glfw.GetPrimaryMonitor()
		if width == 0 || height == 0 {
			vm := monitor.GetVideoMode()
			width, height = vm.Width, vm.Height
		}
	}

	wnd, err := glfw.CreateWindow(width, height, result.title, monitor, nil)
	if err != nil {
		panic(err)
	}
	width, height = wnd.GetSize() // At this time, width and height serve as a "hint" for glfw.CreateWindow, so get actual values from window.

	wnd.MakeContextCurrent()

	result.context = newContext()

	cursorPoint := func(x, y float64) math.Point {
		// HACK: xpos is off by 1 and ypos is off by 3 on OSX.
		// Compensate until real fix is found.
		x -= 1.0
		y -= 3.0
		return math.Point{X: int(x), Y: int(y)}.ScaleS(1 / result.scaling)
	}
	wnd.SetCloseCallback(
		func(*glfw.Window) {
			result.Close()
		},
	)

	wnd.SetPosCallback(
		func(w *glfw.Window, x, y int) {
			result.Lock()
			result.position = math.NewPoint(x, y)
			result.Unlock()
		},
	)

	wnd.SetSizeCallback(
		func(_ *glfw.Window, w, h int) {
			result.Lock()
			result.sizeDipsUnscaled = math.Size{W: w, H: h}
			result.sizeDips = result.sizeDipsUnscaled.ScaleS(1 / result.scaling)
			result.Unlock()
			result.onResize.Fire()
		},
	)

	wnd.SetFramebufferSizeCallback(
		func(_ *glfw.Window, w, h int) {
			result.Lock()
			result.sizePixels = math.Size{W: w, H: h}
			result.Unlock()
			gl.Viewport(0, 0, w, h)
			gl.ClearColor(clearColorR, clearColorG, clearColorB, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT)
		},
	)

	wnd.SetCursorPosCallback(
		func(w *glfw.Window, x, y float64) {
			p := cursorPoint(w.GetCursorPos())
			result.Lock()
			if result.pendingMouseMoveEvent == nil {
				result.pendingMouseMoveEvent = &gxui.MouseEvent{}
				driver.Call(func() {
					result.Lock()
					ev := *result.pendingMouseMoveEvent
					result.pendingMouseMoveEvent = nil
					result.Unlock()
					result.onMouseMove.Fire(ev)
				})
			}
			result.pendingMouseMoveEvent.Point = p
			result.pendingMouseMoveEvent.State = getMouseState(w)
			result.Unlock()
		},
	)

	wnd.SetCursorEnterCallback(
		func(w *glfw.Window, entered bool) {
			p := cursorPoint(w.GetCursorPos())
			ev := gxui.MouseEvent{
				Point: p,
			}
			ev.State = getMouseState(w)
			if entered {
				result.onMouseEnter.Fire(ev)
			} else {
				result.onMouseExit.Fire(ev)
			}
		},
	)

	wnd.SetScrollCallback(
		func(w *glfw.Window, xoff, yoff float64) {
			p := cursorPoint(w.GetCursorPos())
			result.Lock()
			if result.pendingMouseScrollEvent == nil {
				result.pendingMouseScrollEvent = &gxui.MouseEvent{}
				driver.Call(func() {
					result.Lock()
					ev := *result.pendingMouseScrollEvent
					result.pendingMouseScrollEvent = nil
					ev.ScrollX, ev.ScrollY = int(result.scrollAccumX), int(result.scrollAccumY)
					if ev.ScrollX != 0 || ev.ScrollY != 0 {
						result.scrollAccumX -= float64(ev.ScrollX)
						result.scrollAccumY -= float64(ev.ScrollY)
						result.Unlock()
						result.onMouseScroll.Fire(ev)
					} else {
						result.Unlock()
					}
				})
			}
			result.pendingMouseScrollEvent.Point = p
			result.scrollAccumX += xoff * platform.ScrollSpeed
			result.scrollAccumY += yoff * platform.ScrollSpeed
			result.pendingMouseScrollEvent.State = getMouseState(w)
			result.Unlock()
		},
	)

	wnd.SetMouseButtonCallback(
		func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
			p := cursorPoint(w.GetCursorPos())
			ev := gxui.MouseEvent{
				Point:    p,
				Modifier: translateKeyboardModifier(mod),
			}
			ev.Button = translateMouseButton(button)
			ev.State = getMouseState(w)
			if action == glfw.Press {
				result.onMouseDown.Fire(ev)
			} else {
				result.onMouseUp.Fire(ev)
			}
		},
	)

	wnd.SetKeyCallback(
		func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			ev := gxui.KeyboardEvent{
				Key:      translateKeyboardKey(key),
				Modifier: translateKeyboardModifier(mods),
			}
			switch action {
			case glfw.Press:
				result.onKeyDown.Fire(ev)
			case glfw.Release:
				result.onKeyUp.Fire(ev)
			case glfw.Repeat:
				result.onKeyRepeat.Fire(ev)
			}
		},
	)

	wnd.SetCharModsCallback(
		func(w *glfw.Window, char rune, mods glfw.ModifierKey) {
			if !unicode.IsControl(char) &&
				!unicode.IsGraphic(char) &&
				!unicode.IsLetter(char) &&
				!unicode.IsMark(char) &&
				!unicode.IsNumber(char) &&
				!unicode.IsPunct(char) &&
				!unicode.IsSpace(char) &&
				!unicode.IsSymbol(char) {
				return // Weird unicode character. Ignore
			}

			ev := gxui.KeyStrokeEvent{
				Character: char,
				Modifier:  translateKeyboardModifier(mods),
			}
			result.onKeyStroke.Fire(ev)
		},
	)

	wnd.SetRefreshCallback(
		func(w *glfw.Window) {
			if result.canvas != nil {
				result.render()
			}
		},
	)

	fw, fh := wnd.GetFramebufferSize()
	posX, posY := wnd.GetPos()

	// Pre-multiplied alpha blending
	gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Viewport(0, 0, fw, fh)
	gl.Scissor(0, 0, int32(fw), int32(fh))
	gl.ClearColor(clearColorR, clearColorG, clearColorB, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	wnd.SwapBuffers()

	result.window = wnd
	result.driver = driver

	result.onClose = driver.createAppEvent(func() {})
	result.onResize = driver.createAppEvent(func() {})

	result.onMouseMove = gxui.CreateEvent(func(gxui.MouseEvent) {})
	result.onMouseEnter = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseExit = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseDown = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseUp = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseScroll = gxui.CreateEvent(func(gxui.MouseEvent) {})
	result.onKeyDown = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	result.onKeyUp = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	result.onKeyRepeat = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	result.onKeyStroke = driver.createAppEvent(func(gxui.KeyStrokeEvent) {})
	result.onDestroy = driver.createDriverEvent(func() {})

	result.sizeDipsUnscaled = math.Size{W: width, H: height}
	result.sizeDips = result.sizeDipsUnscaled.ScaleS(1 / result.scaling)
	result.sizePixels = math.Size{W: fw, H: fh}
	result.position = math.Point{X: posX, Y: posY}

	return result
}

// Driver methods
// These methods are all called on the driver routine
func (v *ViewportImpl) render() {
	if v.destroyed {
		return
	}

	v.window.MakeContextCurrent()

	ctx := v.context
	ctx.beginDraw(v.SizeDips(), v.SizePixels())

	stack := drawStateStack{
		drawState{
			ClipPixels: v.sizePixels.Rect(),
		},
	}

	v.canvas.draw(ctx, &stack)
	if len(stack) != 1 {
		panic("DrawStateStack count was not 1 after calling Canvas.Draw")
	}

	ctx.apply(stack.head())
	ctx.blitter.commit(ctx)

	if viewportDebugEnabled {
		v.drawFrameUpdate(ctx)
	}

	ctx.endDraw()

	v.window.SwapBuffers()
}

func (v *ViewportImpl) drawFrameUpdate(ctx *context) {
	dx := (ctx.stats.frameCount * 10) & 0xFF
	rect := math.CreateRect(dx-5, 0, dx+5, 3)
	state := &drawState{}
	ctx.blitter.blitRect(ctx, rect, gxui.White, state)
}

// gxui.Viewport compliance
// These methods are all called on the application routine
func (v *ViewportImpl) SetCanvas(newCanvas gxui.Canvas) {
	cnt := atomic.AddUint32(&v.redrawCount, 1)
	childCanvas := newCanvas.(*CanvasImpl)
	v.driver.asyncDriver(func() {
		// Only use the canvas of the most recent SetCanvas call.
		v.window.MakeContextCurrent()
		if atomic.LoadUint32(&v.redrawCount) == cnt {
			v.canvas = childCanvas
			if v.canvas != nil {
				v.render()
			}
		}
	})
}

func (v *ViewportImpl) Scale() float32 {
	v.Lock()
	defer v.Unlock()
	return v.scaling
}

func (v *ViewportImpl) SetScale(newScale float32) {
	v.Lock()
	defer v.Unlock()
	if newScale != v.scaling {
		v.scaling = newScale
		v.sizeDips = v.sizeDipsUnscaled.ScaleS(1 / newScale)
		v.onResize.Fire()
	}
}

func (v *ViewportImpl) SizeDips() math.Size {
	v.Lock()
	defer v.Unlock()
	return v.sizeDips
}

func (v *ViewportImpl) SetSizeDips(size math.Size) {
	v.driver.syncDriver(func() {
		v.sizeDips = size
		v.sizeDipsUnscaled = size.ScaleS(v.scaling)
		v.window.SetSize(v.sizeDipsUnscaled.W, v.sizeDipsUnscaled.H)
	})
}

func (v *ViewportImpl) SizePixels() math.Size {
	v.Lock()
	defer v.Unlock()
	return v.sizePixels
}

func (v *ViewportImpl) Title() string {
	v.Lock()
	defer v.Unlock()
	return v.title
}

func (v *ViewportImpl) SetTitle(title string) {
	v.Lock()
	v.title = title
	v.Unlock()
	v.driver.asyncDriver(func() {
		v.window.SetTitle(title)
	})
}

func (v *ViewportImpl) Position() math.Point {
	v.Lock()
	defer v.Unlock()
	return v.position
}

func (v *ViewportImpl) SetPosition(newPosition math.Point) {
	v.Lock()
	v.position = newPosition
	v.Unlock()
	v.driver.asyncDriver(func() {
		v.window.SetPos(newPosition.X, newPosition.Y)
	})
}

func (v *ViewportImpl) Fullscreen() bool {
	return v.fullscreen
}

func (v *ViewportImpl) Show() {
	v.driver.asyncDriver(func() { v.window.Show() })
}

func (v *ViewportImpl) Hide() {
	v.driver.asyncDriver(func() { v.window.Hide() })
}

func (v *ViewportImpl) Close() {
	v.onClose.Fire()
	v.Destroy()
}

func (v *ViewportImpl) OnResize(f func()) gxui.EventSubscription {
	return v.onResize.Listen(f)
}

func (v *ViewportImpl) OnClose(f func()) gxui.EventSubscription {
	return v.onClose.Listen(f)
}

func (v *ViewportImpl) OnMouseMove(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseMove.Listen(f)
}

func (v *ViewportImpl) OnMouseEnter(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseEnter.Listen(f)
}

func (v *ViewportImpl) OnMouseExit(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseExit.Listen(f)
}

func (v *ViewportImpl) OnMouseDown(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseDown.Listen(f)
}

func (v *ViewportImpl) OnMouseUp(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseUp.Listen(f)
}

func (v *ViewportImpl) OnMouseScroll(f func(gxui.MouseEvent)) gxui.EventSubscription {
	return v.onMouseScroll.Listen(f)
}

func (v *ViewportImpl) OnKeyDown(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyDown.Listen(f)
}

func (v *ViewportImpl) OnKeyUp(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyUp.Listen(f)
}

func (v *ViewportImpl) OnKeyRepeat(f func(gxui.KeyboardEvent)) gxui.EventSubscription {
	return v.onKeyRepeat.Listen(f)
}

func (v *ViewportImpl) OnKeyStroke(f func(gxui.KeyStrokeEvent)) gxui.EventSubscription {
	return v.onKeyStroke.Listen(f)
}

func (v *ViewportImpl) Destroy() {
	v.driver.asyncDriver(func() {
		if !v.destroyed {
			v.window.MakeContextCurrent()
			v.canvas = nil
			v.context.destroy()
			v.window.Destroy()
			v.onDestroy.Fire()
			v.destroyed = true
		}
	})
}
