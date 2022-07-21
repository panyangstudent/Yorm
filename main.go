package main

import (
	"OrmIsImplementedIn7Days/yorm"
	"fmt"
	"log"
)

func main() {
	log.Printf("start")
	// 数据库初始化
	err := yorm.MysqlInit()
	if err != nil {
		log.Printf(fmt.Sprintf("[main] MysqlInit error : %v", err))
		return
	}
}