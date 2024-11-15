package router

import (
	"gin-starter/internal/api"
	"gin-starter/internal/api/v1/auth"
	"gin-starter/internal/middle"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

func addBaseRouter(r *gin.Engine) *gin.Engine {

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
	// 测试demo路由

	return r
}
