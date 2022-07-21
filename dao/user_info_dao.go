package dao

import (
	"OrmIsImplementedIn7Days/model"
	"OrmIsImplementedIn7Days/yorm"
	"fmt"
	"log"
)

const TableNameUserInfo = "user_info"

func InsertUser(info model.UserInfo) (id int64, err error) {
	dbInstance := yorm.NewDBInstance()
	if id, err = dbInstance.SetTable(TableNameUserInfo).Insert(info); err != nil {
		log.Println(fmt.Sprintf("InsertUser Insert error : %v info : %v", err, info))
		return
	}
	return
}

func Where()  {
	dbInstance := yorm.NewDBInstance()
	dbInstance.Where("id", "=", 1).Where("id", "in", []int64{1, 2, 3, 4})
}
func SelectUsers() (userList []model.UserInfo, err error) {
	return
}

func UpdateUserInfo() (err error) {
	return
}
func Order() (err error) {
	dbInstance := yorm.NewDBInstance()
	dbInstance.Order("id","desc","uid","acs")
	return
}