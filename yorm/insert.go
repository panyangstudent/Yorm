package yorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// 结构体入参，按照结构体的入参进行数据插入
func (y *YEngine) Insert(param interface{}) (id int64, err error) {
	// 初始化字段名，占位符数组
	fieldName := make([]string, 0)
	placeHolder := make([]string, 0)

	// 反射当前参数的类型和值
	v := reflect.ValueOf(param)
	t := reflect.TypeOf(param)

	for i := 0; i < t.NumField(); i++ {
		// 判断是否可以反射
		if !v.Field(i).CanInterface() {
			continue
		}

		// 包含sql tag的认为是数据库字段
		sqlTag := t.Field(i).Tag.Get("sql")
		if sqlTag != "" {
			// 跳过自增字段, 查看sqlTag中是否包含auto_increment字段
			if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {
				continue
			} else  {
				fieldName = append(fieldName, strings.Split(sqlTag, ",")[0])
			}
		} else {
			fieldName = append(fieldName, strings.ToLower(t.Field(i).Name))
		}
		placeHolder = append(placeHolder, "?")
		// 获取字段值
		y.allExec = append(y.allExec, v.Field(i).Interface())
	}
	// 拼接表，字段名，占位符等
	y.prePare = "insert into " + y.GetTableName() + "("  + strings.Join(fieldName, ",") + ") values(" + strings.Join(placeHolder, ",") + ")"

	fmt.Println(y.prePare)
	fmt.Println(y.allExec)
	//执行
	stmt, err  := y.DB.Prepare(y.prePare)
	if err != nil {
		err = errors.New(fmt.Sprintf("[db] prepare error : %v", err))
		return
	}

	// 执行exec
	result, err := stmt.Exec(y.allExec...)
	if err != nil {
		err = errors.New(fmt.Sprintf("[db] exec error : %v", err))
		return
	}
	id, _ = result.LastInsertId()
	return
}
