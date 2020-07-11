package dbutil

import "testing"

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

func TestInsertEngine_InsertMany(t *testing.T) {
	var user1 UserModel

	user1.UserId = "111"
	user1.Username = "test1"
	user1.Password = "123"
	user1.Realname = "测试1"
	user1.CreateId = "1"
	user1.CreateTime = "2020-07-09 23:38:38"

	var user2 UserModel

	user2.UserId = "333"
	user2.Username = "test2"
	user2.Password = "123"
	user2.Realname = "测试2"
	user2.CreateId = "1"
	user2.CreateTime = "2020-07-09 23:38:38"

	err := NewBatchInsertEngine(user1, user2).InsertMany()

	t.Log(err)
}
