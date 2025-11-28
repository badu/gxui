package cgo

import imageFont "golang.org/x/image/font"

type glyphTable struct {
	face  imageFont.Face
	index map[rune]int
	pages []*glyphPage
}

func newGlyphTable(face imageFont.Face) *glyphTable {
	return &glyphTable{face: face, index: make(map[rune]int)}
}

func (t *glyphTable) get(whichRune rune) *glyphPage {
	index, found := t.index[whichRune]
	if found {
		return t.pages[index]
	}

	if len(t.pages) == 0 {
		t.pages = append(t.pages, newGlyphPage(t.face, whichRune))
	} else {
		page := t.pages[len(t.pages)-1]
		if !page.add(t.face, whichRune) {
			page = newGlyphPage(t.face, whichRune)
			t.pages = append(t.pages, page)
		}
	}

	index = len(t.pages) - 1
	t.index[whichRune] = index
	return t.pages[index]
}
