package widget

import (
	"github.com/rivo/tview"
)

// ゲームのタイトルを表示するヘッダ・ウィジェット
type header struct {
	Section *tview.TextView // tview
}

// headerウィジェットを生成する。
func NewHeader() *header {
	section := tview.NewTextView()

	h := &header{
		Section: section,
	}

	h.Section.
		SetTextAlign(tview.AlignCenter).
		SetText(withPadding("Code Breaker")).
		SetBorder(true)

	return h
}
