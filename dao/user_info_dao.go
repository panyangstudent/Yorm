package dao

import (
	"OrmIsImplementedIn7Days/model"
	"OrmIsImplementedIn7Days/yorm"
	"fmt"
	"log"
)

const TableNameUserInfo = "user_info"



func SelectTableForTemplate(tableName string) (columnTypeAndName map[string]string, err error) {
	dbInterface := yorm.NewDBInstance()
	columnTypeAndName, err = dbInterface.SetTable(tableName).Limit(1).FindForTemplate()
	return
}


func InsertUser(info model.UserInfo) (id int64, err error) {
	dbInstance := yorm.NewDBInstance()
	if id, err = dbInstance.SetTable(TableNameUserInfo).Insert(info); err != nil {
		log.Println(fmt.Sprintf("InsertUser Insert error : %v info : %v", err, info))
		return
	}
	return
}
