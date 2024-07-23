package gorm

import (
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	logging "gin-starter/internal/utils/logger"
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
	logging.Logger.Info("数据库连接成功...")
}

func NewGormLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             3 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,   // Log level
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
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	default:
		return GormPgSql()
	}
}
func GormMysql() *gorm.DB {
	DbConfig := config.GloConfig.DB

	mysqlConfig := mysql.Config{
		DSN:                       DbConfig.DSN, // DSN data source name
		DefaultStringSize:         191,          // string 类型字段的默认长度
		SkipInitializeWithVersion: false,        // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务，提高性能
		PrepareStmt:            true, // 缓存预编译语句
		Logger:                 NewGormLogger(),
	})
	if err != nil {
		panic(err)
	}
	db.InstanceSet("gorm:table_options", "ENGINE=innodb")
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(DbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DbConfig.MaxOpenConns)
	return db
}

func GormPgSql() *gorm.DB {
	DbConfig := config.GloConfig.DB

	pgsqlConfig := postgres.Config{
		DSN:                  DbConfig.DSN, // DSN data source name
		PreferSimpleProtocol: false,
	}
	db, err := gorm.Open(postgres.New(pgsqlConfig), &gorm.Config{
		SkipDefaultTransaction: true,            // 跳过默认事务，提高性能
		PrepareStmt:            true,            // 缓存预编译语句
		Logger:                 NewGormLogger(), //拦截、接管 gorm v2 自带日志
	})

	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(time.Second * 30)
	sqlDB.SetConnMaxLifetime(60 * time.Second)
	sqlDB.SetMaxIdleConns(DbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(DbConfig.MaxOpenConns)
	return db

}
