package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/quickguard-oss/code-breaker/widget"
	"github.com/rivo/tview"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// ゲームを開始する。
func run() error {
	length, maxNumber, err := parseFlags()

	if err != nil {
		return err
	}

	app := tview.NewApplication()

	gamestate.Initialize(length, maxNumber)

	header := widget.NewHeader()

	rule := widget.NewRule()

	answer := widget.NewAnswer()

	judge := widget.NewJudge()

	hint := widget.NewHint()

	try := widget.NewTry([]widget.Updatable{answer, judge, hint})

	menu := widget.NewMenu(app, []widget.Resettable{try, answer, judge, hint})

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header.Section, 3, 0, false).
		AddItem(
			tview.NewFlex().
				AddItem(try.Section, 0, 4, true).
				AddItem(rule.Section, 0, 3, false).
				AddItem(menu.Section, 0, 3, false),
			7, 0, true,
		).
		AddItem(
			tview.NewFlex().
				AddItem(answer.Section, 0, 2, false).
				AddItem(judge.Section, 0, 1, false).
				AddItem(hint.Section, 0, 2, false),
			0, 1, false,
		)

	app.
		SetRoot(root, true).
		SetFocus(root).
		EnableMouse(true)

	if err := app.Run(); err != nil {
		return err
	}

	return nil
}

// コマンドライン引数をパースする。
func parseFlags() (int, int, error) {
	length := flag.Int("l", 3, "Number of digits. [1-9]")
	maxNumber := flag.Int("n", 5, "Maximum number to be used, ranging from 1 up to this value. [1-9]")

	flag.Parse()

	if (*length < 1) || (9 < *length) {
		return 0, 0, fmt.Errorf("number of digits must be between 1 and 9")
	}

	if (*maxNumber < 1) || (9 < *maxNumber) {
		return 0, 0, fmt.Errorf("maximum number must be between 1 and 9")
	}

	if *maxNumber < *length {
		return 0, 0, fmt.Errorf("number of digits must be less than or equal to maximum number")
	}

	return *length, *maxNumber, nil
}
