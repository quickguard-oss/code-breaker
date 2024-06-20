package widget

import (
	"github.com/gdamore/tcell/v2"
	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/rivo/tview"
)

// これまでのヒット数とブロー数を表示するウィジェット
type judge struct {
	Section *tview.Table // tview
	history [][]int      // ヒット数とブロー数の履歴
}

// judgeウィジェットを生成する。
func NewJudge() *judge {
	section := tview.NewTable()

	j := &judge{
		Section: section,
	}

	j.Section.
		SetBorders(true).
		SetBorder(true).
		SetTitle(withPadding("Judge"))

	j.reset()

	return j
}

// judgeウィジェットを初期化する。
func (j *judge) reset() {
	j.history = [][]int{}

	j.render()
}

// これまでのヒット数とブロー数を表示する。
func (j *judge) render() {
	j.Section.Clear()

	j.Section.SetCell(0, 0, tview.NewTableCell(withPadding("  ")))

	for c, v := range []string{"Hit", "Blow"} {
		j.Section.SetCell(0, c+1,
			tview.NewTableCell(withPadding(v)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter),
		)
	}

	tries := len(j.history)

	for r := range tries {
		j.Section.SetCell(r+1, 0,
			tview.NewTableCell(withPadding("@", r+1)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignLeft),
		)

		for c, v := range j.history[r] {
			j.Section.SetCell(r+1, c+1,
				tview.NewTableCell(withPadding(v)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignRight),
			)
		}
	}
}

// 新たな解答を受け取り、ヒット数とブロー数を算出して履歴に追加する。
func (j *judge) update(guess []int) {
	hit, blow := hitAndBlow(guess, gamestate.GetState().Code)

	j.history = append(j.history, []int{hit, blow})

	j.render()
}
