package widget

// 解答を受け取って自身の状態を更新するウィジェットを表すインタフェイス
type Updatable interface {
	update(guess []int)
}
