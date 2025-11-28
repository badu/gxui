package cgo

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/badu/gxui/pkg/math"
	imageFont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	dumpGlyphPages     = false
	glyphPageWidth     = 512
	glyphPageHeight    = 512
	glyphSizeAlignment = 8
	glyphPadding       = 1
)

type glyphEntry struct {
	offset math.Point
	bounds math.Rect
}

type glyphPage struct {
	image     *image.Alpha
	entries   map[rune]glyphEntry
	tex       *TextureImpl
	size      math.Size // in pixels
	nextPoint math.Point
	rowHeight int
}

func point26_6toPoint(point fixed.Point26_6) math.Point {
	return math.Point{X: int(point.X) >> 6, Y: int(point.Y) >> 6}
}

func rectangle26_6toRect(point fixed.Rectangle26_6) math.Rect {
	return math.Rect{Min: point26_6toPoint(point.Min), Max: point26_6toPoint(point.Max)}
}

func align(width, size int) int {
	return (width + size - 1) & ^(size - 1)
}

func newGlyphPage(face imageFont.Face, whichRune rune) *glyphPage {
	// Start the page big enough to hold the initial rune.
	glyphBounds, _, _ := face.GlyphBounds(whichRune)
	bounds := rectangle26_6toRect(glyphBounds)
	size := math.Size{Width: glyphPageWidth, Height: glyphPageHeight}.Max(bounds.Size())
	size.Width = align(size.Width, glyphSizeAlignment)
	size.Height = align(size.Height, glyphSizeAlignment)

	page := &glyphPage{
		image:     image.NewAlpha(image.Rect(0, 0, size.Width, size.Height)),
		size:      size,
		entries:   make(map[rune]glyphEntry),
		rowHeight: 0,
	}
	page.add(face, whichRune)
	return page
}

func (p *glyphPage) commit() {
	if p.tex != nil {
		return
	}

	p.tex = NewTexture(p.image, 1.0)
	if dumpGlyphPages {
		f, _ := os.Create("glyph-page.png")
		defer func() {
			err := f.Close()
			if err != nil {
				fmt.Println("error closing glyph-page.png:" + err.Error())
			}
		}()
		if err := png.Encode(f, p.image); err != nil {
			fmt.Println("error encoding glyph-page.png:" + err.Error())
		}
	}
}

func (p *glyphPage) add(face imageFont.Face, whichRune rune) bool {
	if _, found := p.entries[whichRune]; found {
		panic("Glyph already added to glyph page")
	}

	glyphBounds, mask, maskp, _, _ := face.Glyph(fixed.Point26_6{}, whichRune)
	bounds := math.CreateRect(glyphBounds.Min.X, glyphBounds.Min.Y, glyphBounds.Max.X, glyphBounds.Max.Y)

	w, h := bounds.Size().WH()
	x, y := p.nextPoint.X, p.nextPoint.Y

	if x+w > p.size.Width {
		// Row full, start new line
		x = 0
		y += p.rowHeight + glyphPadding
		p.rowHeight = 0
	}

	if y+h > p.size.Height {
		return false // Page full
	}

	draw.Draw(p.image, image.Rect(x, y, x+w, y+h), mask, maskp, draw.Src)

	p.entries[whichRune] = glyphEntry{
		offset: math.Point{X: x, Y: y}.Sub(bounds.Min),
		bounds: bounds,
	}
	p.nextPoint = math.Point{X: x + w + glyphPadding, Y: y}
	if h > p.rowHeight {
		p.rowHeight = h
	}
	p.tex = nil

	return true
}

func (p *glyphPage) texture() *TextureImpl {
	if p.tex == nil {
		p.commit()
	}
	return p.tex
}

func (p *glyphPage) get(rune rune) glyphEntry {
	return p.entries[rune]
}
