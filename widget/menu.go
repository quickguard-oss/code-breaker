package widget

import (
	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/rivo/tview"
)

// ゲームのメニューを表示するウィジェット
type menu struct {
	Section *tview.Form // tview
}

// menuウィジェットを生成する。
func NewMenu(app *tview.Application, widgets []Resettable) *menu {
	section := tview.NewForm()

	m := &menu{
		Section: section,
	}

	m.Section.
		AddButton("Reset", func() {
			gamestate.GetState().Reset()

			for _, w := range widgets {
				w.reset()
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetBorder(true).
		SetTitle(withPadding("Menu"))

	return m
}
