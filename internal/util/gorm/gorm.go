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

type GormDB interface {
	NewDB() *gorm.DB
}

func InitDB() {
	if config.GloConfig.DB.Enable {
		global.DB = InitGorm(config.GloConfig).NewDB()
	} else {
		glog.Log.Warn("数据库未启用...")
	}
}

func InitGorm(gloConfig *config.Config) GormDB {
	db := &gloConfig.DB
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务，提高性能
		PrepareStmt:            true, // 缓存预编译语句
		Logger:                 NewGormLogger(gloConfig),
	}
	switch gloConfig.DB.Type {
	case "mysql":
		return &Mysql{db, gormConfig}
	case "pgsql":
		return &PgSql{db, gormConfig}
	default:
		return &PgSql{db, gormConfig}
	}
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
