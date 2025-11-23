// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"image"

	"github.com/badu/gxui/pkg/math"
)

// A Font represents a TrueType font loaded by the GXUI driver.
type Font interface {
	LoadGlyphs(first, last rune)
	Size() int
	GlyphMaxSize() math.Size
	Measure(*TextBlock) math.Size
	Layout(*TextBlock) (offsets []math.Point)
}

// TextBlock is a sequence of runes to be laid out.
type TextBlock struct {
	Runes     []rune
	AlignRect math.Rect
	H         HAlign
	V         VAlign
}

type Canvas interface {
	Size() math.Size
	IsComplete() bool
	Complete()
	Push()
	Pop()
	AddClip(rect math.Rect)
	Clear(color Color)
	DrawCanvas(canvas Canvas, position math.Point)
	DrawTexture(texture Texture, bounds math.Rect)
	DrawRunes(font Font, runes []rune, points []math.Point, color Color)
	DrawLines(polygon Polygon, pen Pen)
	DrawPolygon(polygon Polygon, pen Pen, brush Brush)
	DrawRect(rect math.Rect, brush Brush)
	DrawRoundedRect(rect math.Rect, tl, tr, bl, br float32, p Pen, b Brush)
}

type Viewport interface {
	// SizeDips returns the size of the viewport in device-independent pixels.
	// The ratio of pixels to DIPs is based on the screen density and scale
	// adjustments made with the SetScale method.
	SizeDips() math.Size

	// SetSizeDips sets the size of the viewport in device-independent pixels.
	// The ratio of pixels to DIPs is based on the screen density and scale
	// adjustments made with the SetScale method.
	SetSizeDips(newSize math.Size)

	// SizePixels returns the size of the viewport in pixels.
	SizePixels() math.Size

	// Scale returns the display scaling for this viewport.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	Scale() float32

	// SetScale alters the display scaling for this viewport.
	// A scale of 1 is unscaled, 2 is twice the regular scaling.
	SetScale(scale float32)

	// Fullscreen returns true if the viewport was created full-screen.
	Fullscreen() bool

	// Title returns the title of the window.
	// This is usually the text displayed at the top of the window.
	Title() string

	// SetTitle changes the title of the window.
	SetTitle(title string)

	// Position returns position of the window.
	Position() math.Point

	// SetPosition changes position of the window.
	SetPosition(newPosition math.Point)

	// Show makes the window visible.
	Show()

	// Hide makes the window invisible.
	Hide()

	// Close destroys the window.
	// Once the window is closed, no further calls should be made to it.
	Close()

	// SetCanvas changes the displayed content of the viewport to the specified
	// Canvas. As canvases are immutable once completed, every visual update of a
	// viewport will require a call to SetCanvas.
	SetCanvas(canvas Canvas)

	// OnClose subscribes f to be called when the viewport closes.
	OnClose(callback func()) EventSubscription

	// OnResize subscribes f to be called whenever the viewport changes size.
	OnResize(callback func()) EventSubscription

	// OnMouseMove subscribes f to be called whenever the mouse cursor moves over
	// the viewport.
	OnMouseMove(callback func(MouseEvent)) EventSubscription

	// OnMouseEnter subscribes f to be called whenever the mouse cursor enters the
	// viewport.
	OnMouseEnter(callback func(MouseEvent)) EventSubscription

	// OnMouseEnter subscribes f to be called whenever the mouse cursor leaves the
	// viewport.
	OnMouseExit(callback func(MouseEvent)) EventSubscription

	// OnMouseDown subscribes f to be called whenever a mouse button is pressed
	// while the cursor is inside the viewport.
	OnMouseDown(callback func(MouseEvent)) EventSubscription

	// OnMouseUp subscribes f to be called whenever a mouse button is released
	// while the cursor is inside the viewport.
	OnMouseUp(callback func(MouseEvent)) EventSubscription

	// OnMouseScroll subscribes f to be called whenever the mouse scroll wheel
	// turns while the cursor is inside the viewport.
	OnMouseScroll(callback func(MouseEvent)) EventSubscription

	// OnKeyDown subscribes f to be called whenever a keyboard key is pressed
	// while the viewport has focus.
	OnKeyDown(callback func(KeyboardEvent)) EventSubscription

	// OnKeyUp subscribes f to be called whenever a keyboard key is released
	// while the viewport has focus.
	OnKeyUp(callback func(KeyboardEvent)) EventSubscription

	// OnKeyRepeat subscribes f to be called whenever a keyboard key-repeat event
	// is raised while the viewport has focus.
	OnKeyRepeat(callback func(KeyboardEvent)) EventSubscription

	// OnKeyStroke subscribes f to be called whenever a keyboard key-stroke event
	// is raised while the viewport has focus.
	OnKeyStroke(callback func(KeyStrokeEvent)) EventSubscription
}

type Driver interface {
	// Call queues f to be run on the UI go-routine, returning before f may have been called.
	// Call returns false if the driver has been terminated, in which case f may not be called.
	Call(callback func()) bool

	// CallSync queues and then blocks for f to be run on the UI go-routine.
	// Call returns false if the driver has been terminated, in which case f may not be called.
	CallSync(callback func()) bool

	Terminate()
	SetClipboard(content string)
	GetClipboard() (content string, err error)

	// CreateFont loads a font from the provided TrueType bytes.
	CreateFont(data []byte, size int) (Font, error)

	// CreateWindowedViewport creates a new windowed Viewport with the specified width and height in device independent pixels.
	CreateWindowedViewport(width, height int, name string) Viewport

	// CreateFullscreenViewport creates a new fullscreen Viewport with the specified width and height in device independent pixels.
	// If width or height is 0, then the viewport adopts the current screen resolution.
	CreateFullscreenViewport(width, height int, name string) Viewport

	CreateCanvas(size math.Size) Canvas

	CreateTexture(img image.Image, pixelsPerDip float32) Texture

	// Debug function used to verify that the caller is executing on the UI go-routine. If the caller is not on the UI go-routine then the function panics.
	AssertUIGoroutine()
}
