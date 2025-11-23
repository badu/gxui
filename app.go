// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"time"

	"github.com/badu/gxui/pkg/math"
)

type StyleDefs struct {
	DefaultFont          Font
	DefaultMonospaceFont Font

	BubbleOverlayStyle Style

	ButtonDefaultStyle Style
	ButtonOverStyle    Style
	ButtonPressedStyle Style

	CodeSuggestionListStyle Style
	CodeEditorStyle         Style

	DropDownListDefaultStyle Style
	DropDownListOverStyle    Style

	FocusedStyle   Style
	HighlightStyle Style

	LabelStyle Style

	PanelBackgroundStyle Style

	ScrollBarBarDefaultStyle  Style
	ScrollBarBarOverStyle     Style
	ScrollBarRailDefaultStyle Style
	ScrollBarRailOverStyle    Style

	SplitterBarDefaultStyle Style
	SplitterBarOverStyle    Style

	TabActiveHighlightStyle Style
	TabDefaultStyle         Style
	TabOverStyle            Style
	TabPressedStyle         Style

	TextBoxDefaultStyle Style
	TextBoxOverStyle    Style

	ScreenWidth  int
	ScreenHeight int
	FontSize     int

	WindowBackground Color
}

func CreateBubbleOverlay(driver Driver, styles *StyleDefs) *BubbleOverlay {
	result := &BubbleOverlay{}
	result.Init(result, driver)
	result.SetMargin(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	result.SetPadding(math.Spacing{Left: 5, Top: 5, Right: 5, Bottom: 5})
	result.SetPen(styles.BubbleOverlayStyle.Pen)
	result.SetBrush(styles.BubbleOverlayStyle.Brush)
	return result
}

func CreateButton(driver Driver, styles *StyleDefs) *Button {
	result := &Button{}
	result.Init(result, driver, styles)

	result.OnMouseEnter(func(event MouseEvent) { result.Redraw() })
	result.OnMouseExit(func(event MouseEvent) { result.Redraw() })
	result.OnMouseDown(func(event MouseEvent) { result.Redraw() })
	result.OnMouseUp(func(event MouseEvent) { result.Redraw() })
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	return result
}

func CreateCodeEditor(driver Driver, styles *StyleDefs) *AppCodeEditor {
	result := &AppCodeEditor{}
	result.Init(result, driver, styles)
	result.SetTextColor(styles.TextBoxDefaultStyle.FontColor)
	result.SetMargin(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	result.SetPadding(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	result.SetBorderPen(TransparentPen)
	return result
}

func CreateDropDownList(driver Driver, styles *StyleDefs) *DropDownList {
	result := &DropDownList{}
	result.Init(result, driver, styles)
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	result.List().OnAttach(result.Redraw)
	result.List().OnDetach(result.Redraw)
	result.OnMouseEnter(
		func(event MouseEvent) {
			result.SetBorderPen(styles.DropDownListOverStyle.Pen)
		},
	)
	result.OnMouseExit(
		func(event MouseEvent) {
			result.SetBorderPen(styles.DropDownListDefaultStyle.Pen)
		},
	)
	result.SetPadding(math.CreateSpacing(2))
	result.SetBorderPen(styles.DropDownListDefaultStyle.Pen)
	result.SetBackgroundBrush(styles.DropDownListDefaultStyle.Brush)
	return result
}

func CreateImage(driver Driver, styles *StyleDefs) *Image {
	result := &Image{}
	result.Init(result, driver)
	return result
}

func CreateLabel(driver Driver, styles *StyleDefs) *Label {
	result := &Label{}
	result.Init(result, driver, styles)
	result.SetMargin(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	return result
}

func CreateLinearLayout(driver Driver, styles *StyleDefs) *LinearLayoutImpl {
	result := &LinearLayoutImpl{}
	result.Init(result, driver)
	return result
}

func CreateList(driver Driver, styles *StyleDefs) *ListImpl {
	result := &ListImpl{}
	result.Init(result, driver, styles)
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	result.SetPadding(math.CreateSpacing(2))
	result.SetBorderPen(TransparentPen)
	return result
}

func CreatePanelHolder(driver Driver, styles *StyleDefs) *AppPanelHolder {
	result := &AppPanelHolder{}
	result.Init(result, driver, styles)
	result.SetMargin(math.Spacing{Left: 0, Top: 2, Right: 0, Bottom: 0})
	return result
}

type AppPanelHolder struct {
	PanelHolderImpl
}

func (p *AppPanelHolder) CreatePanelTab() PanelTab {
	result := &AppPanelTab{}
	result.Button.Init(result, p.driver, p.styles)
	result.Button.SetType(ToggleButton) // TODO : @Badu - setting ToggleButton requires POST Init() call (in Init() we set it to PushButton
	result.SetPadding(math.Spacing{Left: 5, Top: 3, Right: 5, Bottom: 3})
	result.OnMouseEnter(func(MouseEvent) { result.Redraw() })
	result.OnMouseExit(func(MouseEvent) { result.Redraw() })
	result.OnMouseDown(func(MouseEvent) { result.Redraw() })
	result.OnMouseUp(func(MouseEvent) { result.Redraw() })
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	return result
}

func CreateProgressBar(driver Driver, styles *StyleDefs) *AppProgressBar {
	result := &AppProgressBar{}
	result.Init(result, driver, styles)
	result.chevronWidth = 10

	result.OnAttach(
		func() {

			result.ticker = time.NewTicker(time.Millisecond * 50)
			go func() {
				for _ = range result.ticker.C {
					if !driver.Call(result.animationTick) {
						return
					}
				}
			}()
		},
	)

	result.OnDetach(
		func() {
			if result.chevrons != nil {
				result.chevrons = nil
				result.ticker.Stop()
				result.ticker = nil
			}
		},
	)

	result.SetBackgroundBrush(CreateBrush(Gray10))
	result.SetBorderPen(CreatePen(1, Gray40))
	return result
}

func CreateScrollBar(driver Driver, styles *StyleDefs) *ScrollBarImpl {
	result := &ScrollBarImpl{}
	result.Init(result, driver)
	result.SetBarBrush(styles.ScrollBarBarDefaultStyle.Brush)
	result.SetBarPen(styles.ScrollBarBarDefaultStyle.Pen)
	result.SetRailBrush(styles.ScrollBarRailDefaultStyle.Brush)
	result.SetRailPen(styles.ScrollBarRailDefaultStyle.Pen)
	updateColors := func() {
		switch {
		case result.IsMouseOver():
			result.SetBarBrush(styles.ScrollBarBarOverStyle.Brush)
			result.SetBarPen(styles.ScrollBarBarOverStyle.Pen)
			result.SetRailBrush(styles.ScrollBarRailOverStyle.Brush)
			result.SetRailPen(styles.ScrollBarRailOverStyle.Pen)
		default:
			result.SetBarBrush(styles.ScrollBarBarDefaultStyle.Brush)
			result.SetBarPen(styles.ScrollBarBarDefaultStyle.Pen)
			result.SetRailBrush(styles.ScrollBarRailDefaultStyle.Brush)
			result.SetRailPen(styles.ScrollBarRailDefaultStyle.Pen)
		}
		result.Redraw()
	}
	result.OnMouseEnter(func(event MouseEvent) { updateColors() })
	result.OnMouseExit(func(event MouseEvent) { updateColors() })
	return result
}

func CreateScrollLayout(driver Driver, styles *StyleDefs) *ScrollLayoutImpl {
	result := &ScrollLayoutImpl{}
	result.Init(result, driver, styles)
	return result
}

func CreateSplitterLayout(driver Driver, styles *StyleDefs) *AppSplitterLayout {
	result := &AppSplitterLayout{}
	result.Init(result, driver, styles)
	return result
}

func CreateTableLayout(driver Driver, styles *StyleDefs) *TableLayoutImpl {
	result := &TableLayoutImpl{}
	result.Init(result, driver)
	return result
}

func CreateTextBox(driver Driver, styles *StyleDefs) *TextBox {
	result := &TextBox{}
	result.Init(result, driver, styles, styles.DefaultFont)
	result.SetTextColor(styles.TextBoxDefaultStyle.FontColor)
	result.SetMargin(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	result.SetPadding(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	result.SetBackgroundBrush(styles.TextBoxDefaultStyle.Brush)
	result.SetBorderPen(styles.TextBoxDefaultStyle.Pen)

	result.OnMouseEnter(
		func(event MouseEvent) {
			result.SetBackgroundBrush(styles.TextBoxOverStyle.Brush)
			result.SetBorderPen(styles.TextBoxOverStyle.Pen)
		},
	)

	result.OnMouseExit(
		func(event MouseEvent) {
			result.SetBackgroundBrush(styles.TextBoxDefaultStyle.Brush)
			result.SetBorderPen(styles.TextBoxDefaultStyle.Pen)
		},
	)
	return result
}

func CreateTree(driver Driver, styles *StyleDefs) *AppTree {
	result := &AppTree{}
	result.Init(result, driver, styles)
	result.SetPadding(math.Spacing{Left: 3, Top: 3, Right: 3, Bottom: 3})
	result.SetBorderPen(TransparentPen)
	result.SetControlCreator(treeControlCreator{})
	return result
}

func CreateWindow(driver Driver, styles *StyleDefs, width, height int, title string) *WindowImpl {
	result := &WindowImpl{}
	result.Init(result, driver, width, height, title)
	result.SetBackgroundBrush(CreateBrush(styles.WindowBackground))
	return result
}

type Style struct {
	Font      Font
	FontColor Color
	Brush     Brush
	Pen       Pen
	VAlign    VAlign
	HAlign    HAlign
}

func CreateStyle(fontColor, brushColor, penColor Color, penWidth float32, font Font) Style {
	return Style{
		FontColor: fontColor,
		Pen:       CreatePen(penWidth, penColor),
		Brush:     CreateBrush(brushColor),
		Font:      font,
	}
}

type AppCodeEditor struct {
	CodeEditor
}

// mixins.CodeEditor overrides
func (t *AppCodeEditor) Paint(canvas Canvas) {
	t.CodeEditor.Paint(canvas)

	if t.HasFocus() {
		rect := t.Size().Rect()
		canvas.DrawRoundedRect(rect, 3, 3, 3, 3, t.styles.FocusedStyle.Pen, t.styles.FocusedStyle.Brush)
	}
}

func (t *AppCodeEditor) CreateSuggestionList() *ListImpl {
	result := CreateList(t.driver, t.styles)
	result.SetBackgroundBrush(t.styles.CodeSuggestionListStyle.Brush)
	result.SetBorderPen(t.styles.CodeSuggestionListStyle.Pen)
	result.SetPadding(math.CreateSpacing(10))
	result.SetBorderPen(WhitePen)
	return result
}

type AppPanelTab struct {
	Button
	active bool
}

func (t *AppPanelTab) SetActive(active bool) {
	t.active = active
	t.Redraw()
}

func (t *AppPanelTab) Paint(canvas Canvas) {
	size := t.Size()
	var style Style
	switch {
	case t.IsMouseDown(MouseButtonLeft) && t.IsMouseOver():
		style = t.styles.TabPressedStyle
	case t.IsMouseOver():
		style = t.styles.TabOverStyle
	default:
		style = t.styles.TabDefaultStyle
	}
	if l := t.Label(); l != nil {
		l.SetColor(style.FontColor)
	}

	canvas.DrawRoundedRect(size.Rect(), 5.0, 5.0, 0.0, 0.0, style.Pen, style.Brush)

	if t.HasFocus() {
		style = t.styles.FocusedStyle
		r := math.CreateRect(1, 1, size.Width-1, size.Height-1)
		canvas.DrawRoundedRect(r, 4.0, 4.0, 0.0, 0.0, style.Pen, style.Brush)
	}

	if t.active {
		style = t.styles.TabActiveHighlightStyle
		r := math.CreateRect(1, size.Height-1, size.Width-1, size.Height)
		canvas.DrawRect(r, style.Brush)
	}

	t.Button.Paint(canvas)
}

type AppProgressBar struct {
	ProgressBarImpl
	chevrons     Canvas
	ticker       *time.Ticker
	chevronWidth int
	scroll       int
}

func (b *AppProgressBar) animationTick() {
	if b.Attached() {
		b.scroll = (b.scroll + 1) % (b.chevronWidth * 2)
		b.Redraw()
	}
}

func (b *AppProgressBar) SetSize(size math.Size) {
	b.ProgressBarImpl.SetSize(size)

	b.chevrons = nil
	if size.Area() > 0 {
		b.chevrons = b.ControlBase.driver.CreateCanvas(size)
		b.chevronWidth = size.Height / 2
		cw := b.chevronWidth
		for x := -cw * 2; x < size.Width; x += cw * 2 {
			// x0    x2
			// |  x1 |  x3
			//    |     |
			// A-----B    - y0
			//  \     \
			//   \     \
			//    F     C - y1
			//   /     /
			//  /     /
			// E-----D    - y2
			y0, y1, y2 := 0, size.Height/2, size.Height
			x0, x1 := x, x+cw/2
			x2, x3 := x0+cw, x1+cw
			var chevron = Polygon{
				/* A */ PolygonVertex{Position: math.Point{X: x0, Y: y0}},
				/* B */ PolygonVertex{Position: math.Point{X: x2, Y: y0}},
				/* C */ PolygonVertex{Position: math.Point{X: x3, Y: y1}},
				/* D */ PolygonVertex{Position: math.Point{X: x2, Y: y2}},
				/* E */ PolygonVertex{Position: math.Point{X: x0, Y: y2}},
				/* F */ PolygonVertex{Position: math.Point{X: x1, Y: y1}},
			}
			b.chevrons.DrawPolygon(chevron, TransparentPen, CreateBrush(Gray30))
		}
		b.chevrons.Complete()
	}
}

func (b *AppProgressBar) PaintProgress(canvas Canvas, rect math.Rect, frac float32) {
	rect.Max.X = math.Lerp(rect.Min.X, rect.Max.X, frac)
	canvas.DrawRect(rect, CreateBrush(Gray50))
	canvas.Push()
	canvas.AddClip(rect)
	canvas.DrawCanvas(b.chevrons, math.Point{X: b.scroll})
	canvas.Pop()
}

type AppSplitterLayout struct {
	SplitterLayoutImpl
}

// mixins.SplitterLayoutImpl overrides
func (l *AppSplitterLayout) CreateSplitterBar() Control {
	result := &SplitterBar{}
	result.Init(result, l.driver, l.styles)
	result.BackgroundColor = l.styles.SplitterBarDefaultStyle.Brush.Color
	result.ForegroundColor = l.styles.SplitterBarDefaultStyle.Pen.Color
	result.OnSplitterDragged(func(wndPnt math.Point) { l.SplitterDragged(result, wndPnt) })
	updateForegroundColor := func() {
		switch {
		case result.IsDragging:
			result.ForegroundColor = l.styles.HighlightStyle.Pen.Color
		case result.IsMouseOver():
			result.ForegroundColor = l.styles.SplitterBarOverStyle.Pen.Color
		default:
			result.ForegroundColor = l.styles.SplitterBarDefaultStyle.Pen.Color
		}
		result.Redraw()
	}
	result.OnDragStart(func(event MouseEvent) { updateForegroundColor() })
	result.OnDragEnd(func(event MouseEvent) { updateForegroundColor() })
	result.OnDragStart(func(event MouseEvent) { updateForegroundColor() })
	result.OnMouseEnter(func(event MouseEvent) { updateForegroundColor() })
	result.OnMouseExit(func(event MouseEvent) { updateForegroundColor() })
	return result
}

type AppTree struct {
	TreeImpl
}

// mixins.TreeImpl overrides
func (t *AppTree) Paint(canvas Canvas) {
	rect := t.Size().Rect()

	t.TreeImpl.Paint(canvas)

	if t.HasFocus() {
		style := t.styles.FocusedStyle
		canvas.DrawRoundedRect(rect, 3, 3, 3, 3, style.Pen, style.Brush)
	}
}

func (t *AppTree) PaintMouseOverBackground(canvas Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, TransparentPen, CreateBrush(Gray15))
}

// mixins.ListImpl overrides
func (t *AppTree) PaintSelection(canvas Canvas, rect math.Rect) {
	style := t.styles.HighlightStyle
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, style.Pen, style.Brush)
}

type treeControlCreator struct{}

func (treeControlCreator) Create(driver Driver, styles *StyleDefs, control Control, node *TreeToListNode) Control {
	img := CreateImage(driver, styles)
	imgSize := math.Size{Width: 10, Height: 10}

	layout := CreateLinearLayout(driver, styles)
	layout.SetDirection(LeftToRight)

	btn := CreateButton(driver, styles)
	btn.SetBackgroundBrush(TransparentBrush)
	btn.SetBorderPen(CreatePen(1, Gray30))
	btn.SetMargin(math.Spacing{Left: 1, Right: 1, Top: 1, Bottom: 1})
	btn.OnClick(func(ev MouseEvent) {
		if ev.Button == MouseButtonLeft {
			node.ToggleExpanded()
		}
	})
	btn.AddChild(img)

	var expandedPoly = Polygon{
		PolygonVertex{Position: math.Point{X: 2, Y: 3}},
		PolygonVertex{Position: math.Point{X: 8, Y: 3}},
		PolygonVertex{Position: math.Point{X: 5, Y: 8}},
	}

	var collapsedPoly = Polygon{
		PolygonVertex{Position: math.Point{X: 3, Y: 2}},
		PolygonVertex{Position: math.Point{X: 8, Y: 5}},
		PolygonVertex{Position: math.Point{X: 3, Y: 8}},
	}
	update := func() {
		expanded := node.IsExpanded()
		canvas := driver.CreateCanvas(imgSize)
		btn.SetVisible(!node.IsLeaf())
		switch {
		case !btn.IsMouseDown(MouseButtonLeft) && expanded:
			canvas.DrawPolygon(expandedPoly, TransparentPen, CreateBrush(Gray70))
		case !btn.IsMouseDown(MouseButtonLeft) && !expanded:
			canvas.DrawPolygon(collapsedPoly, TransparentPen, CreateBrush(Gray70))
		case expanded:
			canvas.DrawPolygon(expandedPoly, TransparentPen, CreateBrush(Gray30))
		case !expanded:
			canvas.DrawPolygon(collapsedPoly, TransparentPen, CreateBrush(Gray30))
		}
		canvas.Complete()
		img.SetCanvas(canvas)
	}
	btn.OnMouseDown(func(event MouseEvent) { update() })
	btn.OnMouseUp(func(event MouseEvent) { update() })
	update()

	WhileAttached(btn, node.OnChange, update)

	layout.AddChild(btn)
	layout.AddChild(control)
	layout.SetPadding(math.Spacing{Left: 16 * node.Depth()})
	return layout
}

func (treeControlCreator) Size(styles *StyleDefs, treeControlSize math.Size) math.Size {
	return treeControlSize
}
