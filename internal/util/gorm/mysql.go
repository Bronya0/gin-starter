package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/util/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Mysql struct {
	*config.DB
	GormConfig *gorm.Config
}

func (m *Mysql) NewDB() *gorm.DB {

	mysqlConfig := mysql.Config{
		DSN:                       m.DSN, // DSN data source name
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), m.GormConfig)
	if err != nil {
		glog.Log.Error(err)
	} else {
		glog.Log.Info("数据库连接成功...")
	}
	db.InstanceSet("gorm:table_options", "ENGINE=innodb")

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdletime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifetime))
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db
}
