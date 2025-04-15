package models

import (
	"fmt"

	"github.com/tiamxu/kit/sql"
)

// DB 全局数据库连接
var DB *sql.DB

// Init 初始化数据库连接
func Init(cfg *sql.Config) error {
	if cfg == nil {
		return fmt.Errorf("database config is nil")
	}

	db := sql.NewPreDB()
	if err := db.Init(cfg); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	DB = db.DB
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
