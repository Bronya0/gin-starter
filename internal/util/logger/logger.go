package logger

import (
	"gin-starter/internal/config"
	logging "github.com/donnie4w/go-logger/logger"
)

var (
	Logger = InitLogger(config.GloConfig.Logs.Path)
)

// InitLogger pathFile: 日志全路径
func InitLogger(path string) *logging.Logging {

	logger := logging.NewLogger()
	logger.SetOption(&logging.Option{
		Level:     logging.LEVEL_INFO,
		Console:   true, // 控制台输出
		Format:    logging.FORMAT_LEVELFLAG | logging.FORMAT_SHORTFILENAME | logging.FORMAT_DATE | logging.FORMAT_MICROSECNDS,
		Formatter: "{level} [{time}] {file}: {message}\n",
		// size或者time模式
		FileOption: &logging.FileTimeMode{ // 这里用时间切割
			Filename:   path,             // 日志文件路径
			Timemode:   logging.MODE_DAY, // 按天
			Maxbuckup:  180,              // 最多备份日志文件数
			IsCompress: false,            // 是否压缩
		},
	})

	return logger
}
