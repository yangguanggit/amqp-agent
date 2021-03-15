package model

import (
	"amqp-agent/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

/**
 * 获取数据库连接
 */
func GetDB() *gorm.DB {
	return db
}

/**
 * 连接数据库
 */
func ConnectDB(dbConfig config.Database) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database, dbConfig.Charset)
	db, err = gorm.Open(dbConfig.Dialect, dsn)
	if err != nil {
		panic(fmt.Errorf("fatal error: connect database: %s\n", err))
	}

	db.DB().SetMaxIdleConns(dbConfig.MaxIdleConnNum)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConnNum)
	db.BlockGlobalUpdate(true)
	db.InstantSet("gorm:save_associations", false)
	db.InstantSet("gorm:association_save_reference", false)
}

/**
 * 关闭数据库连接
 */
func CloseDB() {
	if db != nil {
		_ = db.Close()
	}
}
