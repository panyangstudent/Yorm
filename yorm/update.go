package yorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/**
u.Table("userinfo").Where("uid", 123).Update("status", 1)
u.Table("userinfo").Where("uid", 123).Update(user2)
 */

func (y *YEngine) Update(data ...interface{}) (count int64, err error) {
	var (
		dataType int
	)
	switch len(data) {
	case 0:
		err = errors.New("参数个数错误")
		return 0, y.setError(err)
	case 1:
		dataType = 1
	case 2:
		dataType = 2
	default:
		err = errors.New("参数个数错误")
		return 0, y.setError(err)
	}
	// 如果是结构体
	if dataType == 1 {
		t := reflect.TypeOf(data[0])
		v := reflect.ValueOf(data[0])
		fieldName := make([]string, 0)
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).CanInterface() {
				continue
			}
			//解析tag,找出真实的sql字段名
			sqlTag := t.Field(i).Tag.Get("sql")
			if sqlTag != "" {
				fieldName = append(fieldName, strings.Split(sqlTag, ",")[0]+"=?")
			} else {
				fieldName = append(fieldName, t.Field(i).Name+"=?")
			}
			y.updateExec = append(y.updateExec, v.Field(i).Interface())
		}
		y.upDateParam += strings.Join(fieldName, ",")
	} else {
		y.upDateParam += data[0].(string) + "=?"
		y.updateExec = append(y.updateExec, data[1])
	}
	// 拼接sql
	y.prePare = "update " + y.GetTableName() + "set" + y.upDateParam

	// 如果whereParam不为空
	if y.whereParam != "" {
		y.prePare += "where " + y.whereParam
	}

	// limit不为空
	if y.limitParam != 0 {
		y.prePare += fmt.Sprintf(" %v %d ", " limit ", y.limitParam)
	}

	// Prepare
	var stmt *sql.Stmt
	stmt, err = y.DB.Prepare(y.prePare)
	if err != nil {
		return 0, y.setError(err)
	}
	// 合并UpdateExec和WhereExec
	if y.whereExec != nil {
		y.allExec = append(y.allExec, y.whereExec...)
	}
	// 执行exec，
	result, err := stmt.Exec(y.allExec...)
	if err != nil {
		return 0, y.setError(err)
	}
	// 影响行数
	id, _ := result.RowsAffected()
	return id, nil
}