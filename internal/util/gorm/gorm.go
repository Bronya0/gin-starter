package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	"gin-starter/internal/util/glog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type GormDB interface {
	NewDB() *gorm.DB
}

func InitDB() {
	if config.GloConfig.DB.Enable {
		global.DB = InitGorm(&config.GloConfig.DB).NewDB()
	} else {
		glog.Log.Warn("数据库未启用...")
	}
}

func InitGorm(db *config.DB) GormDB {

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务，提高性能
		PrepareStmt:            true, // 缓存预编译语句
		Logger:                 NewGormLogger(config.GloConfig.Logs.DbLog),
	}
	switch db.Type {
	case "mysql":
		return &Mysql{db, gormConfig}
	case "pgsql":
		return &PgSql{db, gormConfig}
	default:
		return &PgSql{db, gormConfig}
	}
}

func NewGormLogger(logFile string) logger.Interface {

	var writer logger.Writer
	if config.GloConfig.Server.Debug {
		writer = log.New(os.Stdout, "\r\n", log.LstdFlags)
	} else {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			glog.Log.Fatalf("无法创建gorm日志: %v", err)
		}
		defer file.Close() // 文件将在函数退出时自动关闭
		writer = log.New(file, "\r\n", log.LstdFlags)
	}

	// 使用os.File作为io.Writer
	newLogger := logger.New(
		writer,
		logger.Config{
			SlowThreshold:             3 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn,     // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,            // Don't include params in the SQL log
			Colorful:                  false,           // Disable color
		},
	)
	return newLogger
}
