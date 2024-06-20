package widget

import (
	"fmt"
)

// 値を空白で囲む。
func withPadding(a ...any) string {
	return " " + fmt.Sprint(a...) + " "
}

// ヒット数とブロー数を算出する。
func hitAndBlow(guess []int, code []int) (int, int) {
	hit := 0
	blow := 0

	digitMap := map[int]int{}

	for i, v := range code {
		digitMap[v] = i
	}

	for i, v := range guess {
		if pos, ok := digitMap[v]; ok {
			if i == pos {
				hit++
			} else {
				blow++
			}
		}
	}

	return hit, blow
}
