package model

import (
	"github.com/jinzhu/gorm"
)

type Base struct {
	db *gorm.DB
}

/**
 * 获取数据库连接
 */
func (b *Base) GetDB() *gorm.DB {
	if b.db == nil {
		b.SetDB(GetDB())
	}
	return b.db
}

/**
 * 设置数据库连接（仅当使用事物需要变更操作句柄时使用）
 */
func (b *Base) SetDB(db *gorm.DB) {
	b.db = db
}
