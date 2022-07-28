package yorm
/**
设置sql的table name,将table name挂载到YEngine的tableName字段上
 */
func (y *YEngine) SetTable(tableName string) *YEngine {
	if tableName == "" {
		panic("[db] error table is empty")
	}
	y.tableName = tableName
	y.RsetYEngine()
	return y
}

// 清空当前Yengine上的挂载数据，在此之前应该在NewDbInterface时就生成一个新的DB链接，这个后续来实现
func (y *YEngine) RsetYEngine() {
	n := &YEngine{
		DB:            y.DB,
		tableName:     y.tableName,
		prePare:       "",
		allExec:       nil,
		sql:           "",
		whereExec:     nil,
		whereParam:    FirstWhere,
		limitParam:    0,
		orderParam:    "",
		upDateParam:   "",
		updateExec:    nil,
		fieldParam:    "*",
		tx:            nil,
		groupParam:    "",
	}
	y = n
}

func (y *YEngine) GetTableName() string {
	return y.tableName
}