package logger

import (
	"gin-starter/internal/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

// InitLogger pathFile: 日志全路径
func InitLogger() {
	// 文件写入器配置
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   global.GloConfig.Logs.Path,
		MaxSize:    global.GloConfig.Logs.MaxSize,
		MaxBackups: global.GloConfig.Logs.MaxBackups,
		MaxAge:     global.GloConfig.Logs.MaxAge,
		Compress:   global.GloConfig.Logs.Compress,
	})

	// 控制台写入器
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建两个编码器，一个用于文件，另一个用于控制台（这里为了演示使用了相同的配置，实际可以根据需求定制）
	fileEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 创建两个核心，分别对应文件和控制台输出。指定日志下限
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, zap.InfoLevel)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zap.InfoLevel)

	// 使用zapcore.NewTee将两个核心合并
	core := zapcore.NewTee(fileCore, consoleCore)

	// 构建logger实例
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)).Sugar() // error及以上的级别增加堆栈;

	// 设置全局logger
	global.Logger = logger
	global.Logger.Info("日志初始化成功...")
}
