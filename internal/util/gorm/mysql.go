package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/util/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Mysql struct{}

func (m *Mysql) NewDB() *gorm.DB {
	DbConfig := config.GloConfig.DB

	mysqlConfig := mysql.Config{
		DSN:                       DbConfig.DSN, // DSN data source name
		DefaultStringSize:         191,          // string 类型字段的默认长度
		SkipInitializeWithVersion: false,        // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务，提高性能
		PrepareStmt:            true, // 缓存预编译语句
		Logger:                 NewGormLogger(config.GloConfig.Logs.DbLog),
	})
	if err != nil {
		glog.Log.Error(err)
	} else {
		glog.Log.Info("数据库连接成功...")
	}
	db.InstanceSet("gorm:table_options", "ENGINE=innodb")

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(DbConfig.MaxIdletime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(DbConfig.MaxLifetime))
	sqlDB.SetMaxIdleConns(DbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DbConfig.MaxOpenConns)
	return db
}
