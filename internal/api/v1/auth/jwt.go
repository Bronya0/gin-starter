package auth

import (
	"gin-starter/internal/model/auth"
	"gin-starter/internal/model/resp"
	"gin-starter/internal/service"
	"gin-starter/internal/util/glog"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// 用户发送用户名和密码过来
	var login auth.Login
	err := c.ShouldBindJSON(&login)
	if err != nil {
		resp.Error(c, "非法参数", nil)
		return
	}
	// 校验用户名和密码是否正确
	// 密码可能是哈希值
	if login.Username == "admin" && login.Password == "admin123" {
		// 生成Token
		tokenString, err := service.GenToken(login.Username)
		if err != nil {
			glog.Log.Error("生成token失败", err)
			return
		}
		resp.Success(c, "登录成功", gin.H{"token": tokenString})
		return
	}
	// 黑名单
	// ...

	resp.ErrorAuth(c)
	return
}
