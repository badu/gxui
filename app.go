// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
	"time"
)

type App interface {
	Driver() Driver
	DefaultFont() Font
	SetDefaultFont(font Font)
	DefaultMonospaceFont() Font
	SetDefaultMonospaceFont(font Font)
	CreateBubbleOverlay() BubbleOverlay
	CreateButton() Button
	CreateCodeEditor() CodeEditor
	CreateDropDownList() DropDownList
	CreateImage() Image
	CreateLabel() Label
	CreateLinearLayout() LinearLayout
	CreateList() List
	CreatePanelHolder() PanelHolder
	CreateProgressBar() ProgressBar
	CreateScrollBar() ScrollBar
	CreateScrollLayout() ScrollLayout
	CreateSplitterLayout() SplitterLayout
	CreateTableLayout() TableLayout
	CreateTextBox() TextBox
	CreateTree() Tree
	CreateWindow(width, height int, title string) Window
	DisplayWidth() int
	DisplayHeight() int
	DefaultFontSize() int
}

type DefaultApp struct {
	DriverInfo               Driver
	DefaultFontInfo          Font
	DefaultMonospaceFontInfo Font

	WindowBackground Color

	BubbleOverlayStyle Style

	ButtonDefaultStyle Style
	ButtonOverStyle    Style
	ButtonPressedStyle Style

	CodeSuggestionListStyle Style

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
}

// gxui.App compliance
func (a *DefaultApp) Driver() Driver {
	return a.DriverInfo
}

func (a *DefaultApp) DefaultFont() Font {
	return a.DefaultFontInfo
}

func (a *DefaultApp) SetDefaultFont(f Font) {
	a.DefaultFontInfo = f
}

func (a *DefaultApp) DefaultMonospaceFont() Font {
	return a.DefaultMonospaceFontInfo
}

func (a *DefaultApp) SetDefaultMonospaceFont(font Font) {
	a.DefaultMonospaceFontInfo = font
}

func (a *DefaultApp) CreateBubbleOverlay() BubbleOverlay {
	result := &BubbleOverlayImpl{}
	result.Init(result, a)
	result.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetPadding(math.Spacing{L: 5, T: 5, R: 5, B: 5})
	result.SetPen(a.BubbleOverlayStyle.Pen)
	result.SetBrush(a.BubbleOverlayStyle.Brush)
	return result
}

func (a *DefaultApp) CreateButton() Button {
	result := &AppButton{}
	result.Init(result, a)
	result.app = a
	result.SetPadding(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetBackgroundBrush(result.app.ButtonDefaultStyle.Brush)
	result.SetBorderPen(result.app.ButtonDefaultStyle.Pen)
	result.OnMouseEnter(func(event MouseEvent) { result.Redraw() })
	result.OnMouseExit(func(event MouseEvent) { result.Redraw() })
	result.OnMouseDown(func(event MouseEvent) { result.Redraw() })
	result.OnMouseUp(func(event MouseEvent) { result.Redraw() })
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	return result
}

func (a *DefaultApp) CreateCodeEditor() CodeEditor {
	result := &AppCodeEditor{}
	result.app = a
	result.Init(result, a.Driver(), a, a.DefaultMonospaceFont())
	result.SetTextColor(a.TextBoxDefaultStyle.FontColor)
	result.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetPadding(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetBorderPen(TransparentPen)
	return result
}

func (a *DefaultApp) CreateDropDownList() DropDownList {
	result := &AppDropDownList{}
	result.Init(result, a)
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	result.List().OnAttach(result.Redraw)
	result.List().OnDetach(result.Redraw)
	result.OnMouseEnter(
		func(event MouseEvent) {
			result.SetBorderPen(a.DropDownListOverStyle.Pen)
		},
	)
	result.OnMouseExit(
		func(event MouseEvent) {
			result.SetBorderPen(a.DropDownListDefaultStyle.Pen)
		},
	)
	result.SetPadding(math.CreateSpacing(2))
	result.SetBorderPen(a.DropDownListDefaultStyle.Pen)
	result.SetBackgroundBrush(a.DropDownListDefaultStyle.Brush)
	result.app = a
	return result
}

func (a *DefaultApp) CreateImage() Image {
	result := &ImageImpl{}
	result.Init(result, a)
	return result
}

func (a *DefaultApp) CreateLabel() Label {
	result := &LabelImpl{}
	result.Init(result, a, a.DefaultFont(), a.LabelStyle.FontColor)
	result.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	return result
}

func (a *DefaultApp) CreateLinearLayout() LinearLayout {
	result := &LinearLayoutImpl{}
	result.Init(result, a)
	return result
}

func (a *DefaultApp) CreateList() List {
	result := &AppList{}
	result.Init(result, a)
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	result.SetPadding(math.CreateSpacing(2))
	result.SetBorderPen(TransparentPen)
	result.app = a
	return result
}

func (a *DefaultApp) CreatePanelHolder() PanelHolder {
	result := &AppPanelHolder{}
	result.PanelHolderImpl.Init(result, a)
	result.app = a
	result.SetMargin(math.Spacing{L: 0, T: 2, R: 0, B: 0})
	return result
}

func (a *DefaultApp) CreateProgressBar() ProgressBar {
	result := &AppProgressBar{}
	result.Init(result, a)
	result.app = a
	result.chevronWidth = 10

	result.OnAttach(
		func() {
			driver := a.Driver()
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

func (a *DefaultApp) CreateScrollBar() ScrollBar {
	result := &ScrollBarImpl{}
	result.Init(result, a)
	result.SetBarBrush(a.ScrollBarBarDefaultStyle.Brush)
	result.SetBarPen(a.ScrollBarBarDefaultStyle.Pen)
	result.SetRailBrush(a.ScrollBarRailDefaultStyle.Brush)
	result.SetRailPen(a.ScrollBarRailDefaultStyle.Pen)
	updateColors := func() {
		switch {
		case result.IsMouseOver():
			result.SetBarBrush(a.ScrollBarBarOverStyle.Brush)
			result.SetBarPen(a.ScrollBarBarOverStyle.Pen)
			result.SetRailBrush(a.ScrollBarRailOverStyle.Brush)
			result.SetRailPen(a.ScrollBarRailOverStyle.Pen)
		default:
			result.SetBarBrush(a.ScrollBarBarDefaultStyle.Brush)
			result.SetBarPen(a.ScrollBarBarDefaultStyle.Pen)
			result.SetRailBrush(a.ScrollBarRailDefaultStyle.Brush)
			result.SetRailPen(a.ScrollBarRailDefaultStyle.Pen)
		}
		result.Redraw()
	}
	result.OnMouseEnter(func(event MouseEvent) { updateColors() })
	result.OnMouseExit(func(event MouseEvent) { updateColors() })
	return result
}

func (a *DefaultApp) CreateScrollLayout() ScrollLayout {
	result := &ScrollLayoutImpl{}
	result.Init(result, a)
	return result
}

func (a *DefaultApp) CreateSplitterLayout() SplitterLayout {
	result := &AppSplitterLayout{}
	result.app = a
	result.Init(result, a)
	return result
}

func (a *DefaultApp) CreateTableLayout() TableLayout {
	result := &TableLayoutImpl{}
	result.Init(result, a)
	return result
}

func (a *DefaultApp) CreateTextBox() TextBox {
	result := &AppTextBox{}
	result.Init(result, a.Driver(), a, a.DefaultFont())
	result.SetTextColor(a.TextBoxDefaultStyle.FontColor)
	result.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetPadding(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetBackgroundBrush(a.TextBoxDefaultStyle.Brush)
	result.SetBorderPen(a.TextBoxDefaultStyle.Pen)

	result.OnMouseEnter(
		func(event MouseEvent) {
			result.SetBackgroundBrush(a.TextBoxOverStyle.Brush)
			result.SetBorderPen(a.TextBoxOverStyle.Pen)
		},
	)

	result.OnMouseExit(
		func(event MouseEvent) {
			result.SetBackgroundBrush(a.TextBoxDefaultStyle.Brush)
			result.SetBorderPen(a.TextBoxDefaultStyle.Pen)
		},
	)

	result.app = a

	return result
}

func (a *DefaultApp) CreateTree() Tree {
	result := &AppTree{}
	result.Init(result, a)
	result.SetPadding(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	result.SetBorderPen(TransparentPen)
	result.app = a
	result.SetControlCreator(treeControlCreator{})
	return result
}

func (a *DefaultApp) CreateWindow(width, height int, title string) Window {
	result := &WindowImpl{}
	result.Init(result, a.Driver(), width, height, title)
	result.SetBackgroundBrush(CreateBrush(a.WindowBackground))
	return result
}

func (a *DefaultApp) DisplayWidth() int {
	return a.ScreenWidth
}

func (a *DefaultApp) DisplayHeight() int {
	return a.ScreenHeight
}

func (a *DefaultApp) DefaultFontSize() int {
	return a.FontSize
}

type Style struct {
	FontColor Color
	Brush     Brush
	Pen       Pen
}

func CreateStyle(fontColor, brushColor, penColor Color, penWidth float32) Style {
	return Style{
		FontColor: fontColor,
		Pen:       CreatePen(penWidth, penColor),
		Brush:     CreateBrush(brushColor),
	}
}

type AppButton struct {
	ButtonImpl
	app *DefaultApp
}

// Button internal overrides
func (b *AppButton) Paint(canvas Canvas) {
	pen := b.ButtonImpl.BorderPen()
	brush := b.ButtonImpl.BackgroundBrush()
	fontColor := b.app.ButtonDefaultStyle.FontColor

	switch {
	case b.IsMouseDown(MouseButtonLeft) && b.IsMouseOver():
		pen = b.app.ButtonPressedStyle.Pen
		brush = b.app.ButtonPressedStyle.Brush
		fontColor = b.app.ButtonPressedStyle.FontColor
	case b.IsMouseOver():
		pen = b.app.ButtonOverStyle.Pen
		brush = b.app.ButtonOverStyle.Brush
		fontColor = b.app.ButtonOverStyle.FontColor
	}

	if label := b.Label(); label != nil {
		label.SetColor(fontColor)
	}

	rect := b.Size().Rect()

	canvas.DrawRoundedRect(rect, 2, 2, 2, 2, TransparentPen, brush)

	b.PaintChildrenPart.Paint(canvas)

	canvas.DrawRoundedRect(rect, 2, 2, 2, 2, pen, TransparentBrush)

	if b.IsChecked() {
		pen = b.app.HighlightStyle.Pen
		brush = b.app.HighlightStyle.Brush
		canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, pen, brush)
	}

	if b.HasFocus() {
		pen = b.app.FocusedStyle.Pen
		brush = b.app.FocusedStyle.Brush
		canvas.DrawRoundedRect(rect.ContractI(int(pen.Width)), 3.0, 3.0, 3.0, 3.0, pen, brush)
	}
}

type AppCodeEditor struct {
	CodeEditorImpl
	app *DefaultApp
}

// mixins.CodeEditorImpl overrides
func (t *AppCodeEditor) Paint(canvas Canvas) {
	t.CodeEditorImpl.Paint(canvas)

	if t.HasFocus() {
		rect := t.Size().Rect()
		canvas.DrawRoundedRect(rect, 3, 3, 3, 3, t.app.FocusedStyle.Pen, t.app.FocusedStyle.Brush)
	}
}

func (t *AppCodeEditor) CreateSuggestionList() List {
	result := t.app.CreateList()
	result.SetBackgroundBrush(t.app.CodeSuggestionListStyle.Brush)
	result.SetBorderPen(t.app.CodeSuggestionListStyle.Pen)
	return result
}

type AppDropDownList struct {
	DropDownListImpl
	app *DefaultApp
}

// mixin.ListImpl overrides
func (l *AppDropDownList) Paint(canvas Canvas) {
	l.DropDownListImpl.Paint(canvas)
	if l.HasFocus() || l.ListShowing() {
		r := l.Size().Rect().ContractI(1)
		canvas.DrawRoundedRect(r, 3.0, 3.0, 3.0, 3.0, l.app.FocusedStyle.Pen, l.app.FocusedStyle.Brush)
	}
}

func (l *AppDropDownList) DrawSelection(c Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, l.app.HighlightStyle.Pen, l.app.HighlightStyle.Brush)
}

type AppList struct {
	ListImpl
	app *DefaultApp
}

// mixin.ListImpl overrides
func (l *AppList) Paint(canvas Canvas) {
	l.ListImpl.Paint(canvas)
	if l.HasFocus() {
		rect := l.Size().Rect().ContractI(1)
		canvas.DrawRoundedRect(rect, 3.0, 3.0, 3.0, 3.0, l.app.FocusedStyle.Pen, l.app.FocusedStyle.Brush)
	}
}

func (l *AppList) PaintSelection(c Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, l.app.HighlightStyle.Pen, l.app.HighlightStyle.Brush)
}

func (l *AppList) PaintMouseOverBackground(c Canvas, r math.Rect) {
	c.DrawRoundedRect(r, 2.0, 2.0, 2.0, 2.0, TransparentPen, CreateBrush(Gray15))
}

type AppPanelHolder struct {
	PanelHolderImpl
	app *DefaultApp
}

func (p *AppPanelHolder) CreatePanelTab() PanelTab {
	result := &AppPanelTab{}
	result.ButtonImpl.Init(result, p.app)
	result.app = p.app
	result.SetPadding(math.Spacing{L: 5, T: 3, R: 5, B: 3})
	result.OnMouseEnter(func(MouseEvent) { result.Redraw() })
	result.OnMouseExit(func(MouseEvent) { result.Redraw() })
	result.OnMouseDown(func(MouseEvent) { result.Redraw() })
	result.OnMouseUp(func(MouseEvent) { result.Redraw() })
	result.OnGainedFocus(result.Redraw)
	result.OnLostFocus(result.Redraw)
	return result
}

func (p *AppPanelHolder) Paint(c Canvas) {
	panel := p.SelectedPanel()
	if panel != nil {
		bounds := p.Children().Find(panel).Bounds()
		c.DrawRoundedRect(bounds, 0.0, 0.0, 3.0, 3.0, p.app.PanelBackgroundStyle.Pen, p.app.PanelBackgroundStyle.Brush)
	}
	p.PanelHolderImpl.Paint(c)
}

type AppPanelTab struct {
	ButtonImpl
	app    *DefaultApp
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
		style = t.app.TabPressedStyle
	case t.IsMouseOver():
		style = t.app.TabOverStyle
	default:
		style = t.app.TabDefaultStyle
	}
	if l := t.Label(); l != nil {
		l.SetColor(style.FontColor)
	}

	canvas.DrawRoundedRect(size.Rect(), 5.0, 5.0, 0.0, 0.0, style.Pen, style.Brush)

	if t.HasFocus() {
		style = t.app.FocusedStyle
		r := math.CreateRect(1, 1, size.W-1, size.H-1)
		canvas.DrawRoundedRect(r, 4.0, 4.0, 0.0, 0.0, style.Pen, style.Brush)
	}

	if t.active {
		style = t.app.TabActiveHighlightStyle
		r := math.CreateRect(1, size.H-1, size.W-1, size.H)
		canvas.DrawRect(r, style.Brush)
	}

	t.ButtonImpl.Paint(canvas)
}

type AppProgressBar struct {
	ProgressBarImpl
	app          *DefaultApp
	ticker       *time.Ticker
	chevrons     Canvas
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
		b.chevrons = b.app.Driver().CreateCanvas(size)
		b.chevronWidth = size.H / 2
		cw := b.chevronWidth
		for x := -cw * 2; x < size.W; x += cw * 2 {
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
			y0, y1, y2 := 0, size.H/2, size.H
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
	app *DefaultApp
}

// mixins.SplitterLayoutImpl overrides
func (l *AppSplitterLayout) CreateSplitterBar() Control {
	result := &SplitterBar{}
	result.Init(result, l.app)
	result.SetBackgroundColor(l.app.SplitterBarDefaultStyle.Brush.Color)
	result.SetForegroundColor(l.app.SplitterBarDefaultStyle.Pen.Color)
	result.OnSplitterDragged(func(wndPnt math.Point) { l.SplitterDragged(result, wndPnt) })
	updateForegroundColor := func() {
		switch {
		case result.IsDragging():
			result.SetForegroundColor(l.app.HighlightStyle.Pen.Color)
		case result.IsMouseOver():
			result.SetForegroundColor(l.app.SplitterBarOverStyle.Pen.Color)
		default:
			result.SetForegroundColor(l.app.SplitterBarDefaultStyle.Pen.Color)
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

type AppTextBox struct {
	TextBoxImpl
	app *DefaultApp
}

// mixins.TextBoxImpl overrides
func (t *AppTextBox) Paint(canvas Canvas) {
	t.TextBoxImpl.Paint(canvas)

	if t.HasFocus() {
		rect := t.Size().Rect()
		style := t.app.FocusedStyle
		canvas.DrawRoundedRect(rect, 3, 3, 3, 3, style.Pen, style.Brush)
	}
}

type AppTree struct {
	TreeImpl
	app *DefaultApp
}

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

// mixins.TreeImpl overrides
func (t *AppTree) Paint(canvas Canvas) {
	rect := t.Size().Rect()

	t.TreeImpl.Paint(canvas)

	if t.HasFocus() {
		style := t.app.FocusedStyle
		canvas.DrawRoundedRect(rect, 3, 3, 3, 3, style.Pen, style.Brush)
	}
}

func (t *AppTree) PaintMouseOverBackground(canvas Canvas, rect math.Rect) {
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, TransparentPen, CreateBrush(Gray15))
}

// mixins.ListImpl overrides
func (t *AppTree) PaintSelection(canvas Canvas, rect math.Rect) {
	style := t.app.HighlightStyle
	canvas.DrawRoundedRect(rect, 2.0, 2.0, 2.0, 2.0, style.Pen, style.Brush)
}

type treeControlCreator struct{}

func (treeControlCreator) Create(app App, control Control, node *TreeToListNode) Control {
	img := app.CreateImage()
	imgSize := math.Size{W: 10, H: 10}

	layout := app.CreateLinearLayout()
	layout.SetDirection(LeftToRight)

	btn := app.CreateButton()
	btn.SetBackgroundBrush(TransparentBrush)
	btn.SetBorderPen(CreatePen(1, Gray30))
	btn.SetMargin(math.Spacing{L: 1, R: 1, T: 1, B: 1})
	btn.OnClick(func(ev MouseEvent) {
		if ev.Button == MouseButtonLeft {
			node.ToggleExpanded()
		}
	})
	btn.AddChild(img)

	update := func() {
		expanded := node.IsExpanded()
		canvas := app.Driver().CreateCanvas(imgSize)
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
	layout.SetPadding(math.Spacing{L: 16 * node.Depth()})
	return layout
}

func (treeControlCreator) Size(app App, treeControlSize math.Size) math.Size {
	return treeControlSize
}
