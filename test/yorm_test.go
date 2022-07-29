package test

import (
	"OrmIsImplementedIn7Days/dao"
	"OrmIsImplementedIn7Days/model"
	"OrmIsImplementedIn7Days/yorm"
	"fmt"
	"log"
	"testing"
)
var err error
// 初始化数据库链接
func init() {
	err = yorm.MysqlInit()
	if err != nil {
		log.Println(fmt.Sprintf("[main] MysqlInit error : %v", err))
		return
	}
}

func TestInsertUser(t *testing.T) {
	insertData := model.UserInfo{
		UserName:   "yanglei",
		DepartName: "司补",
		Status:     1,
	}
	id, err := dao.InsertUser(insertData)
	log.Println(fmt.Sprintf("TestInsertUser InsertUser resp id : %v err : %v", id, err))
}

func TestFindForTemplate(t *testing.T)  {
	columnTypeAndName, err := dao.SelectTableForTemplate(dao.TableNameUserInfo)
	log.Println(fmt.Sprintf("TestInsertUser InsertUser resp id : %v err : %v", columnTypeAndName, err))
}
