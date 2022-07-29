package common

import (
	"OrmIsImplementedIn7Days/dao"
	"log"
	"os"
	"strings"
	"text/template"
)



func Generate(tableName string) (err error)  {
	toUpTableName := ToUpperStr(tableName)
	//生成表头
	tmpl1, err1 := template.New("test1").Parse(`
	type {{.tableName}} struct {`)
	if err1 != nil {
		panic(err)
	}
	err = tmpl1.Execute(os.Stdout, map[string]string{"tableName":toUpTableName})
	if err != nil {
		panic(err)
	}

	columnTypeAndName, err := dao.SelectTableForTemplate(tableName)
	if err != nil {
		log.Println("generate SelectTableForTemplate error : %v", err)
		return
	}
	// 生成中间结构体
	result := disposeColumn(columnTypeAndName)
	tmpl, err := template.New("test").Parse(`
	{{range $idx, $value := .val}}
		{{$idx}}  {{$value}}{{end}}
	}`)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, result)
	if err != nil {
		panic(err)
	}
	return
}

// 预处理
func disposeColumn(columnTypeAndName map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	result["val"] = make(map[string]string)
	for columnName, columnType := range columnTypeAndName {
		resultColumnName := ToUpperStr(columnName)
		columnType += " `sql:" + columnName + ", json:" + columnName + "`"
		result["val"][resultColumnName] = columnType
	}
	return result
}

func ToUpperStr(oldStr string)  (resultStr string) {
	strSplit := strings.Split(oldStr, "_")
	if len(strSplit) >= 1 {
		for _, s := range strSplit {
			resultStr += strings.ToUpper(s[:1]) + s[1:]
		}
	} else {
		resultStr = strings.ToUpper(strSplit[0][:1]) + strSplit[0][1:]
	}
	return
}