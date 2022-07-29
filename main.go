package main

import (
	"OrmIsImplementedIn7Days/common"
	"OrmIsImplementedIn7Days/yorm"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Printf("start")
	// 数据库初始化
	err := yorm.MysqlInit()
	if err != nil {
		log.Printf(fmt.Sprintf("[main] MysqlInit error : %v", err))
		return
	}

	// 获取命令行输入的table name
	params := os.Args
	if len(params) <= 1 {
		fmt.Println("请追加表名")
		return
	}

	// 生成模板
	err  = common.Generate(params[1])
	if err != nil {
		return
	 }
	fmt.Println("\n\n\n\nsuccess\nstruct模板已生成")
}