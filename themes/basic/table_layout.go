package basic

import (
	"github.com/badu/gxui"
)

func CreateTableLayout(theme *Theme) gxui.TableLayout {
	l := &gxui.TableLayoutImpl{}
	l.Init(l, theme)
	return l
}
