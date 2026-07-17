package data

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化 MySQL 连接
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 打印所有的 SQL 语句，方便调试
	})

	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 配置底层连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层的 sql.DB 失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)           // 空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接可复用的最大时间

	log.Println("MySQL 连接初始化成功！")
	return db
}
