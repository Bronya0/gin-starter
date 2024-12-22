package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/util/glog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type PgSql struct{}

func (p *PgSql) NewDB() *gorm.DB {
	DbConfig := config.GloConfig.DB

	pgsqlConfig := postgres.Config{
		DSN:                  DbConfig.DSN, // DSN data source name
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{
		SkipDefaultTransaction: true,                                       // 跳过默认事务，提高性能
		PrepareStmt:            true,                                       // 缓存预编译语句
		Logger:                 NewGormLogger(config.GloConfig.Logs.DbLog), //拦截、接管 gorm v2 自带日志
	})

	if err != nil {
		glog.Log.Error(err)
	} else {
		glog.Log.Info("数据库连接成功...")
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(DbConfig.MaxIdletime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(DbConfig.MaxLifetime))
	sqlDB.SetMaxIdleConns(DbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DbConfig.MaxOpenConns)
	return db

}
