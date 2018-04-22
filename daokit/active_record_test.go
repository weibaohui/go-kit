package daokit

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func check() {
	entity := New().SetTableName("auth_config")
	i, err := entity.Set("app_id", "x23").Set("key", "y").Set("token", nil).Insert()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(i)
}

type NotknownType struct {
	s1, s2, s3 string
	ff         int
	fff64      int64
}

var secret interface{} = NotknownType{"Ada", "Go", "Oberon", 5, 55}

func TestEntity_Insert(t *testing.T) {
	value := reflect.ValueOf(secret)
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		fieldValue := value.Field(i)
		fmt.Printf("Field %d: %s= %v \n", i, fieldName, fieldValue)
	}
}
func TestEntity_Delete(t *testing.T) {
	db := initOrm("root", "root", "127.0.0.1", "4000", "todo", "")
	entity := New().SetDB(db).SetTableName("auth_config")
	where := New().Set("app_id", "x").Set("token", "xxxx")
	entity.Delete(where)
}

func TestEntity_Insert2(t *testing.T) {
	db := initOrm("root", "root", "127.0.0.1", "4000", "todo", "")
	New().SetDB(db).SetTableName("auth_config").
		Set("app_id", "bpm").
		Set("token", "cYl8R3UY8a3i3YPC").
		Set("key", "94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn").
		Insert()
}

func TestEntity_Update(t *testing.T) {
	db := initOrm("root", "root", "127.0.0.1", "4000", "todo", "")
	where := New().Set("app_id", "bpm2")
	New().SetDB(db).SetTableName("auth_config").
		Set("app_id", "bpm").
		Set("token", "cYl8R3UY8a3i3YPC").
		Set("key", "94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn").
		Update(where)
}
func TestEntity_QueryNumber(t *testing.T) {
	db := initOrm("root", "root", "127.0.0.1", "4000", "todo", "")
	count, err := New().SetDB(db).QueryNumber("select count(*) from auth_config")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("count=" + strconv.FormatInt(count, 10))
	t.Log(count)
}

func TestEntity_QueryString(t *testing.T) {
	db := initOrm("root", "root", "127.0.0.1", "4000", "todo", "")
	appID, err := New().SetDB(db).QueryString("select app_id from auth_config order by app_id asc limit 1")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("appID=" + appID)
	t.Log(appID)
}

type config struct {
	AppID string
	Token string
	Key   string
}

func TestEntity_QueryOne(t *testing.T) {
	db := initOrm("root", "root", "127.0.0.1", "4000", "todo", "")
	sql := "select * from auth_config order by app_id asc limit 1"
	c := config{}
	err := New().SetDB(db).QueryOne(&c, sql)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(c)
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		t.Fatal(err.Error())
	}
	for rows.Next() {
		db.ScanRows(rows, &c)
	}
	t.Log(c)
}
