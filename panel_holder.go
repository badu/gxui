// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"github.com/badu/gxui/math"
)

type PanelHolder interface {
	Control
	AddPanel(panel Control, name string)
	AddPanelAt(panel Control, name string, index int)
	RemovePanel(panel Control)
	Select(int)
	PanelCount() int
	PanelIndex(Control) int
	Panel(int) Control
	Tab(int) Control
}

type PanelTab interface {
	Control
	SetText(string)
	SetActive(bool)
}

type PanelTabCreator interface {
	CreatePanelTab() PanelTab
}

type PanelHolderParent interface {
	ContainerBaseNoControlOuter
	PanelHolder
	PanelTabCreator
}

type PanelHolderImpl struct {
	ContainerBase
	parent    PanelHolderParent
	driver    Driver
	styles    *StyleDefs
	tabLayout LinearLayout
	entries   []PanelEntry
	selected  PanelEntry
}

func insertIndex(holder PanelHolder, at math.Point) int {
	count := holder.PanelCount()
	bestIndex := count
	bestScore := float32(1e20)
	score := func(point math.Point, index int) {
		score := point.Sub(at).Len()
		if score < bestScore {
			bestIndex = index
			bestScore = score
		}
	}
	for i := 0; i < holder.PanelCount(); i++ {
		tab := holder.Tab(i)
		size := tab.Size()
		ml := math.Point{Y: size.H / 2}
		mr := math.Point{Y: size.H / 2, X: size.W}
		score(TransformCoordinate(ml, tab, holder), i)
		score(TransformCoordinate(mr, tab, holder), i+1)
	}
	return bestIndex
}

func beginTabDragging(holder PanelHolder, panel Control, name string, window Window) {
	var mms, mos EventSubscription
	mms = window.OnMouseMove(
		func(event MouseEvent) {
			for _, c := range TopControlsUnder(event.WindowPoint, event.Window) {
				if over, ok := c.Control.(PanelHolder); ok {
					insertAt := insertIndex(over, c.Point)
					if over == holder {
						if insertAt > over.PanelIndex(panel) {
							insertAt--
						}
					}
					holder.RemovePanel(panel)
					holder = over
					holder.AddPanelAt(panel, name, insertAt)
					holder.Select(insertAt)
				}
			}
		},
	)
	mos = window.OnMouseUp(
		func(event MouseEvent) {
			mms.Forget()
			mos.Forget()
		},
	)
}

func (p *PanelHolderImpl) Init(parent PanelHolderParent, driver Driver, styles *StyleDefs) {
	p.ContainerBase.Init(parent, driver)

	p.parent = parent
	p.driver = driver
	p.styles = styles

	p.tabLayout = CreateLinearLayout(driver, styles)
	p.tabLayout.SetDirection(LeftToRight)
	p.ContainerBase.AddChild(p.tabLayout)
	p.SetMargin(math.Spacing{L: 1, T: 2, R: 1, B: 1})
	p.SetMouseEventTarget(true) // For drag-drop targets
}

func (p *PanelHolderImpl) LayoutChildren() {
	size := p.Size()

	tabHeight := p.tabLayout.DesiredSize(math.ZeroSize, size).H
	panelRect := math.CreateRect(0, tabHeight, size.W, size.H).Contract(p.Padding())

	for _, child := range p.Children() {
		if child.Control == p.tabLayout {
			child.Control.SetSize(math.Size{W: size.W, H: tabHeight})
			child.Offset = math.ZeroPoint
		} else {
			rect := panelRect.Contract(child.Control.Margin())
			child.Control.SetSize(rect.Size())
			child.Offset = rect.Min
		}
	}
}

func (p *PanelHolderImpl) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (p *PanelHolderImpl) SelectedPanel() Control {
	return p.selected.Panel
}

// gxui.PanelHolder compliance
func (p *PanelHolderImpl) AddPanel(panel Control, name string) {
	p.AddPanelAt(panel, name, len(p.entries))
}

func (p *PanelHolderImpl) AddPanelAt(panel Control, name string, index int) {
	if index < 0 || index > p.PanelCount() {
		panic(fmt.Errorf("index %d is out of bounds. Acceptable range: [%d - %d]", index, 0, p.PanelCount()))
	}
	tab := p.parent.CreatePanelTab()
	tab.SetText(name)
	mds := tab.OnMouseDown(
		func(ev MouseEvent) {
			p.Select(p.PanelIndex(panel))
			beginTabDragging(p.parent, panel, name, ev.Window)
		},
	)

	p.entries = append(p.entries, PanelEntry{})
	copy(p.entries[index+1:], p.entries[index:])
	p.entries[index] = PanelEntry{
		Panel:                 panel,
		Tab:                   tab,
		MouseDownSubscription: mds,
	}
	p.tabLayout.AddChildAt(index, tab)

	if p.selected.Panel == nil {
		p.Select(index)
	}
}

func (p *PanelHolderImpl) RemovePanel(panel Control) {
	index := p.PanelIndex(panel)
	if index < 0 {
		panic("PanelHolderImpl does not contain panel")
	}

	entry := p.entries[index]
	entry.MouseDownSubscription.Forget()
	p.entries = append(p.entries[:index], p.entries[index+1:]...)
	p.tabLayout.RemoveChildAt(index)

	if panel == p.selected.Panel {
		if p.PanelCount() > 0 {
			p.Select(math.Max(index-1, 0))
		} else {
			p.Select(-1)
		}
	}
}

func (p *PanelHolderImpl) Select(index int) {
	if index >= p.PanelCount() {
		panic(fmt.Errorf("index %d is out of bounds. Acceptable range: [%d - %d]", index, -1, p.PanelCount()-1))
	}

	if p.selected.Panel != nil {
		p.selected.Tab.SetActive(false)
		p.ContainerBase.RemoveChild(p.selected.Panel)
	}

	if index >= 0 {
		p.selected = p.entries[index]
	} else {
		p.selected = PanelEntry{}
	}

	if p.selected.Panel != nil {
		p.ContainerBase.AddChild(p.selected.Panel)
		p.selected.Tab.SetActive(true)
	}
}

func (p *PanelHolderImpl) PanelCount() int {
	return len(p.entries)
}

func (p *PanelHolderImpl) PanelIndex(panel Control) int {
	for i, e := range p.entries {
		if e.Panel == panel {
			return i
		}
	}
	return -1
}

func (p *PanelHolderImpl) Panel(index int) Control {
	return p.entries[index].Panel
}

func (p *PanelHolderImpl) Tab(index int) Control {
	return p.entries[index].Tab
}

type PanelEntry struct {
	Tab                   PanelTab
	Panel                 Control
	MouseDownSubscription EventSubscription
}
