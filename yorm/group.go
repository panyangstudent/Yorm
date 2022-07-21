package yorm

import "strings"

func (y *YEngine)Group(param ...string) *YEngine {
	if len(param) != 0 {
		y.groupParam = strings.Join(param, ",")
	}
	return y
}