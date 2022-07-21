package yorm

func (y * YEngine) Generate(tableName string)  {
	y.prePare = "select * from " + tableName + " limit 1"
	y.DB.QueryRow(y.prePare)
}
