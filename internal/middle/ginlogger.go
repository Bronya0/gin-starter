package middle

import (
	"fmt"
	"gin-starter/internal/util/logger"
	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取请求信息
		method := c.Request.Method
		url := c.Request.RequestURI
		statusCode := c.Writer.Status()

		// 记录日志
		if statusCode >= 500 {
			logger.Logger.Error(fmt.Sprintf("【GinLogger】Method: %s, URL: %s, Status: %d", method, url, statusCode))
		} else if statusCode >= 400 {
			logger.Logger.Warn(fmt.Sprintf("【GinLogger】Method: %s, URL: %s, Status: %d", method, url, statusCode))
		} else {
			logger.Logger.Info(fmt.Sprintf("【GinLogger】Method: %s, URL: %s, Status: %d", method, url, statusCode))
		}

		// 继续执行后续的中间件和路由
		c.Next()
	}
}
