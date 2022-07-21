package yorm

/**
y.Table("userinfo").Where("status", 2).Field("uid,status").Select()
 */
func (y *YEngine) Field(param string) *YEngine {
	y.fieldParam = param
	return y
}
