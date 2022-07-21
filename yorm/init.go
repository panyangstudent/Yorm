package yorm

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //引入数据库驱动注册及初始化
	"log"
)

var baseDBInstance  *YEngine

//初始化mysql链接
func MysqlInit() (err error ){
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	fmt.Println("db init start")
	// 链接数据库
	baseDBInstance = &YEngine{
		whereParam: FirstWhere,
		fieldParam: "*",
	}
	if baseDBInstance.DB , err = sql.Open("mysql","root:12345678@tcp(127.0.0.1)/user_test?charset=utf8"); err != nil {
		log.Println(fmt.Sprintf("mysql conn error : %v", err))
		return
	}

	// 检查是否可以正常连接
	if err = baseDBInstance.DB.Ping(); err != nil {
		log.Printf(fmt.Sprintf("db ping error : %v", err))
		return
	}
	// 设置最大打开链接数(空闲+使用中)
	// 设置同时打开的连接数(使用中+空闲)
	// 设为5。将此值设置为小于或等于0表示没有限制
	baseDBInstance.DB.SetMaxOpenConns(5)

	// 在连接池中保持最大3个空闲链接
	// 理论上保持更多的空闲链接将提高性能，降低从头创建新链接的可能性
	baseDBInstance.DB.SetMaxIdleConns(3)

	// 连接池中空闲链接的最大存活时间，超过当前时间将会从连接池中移除
	baseDBInstance.DB.SetConnMaxIdleTime(3)

	// 设置链接空闲最大保持时间，超过便会断开
	// 链接将在第一次被创建之后，超过该时间断开
	baseDBInstance.DB.SetConnMaxLifetime(3)
	return
}

func NewDBInstance() *YEngine {
	return baseDBInstance
}
