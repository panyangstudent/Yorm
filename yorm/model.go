package yorm

import (
	"database/sql"
	"errors"
	"fmt"
)

type YEngine struct {
	DB            *sql.DB //数据库链接
	tableName     string
	prePare       string
	allExec       []interface{} //入参
	sql           string

	whereExec     []interface{} //where条件的入参
	whereParam    string

	limitParam    int
	orderParam    string
	groupParam    string

	upDateParam   string
	updateExec    []interface{} //update条件的入参

	fieldParam    string
	tx            *sql.Tx
}

const FirstWhere = "where ( "

func (y *YEngine) setError(err error) error {
	return errors.New(fmt.Sprintf("db error : %v", err))
}