package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/util/glog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgSql struct {
	*config.DB
	GormConfig *gorm.Config
}

func (p *PgSql) NewDB() *gorm.DB {

	pgsqlConfig := postgres.Config{
		DSN:                  p.DSN, // DSN data source name
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), p.GormConfig)

	if err != nil {
		glog.Log.Error(err)
		panic(err)
	} else {
		glog.Log.Info("数据库连接成功...")
	}

	return db

}
