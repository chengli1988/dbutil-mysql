package dbutil

import "testing"

func TestDeleteEngine_Delete(t *testing.T) {
	var user UserModel

	user.UserId = "444"

	deleteEngine := NewDeleteEngine(user)

	deleteEngine.WhereEqs("userId")

	t.Log(deleteEngine.Delete())
}
