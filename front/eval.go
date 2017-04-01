package front

type EvalItem struct {
	Eval     string
	EvalName string // gen by server
	EvalAt   int64  // gen by server
	RateStar uint
}
