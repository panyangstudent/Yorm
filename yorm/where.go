package yorm

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// 单个参数调用
func (y *YEngine) Where(field interface{}, opt interface{}, param interface{}) *YEngine {
	if field.(string) == "" || opt.(string) == "" {
		log.Println(fmt.Sprintf("Where param error"))
		return y
	}
	if y.whereParam != FirstWhere {
		y.whereParam += " and "
	}
	opt = strings.ToLower(opt.(string))
	// 针对in， 和not in的情况单独处理
	if opt == "in" || opt == "not in" {
		paramType := reflect.TypeOf(param).Kind()
		if paramType != reflect.Slice && paramType != reflect.Array {
			log.Println(fmt.Sprintf("param type error"))
			return y
		}
		v := reflect.ValueOf(param)
		vLen := v.Len()
		ps := make([]string, vLen)
		for i := 0; i < v.NumField(); i++ {
			ps[i] = "?"
			y.whereExec = append(y.whereExec, v.Field(i).Interface())
		}
		y.whereParam += field.(string) + " " + opt.(string) + " (" + strings.Join(ps, ",") + ") "
	} else {
		y.whereParam += field.(string) + " " + opt.(string) + " ？"
		y.whereExec = append(y.whereExec, param)
	}
	return y
}