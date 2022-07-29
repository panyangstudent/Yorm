package yorm

import (
	"database/sql"
	"fmt"
)

func (y *YEngine) FindForTemplate() (columnTypeAndName map[string]string, err error) {
	y.prePare = "select " + y.fieldParam + " from " + y.GetTableName()

	if y.whereParam != FirstWhere {
		y.prePare += y.whereParam + " ) "
	}

	if y.limitParam != 0 {
		y.prePare += fmt.Sprintf(" %v %d ", " limit ", y.limitParam)
	}

	rows, err := y.DB.Query(y.prePare)
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
	columnTypeAndName = make(map[string]string)
	for _, columnType := range dest {
		switch columnType.DatabaseTypeName() {
		case "FLOAT", "DOUBLE", "DECIMAL":
			columnTypeAndName[columnType.Name()] = "float64"
		case "VARCHAR", "TEXT", "NVARCHAR","CHAR":
			columnTypeAndName[columnType.Name()] = "string"
		case "INT", "BIGINT", "TINYINT":
			columnTypeAndName[columnType.Name()] = "int64"
		case "BOOL":
			columnTypeAndName[columnType.Name()] = "bool"
		case "DATE", "TIME", "TIMESTAMP", "DATETIME":
			columnTypeAndName[columnType.Name()] = "time.Time"
		}
	}
	return
}