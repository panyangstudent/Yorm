package model

const TABLE_NAME_USER_INFO = "user_info"

type UserInfo struct {
	Id         int64  `sql:"id, auto_increment"`
	UserName   string `sql:"user_name"`
	DepartName string `sql:"depart_name"`
	Status     int64  `sql:"status"`
}