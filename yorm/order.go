package yorm

import (
	"log"
	"strings"
)

/**
// 查询结果按照多个参数，不同排序
select * from user_info where uid >= 10888 order by uid asc, status desc
参数之间使用，分割。例如：id，desc，uid，acs
 */

func (y *YEngine) Order(param ...string) *YEngine {
	orderLen := len(param)
	if orderLen%2 != 0 {
		log.Println("order by 参数数量错误，需要偶数个")
		return y
	}
	// 可能存在多次调用的情况
	if y.orderParam != "" {
		y.orderParam += ","
	}
	for i := 0; i < orderLen/2; i++ {
		order := strings.ToLower(param[2*i+1])
		if order != "desc" && order != "acs" {
			log.Println("排序关键字错误")
			return y
		}
		if i < orderLen/2-1 {
			y.orderParam += param[i*2] + " " + param[i*2+1] + ","
		} else {
			y.orderParam += param[i*2] + " " + param[i*2+1]
		}
	}
	return y
}