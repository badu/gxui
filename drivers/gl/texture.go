// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"image"

	"github.com/badu/gxui/pkg/math"
	"github.com/goxjs/gl"
)

type TextureImpl struct {
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func NewTexture(fromImage image.Image, pixelsPerDip float32) *TextureImpl {
	result := &TextureImpl{
		image:        fromImage,
		pixelsPerDip: pixelsPerDip,
	}
	return result
}

// gxui.Texture compliance
func (t *TextureImpl) Image() image.Image {
	return t.image
}

func (t *TextureImpl) Size() math.Size {
	return t.SizePixels().ScaleS(1.0 / t.pixelsPerDip)
}

func (t *TextureImpl) SizePixels() math.Size {
	s := t.image.Bounds().Size()
	return math.Size{Width: s.X, Height: s.Y}
}

func (t *TextureImpl) FlipY() bool {
	return t.flipY
}

func (t *TextureImpl) SetFlipY(flipY bool) {
	t.flipY = flipY
}

func (t *TextureImpl) newContext() *textureContext {
	var format gl.Enum
	var data []byte
	var pma bool

	switch imageType := t.image.(type) {
	case *image.RGBA:
		format = gl.RGBA
		data = imageType.Pix
		pma = true
	case *image.NRGBA:
		format = gl.RGBA
		data = imageType.Pix
	case *image.Alpha:
		format = gl.ALPHA
		data = imageType.Pix
	default:
		panic("Unsupported image type")
	}

	glTexture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, glTexture)
	w, h := t.SizePixels().WH()
	gl.TexImage2D(gl.TEXTURE_2D, 0, w, h, format, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture{})
	checkError()

	globalStats.textureContextCount.inc()
	return &textureContext{
		texture:    glTexture,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

type textureContext struct {
	contextResource
	sizePixels math.Size
	texture    gl.Texture
	flipY      bool
	pma        bool
}

func (c *textureContext) destroy() {
	globalStats.textureContextCount.dec()
	gl.DeleteTexture(c.texture)
	c.texture = gl.Texture{}
}
