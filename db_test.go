package dbutil

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func init() {
	InitPool("root", "root", "127.0.0.1", 3306, "demo", "utf8mb4")
}

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

func (user UserModel) GetTableName() string {
	return "sys_user"
}

// 新增例子
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

//  查询例子：查询一行记录
func TestSelectEngine_SelectOne(t *testing.T) {
	var user UserModel

	user.UserId = "123"

	selectEngine := NewSelectEngine(user)

	selectEngine.WhereEqs("userId")

	rows, _ := selectEngine.SelectOne()

	t.Log("rows:", rows)

	var user2 UserModel

	mapToStruct(rows, &user2)

	log.Println(user2)
	t.Log(user2.UpdateTime)
}

// map数据转成Struct
func mapToStruct(mapResult map[string]interface{}, targetStruct interface{}) {

	jsonResult, err := json.Marshal(mapResult)

	if err != nil {
		log.Println(err)
	} else {
		json.Unmarshal(jsonResult, targetStruct)
	}

}

// 修改例子
func TestUpdateEngine_Update(t *testing.T) {
	var user UserModel

	user.UserId = "123"
	user.Username = "test1111111"
	user.Password = "123456"
	user.Realname = "测试11111111"
	user.CreateId = "1"
	user.CreateTime = "2020-07-09 23:38:38"
	user.UpdateId = "1"
	user.UpdateTime = LocalTime(time.Now())

	updateEngine := NewUpdateEngine(user)
	updateEngine.WhereEqs("userId")

	t.Log(updateEngine.Update())
}

// 删除例子
func TestDeleteEngine_Delete(t *testing.T) {
	var user UserModel

	user.UserId = "123"

	deleteEngine := NewDeleteEngine(user)

	deleteEngine.WhereEqs("userId")

	t.Log(deleteEngine.Delete())
}

// 查询符合条件的所有数据
func TestSelectEngine_selectAll(t *testing.T) {
	var user UserModel
	user.Username = "lisi759"
	selectEngine := NewSelectEngine(user)
	selectEngine.WhereLeftLikes("username")
	selectEngine.OrderByAsc("username")
	t.Log(selectEngine.SelectAll())
}

// 查询符合条件的分页数据
func TestSelectEngin_SelectPage(t *testing.T) {
	var user UserModel

	user.Username = "lisi7"

	selectEngine := NewSelectEngine(user)
	selectEngine.WhereLeftLikes("username")
	selectEngine.OrderByAsc("username").OrderByDesc("createTime").Limit(1, 2)

	t.Log(selectEngine.SelectPage())
}

// 插入多条数据
func TestInsertEngine_InsertMany(t *testing.T) {
	var user1 UserModel

	user1.UserId = "111"
	user1.Username = "test1"
	user1.Password = "123"
	user1.Realname = "测试1"
	user1.CreateId = "1"
	user1.CreateTime = "2020-07-09 23:38:38"

	var user2 UserModel

	user2.UserId = "222"
	user2.Username = "test2"
	user2.Password = "123"
	user2.Realname = "测试2"
	user2.CreateId = "1"
	user2.CreateTime = "2020-07-09 23:38:38"

	users := make([]Model, 0)

	users = append(users, user1)
	users = append(users, user2)

	err := NewBatchInsertEngine(users).InsertMany()

	t.Log(err)
}

// 替换数据
func TestUpdateEngine_ReplaceIntoMany(t *testing.T) {
	var user1 UserModel

	user1.UserId = "111"
	user1.Username = "test111"
	user1.Password = "123"
	user1.Realname = "测试1111"
	user1.CreateId = "1"
	user1.CreateTime = "2020-07-09 23:38:38"

	var user2 UserModel

	user2.UserId = "333"
	user2.Username = "test222"
	user2.Password = "123"
	user2.Realname = "测试2222"
	user2.CreateId = "1"
	user2.CreateTime = "2020-07-09 23:38:38"

	users := make([]Model, 0)

	users = append(users, user1)
	users = append(users, user2)

	err := NewBatchUpdateEngine(users).ReplaceIntoMany()

	t.Log(err)
}
