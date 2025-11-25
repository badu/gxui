package purego

import (
	"image"

	"github.com/badu/gxui/pkg/math"
)

type textureContext struct {
	contextResource
	fn         *Functions
	sizePixels math.Size
	texture    Texture
	flipY      bool
	pma        bool
}

type TextureImpl struct {
	fn           *Functions
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func NewTexture(fn *Functions, fromImage image.Image, pixelsPerDip float32) *TextureImpl {
	result := &TextureImpl{
		fn:           fn,
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
	var format Enum
	var data []byte
	var pma bool

	switch imageType := t.image.(type) {
	case *image.RGBA:
		format = RGBA
		data = imageType.Pix
		pma = true
	case *image.NRGBA:
		format = RGBA
		data = imageType.Pix
	case *image.Alpha:
		format = ALPHA
		data = imageType.Pix
	default:
		panic("Unsupported image type")
	}

	glTexture := t.fn.CreateTexture()
	t.fn.BindTexture(TEXTURE_2D, glTexture)
	w, h := t.SizePixels().WH()
	t.fn.TexImage2D(TEXTURE_2D, 0, w, h, format, UNSIGNED_BYTE, data)
	t.fn.TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)
	t.fn.TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR)
	t.fn.BindTexture(TEXTURE_2D, Texture{})
	checkError(t.fn)

	globalStats.textureContextCount.inc()
	return &textureContext{
		fn:         t.fn,
		texture:    glTexture,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

func (c *textureContext) destroy() {
	globalStats.textureContextCount.dec()
	c.fn.DeleteTexture(c.texture)
	c.texture = Texture{}
}
