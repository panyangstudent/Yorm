package yorm

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (y *YEngine) Find(result interface{}) (err error) {
	// 判断入参是否是指针
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		return errors.New("入参需要是指针类型")
	}

	if reflect.ValueOf(result).IsNil() {
		return errors.New("入参不能是nil")
	}

	// 拼接sql
	y.prePare = "select " + y.fieldParam + " from " + y.GetTableName()

	if y.whereParam != FirstWhere {
		y.prePare += y.whereParam + " ) "
	}

	if y.orderParam != "" {
		y.prePare += " order by " + y.orderParam
	}

	if y.groupParam != "" {
		y.prePare += " group by " + y.groupParam
	}

	if y.limitParam != 0 {
		y.prePare += fmt.Sprintf(" %v %d ", " limit ", y.limitParam)
	}


	// 读出每个列
	rows, err := y.DB.Query(y.prePare, y.allExec)
	if err != nil {
		return y.setError(err)
	}

	// 读出查询出的列字段名
	column, err := rows.Columns()
	if err != nil {
		return y.setError(err)
	}

	// values是每个列的值, 这里获取到byte里
	values := make([][]byte, len(column))

	//因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	scans := make([]interface{}, len(column))

	//原始struct的slice，当前destSlice也是指针类型
	destSlice := reflect.ValueOf(result).Elem()

	//原始单个struct的各个字段
	destType := destSlice.Type().Elem()
	for i := range values {
		scans[i] = &values[i]
	}

	// 循环遍历
	var fieldName string
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		if err = rows.Scan(scans...); err != nil {
			return y.setError(err)
		}

		// 遍历一行数据的各个字段
		for k, v := range values {
			key := column[k]
			value := string(v)
			for i := 0; i < destType.NumField(); i++ {
				// 看下是否有sql别名
				sqlTag := destType.Field(i).Tag.Get("sql")
				if sqlTag != "" {
					fieldName = strings.Split(sqlTag, ",")[0]
				} else {
					fieldName = destType.Field(i).Name
				}
				// 判断字段名是否相同
				if key != fieldName {
					continue
				}
				// 反射赋值
				if err = y.reflectSet(dest, i, value); err != nil {
					return err
				}
			}
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return
}

func (y *YEngine) reflectSet(dest reflect.Value, i int, value string) error {
	switch dest.Field(i).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return y.setError(err)
		}
		dest.Field(i).SetInt(res)
	case reflect.String:
		dest.Field(i).SetString(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return y.setError(err)
		}
		dest.Field(i).SetUint(res)
	case reflect.Float32:
		res, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return y.setError(err)
		}
		dest.Field(i).SetFloat(res)
	case reflect.Float64:
		res, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return y.setError(err)
		}
		dest.Field(i).SetFloat(res)
	case reflect.Bool:
		res, err := strconv.ParseBool(value)
		if err != nil {
			return y.setError(err)
		}
		dest.Field(i).SetBool(res)
	}
	return nil
}