package logger

import (
	"gin-starter/internal/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

// InitLogger pathFile: 日志全路径
func InitLogger() {

	//写入器
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   global.GloConfig.Logs.Path,       //日志文件的位置 /xxx.log
		MaxSize:    global.GloConfig.Logs.MaxSize,    //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: global.GloConfig.Logs.MaxBackups, //保留旧文件的最大个数
		MaxAge:     global.GloConfig.Logs.MaxAge,     //保留旧文件的最大天数
		Compress:   global.GloConfig.Logs.Compress,   //是否压缩/归档旧文件
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000")) // 时间格式
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	var encoder = zapcore.NewConsoleEncoder(encoderConfig) // 普通模式，还有json模式
	writer = zapcore.AddSync(writer)
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)                              // 日志等级下限
	_logger := zap.New(zapCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)).Sugar() // error及以上的级别增加堆栈; sugar允许使用f方法
	global.Logger = _logger
	global.Logger.Info("日志初始化成功...")
}
