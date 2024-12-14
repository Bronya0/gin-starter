package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// addAuthRouter 需要认证的路由
func addAuthRouter(r *gin.Engine) *gin.Engine {
	authApi := r.Group("/api/v1")
	{
		authApi.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"love you": time.Now().Format(time.DateTime),
			})
		})
	}

	return r
}
