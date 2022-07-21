package yorm

import (
	"fmt"
	"strconv"
	"strings"
)

func (y *YEngine) Sum(param string)(int64, error) {
	sum, err := y.aggregationQuery("sum", param)
	if err != nil {
		return 0, y.setError(err)
	}
	return sum.(int64), err
}

func (y *YEngine) Avg(param string) (string, error) {
	avg, err := y.aggregationQuery("avg", param)
	if err != nil {
		return "0", y.setError(err)
	}
	return string(avg.([]byte)), nil
}

func (y *YEngine) Max(param string) (string, error) {
	max, err := y.aggregationQuery("max", param)
	if err != nil {
		return "0", y.setError(err)
	}
	return string(max.([]byte)), nil
}

func (y YEngine) Min(param string) (string, error) {
	min, err := y.aggregationQuery("min", param)
	if err != nil {
		return "0", err
	}
	return string(min.([]byte)), nil
}

func (y *YEngine) Count() (int64, error) {
	count, err := y.aggregationQuery("count", "*")
	if err != nil {
		return 0, y.setError(err)
	}
	return count.(int64), err
}

func (y *YEngine) aggregationQuery(name, param string) (cnt interface{}, err error) {
	// 拼接sql
	y.prePare = "select " + name + "(" + param + ") as cnt from " + y.GetTableName()

	// 拼接where
	if y.whereParam != "" {
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

	y.allExec = y.whereExec
	// 执行绑定
	y.generateSql()
	//查询
	err = y.DB.QueryRow(y.prePare, y.allExec...).Scan(&cnt)
	return
}

func (y *YEngine) generateSql() {
	y.sql = y.prePare
	for _, v := range y.allExec {
		switch v.(type) {
		case int:
			y.sql = strings.Replace(y.sql, "?", strconv.Itoa(v.(int)), 1)
		case int64:
			y.sql = strings.Replace(y.sql, "?", strconv.FormatInt(v.(int64), 10), 1)
		case bool:
			y.sql = strings.Replace(y.sql,"?", strconv.FormatBool(v.(bool)),1)
		default:
			y.sql = strings.Replace(y.sql, "?", "'"+v.(string)+"'", 1)
		}
	}
}