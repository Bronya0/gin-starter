package router

import (
	"fmt"
	"gin-starter/internal/config"
	"gin-starter/internal/middle"
	"gin-starter/internal/util/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

// InitServer 加载配置文件的端口，启动gin服务，同时初始化路由
func InitServer() {

	engine := CreateEngine()

	cfg := config.GloConfig.Server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  300 * time.Second,
	}
	logger.Logger.Info("欢迎主人！服务运行地址：http://", addr)
	printRegisteredRoutes(engine)
	logger.Logger.Error(srv.ListenAndServe().Error())

}

// printRegisteredRoutes 打印注册的路由信息
func printRegisteredRoutes(r *gin.Engine) {
	// 遍历注册的路由
	for _, route := range r.Routes() {
		// 输出路由信息
		fmt.Printf("%s %s, ", route.Method, route.Path)
	}
	logger.Logger.Info("")
}

// CreateEngine 注册通用的路由
func CreateEngine() *gin.Engine {
	engine := Engine()
	// 中间件
	addMiddleware(engine)
	// 通用路由
	addBaseRouter(engine)
	// 自定义路由...

	return engine
}

func Engine() *gin.Engine {
	var engine *gin.Engine

	// 根据配置文件的debug初始化gin路由
	if config.GloConfig.Server.Debug == false {
		//【生产模式】
		// 禁用 gin 记录接口访问日志，
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		engine = gin.New()
	} else {
		// 【调试模式】
		// 开启 pprof 包，便于开发阶段分析程序性能
		engine = gin.Default()
		gin.ForceConsoleColor()

		// pprof
		//http://localhost:8001/debug/pprof
		//pprof.Register(r)
	}
	return engine

}

func addMiddleware(r *gin.Engine) {

	// 前置通用中间件
	r.Use(
		//middle.GinLogger(),
		middle.CustomRecovery(),
		gzip.Gzip(gzip.DefaultCompression),
		middle.SlowTimeMiddleware(),
	)

}
