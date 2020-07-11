package dbutil

import (
	"testing"
)

type UserModel struct {
	UserId     string `json:"userId" db:"user_id" dbField:"true"`
	Username   string `json:"username" db:"username" dbField:"true"`
	Realname   string `json:"realname" db:"realname" dbField:"true"`
	Password   string `json:"password" db:"password" dbField:"true"`
	Remark     string `json:"remark" db:"remark" dbField:"true"`
	CreateId   string `json:"createId" db:"create_id" dbField:"true"`
	CreateTime string `json:"createTime" db:"create_time" dbField:"true"`
	UpdateId   string `json:"updateId" db:"update_id" dbField:"true"`
	UpdateTime string `json:"updateTime" db:"update_time" dbField:"true"`
	UserIds    string `json:"userIds" db:"user_id"`
}

func (user UserModel) GetTableName() string {
	return "sys_user"
}

func TestSelectEngine_SelectOne(t *testing.T) {
	var user UserModel

	user.UserId = "00000450AC2811E98CF67EB0F7C6F98E"

	selectEngine := NewSelectEngine(user)

	selectEngine.WhereEqs("userId")

	t.Log(selectEngine.SelectOne())
}

func TestSelectEngine_selectAll(t *testing.T) {
	var user UserModel
	user.Username = "lisi759"
	selectEngine := NewSelectEngine(user)
	selectEngine.WhereLeftLikes("username")
	selectEngine.OrderByAsc("username")
	t.Log(selectEngine.SelectAll())
}

func TestSelectEngin_SelectPage(t *testing.T) {
	var user UserModel

	user.Username = "lisi7"

	selectEngine := NewSelectEngine(user)
	selectEngine.WhereLeftLikes("username")
	selectEngine.OrderByAsc("username").Limit(1, 2)

	t.Log(selectEngine.SelectPage())
}
