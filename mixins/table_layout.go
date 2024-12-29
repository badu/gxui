package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type Cell struct {
	x, y, w, h int
}

func (c Cell) AtColumn(x int) bool {
	return c.x <= x && c.x+c.w >= x
}

func (c Cell) AtRow(y int) bool {
	return c.y <= y && c.y+c.h >= y
}

type TableLayout struct {
	ContainerBase
	outer   ContainerBaseOuter
	grid    map[gxui.Control]Cell
	rows    int
	columns int
}

func (l *TableLayout) Init(outer ContainerBaseOuter, theme gxui.Theme) {
	l.ContainerBase.Init(outer, theme)
	l.outer = outer
	l.grid = make(map[gxui.Control]Cell)
}

func (l *TableLayout) LayoutChildren() {
	size := l.outer.Size().Contract(l.outer.Padding())
	offset := l.outer.Padding().LT()

	columnWidth, columnHeight := size.W/l.columns, size.H/l.rows

	var childRect math.Rect
	for _, child := range l.outer.Children() {
		childMargin := child.Control.Margin()
		cell := l.grid[child.Control]

		x, y := cell.x*columnWidth, cell.y*columnHeight
		w, h := x+cell.w*columnWidth, y+cell.h*columnHeight

		childRect = math.CreateRect(x+childMargin.L, y+childMargin.T, w-childMargin.R, h-childMargin.B)

		child.Layout(childRect.Offset(offset).Canon())
	}
}

func (l *TableLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *TableLayout) SetGrid(columns, rows int) {
	if l.columns != columns {
		if l.columns > columns {
			for c := l.columns; c > columns; c-- {
				for _, cell := range l.grid {
					if cell.AtColumn(c) {
						panic("Can't remove column with cells")
					}
				}
				l.columns--
			}
		} else {
			l.columns = columns
		}
	}

	if l.rows != rows {
		if l.rows > rows {
			for r := l.rows; r > rows; r-- {
				for _, cell := range l.grid {
					if cell.AtRow(r) {
						panic("Can't remove row with cells")
					}
				}
				l.rows--
			}
		} else {
			l.rows = rows
		}
	}

	if l.rows != rows || l.columns != columns {
		l.LayoutChildren()
	}
}

func (l *TableLayout) SetChildAt(x, y, w, h int, child gxui.Control) *gxui.Child {
	if x+w > l.columns || y+h > l.rows {
		panic("Cell is out of grid")
	}

	for _, c := range l.grid {
		if c.x+c.w > x && c.x < x+w && c.y+c.h > y && c.y < y+h {
			panic("Cell already has a child")
		}
	}

	l.grid[child] = Cell{x, y, w, h}
	return l.ContainerBase.AddChild(child)
}

func (l *TableLayout) RemoveChild(child gxui.Control) {
	delete(l.grid, child)
	l.ContainerBase.RemoveChild(child)
}
