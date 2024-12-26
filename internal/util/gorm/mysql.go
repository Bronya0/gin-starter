package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/util/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		panic(err)
	} else {
		glog.Log.Info("数据库连接成功...")
	}
	db.InstanceSet("gorm:table_options", "ENGINE=innodb")

	return db
}
