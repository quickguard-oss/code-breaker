package widget

import (
	"strconv"

	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/rivo/tview"
)

// プレイヤーの解答を入力するウィジェット
type try struct {
	Section *tview.Form // tview
}

// tryウィジェットを生成する。
func NewTry(widgets []Updatable) *try {
	section := tview.NewForm()

	t := &try{
		Section: section,
	}

	s := gamestate.GetState()

	t.Section.
		AddInputField("Numbers:", "", s.Length+1, func(textToCheck string, lastChar rune) bool {
			if s.Length < len(textToCheck) {
				return false
			}

			if (lastChar < '1') || (rune('0'+s.MaxNumber) < lastChar) {
				return false
			}

			return true
		}, nil).
		AddButton("OK", func() {
			if s.IsRevealed {
				return
			}

			if len(t.getInputField().GetText()) != s.Length {
				return
			}

			s.Try()

			guess := t.getAnswer()

			for _, w := range widgets {
				w.update(guess)
			}

			t.reset()
		}).
		AddButton("Clear", func() {
			t.reset()
		}).
		SetBorder(true).
		SetTitle(withPadding("Try"))

	return t
}

// tryウィジェットを初期化する。
func (t *try) reset() {
	inputField := t.getInputField()

	inputField.SetText("")
}

// 入力フィールドのインスタンスを取得する。
func (t *try) getInputField() *tview.InputField {
	return t.Section.GetFormItemByLabel("Numbers:").(*tview.InputField)
}

// 入力された解答を取得する。
func (t *try) getAnswer() []int {
	var numbers []int

	text := t.getInputField().GetText()

	for _, char := range text {
		digit, _ := strconv.Atoi(string(char))

		numbers = append(numbers, digit)
	}

	return numbers
}
