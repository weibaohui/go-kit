package daokit

import "strings"

type Entity struct {
	tableName string
	data      map[string]interface{}
}

func M(name string) *Entity {
	entity := &Entity{
		data: make(map[string]interface{}),
	}
	entity.tableName = name
	return entity
}
func (e *Entity) Set(key string, val interface{}) *Entity {
	e.data[key] = val
	return e
}
func (e *Entity) QueryNumber(sql string, args ...interface{}) (int64, error) {
	var number int64
	db := GetOrm().Raw(sql, args...).Scan(&number)
	return number, db.Error
}
func (e *Entity) QueryString(sql string, args ...interface{}) (string, error) {
	var str string
	db := GetOrm().Raw(sql, args...).Scan(&str)
	return str, db.Error
}
func (e *Entity) Insert() (int64, error) {
	sql := "insert into " + e.tableName + "("
	var vals []interface{}
	for k, v := range e.data {
		sql = sql + "`" + k + "`,"
		vals = append(vals, v)
	}
	sql = strings.TrimRight(sql, ",")
	sql = sql + ") values ("
	for range e.data {
		sql = sql + "?,"
	}
	sql = strings.TrimRight(sql, ",")
	sql = sql + ")"
	db := GetOrm().Exec(sql, vals...)
	return db.RowsAffected, db.Error
}
