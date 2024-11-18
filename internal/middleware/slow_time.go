package middleware

import (
	"fmt"
	"gin-starter/internal/util/glog"
	"github.com/gin-gonic/gin"
	"time"
)

// SlowTimeMiddleware 检查接口响应耗时，将慢接口的具体信息打印到日志里
func SlowTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		if latency.Seconds() > 1 { // 设置阈值，超过1秒则认为是慢接口
			glog.Log.Warn(fmt.Sprintf("【慢接口】%v %v %v", c.Request.Method, c.Request.URL, latency))
		}
	}
}
