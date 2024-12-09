package router

import (
	"gin-starter/internal/api"
	"gin-starter/internal/api/v1/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func addAccessRouter(r *gin.Engine) *gin.Engine {

	//设置跨域，真正的跨域保护应该在网关层做
	//r.Use(middle.AccessCors())

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

	// 登录获取JWT
	r.POST("/login", auth.Login)

	// 注册
	//r.POST("/register", auth.Register)

	return r
}

func addRouter(r *gin.Engine) *gin.Engine {

	r.GET("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"love you": time.Now().Format(time.DateTime),
		})
	})

	return r
}
