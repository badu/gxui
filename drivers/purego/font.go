package purego

import (
	"fmt"
	"unicode"

	"github.com/badu/gxui"
	"github.com/badu/gxui/pkg/math"
	"github.com/golang/freetype/truetype"
	imageFont "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type font struct {
	ttf              *truetype.Font
	resolutions      map[resolution]*glyphTable
	glyphAdvanceDips map[rune]int
	glyphMaxSizeDips math.Size
	size             int
	ascentDips       int
	scale            fixed.Int26_6
}

func newFont(data []byte, size int) (*font, error) {
	ttf, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}

	scale := fixed.Int26_6(size << 6)
	bounds := rectangle26_6toRect(ttf.Bounds(scale))
	ascentDips := bounds.Max.Y

	return &font{
		size:             size,
		scale:            scale,
		glyphMaxSizeDips: bounds.Size(),
		ascentDips:       ascentDips,
		ttf:              ttf,
		resolutions:      make(map[resolution]*glyphTable),
		glyphAdvanceDips: make(map[rune]int),
	}, nil
}

func (f *font) advanceDips(ofRune rune) int {
	if g, found := f.glyphAdvanceDips[ofRune]; found {
		return g
	}

	idx := f.ttf.Index(ofRune)
	buffer := &truetype.GlyphBuf{}
	err := buffer.Load(f.ttf, f.scale, idx, imageFont.HintingFull)
	if err != nil {
		panic(err)
	}

	advance := int((buffer.AdvanceWidth + 0x3f) >> 6)
	f.glyphAdvanceDips[ofRune] = advance
	return advance
}

func (f *font) glyphTable(fn *Functions, resolution resolution) *glyphTable {
	result, found := f.resolutions[resolution]
	if !found {
		opt := truetype.Options{
			Size:              float64(f.size),
			DPI:               float64(resolution.intDipsToPixels(72)),
			Hinting:           imageFont.HintingFull,
			GlyphCacheEntries: 1,
			SubPixelsX:        1,
			SubPixelsY:        1,
		}
		result = newGlyphTable(fn, truetype.NewFace(f.ttf, &opt))
		f.resolutions[resolution] = result
	}
	return result
}

func (f *font) align(rect math.Rect, size math.Size, ascent int, horizontalAlignment gxui.HAlign, verticalAlignment gxui.VAlign) math.Point {
	var origin math.Point

	switch horizontalAlignment {
	case gxui.AlignLeft:
		origin.X = rect.Min.X
	case gxui.AlignCenter:
		origin.X = rect.Middle().X - (size.Width / 2)
	case gxui.AlignRight:
		origin.X = rect.Max.X - size.Width
	}

	switch verticalAlignment {
	case gxui.AlignTop:
		origin.Y = rect.Min.Y + ascent
	case gxui.AlignMiddle:
		origin.Y = rect.Middle().Y - (size.Height / 2) + ascent
	case gxui.AlignBottom:
		origin.Y = rect.Max.Y - size.Height + ascent
	}

	return origin
}

func (f *font) DrawRunes(fn *Functions, ctx *context, runes []rune, offsets []math.Point, color gxui.Color, state *drawState) {
	if len(runes) != len(offsets) {
		panic(fmt.Errorf("there must be the same number of runes to offsets. Got %d runes and %d offsets", len(runes), len(offsets)))
	}

	atResolution := ctx.resolution
	table := f.glyphTable(fn, atResolution)

	for runeIdx, curRune := range runes {
		if unicode.IsSpace(curRune) {
			continue
		}

		page := table.get(curRune)
		glyphTexture := page.texture()
		entry := page.get(curRune)
		srcRect := entry.bounds.Offset(entry.offset)
		dstRect := entry.bounds.Offset(atResolution.pointDipsToPixels(offsets[runeIdx]))
		textureCtx := ctx.getOrCreateTextureContext(glyphTexture)
		ctx.blitter.blitGlyph(ctx, textureCtx, color, srcRect, dstRect, state)
	}
}

func (f *font) Size() int {
	return f.size
}

func (f *font) Measure(textBlock *gxui.TextBlock) math.Size {
	size := math.Size{Width: 0, Height: f.glyphMaxSizeDips.Height}
	var offset math.Point
	for _, curRune := range textBlock.Runes {
		if curRune == '\n' {
			offset.X = 0
			offset.Y += f.glyphMaxSizeDips.Height
			continue
		}

		offset.X += f.advanceDips(curRune)
		size = size.Max(math.Size{Width: offset.X, Height: offset.Y + f.glyphMaxSizeDips.Height})
	}
	return size
}

func (f *font) Layout(textBlock *gxui.TextBlock) []math.Point {
	sizeDips := math.Size{}
	offsets := make([]math.Point, len(textBlock.Runes))
	var offset math.Point
	for i, r := range textBlock.Runes {
		if r == '\n' {
			offset.X = 0
			offset.Y += f.glyphMaxSizeDips.Height
			continue
		}

		offsets[i] = offset
		offset.X += f.advanceDips(r)
		sizeDips = sizeDips.Max(math.Size{Width: offset.X, Height: offset.Y + f.glyphMaxSizeDips.Height})
	}

	origin := f.align(textBlock.AlignRect, sizeDips, f.ascentDips, textBlock.H, textBlock.V)
	for i, p := range offsets {
		offsets[i] = p.Add(origin)
	}

	return offsets
}

func (f *font) LoadGlyphs(first, last rune) {
	if first > last {
		first, last = last, first
	}
	for r := first; r < last; r++ {
		f.advanceDips(r)
	}
}

func (f *font) GlyphMaxSize() math.Size {
	return f.glyphMaxSizeDips
}
