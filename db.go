package upms

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db_upms *gorm.DB
)

func InitDb(dbconfig string) {
	db, err := gorm.Open("mysql", dbconfig)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database, config:%s, error:%s", dbconfig, err.Error()))
	}
	SetDb(db)
}

func SetDb(db *gorm.DB) {
	db_upms = db
}

func SetDebug() {
	db_upms = db_upms.Debug()
}
