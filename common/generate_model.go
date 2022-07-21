package common

import (
	"OrmIsImplementedIn7Days/yorm"
	"os"
	"text/template"
)

func selectTable(tableName string)  {
	dbInterface := yorm.NewDBInstance()
	dbInterface.SetTable(tableName)
}
func Generate(tableName string)  {
	tmpl , err := template.New("test").Parse(`hello {{.name}}!obj:{{.}}`)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"name": "world", "age": 18})
	if err != nil {
		panic(err)
	}
}