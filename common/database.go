package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 首先下载gorm，然后执行go get -u gorm.io/gorm
// 在下载数据库驱动，然后执行go get -u gorm.io/driver/sqlite  或者 go get -u gorm.io/driver/mysql 或者 go get -u gorm.io/driver/postgre 或者 go get -u gorm.io/driver/sqlserver

import (
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *gorm.DB
)

func DBInit() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("connect is db failed: " + err.Error())
	}
	Db = db
	return
}
