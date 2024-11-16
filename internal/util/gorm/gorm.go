package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	logging "gin-starter/internal/util/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitDB() {
	global.DB = InitGorm(config.GloConfig.DB.Type)
	logging.Log.Info("数据库连接成功...")
}

func NewGormLogger(logFile string) logger.Interface {

	var writer logger.Writer
	if config.GloConfig.Server.Debug {
		writer = log.New(os.Stdout, "\r\n", log.LstdFlags)
	} else {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("无法创建gorm日志: %v", err)
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

func InitGorm(DbType string) *gorm.DB {
	switch DbType {
	case "mysql":
		return Mysql()
	case "pgsql":
		return PgSql()
	default:
		return PgSql()
	}
}
func Mysql() *gorm.DB {
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
		logging.Log.Error(err)
	}
	db.InstanceSet("gorm:table_options", "ENGINE=innodb")

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(DbConfig.MaxIdletime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(DbConfig.MaxLifetime))
	sqlDB.SetMaxIdleConns(DbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DbConfig.MaxOpenConns)
	return db
}

func PgSql() *gorm.DB {
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
		logging.Log.Error(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(DbConfig.MaxIdletime))
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(DbConfig.MaxLifetime))
	sqlDB.SetMaxIdleConns(DbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DbConfig.MaxOpenConns)
	return db

}
