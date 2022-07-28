package yorm

import (
	"database/sql"
	"fmt"
	"reflect"
)

func (y *YEngine) FindForTemplate() (columnTypeAndName map[string]string, err error) {
	y.prePare = "select " + y.fieldParam + " from " + y.GetTableName()

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

	rows, err := y.DB.Query(y.prePare, y.allExec)
	if err != nil {
		return nil, y.setError(err)
	}

	// 读出查询出的列字段名和对应类型
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, y.setError(err)
	}

	columnTypeAndName = y.reflectColumnTypeToGoStruct(columnTypes)
	return
}

func (y *YEngine) reflectColumnTypeToGoStruct(dest []*sql.ColumnType) (columnTypeAndName map[string]string) {
	for _, columnType := range dest {
		switch columnType.ScanType().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			columnTypeAndName[columnType.Name()] = "int64"
		case reflect.String:
			columnTypeAndName[columnType.Name()] = "string"
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			columnTypeAndName[columnType.Name()] = "float64"
		case reflect.Bool:
			columnTypeAndName[columnType.Name()] = "bool"
		}
	}
	return
}