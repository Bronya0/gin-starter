package router

import (
	"gin-starter/internal/api"
	"gin-starter/internal/api/v1/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// addPublicRouter  公开的路由
func addPublicRouter(r *gin.Engine) *gin.Engine {

	// 设置跨域，真正的跨域保护应该在网关层做
	// r.Use(middle.AccessCors())

	// 代理静态文件
	// http.Handle("/front/", http.FileServer(http.Dir("front/")))
	// r.LoadHTMLGlob("front/*.tmpl")
	// r.Static("front", "front")

	publicApi := r.Group("/api/public")
	{
		publicApi.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"love you": time.Now().Format(time.DateTime),
			})
		})

		publicApi.GET("/father", api.TestGorm)

		// 登录获取JWT
		publicApi.POST("/jwt-login", auth.JwtLogin)
		// 注册
		// r.POST("/register", auth.Register)

	}

	return r
}
