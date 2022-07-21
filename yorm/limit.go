package yorm

func (y *YEngine) Limit(param int) *YEngine {
	y.limitParam = param
	return y
}