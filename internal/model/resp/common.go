package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  msg,
		"data": data,
	})
	c.Abort()
}
func Error(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 5000,
		"msg":  msg,
		"data": data,
	})
	c.Abort()
}
func ErrorAuth(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": 4000,
		"msg":  "认证不通过",
		"data": nil,
	})
	c.Abort()
}
