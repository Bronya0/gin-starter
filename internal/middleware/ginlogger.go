package middleware

import (
	"fmt"
	"gin-starter/internal/util/glog"
	"github.com/gin-gonic/gin"
)

func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求信息
		method := c.Request.Method
		url := c.Request.RequestURI
		statusCode := c.Writer.Status()

		// 记录日志
		if statusCode >= 500 {
			glog.Log.Error(fmt.Sprintf("【ErrorLogger】Method: %s, URL: %s, Status: %d", method, url, statusCode))
		}

		// 继续执行后续的中间件和路由
		c.Next()
	}
}
