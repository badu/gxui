// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type PanelHolder struct {
	gxui.PanelHolderImpl
	theme *Theme
}

func CreatePanelHolder(theme *Theme) gxui.PanelHolder {
	p := &PanelHolder{}
	p.PanelHolderImpl.Init(p, theme)
	p.theme = theme
	p.SetMargin(math.Spacing{L: 0, T: 2, R: 0, B: 0})
	return p
}

func (p *PanelHolder) CreatePanelTab() gxui.PanelTab {
	return CreatePanelTab(p.theme)
}

func (p *PanelHolder) Paint(c gxui.Canvas) {
	panel := p.SelectedPanel()
	if panel != nil {
		bounds := p.Children().Find(panel).Bounds()
		c.DrawRoundedRect(bounds, 0.0, 0.0, 3.0, 3.0, p.theme.PanelBackgroundStyle.Pen, p.theme.PanelBackgroundStyle.Brush)
	}
	p.PanelHolderImpl.Paint(c)
}
