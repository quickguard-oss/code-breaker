package widget

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/quickguard-oss/code-breaker/gamestate"
	"github.com/rivo/tview"
)

// ヒントの凡例
const (
	symbolHit           = '@' // ヒット
	symbolExist         = 'O' // 正解には当該数字が含まれる
	symbolExcluded      = 'X' // 正解には当該数字が含まれない
	symbolNotDetermined = ' ' // 未確定 (= 情報が足りず、判断できない)
)

// これまでの試行から得られた情報をヒントとして表示するウィジェット
type hint struct {
	Section    *tview.Table // tview
	matrix     [][]rune     // ヒント表
	candidates [][]int      // これまでの結果に合致する、正解の候補群 (ヒントの導出に使用する)
}

// hintウィジェットを生成する。
func NewHint() *hint {
	section := tview.NewTable()

	h := &hint{
		Section: section,
	}

	h.Section.
		SetBorders(true).
		SetBorder(true).
		SetTitle(withPadding("Hint"))

	h.reset()

	return h
}

// hintウィジェットを初期化する。
func (h *hint) reset() {
	h.initializeCandidates()

	h.initializeMatrix()

	h.render()
}

// 正解候補群に初期値 (= 使用可能な数字から成る順列) をセットする。
func (h *hint) initializeCandidates() {
	s := gamestate.GetState()

	digits := make([]int, s.MaxNumber)

	for i := range digits {
		digits[i] = i + 1
	}

	h.candidates = h.generatePermutations([]int{}, digits)
}

// 正解の候補となる順列を生成する。
func (h *hint) generatePermutations(permutation []int, rest []int) [][]int {
	result := [][]int{}

	s := gamestate.GetState()

	if len(permutation) == s.Length {
		candidate := make([]int, len(permutation))

		copy(candidate, permutation)

		return append(result, candidate)
	}

	for i, v := range rest {
		result = slices.Concat(
			result,
			h.generatePermutations(
				append(permutation, v),
				slices.Concat(rest[:i], rest[i+1:]),
			),
		)
	}

	return result
}

// ヒント表に初期値 (= すべてのマスが未確定) をセットする。
func (h *hint) initializeMatrix() {
	s := gamestate.GetState()

	matrix := make([][]rune, s.MaxNumber)

	for r := range s.MaxNumber {
		matrix[r] = make([]rune, s.Length)

		for c := range s.Length {
			matrix[r][c] = symbolNotDetermined
		}
	}

	h.matrix = matrix
}

// ヒントを表示する。
func (h *hint) render() {
	s := gamestate.GetState()

	h.Section.Clear()

	h.Section.SetCell(0, 0, tview.NewTableCell(withPadding("")))

	for c := range s.Length {
		h.Section.SetCell(0, c+1,
			tview.NewTableCell(withPadding("#", c+1)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignCenter),
		)
	}

	for r := range s.MaxNumber {
		h.Section.SetCell(r+1, 0,
			tview.NewTableCell(withPadding(r+1)).
				SetTextColor(tcell.ColorBlue).
				SetAlign(tview.AlignLeft),
		)

		for c := range s.Length {
			cellValue := h.matrix[r][c]

			cell := tview.NewTableCell(withPadding(string(cellValue))).
				SetAlign(tview.AlignRight)

			switch cellValue {
			case symbolHit:
				cell.SetTextColor(tcell.ColorLightGreen)
			case symbolExist:
				cell.SetTextColor(tcell.ColorYellow)
			case symbolExcluded:
				cell.SetTextColor(tcell.ColorRed)
			default:
				cell.SetTextColor(tcell.ColorWhite)
			}

			h.Section.SetCell(r+1, c+1, cell)
		}
	}
}

// 新たな解答とその判定結果を受け取り、ヒントを更新する。
func (h *hint) update(guess []int) {
	h.updateCandidates(guess)

	h.updateMatrix()

	h.render()
}

// 正解候補群を更新する。
// (受け取った解答と判定結果に合致する候補のみを残す。)
func (h *hint) updateCandidates(guess []int) {
	newCandidates := [][]int{}

	hitA, blowA := hitAndBlow(guess, gamestate.GetState().Code)

	for _, candidate := range h.candidates {
		hitB, blowB := hitAndBlow(guess, candidate)

		if hitA == hitB && blowA == blowB {
			newCandidates = append(newCandidates, candidate)
		}
	}

	h.candidates = newCandidates
}

// ヒント表を更新する。
func (h *hint) updateMatrix() {
	s := gamestate.GetState()

	digitPossibilities := make([]map[int]struct{}, s.Length)

	for i := range s.Length {
		digitPossibilities[i] = map[int]struct{}{}
	}

	exists := 0

	for i := range s.MaxNumber {
		exists = exists + (1 << i)
	}

	for _, candidate := range h.candidates {
		bits := 0

		for i, v := range candidate {
			digitPossibilities[i][v] = struct{}{}

			bits = bits + (1 << (v - 1))
		}

		exists = exists & bits
	}

	for c, possibility := range digitPossibilities {
		possibles := []int{}

		for r := range s.MaxNumber {
			if _, ok := possibility[r+1]; ok {
				possibles = append(possibles, r)
			} else {
				// ルール1: 当該桁において候補に含まれない数字は、正解から除外できる。
				h.matrix[r][c] = symbolExcluded
			}
		}

		if len(possibles) == 1 {
			// ルール2: 当該桁において候補数が唯一つの数字は、ヒットであることが確定する。
			h.matrix[possibles[0]][c] = symbolHit
		}
	}

	for v := range s.MaxNumber {
		if exists&(1<<v) != 0 {
			for c := range s.Length {
				if h.matrix[v][c] == symbolNotDetermined {
					// ルール3: 正解候補群のすべてに含まれる数字は、正解のいずれかの桁に存在する。
					h.matrix[v][c] = symbolExist
				}
			}
		}
	}
}
