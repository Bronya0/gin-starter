package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func addRouter(r *gin.Engine) *gin.Engine {

	r.GET("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"love you": time.Now().Format(time.DateTime),
		})
	})

	return r
}
