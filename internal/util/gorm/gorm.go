package gorm

import (
	"fmt"
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	"gin-starter/internal/util/glog"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"strings"
	"time"
)

type IGormDB interface {
	NewDB() *gorm.DB
}

func InitDB() {
	if config.GloConfig.DB.Enable {
		global.DB = InitGorm(config.GloConfig)
	} else {
		glog.Log.Warn("数据库未启用...")
	}
}

func InitGorm(gloConfig *config.Config) *gorm.DB {
	dbConf := &gloConfig.DB
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务，提高性能
		PrepareStmt:            true, // 缓存预编译语句
		Logger:                 NewGormLogger(gloConfig),
	}

	var db IGormDB
	switch gloConfig.DB.Type {
	case "mysql":
		db = &Mysql{dbConf, gormConfig}
	case "pgsql":
		db = &PgSql{dbConf, gormConfig}
	default:
		db = &PgSql{dbConf, gormConfig}
	}

	instance := db.NewDB()
	sqlDB, _ := instance.DB()

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	// 应小于数据库服务器、负载均衡器等设置的超时时间
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(dbConf.MaxIdletime))
	// SetConnMaxIdleTime 设置空闲连接最大存活时间
	// 建议设置以避免空闲连接占用资源
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(dbConf.MaxLifetime))
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	// 根据服务器性能和内存情况来设置，通常设置为 MaxOpenConns 的 25%-50%
	sqlDB.SetMaxIdleConns(dbConf.MaxIdleConns)
	// SetMaxOpenConns 设置打开数据库连接的最大数量
	// 需要考虑数据库服务器的性能和配置，建议不超过数据库支持的最大连接数的75%
	sqlDB.SetMaxOpenConns(dbConf.MaxOpenConns)
	return instance
}

func NewGormLogger(gloConfig *config.Config) gormLog.Interface {
	newLogger := gormLog.New(
		&CustomWriter{},
		gormLog.Config{
			SlowThreshold:             time.Duration(gloConfig.DB.SlowThreshold) * time.Second, // Slow SQL threshold
			LogLevel:                  gormLog.Warn,                                            // Log level
			IgnoreRecordNotFoundError: true,                                                    // 记录日志时会忽略ErrRecordNotFound错误
			ParameterizedQueries:      true,                                                    // 不会在SQL日志中记录参数值，这有助于保护敏感信息不被记录在日志中。
			Colorful:                  true,                                                    // 如果设置为true，日志将以彩色显示，这有助于在终端中快速区分不同级别的日志。
		},
	)

	return newLogger
}

type CustomWriter struct{}

func (l CustomWriter) Printf(strFormat string, args ...interface{}) {
	logRes := fmt.Sprintf(strFormat, args...)
	logFlag := "gorm日志:"
	if strings.HasPrefix(strFormat, "[info]") || strings.HasPrefix(strFormat, "[traceStr]") {
		glog.Log.Info(logRes)
	} else if strings.HasPrefix(strFormat, "[error]") || strings.HasPrefix(strFormat, "[traceErr]") {
		glog.Log.Error(logFlag, logRes)
	} else if strings.HasPrefix(strFormat, "[warn]") || strings.HasPrefix(strFormat, "[traceWarn]") {
		glog.Log.Warn(logFlag, logRes)
	} else {
		fmt.Println(111, strFormat, args)
		glog.Log.Info(logFlag, logRes)
	}

}
