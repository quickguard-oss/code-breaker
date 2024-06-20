package gamestate

import (
	"math/rand"
)

// 最大試行回数
const maxTry = 9

// ゲーム全体にまつわる状態
type state struct {
	Length     int   // 数列の桁数
	MaxNumber  int   // 使用可能な数字 (1〜MaxNumber)
	IsRevealed bool  // 正解を開示するか否か (正解を見破った or 試行回数を使い切った場合に true)
	Code       []int // 正解
	tryCount   int   // 試行した回数 (上限に達したらチャレンジ終了)
}

var stateInstance *state

// stateインスタンスを初期化する。
func Initialize(length int, maxNumber int) {
	stateInstance = &state{
		Length:    length,
		MaxNumber: maxNumber,
	}

	stateInstance.Reset()
}

// stateインスタンスを取得する。
func GetState() *state {
	return stateInstance
}

// stateインスタンスを初期状態にリセットする。
func (s *state) Reset() {
	s.IsRevealed = false

	s.Code = s.generateCode()

	s.tryCount = 0
}

// 試行回数をカウントアップする。
func (s *state) Try() {
	s.tryCount++

	if maxTry <= s.tryCount {
		s.Reveal()
	}
}

// 正解を開示する。
func (s *state) Reveal() {
	s.IsRevealed = true
}

// ランダムな正解を生成する。
func (s *state) generateCode() []int {
	code := make([]int, s.MaxNumber)

	for i := range s.MaxNumber {
		code[i] = i + 1
	}

	rand.Shuffle(
		s.MaxNumber,
		func(i, j int) {
			code[i], code[j] = code[j], code[i]
		},
	)

	return code[0:s.Length]
}
