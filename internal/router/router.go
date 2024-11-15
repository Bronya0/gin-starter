package router

import (
	"fmt"
	"gin-starter/internal/api"
	"gin-starter/internal/api/v1/auth"
	"gin-starter/internal/config"
	"gin-starter/internal/middle"
	"gin-starter/internal/utils/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"net/http"
	"time"
)

// InitServer 加载配置文件的端口，启动gin服务，同时初始化路由
func InitServer() {

	// ===注册路由===
	router := CreateRouter()

	cfg := config.GloConfig.Server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  300 * time.Second,
	}
	logger.Logger.Info("欢迎主人！服务运行地址：http://", addr)
	//printRegisteredRoutes(router)
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

// CreateRouter 注册通用的路由
func CreateRouter() *gin.Engine {

	// 注册通用路由
	r := CommonRouter()

	// 注册中间件
	InitMiddleware(r)

	// 注册自定义路由
	CustomRouter(r)

	return r
}

func CommonRouter() *gin.Engine {
	var r *gin.Engine

	// 根据配置文件的debug初始化gin路由
	if config.GloConfig.Server.Debug == false {
		//【生产模式】
		// 禁用 gin 记录接口访问日志，
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		r = gin.New()
	} else {
		// 【调试模式】
		// 开启 pprof 包，便于开发阶段分析程序性能
		r = gin.Default()
		gin.ForceConsoleColor()

		// pprof
		//http://localhost:8001/debug/pprof
		//pprof.Register(r)
	}
	return r

}

func InitMiddleware(r *gin.Engine) {

	// 前置通用中间件
	r.Use(
		//middle.GinLogger(),
		middle.CustomRecovery(),
		gzip.Gzip(gzip.DefaultCompression),
		middle.SlowTimeMiddleware(),
	)

	//设置跨域，真正的跨域保护应该在网关层做
	//r.Use(middle.AccessCors())

	// swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 代理静态文件
	//http.Handle("/front/", http.FileServer(http.Dir("front/")))
	//r.LoadHTMLGlob("front/*.tmpl")
	//r.Static("front", "front")

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"love you": time.Now().Format(time.DateTime),
		})
	})

	r.GET("/father", api.TestGorm)

	// 注册
	r.POST("/register", auth.Register)

	// 登录接口
	r.POST("/login", auth.Login)

	// JWT认证中间件
	r.Use(middle.CheckTokenAuth())

	r.GET("/auth", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"title": "欢迎主人访问授权接口",
		})
	})
}
