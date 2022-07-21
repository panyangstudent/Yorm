package yorm

import (
	"log"
)

func (y *YEngine) Or () *YEngine {
	if y.whereParam == "" {
		log.Panicln("or必须在Where之后调用")
		return y
	}

	y.whereParam += " ) or ( "
	return y
}