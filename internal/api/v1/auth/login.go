package auth

import (
	"gin-starter/internal/model/req"
	"gin-starter/internal/model/resp"
	"gin-starter/internal/service/auth"
	"gin-starter/internal/util/glog"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// 用户发送用户名和密码过来
	var login req.LoginReq
	err := c.ShouldBindJSON(&login)
	if err != nil {
		resp.Error(c, "非法参数", nil)
		return
	}
	// 验证用户名和密码
	user, err := auth.CheckUser(login.Username, login.Password)
	if err != nil {
		resp.Error(c, "用户名或密码错误", nil)
		return
	}
	// 黑名单。从redis里检索set，存在则返回错误
	isBlacklisted, err := auth.CheckBlacklist(login.Username)
	if err != nil {
		resp.Error(c, "服务器内部错误", nil)
		return
	}
	if isBlacklisted {
		resp.Error(c, "账户已被禁用", nil)
		return
	}
	// 生成JWT
	tokenString, err := auth.GenToken(user.Id)
	if err != nil {
		glog.Log.Error("生成token失败", err)
		return
	}
	resp.Success(c, "登录成功", gin.H{"token": tokenString})
}
