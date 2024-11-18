package middleware

import (
	"gin-starter/internal/model/resp"
	"gin-starter/internal/svc"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			resp.ErrorAuth(c)
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			resp.ErrorAuth(c)
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它；也会自动校验过期时间
		payload, err := svc.ParseToken(parts[1])
		if err != nil {
			resp.ErrorAuth(c)
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", payload.Username)
		c.Next()
	}
}
