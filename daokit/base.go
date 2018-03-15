package daokit

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
	"github.com/weibaohui/go-kit/propkit"
	"github.com/weibaohui/go-kit/uikit"
)

var err error

var globalDB *gorm.DB
var once sync.Once

// GetOrm :set orm  singleton
func GetOrm() *gorm.DB {
	once.Do(func() {
		user := propkit.Init().Get("db.user")
		password := propkit.Init().Get("db.password")
		host := propkit.Init().Get("db.host")
		port := propkit.Init().GetInt64("db.port")
		name := propkit.Init().Get("db.name")
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			user, password, host, port, name)
		fmt.Println(url)
		globalDB, err = gorm.Open("mysql", url)
		if err != nil {
			log.Fatalf("connect to db err: %s", err.Error())
		}
		globalDB.DB().SetMaxOpenConns(15)
		globalDB.DB().SetMaxIdleConns(15)
		if propkit.IsDevMod() {
			globalDB.LogMode(true)
		} else {
			globalDB.LogMode(false)
		}
		fmt.Println("db 初始化")
	})
	return globalDB
}
func SetPageLimitAndCount(sql *gorm.DB, page ...*uikit.Pagination) *gorm.DB {
	if len(page) == 1 {
		var total int
		sql.Count(&total)
		page[0].SetTotal(total)
		sql = sql.Limit(page[0].PageSize).Offset(page[0].PageSize * (page[0].PageIndex - 1))
	}
	return sql
}
