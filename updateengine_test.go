package dbutil

import (
	"testing"
)

func TestUpdateEngine_Update(t *testing.T) {
	var user UserModel

	user.UserId = "111"
	user.Username = "test1111111"
	user.Password = "123456"
	user.Realname = "测试11111111"
	user.CreateId = "1"
	user.CreateTime = "2020-07-09 23:38:38"
	user.UpdateId = "1"
	user.UpdateTime = "2020-07-10 23:38:38"

	updateEngine := NewUpdateEngine(user)
	updateEngine.WhereEqs("userId")

	t.Log(updateEngine.Update())
}

func TestUpdateEngine_ReplaceIntoMany(t *testing.T) {
	var user1 UserModel

	user1.UserId = "111"
	user1.Username = "test1"
	user1.Password = "123"
	user1.Realname = "测试1"
	user1.CreateId = "1"
	user1.CreateTime = "2020-07-09 23:38:38"

	var user2 UserModel

	user2.UserId = "444"
	user2.Username = "test2"
	user2.Password = "123"
	user2.Realname = "测试2"
	user2.CreateId = "1"
	user2.CreateTime = "2020-07-09 23:38:38"

	err := NewBatchUpdateEngine(user1, user2).ReplaceIntoMany()

	t.Log(err)
}
