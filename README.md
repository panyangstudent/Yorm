# Yorm
7天实现orm

# 前言
作为一名程序猿，相信没有人不对orm熟悉。orm作为桥梁，连接了高级编程语言和关系型数据库。在一个全新的项目中当考虑到操作数据库时，相信大多数人都会选择orm。一方面orm为我们带来了便捷，安全等好处，但是大多数时候应该是因为懒和原生sql的开发效率过低问题，哈哈哈哈。
总结下，我们知道了orm是高级编程语言和数据库打交道的中间件，并且更加高效和便捷。

# Golang中原生sql是如何链接和操作mysql的
项目地址https://github.com/panyangstudent/OrmIsImplementedIn7Days
```golang
package main

import (
    "database/sql"
    "fmt"
)
type YEngine struct {
	DB *sql.DB
	TableName string
	PrePare string
	AllExec []interface{}
	Sql string
	Where string
	Limit string
	Order string
	OrWhere string
	AndWhere string
	WhereExec []interface{}
	UpDate string
	UpdateExec []interface{}
	Field string

}
var BaseDBInstance YEngine


type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Ages int64  `json:"ages"`
}

func MysqlInit() (err error ){
    fmt.Println("db init start")

    // 链接数据库
    if BaseDBInstance.DB , err = sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8"); err != nil {
        panic(fmt.Sprintf("mysql conn error : %v", err))
    }
    return
}


func (u UserDao) GetUserById(id int64) (users *User, err error) {
    var  (
        row *sql.Rows
    )
	
    sql1 := fmt.Sprintf("select * from %v where id = %v", model.TABLE_NAME_USER_INFO, id)
    if row, err = common.BaseDBInstance.DB.Query(sql1); err != nil {
        log.Error(fmt.Sprintf("GetUserById error : %v", err))
    }
    for row.Next() {
		if err = row.Scan(users); err != nil {
            return
        }
    }
    return
}

func (u UserDao) InsertUser(userInfo User) (err error) {
        sql1 := fmt.Sprintf("INSERT INTO user_info (name, ages) VALUES (?, ?)")
        if _, err = common.BaseDBInstance.DB.Exec(sql1, userInfo.Name, userInfo.Ages); err != nil {
            log.Error(fmt.Sprintf("[db] InsertUser error : %v", err))
        }
    return
}

```
如上所示首先是载入sql引擎和mysql的驱动，然后再连接Mysql。直接写插入和获取的接口，这一整套下来，开发效率真的非常低。所以orm的作用就凸显出来了，和msyql的交互交给orm，并且orm支持生成复杂的sql，业务层直接调用orm提供的方法就好，函数式编程调用。

# 如何编写一个ORM
orm需要具备的功能：
* 暴露CURD的各类链式调用方法，包括但不限于(select，update，delete，insert， where， and，or，设置tableName，limit，connect，order， group，事务Begin/Commit/Rollback)
* 支持一键生成对应的数据库表的model和常见的dao层方法

### 各类方法的暴露
首先需要构建一个YEngine对象，这个对象会绑定所以对数据库操作的方法和参数，orm的底层本质还是sql的拼接，所以我们需要将各个操作方法生成的数据都保存在这个结构体上，方便最后一步拼接sql。其中需要说明的是DB字段，她是直接用于进行CURD操作的。Tx是*sql.Tx
类型的。她是数据库的事物操作，用于回滚和提交。
```go
type YEngine struct {
  DB            *sql.DB
  TableName     string
  PrePare       string
  AllExec       []interface{}
  Sql           string
  WhereParam    string
  LimitParam    string
  OrderParam    string
  OrWhereParam  string
  WhereParam string
  WhereExec     []interface{}
  UpDateParam   string
  UpdateExec    []interface{}
  FieldParam    string
  Tx            *sql.Tx
  GroupParam    string
}

```

### 各类操作
#### 初始化连接
```go
// 初始化mysql链接
func MysqlInit() (baseDBInstance YEngine, err error){
	fmt.Println("db init start")
	// 链接数据库
	if BaseDBInstance.DB , err = sql.Open("mysql","root:123456789@tcp(127.0.0.1:3306)/test?charset=utf8"); err != nil {
		panic(fmt.Sprintf("mysql conn error : %v", err))
	}

	// 检查是否可以正常连接
	if err = BaseDBInstance.DB.Ping(); err != nil {
		log.Error(fmt.Sprintf("[db] mysql ping error : %v", err))
		return
	}

	// 设置最大打开链接数(空闲+使用中)
	// 设置同时打开的连接数(使用中+空闲)
	// 设为5。将此值设置为小于或等于0表示没有限制
	BaseDBInstance.DB.SetMaxOpenConns(5)

	// 在连接池中保持最大3个空闲链接
	// 理论上保持更多的空闲链接将提高性能，降低从头创建新链接的可能性
	BaseDBInstance.DB.SetMaxIdleConns(3)

	// 连接池中空闲链接的最大存活时间，超过当前时间将会从连接池中移除
	BaseDBInstance.DB.SetConnMaxIdleTime(3)

	// 设置链接空闲最大保持时间，超过便会断开
	// 链接将在第一次被创建之后，超过该时间断开
	BaseDBInstance.DB.SetConnMaxLifetime(3)
	return
}
```
MysqlInit方法来创建一个mysql连接，参数目前都是写死的， 暂时不考虑多数据库，连接池的加入。这两项会在后期进行优化。那么如何实现链式调用呢？应用go引用方法的特性，返回实例本身。

#### 链式调用
链式调用如下所示
```go
func (u *YEngine) Where (userInfo string) (*YEngine) {
	...
	return u
}
```
这样我们就可以链式调用，u.Where().Where().Where()

#### 设置表名
```go
func (u *YEngine) Table (tableName string) (*YEngine) {
	u.TableName = tableName
	
	// 重置YEngine，清除挂载的数据
	u.ResetYEngine()
	return u
}

func (u *YEngine) ResetYEngine () {
	// todo 除过tableName字段，其余字段上的数据清空
	return
}
```
我们每调用一次Table方法设置表名，就会清空YEngine节点上挂载的所有数据，但是这样做如果遇见并发问题，就会存在多个YEngine实例，这个我们放在后面再来优化。

#### insert
我们在本文的最开始看到了InsertUser是原生sql的插入数据方式。一般在使用原生sql时会先执行Prepare方法，做预编译。该方法经过一次编译，多次运行，省去了解析优化等过程，此外预编译语句能防止sql注入。
```go
stmt, err := db.Prepare("INSERT INTO userinfo (username, departname, created) VALUES (?, ?, ?)")

    result2, err := stmt.Exec("zhangsan", "pro", time.Now().Format("2006-01-02"))
```

ok，我们可以看到第一行代码将真正填入的值用？来代替占位。第二行代码会将真正的数值填入，一一和？对应起来。如此我们可以将其进行数据拆分，为了保持方便我们在调用insert方法时可以传入一个map，或者一个struct
。如果是map， 其value都是同一种类型，和数据库中的多类型无法一一对应，所以pass。这里选择struct，满足不动字段可能是不同类型的情况。我们在定义struct时，是先定义一个struct类型，然后初始化赋值。具体如下示例
```go
    type User struct {
        Username   string `sql:"username"`

        Departname string `sql:"departname"`

        Status     int64  `sql:"status"`
    }
    user2 := User{

        Username:   "111",

        Departname: "111", 

        Status:     1,

    }
    id, err := e.Table("userinfo").Insert(user2)
```

我们可以注意到，User结构体中，每个元素后面都有一个sql:"..."的标签，这个叫做tag标签，因为在go中Struct中的元素首字母都是大写表示public。但是sql中的每个字段名一般都是小写字母
所以我们需要该标签来实现转换。现在的问题就是我们怎么转换这个struct到sql语句中，和每个？一一对应。在这里我们拆成两步：
    
* 将sql:"..."标签进行解析和匹配，一次替换城全小写，解析成(username, departname, status)，并且依次生成对应数量
  ```go
    stmt, err := db.Prepare("INSERT INTO userinfo (username, departname, status) VALUES (?, ?, ?)")
  ```
* 将user2的子元素的值拆出来，放到exec中
   ```go
  result2, err := stmt.Exec("111", "111", 1)
  ```
问题是我们如何将user2里面的3个子元素的field，解析成(username，departname，status)呢？golang可以通过反射来推导出传入结构体的变量，他的field是什么，value是多少，类型是什么，tag是什么都可以通过反射推导出来.
我们可以试下reflect.TypeOf和reflect.ValueOf这两个方法
```go
t := reflect.TypeOf(user2)
//反射出这个结构体变量的值

    v := reflect.ValueOf(user2)
	
    fmt.Printf("==== print type ====\n%+v\n", t)
    
    fmt.Printf("==== print value ====\n%+v\n", v)
  
    //输出
    ==== print type ====
    main.User
    ==== print value ====
    {Username:"111" Departname:"111" Status:1}
```
通过上面两个方法我们知道变量user2是User类型，值也是我们初始化时的值。接下来我们就需要循环遍历t.NumField()和t.Field(i)来拆分里面的值
```go
      // 字段名
      var fieldName []string
      // ? 占位符
      var placeholder []string
      for i := 0 ; i< t.NumField(); i++ {
          // 小写开头考，无法反射 todo 这块有点疑惑 CanInterface方法的底层实现。
          if !v.Field(i).CanInterface() {
              continue
          }
          sqlTag := t.Field(i).Tag.Get("sql")
          if sqlTag != "" {
              // 跳过自增字段
              if strings.Contains(strings.ToLower(sqlTag), "auto_increment") {
                  continue
              } else {
                  fieldName =  append(fieldName, strings.Split(sqlTag, ",")[0])
              }
          } else {
              fieldName = append(fieldName, t.Field(i).Name)
          }
          placeholder =  append(placeholder, "?")
          // 字段的值
          u.AllExec = append(u.AllExec, v.Field(i).Interface())
      }
      // 拼接表，字段名，占位符
      u.Prepare =  "insert into " + u.GetTable() + " (" + strings.Join(fieldName, ",") + ") values(" + strings.Join(placeholder, ",") + ")"
```
如上所示，t.NumField()可以获取到这个结构体有多少字段用于for循环，t.Field(i).Tag.Get("sql")可以获取到包含sql:"xxx"的tag值，我们用来sql匹配和替换。
t.Field(i).Name可以获取到字段的field名字。通过v.Field(i).Interface()可以获取到字段的value值。e.GetTable()来获取我们设置的标的名字。通过上面的这一段稍微有点复杂的反射和拼接，我们就完成了Db.Prepare部分。
接下来我们需要获取stmt.exec里面的值的部分，上面我们将所有的值都放入了u.AllExec这属性中了。
  ```go
  // 第一步
  stmt ,err = u.db.Prepare(u.Prepare)
  
  // 第二步 执行exec,注意这是stmt.Exec
  result,err := stmt.Exec(e.AllExec...)
  if err != nil {
  //TODO
  }
  
  //获取自增ID
  id, _ := result.LastInsertId()
  ```
所以整个流程包括如下几步：
* 传入对应的struct类型元素
* 反射算出整个struct有多少元素，这样就好算出value后面需要几个`()`的占位符
* 搞for循环，得出子元素的name和value， 并且生成一个? 占位符
* 将每个name字段放入切片生成 sql语句，同时添加占位符。
* 将value都放入统一的AllExec中

#### where
* 结构体参数调用
  下面我们开始实现where方法的逻辑，这个where主要是为了实现替换sql语句中where后面部分的逻辑，例如原生的sql：
  ```go
  select * from userinfo where status = 1
  delete from userinfo where status = 1 or departname != "aa"
  update userinfo set departname = "bb" where status = 1 and departname = "aa"
  ```
  所以将where后面的数据单独拆出来，改成一个where方法很有必要，大部分的orm也是这样做的。并且通过上述的方法可以看到，where部分也是一样，先用Prepare生成问号占位符，在和exce替换值得方式来操作
  ```go
  stmt, err := db.Prepare("delete from userinfo where uid=?")
  
  result3, err := stmt.Exec("10795")
  
  stmt, err := db.Prepare("update userinfo set username=? where uid=?")

  result, err := stmt.Exec("lisi", 2)
  ```
  所以where的拆分，其实也是分位两部分走，和插入的逻辑基本一致。
  ```go
  type User struct {
  
      Username   string `sql:"username"`
  
      Departname string `sql:"departname"`
  
      Status     int64  `sql:"status"`
  
  }
  user2 := User{
  
      Username:   "111",
  
      Departname: "111", 
  
      Status:     1,
  
  }
  
  result1, err1 := u.Table("userinfo").Where(user2).Delete()
  
  result2, err2 := u.Table("userinfo").Where(user2).Select()
  ```
  我们这次实现的是where部分，这部分不会具体去执行结果，他做的仅仅是将数据拆分出来，用两个新的子元素whereParam和whereExec来暂存数据
  给最后的curd操作方法来使用。具体实现如下：
  ```go
  func (u *YEngine) Where(data interface{}) *YEngine {
  // 反射type和value
  t := reflect.TypeOf(data)
  v := reflect.ValueOf(data)
  
  // 字段名
  var fieldNameArray []string
  for i := 0; i < t.NumField(); i++ {  
    // 检测是否可以反射
    if !v.Field(i).CanInterface() {
        continue
    }
    // 解析tag,寻找真实的sql字段名
    sqlTag := t.Field(i).Tag.Get("sql")
    if sqlTag != nil {
        fieldNameArray = append(fieldNameArray, strings.Split(sqlTag, ",")[0]+"=?")
    } else {
        fieldNameArray = append(fieldNameArray, t.Field(i).Name+"=?")
    } 
    u.WhereExec = append(u.WhereExec, v.Field(i).interface())
  }
  
  // 拼接
  u.WhereParam += strings.Join(fieldNameArray, "and") 
  return u
  }
  ```
  这样，我们就可以调用Where()反复，转换成生成了2个暂存变量。我们打印下这2个值看看：
  ```go
  WhereParam = "username=? and departname=? and Status=?"
  WhereExec = []interface{"111", "111", 1}
  ```
  由于where是中间方法，可以多次调用，所以在第二次调用时，需要拼接上次的调用的结果；代码处理如下：
  ```go
  if u.WhereParam != "" {
    u.WhereParam += " and ("
  } else {
    u.WhereParam += "("
  }
  ```
  但是这里的实现为了简单，都是按照=的逻辑来实现的，传入的参数也只有一个interface。实际上我们在使用时，会有多种比较符，这就涉及了下面单个参数的调用。

* 单个参数的调用

  上面的where方法的参数，其实是我们和insert一样，传入的是个结构体，但是有时候，我们仅仅只需要查询1个字段，如果再去定义&实例化结构体，就显得比较麻烦。所以需要orm更加灵活的支持条件的增加。如下：
  ```go
    where("uid", ">=", 1223)
    where("uid", "=", "fsdfsd")
    where("uid", "in", []interface{1,2,3})
  ```
  这样我们可以使用其他非=的表达方式，比如：!=, like, not in, in等
  针对与这样的表现形式，我们在实现的时候对比下结构体，可以知道方法需要3个入参，第一个是需要查询字段，第二个是比较符号，第三个是查询的值。具体实现如下：
  ```go
  fun (u *YEngine) Where(fieldName string, opt string, fieldValue string) *YEngine {
    lowerOpt := strings.Trim(strings.ToLower(fieldName.(string)))
    
    if lowerOpt == "in" || lowerOpt == "not in" {
        // 判断传入的是否是切片
        fieldType := reflect.TypeOf(fieldValue).Kind()
        if fieldType != reflect.Slice && fieldType != reflect.Array {
            panic("in / not in, fieldValue is not slice or array")
        }
        v := reflect.ValueOf(filedValue)
        dataNum := v.Len()
        // 占位符
        ps := make([]string, dataNum)
        for i:=0; i< dataNum; i++ {
            ps[i] = "?"
            u.WhereExec = append(u.WhereExec, v.index(i).interface())
        } 
        //拼接
        u.WhereParam += fieldName.(string) + " " + lowerOpt + " (" + strings.Join(ps, ",") + ")"
    } else {
        u.WhereParam += fieldName.(string) + " " + lowerOpt + " ? "
        e.WhereExec = append(e.WhereExec, fieldValue)
    }
    return e
  }
  ```
  上面的代码唯一需要注意的就是针对in操作的特别处理。

#### orWhere
上面的where方法，生成的数据块之间都是and的关系，但是我们的sql有一些是or的关系，比如：
```go
where uid >= 123  or name = "vv"
where (uid = 123 and name = 'vv') or (uid = 456 and name = 'bb')
```
这种情况只需要新加一个orWhereParam参数，替代上面的方法中的whereParam即可，whereExec不需要变化，然后把拼接关系改成or，其他的代码基本是一致的
```go
fun (u *YEngine)  OrWhere (fieldName string, opt string, fieldValue string) *YEngine {
    ....
    if u.WhereParam == "" {
        panic("orWhere 必须在 Where之后调用")
    }
    ...
    u.OrWhereParam += "or ("
    ...
	return u
}
```
在这块我的想法是想将and， or这样的关系连接符剥离出来，而非大多数的orm的实现方式。具体实现如下：
```go
fun (u *YEngine) Or() *YEngine {
    if u.WhereParam == "" {
        panic("or必须在Where之后调用") 	
    }
    u.WhereParam += " or ("
	return u
}

func (u *YEngine) And() *YEngine {
    if u.WhereParam == "" {
        panic("and 必须在where之后调用")
    }
    u.WhereParam += " and ("
	return u
}
```

#### delete
删除也是我们经常会用到的结果，当我们完成了前面的where和or、and的数据逻辑绑定后。其实写delete方法就很简单了。
方法实现如下：
```go
func (u *YEngine) Del() (count int64, err error) {
    // 初始化声明
    var (
        stmt *sql.stmt
    )
	u.Prepare = "delete from "+ u.GetTable()
	// 如果where不为空
    if u.WhereParam != "" {
        u.Prepare += "where " + u.WhereParm
    }
	// order by 不为空
    if u.OrderParam != "" {
        u.Prepare += " order by " + u.OrderParam 	
    }
	// limit 不为空
    if u.limitParam != "" {
        u.Prepare += " limit " + u.LimitParam	
    }
	// prepare
	if stmt, err = e.Db.Prepare(e.Prepare); err != nil {
	    return 	
    }
	u.AllExec = u.WhereExec
	
	// 执行exec，这块是stmt.Exec
    result, err := stmt.Exec(e.AllExec)
    if err != nil {
        return 
    }
	
    count, err = result.RowsAffected()
	return 
}
```
我们在这里只需要拼接所有的参数， 并且执行对应逻辑就可以，这里的返回一般的实现会返回对应影响的行数，和select， update等还是不太一样的
具体实践如下：
```go
rowsAffected, err := e.Table("userinfo").Where("uid", ">=", 123).Delete()
```

#### update
修改数据这部分和delete也是基本类似的，但是在入参的时候可以入参两个，也可以入参一个结构体，具体如下：
```go
u.Table("userinfo").Where("uid", 123).Update("status", 1)
u.Table("userinfo").Where("uid", 123).Update(user2)
```
所以update需要支持既能支持一个参数，也可以支持两个参数，具体实现如下：
```go
func (u *YEngine) Update(data interface{}) (count int64 ,err error) {
    var (
        dataType int
    )
    switch len(data) {
    case 0 :
        err = errors.New("参数个数错误")
        return
    case 1 :
        dataType = 1
    case 2 :
        dataType = 2
    default:
        err = errors.New("参数个数错误")
        return 
    } 
    // 如果是结构体
    if dataType == 1 {
        t := reflect.TypeOf(data[0])
        v := reflect.ValueOf(data[0])
        fieldName := make([]string, 0)
        for i := 0 ; i< t.NumField(); i++ {
            if v.Field(i).CanInterface() {
                continue
            }
            //解析tag,找出真实的sql字段名
            sqlTag := t.Field(i).Tag.Get("sql")
            if sqlTag != "" {
                fieldName = append(fieldName, strings.Split(sqlTag, ",")[0] + "=?")
            } else {
                fieldName = append(fieldName, t.Field(i).Name+"=?")
            }
            u.UpdateExec = append(u.UpdateExec, v.Field(i).Interface())
        }
        u.UpdateParam += strings.Join(fieldName, ",")
    } else {
        u.UpdateParam += data[0].(string) + "=?"
        u.UpdateExec = append(u.UpdateExec, data[1])
    }
    // 拼接sql
	u.Prepare = "update " + u.GetTable() + "set" + u.UpdateParam
	
	
    // 如果whereParam不为空
    if u.WhereParam != "" || u.OrWhereParam != "" {
        u.Prepare += "where " + u.WhereParam + u.OrWhereParam 
    }
	
    // limit不为空
    if u.LimitParam != "" {
        u.Prepare += "limit " + u.LimitParam 
    }
    
    // Prepare
    var stmt *sql.stmt    
    var err error
    stmt, err := u.Db.Prepare(u.Prepare)
    if err != nil {
        return 0, u.setErrorInfo(err)		
    }
    // 合并UpdateExec和WhereExec   
    if u.WhereExec != nil {
        u.AllExec = append(u.AllExec, u.WhereExec...)
    }
    // 执行exec，
    result, err := stmt.Exec(e.AllExec...) 
    if err != nil {
        return 0, u.setErrorInfo(err)
    }
    // 影响行数
    id,_ := result.RowsAffected()
    return id , nil
}
```


#### select, 返回值为map切片
查询也是平时使用更多的地方，上述的几个方法我们实现了增删改，也熟悉了对应的写法。查询比较奇葩的地方在于使用了QueryRow和Query方法来获取数据，具体示例如下：
```go
// 单条数据查询
var usernmae, departname, status string
err := db.QueryRow("select username, departname, status from userinfo where uid=?", 4).Scan(&username, &departname, &status)
if err != nil {
	fmt.Println("QueryRow error:", err.Error())
}
fmt.Println("username: ", username, "departname: ", departname, "status: ", status)

```

多条查询一般是将数据序列化到一个多维的结构体中,通过for循环去给对应的结构体进行赋值,需要注意的是select出来的字段数需要和结构体的字段数保持一致,不然会丢失。
```go
// 多条数据查询
rows, err := db.Query("select username, departname, created from userinfo where username=?", "yang")
if err != nil {
	fmt.Println("queryRow error :", err.Error())
}
type userInfo struct {
    Username   string `json:"username"`
    Departname string `json:"departname"`
    Created    string `json:"created"`
}
var user []userInfo

for rows.Next() {
    var username1,departname1,created1 string
    if err := rows.Scan(&username1, &departname1,&created1); err != nil {   
        fmt.Println("query error :", err.Error())
    }
    user = append(user, userInfo{Username: username1, Departname: departname1, Created: created1})
}
```
如上操作后，我们需要提前定义结构体，并且还需要区分单条和多条的数据查询。实现上比较麻烦，并且在使用时也是个问题，用户需要区分单条和多条的情况。
所以我们需要整合单条和多条的情况，考虑到要提前初始化一个数据结构，在初始化成一个数组，这样太过于麻烦，简单的做法是直接按照数据库表里的字段名，
直接初始化出一个同名的map切片。具体实现如下：
```go
result, err := u.Table("user_info").Where("status", "=", 1).Select()
// 返回为
[
    map[
        departname:v 
        status:1 
        uid:123 
        username:yang
    ] 
    map[
        departname:n 
        status:0 
        uid:456 
        username:small
    ]
]
```
这种实现方式的前提是我们可以获取表的字段有哪些，才能根据把这些字段映射成一个map。db.query给我们返回了一个columns()方法，他能返回我们本次查询的表的字段名是哪些。具体实现如下：
```go
rows, err := db.query("select uid, username, departname, status from userinfo where username=?", "yang")
if err != nil {
	fmt.Println("query error:", err.Error())
}
column, err := rows.Columns()
if err != nil {
    fmt.Println("query error : ", err.Error())
}
fmt.Println(column)
// 返回输出
[uid username departname stauts]
```
至此我们获取到了数据库表的字段名，我们需要按照这些字段名进行rows.Scan()数据绑定。由于我们没有预先定义数据类型进行绑定，所以这个数据我们只能动态生成。如下代码，原生的for循环的时候在Scan时通过地址来动态引用赋值。所以这4个字段的名字不重要，最后赋值的是这4个变量指向的地址空间。
```go
for rows.Next() {
  var uid1, username1, departname1, status1 string
  rows.Scan(&uid1, &username1, &departname1, &status1)
  fmt.Println(uid1,username1,departname1,status1)
}
```
正是利用了这一点，所以我们可以按照Columns返回的字段个数，来解决映射关。具体实现如下：
```go
// 读出查询出的列字段名
columns, err := rows.Columns()
if err != nil {
    return nil, u.SetErrorInfo(err)
}
// values是每个列的值，这里获取到bytes里
values := make([][]bytes, len(columns))

// 因为每次查询出来的列都是不定长的，用len(columns)定住当次查询的长度
scans := make([]interface{}, len(columns))

for i := range values {
    scans[i] = &values[i]
}
```
这里新建两个切片，第一个切片是values，初始值是空的，scans初始值是一个空接口类型的切片，通过一个for循环，让scans每个元素的值，都是values里的每个值得地址。
这样做的好处在于，我们在使用Scan()方法时，可以传递进去scans。如下
```go
for rows.Next() {
  rows.Scan(scans[0], scans[1],scans[2], scans[3])
}
```
这样，scans[0]对应到上面的uid1，scans[3]对应到上述的status1，由于scans[3]存储的是对应values下标的存储地址空间，所以values也随之更改。
然后我们通过三个切片的下标映射，就能将表字段和值对应起来，拼接成1个map。当然这种方式也是有致命缺点的，就是必须知道columns的数量，如果过多，
那么Scan方法中就要写多个，这种原始方式肯定不能满足我们的诉求。所以就有了如下的优化：
```go
results, err := make([]map[string]string, 0)
for rows.Next() {
    if err := rows.Scan(scans...); err != nil {
        return nil, u.setErrorInfo(err)
    }
	
	// 每行数据
    row := make(map[string]string)
    // 循环values数据，通过相同的下标，取得columns里对应的列名，生成1个新的map
    for k, v := range values {
        key := columns[k]
        row[key] = string(v)
    }
    // 添加到map中
	results = append(results, row)
}
```
这里关键的点在于 rows.Scan(scans...)这个写法，可以将切片的字段铺开，解决了字段过多的问题。完整的select的实现如下：
```go
func(u *YEngine) Select() ([]map[string]string, error) {
    // 拼接sql
    u.Prepare = "select * from " + u.GetTable()
    // 如果whereParams不为空
    if u.WhereParam != "" || u.OrWhereParam != "" {
        u.Prepare += "where " + u.WhereParam + u.OrWhereParam
        u.AllExec = u.WhereExec
    }
    // order by 不为空
    if u.OrderParam != "" {
        u.Prepare += " order by " + u.OrderParam
    }
    if u.LimitParam != "" {
        u.Prepare += "limit " + u.LimitParam
    }
	
    // query
    rows, err := u.Db.Query(e.Prepare, e.AllExec)
    if err != nil {
        return nil , u.setErrorInfo(err)
    }
    // 读出查询出的列字段名
    columns, err := rows.Columns()
    if err != nil {
        return nil , u.setErrorInfo(err)
    }
    // values是每个列的值，这里获取到byte里
    values := make([][]byte, len(colums))
    // 因为每次查询出来的列都不一定定长，用len(columns)定住当前长度
    scans := make([]interface, len(columns))
    for i := range values {
        scans[i] = &values[i]
    }
    for rows.Next() {
        if err := rows.Scan(scans...); err != nil {
            return nil , u.setErrorInfo(err)
        }       
        // 每行数据
        row :=  make(map[string]string)
        for k,v = range values {
            key := colums[k]    
            row[key] = string(v)
        }
        //添加到map切片中
        results = append(results, row)
    }
    return results, nil
}
```
这里我们就可以非常方便的查询数据，这里需要注意的有两个点：
* 最后返回的map切片，里面的key名是数据库字段名，如果要映射成首字母大写的结构，需要我们自己写方法
* select会将数据库中的所有字段的类型转化为字符串类型，如果需要转化成对应的类型，也需要我们写方法

虽然网上有各种查询单条的方法，但是个人觉得没太大必要单独抽一个方法出来，这样给使用人员的选择就造成了一定迷惑。不如通过select方法和limit方法组合的方式来限制单条的做法

#### 查询多条Find() ,返回值为结构体切片
这个方法其实是对原生go的一个简单包装，通过预先定义好数据结构，然后通过引用赋值。目前大多数的orm都是如下实现。
```go
type User struct (
  Uid        int    `sql:"uid,auto_increment"`
  Username   string `sql:"username"`
  Departname string `sql:"departname"`
  Status     int64  `sql:"status"`
)
var user1 []User
// select * from userinfo where status=1
err := u.Table("userinfo").Where("status", 2).Find(&user1)
if err != nil {
    fmt.Println(err.Error())
} else {
    fmt.Println("%#v", user1)
}
// 输出
[]User{smallorm.User2{Uid:131733, Username:"EE2", Departname:"223", Status:2}, smallorm.User{Uid:131734, Username:"EE2", Departname:"223", Status:2}, smallorm.User{Uid:131735, Username:"EE2", Departname:"223", Status:2}}
```
当前的这种实现可以通过如下步骤实现：
* 先定义一个结构体。里面的字段通过tag标签和表的字段进行关联
* 初始化一个空的结构体切片，然后通过&取地址符传给Find()方法
* Find方法内部先获取到表的列名，通过tag关联和各种反射方法，将数据绑定到传入的结构体切片上，给他赋值

总体感受下来这个实现和上面select的实现基本类型，我们分布实现如下：
```go
// 读取查询出的列名字
column, err := rows.Columns()
if err != nil {
    return e.setErrorInfo(err)
}
// values是每个列的值，这里获取到byte里
values := make([][]byte, len(columns))

// 由于每次查询出来的是不定长的，用len(columns)定住当前的查询长度
scans := make([]interface{}, len(columns))
for i := range values{
    scans[i] = &value[i]
}
```
上面的这几步和之前select是类似的，关键是以下几步：
```go
// 原struct的切片值
destSlice := reflect.ValueOf(reslut).Elem()

//原始当个struct的类型
destType := destSlice.Type().Elem()

//打印下
fmt.Printf("%+v\n", destSlice)
fmt.Printf("%+v", destType)
[]
main.User
```
我们通过反射的两个方法，可以得出传入的User结构体切片是什么类型，他的值是什么。接下来我们通过反射继续解析：
```go
// 获取到User结构体的字段数，这里返回4
destType.NumField()

// 获取到User结构体的第i个字段的tag值，比如返回username
destType.Field(i).Tag.Get("sql")

// 获取到User结构体的第i个字段的名字,比如返回:Username
destType.Field(i).Name

再通过以下几个反射赋值

// 根据类型生成1个新的值，返回：{Uid:0 Username: Departname: Status:0}
dest := reflect.New(destType).Elem()
// 给第i个元素，附值，类型是string类型
dest.Field(i).SetString(value)
// 将dest值添加到destSlice切片中。
reflect.Append(destSlice,dest)
// 将最后得到的切片完全赋值给本身
destSlice.Set(reflect.Append(destSlice,dest))
```
接下来我们看下完整的Find方法实现：
```go
//查询多条返回值为struct切片
func (u *YEngine) Find(result interface{}) error {
	if reflect.ValueOf(result).Kind != reflect.Ptr {
        return u.SetErrorInfo(errors.New("参数请传指针变量"))    
    }
    if reflect.ValueIf(result).IsNil() {
        return u.SetErrorInfo(errors.New("参数不能是空指针"))    
    }
    // 拼接sql
    u.Prepare = "select * from " + u.GetTable()
    // 如果whereParams不为空
    if u.WhereParam != "" || u.OrWhereParam != "" {
        u.Prepare += "where " + u.WhereParam + u.OrWhereParam
        u.AllExec = u.WhereExec
    }
    // order by 不为空
    if u.OrderParam != "" {
        u.Prepare += " order by " + u.OrderParam
    }
    // limit 不为空
    if u.LimitParam != "" {
        u.Prepare += "limit " + u.LimitParam
    }
    // query
    rows, err := e.Db.Query(e.Prepare, e.AllExec...)
    if err != nil {
        return e.setErrorInfo(err)
    }
    // 读出查询出的列字段名
    column, err := rows.Columns()
    if err != nil {
        return  e.setErrorInfo(err)   
    }
    // values是每个列的值, 这里获取到byte里
    values := make([][]byte,len(column))

    //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
    scans := make([]interface{}, len(column))	
    
    // 原始struct的切片值
    destSlice := reflect.ValueOf(result).Elem()

    //原始单个struct的类型
    destType := destSlice.Type().Elem()
    for i := range values {
        scans[i] = &values[i]
    }  
    //循环遍历    
    for rows.Next() {
        dest := reflect.New(destType).Elem()
        if err := rows.Scan(scans...); err != nil {
            return e.setErrorInfo(err)
        }
        // 遍历一行数据的各个字段
        for k ,v := range values {
            key := column[k]       
            value := string(v)
            for i:=0; i < destType.NumField();i++{
                // 看下是否有sql别名  	
                sqlTag := destType.Field(i).Tag.Get("sql")
                var fieldName string 
                if sqlTag != "" {
                    fieldName = strings.Split(sqlTag, ",")[0]			
                } else {
                    fieldName = destType.Field(i).Name    
                }   
				if k != fieldName {
                    continue
                }
                // 反射赋值
                if err := u.reflectSet(dest,i,value); err != nil {
                    return err
                }       
            }
        }
        destSlice.Set(reflect.Append(destSlice, dest))
    }
    return nil
}

```
我们在方法前面增加了几个参数校验，都是基于反射的，来判断传进来的值都是指针类型才行。在反射赋值里，我搞了个reflectSet来进行字段类型的匹配，具体实现如下：
```go
func (e *YEngine) reflectSet(dest reflect.Value, i int, value string) error {
  switch dest.Field(i).Kind() {
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    res, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
      return e.setErrorInfo(err)
    }
    dest.Field(i).SetInt(res)
  case reflect.String:
    dest.Field(i).SetString(value)
  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    res, err := strconv.ParseUint(value, 10, 64)
    if err != nil {
      return e.setErrorInfo(err)
    }
    dest.Field(i).SetUint(res)
  case reflect.Float32:
    res, err := strconv.ParseFloat(value, 32)
    if err != nil {
      return e.setErrorInfo(err)
    }
    dest.Field(i).SetFloat(res)
  case reflect.Float64:
    res, err := strconv.ParseFloat(value, 64)
    if err != nil {
      return e.setErrorInfo(err)
    }
    dest.Field(i).SetFloat(res)
  case reflect.Bool:
    res, err := strconv.ParseBool(value)
    if err != nil {
      return e.setErrorInfo(err)
    }
    dest.Field(i).SetBool(res)
  }
  return nil
}
```
这块网上也有单条查询返回结构体的，同样的个人感觉没必要。直接增加limit，按struct数组返回就可以。同样的网上也有针对find设置查询字段的，个人理解也是没必要的，可以全查

#### 设置查询字段Field
设置字段是一个很基础的其实也很重要的功能，因为我们平时查询数据的时候，都喜欢用select * ，这种方式在只需要某几个字段，但是会查询大量数据的情况下，是十分浪费和低效的。
所以，这里像大多数orm一样实现了Field方法，来执行只查找某些字段。
```go
u.Table("userinfo").Where("status", 2).Field("uid,status").Select()
```
由于是链式调用， 再加上field方法参数本事没有数据属性，所以放在链式调用的中间位置就可以，实现逻辑也比较容易，就是给YEngine的FieldParam赋值就行，在select/find中替换*
具体实现如下：
```go

func (u *YEngine) Field (field string) *YEngine {
    u.FieldParam = field
	return u
}

// select/find方法替换*
u.Prepare = "select " + e.FieldParam + " from " + e.GetTable()
```
u.FieldParam的初始值是*，这是在NewMysql初始化中实现的。这里没有对传入的field参数进行校验，后面可以试着优化下(如果有时间的话，嘻嘻嘻)

#### 设置大小Limit
limit一般用来限制每次查询获取的数量， 更多的会用在分页的场景下，比如后台的列表等。通常的实现如下：
```go
u.Table("userinfo").Where("status", 2).Limit(1).Select()
```
所以具体实现如下：
```go
func (u *YEngine) Limit(limit int)  *YEngine {
    u.LimitParam = limit
    return u
}
```

#### 聚合查询t
聚合查询除了count之外，还有svg，max，min等。以count为例，count是用来获取当前查询数据行数，他的实现方式都是将原来的select * 换成select count(*)，或者是count(1)，因此在实现的时候需要将对应的对应的聚合函数名，和参数传入进去，其实这部分的count可以通过select全部数据之后在业务层来进行计算长度，但是为了完整我们部分实现如下：
```go
func (u *YEngine) AggregateQuery(name, param string) (interface{}, error) {
    u.Prepare = "select " + name + "(" + param + ") as cnt from " + e.GetTable()
}
```
我们申明了一个cnt，用他来获取最终的聚合结果值，之所以是接口类型，是因为聚合的对象类型是不定的，可能是int型，也可能是float型。完整的实现如下：
```go
func (u *YEngine) AggregateQuery(name, param string) (interface{}, error) {
    // 拼接sql
    u.Prepare = "select " + name + "(" + param + ") as cnt from " + e.GetTable()
    
    // 如果whereParam不为空，或者OrWhereParam不为空
    if u.WhereParam != "" || u.OrWhereParam != "" {
        u.Prepare += "where " + u.WhereParam + u.OrWhereParam
    }
	
    // 如果limit不为空
    if u.LimitParam != "" {
        u.Prepare += "limit " + u.LimitParam     
    }
        
    u.AllExec = u.WhereExec
    u.generateSql()
    //执行绑定
    var cnt interface{}
    err := e.Db.QueryRow(e.Prepare, e.AllExec...).Scan(&cnt)
    
    if err != nil {
    
    return nil, e.setErrorInfo(err)
}

// 生成完整的sql，这里的处理比较简陋
func (u *YEngine) generateSql() {
    u.sql = u.Prepare
    for _,i2 := range u.AllExec {
        switch i2.(type) {
            case int:
            u.sql = strings.Replace(u.sql, "?", strconv.itoa(i2.(int)), 1)
            case int64:
            u.sql = strings.Replace(u.sql, "?", strconv.FromatInt(i2.(int64), 10), 1)
            case bool:
            u.sql = strings.Replace(u.sql, "?", strconv.FormatBool(i2.(bool)),1)
            default :   
            u.sql = strings.Replace(u.sql, "?", "'" + i2.(string) + "'", 1)
        }
    }
}
```
这样我们总体就完成了函数的编写，但是各个聚合函数之间可能会有一定的差别。比如count，可以使count()，也可以是count(1)，具体我们在实现的时候可以固定下来。所以count的具体实现方式如下：
```go

func (u *YEngine) Count() (int64, error) {
    count, err := u.aggregateQuery("count", "*")
    if err != nil {
        return 0, u.setErrorInfo(err)
    }
    return count.(int64), err
}
```
#### 获取最大值Max
可以使用max()方法来获取一个字段的最大值，返回总数的类型时string类型。她是链式调用结构最后一次操作。之所以返回值是string，
是因为取最大值，有时候不限制在int类型的表字段最大值，有时候也会有时间最大值等，所以返回string是最合适的第一个参数我们传一个
max，第二个字段我们传表的一个字段，具体实现如下：
```go
func (u *YEngine) Max(param string) (string, error) {
    max, err := u.aggregateQuery("max", param)
    if err != nil {
        return "0", u.setErrorInfo(err)    
    }
    return string(max.([]byte)), nil
}
```


#### 获取最小值min
使用min()方法返回一个字段的最小值，基本和上述max类似，具体实现如下：
```go
func (u *YEngine) Min(param string) (string, error) {
    min, err := u.aggregateQuery("min", param)
    if err != nil {
        return "0", u.setErrorInfo(err)
    }
    return string(min.([]byte)), nil
}
```

#### 平均值Avg 
使用avg方法返回一个平均值，具体实现和上述类似：
```go
func (u *YEngine) Avg(param string) (string, error) {
  avg, err := u.aggregateQuery("avg", param)
  if err != nil {
    return "0", u.setErrorInfo(err)
  }
  return string(avg.([]byte)), nil

}
```
#### 获取总和Sum
使用sum方法返回一个总和，具体实现和上述类似：
```go
func (u *YEngine) Avg(param string) (string, error) {
  sum, err := u.aggregateQuery("sum", param)
  if err != nil {
    return "0", u.setErrorInfo(err)
  }
  return string(sum.([]byte)), nil

}
```

#### 排序Order
排序和limit比较类似，我们平时在使用都是如下写法：
```go
// 按照uid降序
select * from user_info where uid >= 10888 order by uid desc

// 按照uid升序
select * from user_info where uid >= 10888 order by uid asc

// 查询结果按照多个参数，不同排序
select * from user_info where uid >= 10888 order by uid asc, status desc
```
平时我们使用时会按照如上三种方式来使用orm，所以我们可以看到这是一个可变长度的入参。具体的处理如下：
```go
// 参数之间使用，分割。例如：id，desc，uid，acs
func (u *YEngine) Order(param ....string) *YEngine {
    orderLen := len(param)
	if orderLen % 2 != 0 {
        panic("order by 参数数量错误，需要偶数个")
    } 
    // 存在多次调用的情况
    if u.OrderParam != "" {
        u.OrderParam += ","
    }
    for i :=0; i < orderLen/2; i++ {
        keyString := strings.ToLower(param[i*2+1])
        if keyString != "desc" || keyString != "acs" {
            panic("排序关键字错误")
        }
        if i < orderLen / 2 -1 {
            u.OrderParam += param[i*2] + " " + param[i*2+1] + ","
        } else {
            u.OrderParam += param[i*2] + " " + param[i*2+1]
        }                    
    }
    return u
}
```

#### 分组Group
分组也是我们平时使用较多的地方，它用于我们对某些数据根据其中的字段进行分组，这个写法也是很简单的。具体实现如下：
```go
func (u *YEngine) Group(group ...string) *YEngine {
    if len(group) != 0  {
        u.GroupParam = strings.Join(group, ",")
    }
    return u
}
```
这样我们需要在select/find中加入对group的判断，如下：
```go

// 需要放在limit之前
if u.GroupParam != "" {
    e.Prepare += " group by " + e.GroupParam
}
```




# mysql代码模板生成