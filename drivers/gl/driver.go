// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gl contains an OpenGL implementation of the gxui.Driver interface.
package gl

import (
	"container/list"
	"image"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

// Maximum time allowed for application to process events on termination.
const maxFlushTime = time.Second * 3

func init() {
	runtime.LockOSThread()
}

type DriverImpl struct {
	pendingDriver chan func()
	pendingApp    chan func()
	terminated    int32 // non-zero represents driver terminations
	viewports     *list.List

	pcs  []uintptr // reusable scratch-buffer for use by runtime.Callers.
	uiPC uintptr   // the program-counter of the applicationLoop function.
}

// StartDriver starts the gl driver with the given appRoutine.
func StartDriver(appRoutine func(driver gxui.Driver)) {
	if runtime.GOMAXPROCS(-1) < 2 {
		runtime.GOMAXPROCS(2)
	}

	if err := glfw.Init(gl.ContextWatcher); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	result := &DriverImpl{
		pendingDriver: make(chan func(), 256),
		pendingApp:    make(chan func(), 256),
		viewports:     list.New(),
		pcs:           make([]uintptr, 256),
	}

	result.pendingApp <- result.discoverUIGoRoutine
	result.pendingApp <- func() { appRoutine(result) }

	go result.applicationLoop()

	result.driverLoop()
}

func (d *DriverImpl) asyncDriver(callback func()) {
	d.pendingDriver <- callback
	d.wake()
}

func (d *DriverImpl) syncDriver(callback func()) {
	done := make(chan bool, 1)
	d.asyncDriver(
		func() { callback(); done <- true },
	)
	<-done
}

func (d *DriverImpl) createDriverEvent(signature interface{}) gxui.Event {
	return gxui.CreateChanneledEvent(signature, d.pendingDriver)
}

func (d *DriverImpl) createAppEvent(signature interface{}) gxui.Event {
	return gxui.CreateChanneledEvent(signature, d.pendingApp)
}

// driverLoop pulls and executes funcs from the pendingDriver chan until chan
// close. If there are no funcs enqueued, the driver routine calls and blocks on
// glfw.WaitEvents. All sends on the pendingDriver chan should be paired with a
// call to wake() so that glfw.WaitEvents can return.
func (d *DriverImpl) driverLoop() {
	for {
		select {
		case ev, open := <-d.pendingDriver:
			if open {
				ev()
			} else {
				return // terminated
			}
		default:
			glfw.WaitEvents()
		}
	}
}

func (d *DriverImpl) wake() {
	glfw.PostEmptyEvent()
}

// applicationLoop pulls and executes funcs from the pendingApp chan until
// the chan is closed.
func (d *DriverImpl) applicationLoop() {
	for ev := range d.pendingApp {
		ev()
	}
}

// gxui.Driver compliance
func (d *DriverImpl) Call(callback func()) bool {
	if callback == nil {
		panic("Function must not be nil")
	}

	if atomic.LoadInt32(&d.terminated) != 0 {
		return false // Driver.Terminate has been called
	}
	d.pendingApp <- callback
	return true
}

func (d *DriverImpl) CallSync(callback func()) bool {
	done := make(chan struct{})
	if d.Call(
		func() { callback(); close(done) },
	) {
		<-done
		return true
	}
	return false
}

func (d *DriverImpl) Terminate() {
	d.asyncDriver(
		func() {
			// Close all viewports. This will notify the application.
			for frontViewport := d.viewports.Front(); frontViewport != nil; frontViewport = frontViewport.Next() {
				frontViewport.Value.(*ViewportImpl).Destroy()
			}

			// Flush all remaining events from the application and driver.
			// This gives the application an opportunity to handle shutdown.
			flushStart := time.Now()
			for time.Since(flushStart) < maxFlushTime {
				done := true

				// Process any application events
				sync := make(chan struct{})
				d.Call(func() {
					select {
					case ev := <-d.pendingApp:
						ev()
						done = false
					default:
					}
					close(sync)
				})

				<-sync

				// Process any driver events
				select {
				case ev := <-d.pendingDriver:
					ev()
					done = false
				default:
				}

				if done {
					break
				}
			}

			// All done.
			atomic.StoreInt32(&d.terminated, 1)

			close(d.pendingApp)
			close(d.pendingDriver)

			d.viewports = nil
		})
}

func (d *DriverImpl) SetClipboard(str string) {
	d.asyncDriver(
		func() {
			frontViewport := d.viewports.Front().Value.(*ViewportImpl)
			frontViewport.window.SetClipboardString(str)
		},
	)
}

func (d *DriverImpl) GetClipboard() (str string, err error) {
	d.syncDriver(
		func() {
			frontViewport := d.viewports.Front().Value.(*ViewportImpl)
			str = frontViewport.window.GetClipboardString()
		},
	)
	return
}

func (d *DriverImpl) CreateFont(data []byte, size int) (gxui.Font, error) {
	return newFont(data, size)
}

func (d *DriverImpl) CreateWindowedViewport(width, height int, name string) gxui.Viewport {
	var v *ViewportImpl
	d.syncDriver(
		func() {
			v = NewViewport(d, width, height, name, false)
			e := d.viewports.PushBack(v)
			v.onDestroy.Listen(func() {
				d.viewports.Remove(e)
			})
		},
	)
	return v
}

func (d *DriverImpl) CreateFullscreenViewport(width, height int, name string) gxui.Viewport {
	var v *ViewportImpl
	d.syncDriver(
		func() {
			v = NewViewport(d, width, height, name, true)
			e := d.viewports.PushBack(v)
			v.onDestroy.Listen(func() {
				d.viewports.Remove(e)
			})
		},
	)
	return v
}

func (d *DriverImpl) CreateCanvas(s math.Size) gxui.Canvas {
	return NewCanvas(s)
}

func (d *DriverImpl) CreateTexture(img image.Image, pixelsPerDip float32) gxui.Texture {
	return NewTexture(img, pixelsPerDip)
}
