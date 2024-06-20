package widget

import (
	"fmt"
	"reflect"

	"github.com/gdamore/tcell/v2"
	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/rivo/tview"
)

// これまでの解答を表示するウィジェット
type answer struct {
	Section *tview.Table // tview
	history [][]int      // 解答履歴
}

// answerウィジェットを生成する。
func NewAnswer() *answer {
	section := tview.NewTable()

	a := &answer{
		Section: section,
	}

	a.Section.
		SetBorders(true).
		SetBorder(true).
		SetTitle(withPadding("Answer"))

	a.reset()

	return a
}

// answerウィジェットを初期化する。
func (a *answer) reset() {
	a.history = [][]int{}

	a.render()
}

// これまでの解答を表示する。
// 正解を見破った場合は正解も表示する。
func (a *answer) render() {
	s := gamestate.GetState()

	a.Section.Clear()

	a.Section.SetCell(0, 0, tview.NewTableCell(withPadding("")))

	for c := range s.Length {
		a.Section.SetCell(0, c+1,
			tview.NewTableCell(withPadding("#", c+1)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter),
		)
	}

	tries := len(a.history)

	for r := range tries {
		a.Section.SetCell(r+1, 0,
			tview.NewTableCell(withPadding("@", r+1)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignLeft),
		)

		for c := range s.Length {
			a.Section.SetCell(r+1, c+1,
				tview.NewTableCell(withPadding(a.history[r][c])).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignRight),
			)
		}
	}

	a.Section.SetCell(tries+1, 0,
		tview.NewTableCell(withPadding("Answer")).
			SetTextColor(tcell.ColorBlue).
			SetAlign(tview.AlignLeft),
	)

	for c := range s.Length {
		symbol := "?"

		if s.IsRevealed {
			symbol = fmt.Sprint(s.Code[c])
		}

		cell := tview.NewTableCell(withPadding(symbol)).
			SetTextColor(tcell.ColorLightGreen).
			SetAlign(tview.AlignRight)

		if s.IsRevealed {
			cell.SetAttributes(tcell.AttrBlink)
		}

		a.Section.SetCell(tries+1, c+1, cell)
	}
}

// 新たな解答を受け取り、履歴に追加する。
func (a *answer) update(guess []int) {
	s := gamestate.GetState()

	a.history = append(a.history, guess)

	if reflect.DeepEqual(guess, s.Code) {
		s.Reveal()
	}

	a.render()
}
