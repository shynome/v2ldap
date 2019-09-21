package v2ray

import (
	"testing"
)

func TestGetDBUsers(t *testing.T) {
	db := getDB()
	defer db.Close()
	v2 := V2ray{
		DB: db,
	}
	type TestCase struct {
		Users []string
		Count int
	}
	u1, u2 := "1@x.y", "2@x.y"
	testCases := []TestCase{
		{
			[]string{u1},
			1,
		},
		{
			[]string{u1, u2},
			2,
		},
	}
	for index, test := range testCases {
		var err error
		if _, err = v2.Sync(test.Users, true); err != nil {
			t.Error(err)
			return
		}
		var users []User
		if users, err = v2.GetDBUsers(); err != nil {
			t.Error(err)
			return
		}
		var count = len(users)
		if count != test.Count {
			t.Errorf("Index %v expect %v , but get %v . the users list is %v", index, test.Count, count, test.Users)
			return
		}
	}
}
