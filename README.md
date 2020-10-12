# dbutil-mysql
 go mysql orm
 
 此库尚未在生产环境中使用！！！
 
# 依赖第三方库
 1、https://github.com/go-sql-driver/mysql v1.5.0

# 使用指南

## 1.struct 定义规则

```go
type UserModel struct {
    UserId     string    `json:"userId" db:"user_id" dbField:"true" dbType:"varchar"`
    Username   string    `json:"username" db:"username" dbField:"true" dbType:"varchar"`
    Realname   string    `json:"realname" db:"realname" dbField:"true" dbType:"varchar"`
    Password   string    `json:"password" db:"password" dbField:"true" dbType:"varchar"`
    Remark     string    `json:"remark" db:"remark" dbField:"true" dbType:"varchar"`
    CreateId   string    `json:"createId" db:"create_id" dbField:"true" dbType:"varchar"`
    CreateTime string    `json:"createTime" db:"create_time" dbField:"true" dbType:"varchar"`
    UpdateId   string    `json:"updateId" db:"update_id" dbField:"true" dbType:"varchar"`
    UpdateTime LocalTime `json:"updateTime" db:"update_time" dbField:"true" dbType:"datetime"`
    UserIds    string    `json:"userIds" db:"user_id" dbType:"varchar"`
}

// UserModel需要实现GetTableName方法，返回UserModel对应的数据库表名
func (user UserModel) GetTableName() string {
    return "sys_user"
}
```
tag说明：

json: struct字段对应前端字段名称  
db: struct字段对应数据库表字段名称  
dbType: 数据库表字段的数据类型  
dbField: struct字段为数据库表字段时，配置为true  

## 2.初始化连接池

InitPool("root", "root", "127.0.0.1", 3306, "demo", "utf8mb4")

## 3.新增操作
```go
func TestInsertEngine_Insert(t *testing.T) {
    var user UserModel
    
    user.UserId = "123"
    user.Username = "test"
    user.Password = "123"
    user.Realname = "测试"
    user.CreateId = "123"
    user.CreateTime = "2020-07-09 23:38:38"

    err := NewInsertEngine(user).Insert()

    t.Log(err)
}
```

未完待续...
