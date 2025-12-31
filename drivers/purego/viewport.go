package purego

import (
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/badu/gxui"
	"github.com/badu/gxui/pkg/math"
)

const ScrollSpeed = 20.0
const viewportDebugEnabled = false

const clearColorR = 0.5
const clearColorG = 0.5
const clearColorB = 0.5

type ViewportImpl struct {
	sync.Mutex
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
	onDestroy               gxui.Event
	driver                  *DriverImpl
	context                 *context
	window                  *Window
	canvas                  *CanvasImpl
	pendingMouseMoveEvent   *gxui.MouseEvent
	pendingMouseScrollEvent *gxui.MouseEvent
	title                   string
	sizeDipsUnscaled        math.Size
	sizeDips                math.Size
	sizePixels              math.Size
	position                math.Point
	scrollAccumX            float64
	scrollAccumY            float64

	scaling     float32
	redrawCount uint32

	fullscreen bool
	destroyed  bool
}

func NewViewport(driver *DriverImpl, width, height int, title string, fullscreen bool) *ViewportImpl {

	result := &ViewportImpl{fullscreen: fullscreen, scaling: 1, title: title}

	DefaultWindowHints()
	WindowHint(Samples, 4)

	var monitor *Monitor
	if fullscreen {
		monitor = GetPrimaryMonitor()
		if width == 0 || height == 0 {
			vm := monitor.GetVideoMode()
			if vm == nil {
				panic("No video mode available on primary monitor")
			}
			width, height = vm.Width, vm.Height
		}
	}

	wnd := CreateWindow(width, height, result.title, monitor, nil)

	width, height = wnd.GetSize() // At this time, width and height serve as a "hint" for CreateWindow, so get actual values from window.

	wnd.MakeContextCurrent()

	result.context = newContext(driver.fn)

	cursorPoint := func(x, y float64) math.Point {
		// HACK: xpos is off by 1 and ypos is off by 3 on OSX.
		// Compensate until real fix is found.
		x -= 1.0
		y -= 3.0
		return math.Point{X: int(x), Y: int(y)}.ScaleS(1 / result.scaling)
	}
	wnd.SetCloseCallback(
		func(*Window) {
			result.Close()
		},
	)

	wnd.SetPosCallback(
		func(w *Window, x, y int32) {
			result.Lock()
			result.position = math.NewPoint(int(x), int(y))
			result.Unlock()
		},
	)

	wnd.SetSizeCallback(
		func(_ *Window, w, h int32) {
			result.Lock()
			result.sizeDipsUnscaled = math.Size{Width: int(w), Height: int(h)}
			result.sizeDips = result.sizeDipsUnscaled.ScaleS(1 / result.scaling)
			result.Unlock()
			result.onResize.Emit()
		},
	)

	wnd.SetFramebufferSizeCallback(
		func(_ *Window, w, h int32) {
			result.Lock()
			result.sizePixels = math.Size{Width: int(w), Height: int(h)}
			result.Unlock()
			driver.fn.Viewport(0, 0, int32(w), int32(h))
			driver.fn.ClearColor(clearColorR, clearColorG, clearColorB, 1.0)
			driver.fn.Clear(COLOR_BUFFER_BIT)
		},
	)

	wnd.SetCursorPosCallback(
		func(w *Window, x, y float64) {
			p := cursorPoint(w.GetCursorPos())
			result.Lock()
			if result.pendingMouseMoveEvent == nil {
				result.pendingMouseMoveEvent = &gxui.MouseEvent{}
				driver.Call(func() {
					result.Lock()
					ev := *result.pendingMouseMoveEvent
					result.pendingMouseMoveEvent = nil
					result.Unlock()
					result.onMouseMove.Emit(ev)
				})
			}
			result.pendingMouseMoveEvent.Point = p
			result.pendingMouseMoveEvent.State = getMouseState(w)
			result.Unlock()
		},
	)

	wnd.SetCursorEnterCallback(
		func(w *Window, entered bool) {
			p := cursorPoint(w.GetCursorPos())
			ev := gxui.MouseEvent{
				Point: p,
			}
			ev.State = getMouseState(w)
			if entered {
				result.onMouseEnter.Emit(ev)
			} else {
				result.onMouseExit.Emit(ev)
			}
		},
	)

	wnd.SetScrollCallback(
		func(w *Window, xoff, yoff float64) {
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
						result.onMouseScroll.Emit(ev)
					} else {
						result.Unlock()
					}
				})
			}
			result.pendingMouseScrollEvent.Point = p
			result.scrollAccumX += xoff * ScrollSpeed
			result.scrollAccumY += yoff * ScrollSpeed
			result.pendingMouseScrollEvent.State = getMouseState(w)
			result.Unlock()
		},
	)

	wnd.SetMouseButtonCallback(
		func(w *Window, button MouseButton, action Action, mod ModifierKey) {
			p := cursorPoint(w.GetCursorPos())
			ev := gxui.MouseEvent{
				Point:    p,
				Modifier: translateKeyboardModifier(mod),
			}
			ev.Button = translateMouseButton(button)
			ev.State = getMouseState(w)
			if action == Press {
				result.onMouseDown.Emit(ev)
			} else {
				result.onMouseUp.Emit(ev)
			}
		},
	)

	wnd.SetKeyCallback(
		func(w *Window, key Key, scancode int32, action Action, mods ModifierKey) {
			ev := gxui.KeyboardEvent{
				Key:      translateKeyboardKey(key),
				Modifier: translateKeyboardModifier(mods),
			}
			switch action {
			case Press:
				result.onKeyDown.Emit(ev)
			case Release:
				result.onKeyUp.Emit(ev)
			case Repeat:
				result.onKeyRepeat.Emit(ev)
			}
		},
	)

	wnd.SetCharModsCallback(
		func(w *Window, char rune, mods ModifierKey) {
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
			result.onKeyStroke.Emit(ev)
		},
	)

	wnd.SetRefreshCallback(
		func(w *Window) {
			if result.canvas != nil {
				result.render()
			}
		},
	)

	fw, fh := wnd.GetFramebufferSize()
	posX, posY := wnd.GetPos()

	// Pre-multiplied alpha blending
	driver.fn.BlendFunc(ONE, ONE_MINUS_SRC_ALPHA)
	driver.fn.Enable(BLEND)
	driver.fn.Enable(SCISSOR_TEST)
	driver.fn.Viewport(0, 0, fw, fh)
	driver.fn.Scissor(0, 0, fw, fh)
	driver.fn.ClearColor(clearColorR, clearColorG, clearColorB, 1.0)
	driver.fn.Clear(COLOR_BUFFER_BIT)
	wnd.SwapBuffers()

	result.window = wnd
	result.driver = driver

	result.onClose = driver.createAppEvent(func() {})
	result.onResize = driver.createAppEvent(func() {})

	result.onMouseMove = gxui.NewListener(func(gxui.MouseEvent) {})
	result.onMouseEnter = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseExit = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseDown = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseUp = driver.createAppEvent(func(gxui.MouseEvent) {})
	result.onMouseScroll = gxui.NewListener(func(gxui.MouseEvent) {})
	result.onKeyDown = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	result.onKeyUp = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	result.onKeyRepeat = driver.createAppEvent(func(gxui.KeyboardEvent) {})
	result.onKeyStroke = driver.createAppEvent(func(gxui.KeyStrokeEvent) {})
	result.onDestroy = driver.createDriverEvent(func() {})

	result.sizeDipsUnscaled = math.Size{Width: width, Height: height}
	result.sizeDips = result.sizeDipsUnscaled.ScaleS(1 / result.scaling)
	result.sizePixels = math.Size{Width: int(fw), Height: int(fh)}
	result.position = math.Point{X: int(posX), Y: int(posY)}

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

// SetCanvas is gxui.Viewport compliance
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
		v.onResize.Emit()
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
		v.window.SetSize(v.sizeDipsUnscaled.Width, v.sizeDipsUnscaled.Height)
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
	v.onClose.Emit()
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
			v.onDestroy.Emit()
			v.destroyed = true
		}
	})
}
