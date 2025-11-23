package gxui

import (
	"github.com/badu/gxui/pkg/math"
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

type TableLayoutImpl struct {
	ContainerBase
	parent  BaseContainerParent
	grid    map[Control]Cell
	rows    int
	columns int
}

func (l *TableLayoutImpl) Init(parent BaseContainerParent, driver Driver) {
	l.ContainerBase.Init(parent, driver)
	l.parent = parent
	l.grid = make(map[Control]Cell)
}

func (l *TableLayoutImpl) LayoutChildren() {
	size := l.parent.Size().Contract(l.parent.Padding())
	offset := l.parent.Padding().TopLeft()

	columnWidth, columnHeight := size.Width/l.columns, size.Height/l.rows

	var childRect math.Rect
	for _, child := range l.parent.Children() {
		childMargin := child.Control.Margin()
		cell := l.grid[child.Control]

		x, y := cell.x*columnWidth, cell.y*columnHeight
		w, h := x+cell.w*columnWidth, y+cell.h*columnHeight

		childRect = math.CreateRect(x+childMargin.Left, y+childMargin.Top, w-childMargin.Right, h-childMargin.Bottom)

		child.Layout(childRect.Offset(offset).Canon())
	}
}

func (l *TableLayoutImpl) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *TableLayoutImpl) SetGrid(columns, rows int) {
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

func (l *TableLayoutImpl) SetChildAt(x, y, w, h int, child Control) *Child {
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

func (l *TableLayoutImpl) RemoveChild(child Control) {
	delete(l.grid, child)
	l.ContainerBase.RemoveChild(child)
}
