package v2ray

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // init sqlite
)

var db *gorm.DB

func getDB() *gorm.DB {
	if db == nil {
		var err error
		os.Remove("test.db")
		db, err = gorm.Open("sqlite3", "test.db")
		if err != nil {
			panic(err)
		}
		db.AutoMigrate(&User{})
	}
	return db
}

func TestSync(t *testing.T) {
	db := getDB()
	defer db.Close()
	v2 := V2ray{
		DB: db,
	}
	u1, u2, u3 := "1@x.y", "2@x.y", "3@x.y"
	type TestResult struct {
		Add    int
		Delete int
	}
	type TestCase struct {
		Users   []string
		Expect  TestResult
		Confirm bool
	}
	testCases := []TestCase{
		{
			[]string{u1},
			TestResult{Add: 1, Delete: 0},
			false,
		},
		{
			[]string{u1, u2},
			TestResult{Add: 2, Delete: 0},
			false,
		},
		{
			[]string{u1},
			TestResult{Add: 1, Delete: 0},
			true,
		},
		{
			[]string{u2},
			TestResult{Add: 1, Delete: 1},
			true,
		},
		{
			[]string{u1, u2, u3},
			TestResult{Add: 2, Delete: 0},
			true,
		},
		{
			[]string{u1, u2, u3},
			TestResult{Add: 0, Delete: 0},
			true,
		},
	}
	for index, test := range testCases {
		var result SyncResponse
		var err error
		if result, err = v2.Sync(test.Users, test.Confirm); err != nil {
			t.Error(err)
			return
		}
		if len(result.Add) != test.Expect.Add {
			t.Errorf("Index %v expect %v , but get %v . the users list is %v", index, test.Expect.Add, len(result.Add), test.Users)
			return
		}
		if len(result.Delete) != test.Expect.Delete {
			t.Errorf("Index %v expect %v , but get %v . the users list is %v", index, test.Expect.Delete, len(result.Delete), test.Users)
			return
		}
	}
}
