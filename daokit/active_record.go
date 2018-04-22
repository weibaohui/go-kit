package daokit

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/weibaohui/go-kit/strkit"
	"reflect"
	"strings"
)

type Entity struct {
	db        *gorm.DB
	tableName string
	data      map[string]interface{}
}

func New() *Entity {
	entity := &Entity{
		data: make(map[string]interface{}),
	}
	return entity
}

// initDB 如果没有指定db，那么取默认的GetOrm()
func (e *Entity) initDB() *Entity {
	if e.db == nil {
		e.db = GetOrm()
	}
	return e
}
func (e *Entity) SetDB(db *gorm.DB) *Entity {
	e.db = db
	return e
}

func (e *Entity) SetTableName(name string) *Entity {
	e.tableName = name
	return e
}
func (e *Entity) Set(key string, val interface{}) *Entity {
	e.data[key] = val
	return e
}
func (e *Entity) Remove(key string) *Entity {
	delete(e.data, key)
	return e
}
func (e *Entity) QueryOne(m interface{}, sql string, args ...interface{}) error {
	e.initDB()
	rows, err := e.db.Raw(sql).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		err = e.db.ScanRows(rows, &m)
		return err
	}
	return nil
}
func (e *Entity) QueryNumber(sql string, args ...interface{}) (int64, error) {
	e.initDB()
	var number int64
	db := e.db.Raw(sql, args...)
	db.Row().Scan(&number)
	return number, db.Error
}
func (e *Entity) QueryString(sql string, args ...interface{}) (string, error) {
	e.initDB()
	var str string
	db := e.db.Raw(sql, args...)
	db.Row().Scan(&str)
	return str, db.Error
}
func (e *Entity) Insert() (int64, error) {
	if len(e.data) == 0 {
		return 0, errors.New("请先设置字段")
	}
	sql := "INSERT INTO " + e.tableName + " ("
	var vals []interface{}
	for k, v := range e.data {
		sql = sql + "`" + k + "`,"
		vals = append(vals, v)
	}
	sql = strings.TrimRight(sql, ",")
	sql = sql + ") VALUES ("
	for range e.data {
		sql = sql + "?,"
	}
	sql = strings.TrimRight(sql, ",")
	sql = sql + ")"
	e.initDB()
	db := e.db.Exec(sql, vals...)
	return db.RowsAffected, db.Error
}

func (e *Entity) Update(where *Entity) (int64, error) {
	if len(e.data) == 0 {
		return 0, errors.New("请先设置字段")
	}
	if len(where.data) == 0 {
		return 0, errors.New("请设置条件")
	}
	sql := "UPDATE  " + e.tableName + " SET "
	var vals []interface{}
	for k, v := range e.data {
		sql = sql + "`" + k + "`= ? ,"
		vals = append(vals, v)
	}
	sql = strings.TrimRight(sql, ",")
	sql = sql + " where 1=1"
	for k, v := range where.data {
		sql = sql + " AND `" + k + "`= ? "
		vals = append(vals, v)
	}
	//去掉第一个 1=1 and
	sql = strings.Replace(sql, "1=1 AND", "", -1)

	e.initDB()
	db := e.db.Exec(sql, vals...)
	return db.RowsAffected, db.Error
}
func (e *Entity) Delete(where *Entity) (int64, error) {
	if len(where.data) == 0 {
		return 0, errors.New("请设置条件")
	}

	sql := "DELETE FROM " + e.tableName + " WHERE 1=1"
	var vals []interface{}
	for k, v := range where.data {
		sql = sql + " AND `" + k + "`= ? "
		vals = append(vals, v)
	}

	//去掉第一个 1=1 and
	sql = strings.Replace(sql, "1=1 AND", "", -1)
	e.initDB()
	db := e.db.Exec(sql, vals...)
	return db.RowsAffected, db.Error
}

// Exec 直接执行原始SQL语句
func (e *Entity) Exec(sql string, args ...interface{}) (int64, error) {
	e.initDB()
	db := e.db.Exec(sql, args...)
	return db.RowsAffected, db.Error
}

// Parse 将实体转换为Entity的data
// param ignoreNull 是否忽略空值
// param toUnderLine 是否将字段名称转换为下划线方式
func (e *Entity) Parse(obj interface{}, toUnderLine bool) *Entity {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		panic("obj 必须为 指针")
	}
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i).Interface()
		field := t.Field(i).Name
		if toUnderLine {
			field = strkit.ToUnderLine(field)
		}
		e.data[field] = value
	}
	return e
}
