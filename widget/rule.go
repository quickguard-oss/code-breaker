package widget

import (
	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/rivo/tview"
)

// ゲームのルールを表示するウィジェット
type rule struct {
	Section *tview.TextView // tview
}

// ruleウィジェットを生成する。
func NewRule() *rule {
	section := tview.NewTextView()

	r := &rule{
		Section: section,
	}

	s := gamestate.GetState()

	r.Section.
		SetText(
			withPadding("Length: ", s.Length) + "\n" +
				withPadding("Range: 1-", s.MaxNumber),
		).
		SetBorder(true).
		SetTitle(withPadding("Rule"))

	return r
}
