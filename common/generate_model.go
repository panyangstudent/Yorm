package common

import (
	"os"
	"text/template"
)

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