package auth

import (
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	"gin-starter/internal/model/response"
	"gin-starter/internal/service/auth"
	"gin-starter/internal/service/auth/jwt"
	"gin-starter/internal/utils/hash"
	"github.com/gin-gonic/gin"
	"time"
)

type Users struct {
}

type UsersCurd struct {
	userModel *auth.UsersModel
}

var usersCrudObj = UsersCurd{}

func (u *UsersCurd) Register(userName, pass, userIp string) bool {
	pass = hash.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return u.userModel.Register(userName, pass, userIp)
}

func (u *UsersCurd) Store(name string, pass string, realName string, phone string, remark string) bool {

	pass = hash.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return u.userModel.Store(name, pass, realName, phone, remark)
}

func (u *UsersCurd) Update(id int, name string, pass string, realName string, phone string, remark string, clientIp string) bool {
	//预先处理密码加密等操作，然后进行更新
	pass = hash.Base64Md5(pass) // 预先处理密码加密，然后存储在数据库
	return u.userModel.Update(id, name, pass, realName, phone, remark, clientIp)
}

// Register 1.用户注册
func Register(c *gin.Context) {
	// 当然也可以通过gin框架的上下文原始方法获取，例如： c.PostForm("user_name") 获取，这样获取的数据格式为文本，需要自己继续转换
	userName := c.PostForm("user_name")
	pass := c.PostForm("pass")
	userIp := c.ClientIP()
	if usersCrudObj.Register(userName, pass, userIp) {
		response.SuccessResponse(c, global.CurdStatusOkMsg, "")
	} else {
		response.ErrorResponse(c, global.CurdRegisterFailMsg, "")
	}
}

// Login 2.用户登录
func Login(c *gin.Context) {
	userName := c.PostForm("user_name")
	pass := c.PostForm("pass")
	phone := c.PostForm("phone")
	userModelFact := auth.CreateUserFactory("")
	userModel := userModelFact.Login(userName, pass)

	if userModel != nil {
		userTokenFactory := jwt.CreateUserFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.Id, userModel.UserName, userModel.Phone, config.GloConfig.Jwt.JwtTokenCreatedExpireAt); err == nil {
			if userTokenFactory.RecordLoginToken(userToken, c.ClientIP()) {
				data := gin.H{
					"userId":     userModel.Id,
					"user_name":  userName,
					"realName":   userModel.RealName,
					"phone":      phone,
					"token":      userToken,
					"updated_at": time.Now().Format(time.DateTime),
				}
				response.SuccessResponse(c, global.CurdStatusOkMsg, data)
				go userModel.UpdateUserloginInfo(c.ClientIP(), userModel.Id)
				return
			}
		}
	}
	response.ErrorResponse(c, global.CurdLoginFailMsg, "")
}

// RefreshToken 刷新用户token
func (u *Users) RefreshToken(c *gin.Context) {
	oldToken := c.GetString(global.ValidatorPrefix + "token")
	if newToken, ok := jwt.CreateUserFactory().RefreshToken(oldToken, c.ClientIP()); ok {
		res := gin.H{
			"token": newToken,
		}
		response.SuccessResponse(c, global.CurdStatusOkMsg, res)
	} else {
		response.ErrorResponse(c, global.CurdRefreshTokenFailMsg, "")
	}
}

// 后面是 curd 部分，自带版本中为了降低初学者学习难度，使用了最简单的方式操作 增、删、改、查
// 在开发企业实际项目中，建议使用我们提供的一整套 curd 快速操作模式
// 参考地址：https://gitee.com/daitougege/GinSkeleton/blob/master/docs/concise.md
// 您也可以参考 Admin 项目地址：https://gitee.com/daitougege/gin-skeleton-admin-backend/ 中， app/model/  提供的示例语法

// Show 3.用户查询（show）
func Show(c *gin.Context) {
	userName := c.GetString(global.ValidatorPrefix + "user_name")
	page := c.GetFloat64(global.ValidatorPrefix + "page")
	limit := c.GetFloat64(global.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit
	counts, showlist := auth.CreateUserFactory("").Show(userName, int(limitStart), int(limit))
	if counts > 0 && showlist != nil {
		response.SuccessResponse(c, global.CurdStatusOkMsg, gin.H{"counts": counts, "list": showlist})
	} else {
		response.ErrorResponse(c, global.CurdSelectFailMsg, "")
	}
}

// Store 4.用户新增(store)
func Store(c *gin.Context) {
	userName := c.GetString(global.ValidatorPrefix + "user_name")
	pass := c.GetString(global.ValidatorPrefix + "pass")
	realName := c.GetString(global.ValidatorPrefix + "real_name")
	phone := c.GetString(global.ValidatorPrefix + "phone")
	remark := c.GetString(global.ValidatorPrefix + "remark")

	if usersCrudObj.Store(userName, pass, realName, phone, remark) {
		response.SuccessResponse(c, global.CurdStatusOkMsg, "")
	} else {
		response.ErrorResponse(c, global.CurdCreatFailMsg, "")
	}
}

// Update 5.用户更新(update)
func Update(c *gin.Context) {
	//表单参数验证中的int、int16、int32 、int64、float32、float64等数字键（字段），请统一使用 GetFloat64() 获取，其他函数无效
	userId := c.GetFloat64(global.ValidatorPrefix + "id")
	userName := c.GetString(global.ValidatorPrefix + "user_name")
	pass := c.GetString(global.ValidatorPrefix + "pass")
	realName := c.GetString(global.ValidatorPrefix + "real_name")
	phone := c.GetString(global.ValidatorPrefix + "phone")
	remark := c.GetString(global.ValidatorPrefix + "remark")
	userIp := c.ClientIP()

	// 检查正在修改的用户名是否被其他人使用
	if auth.CreateUserFactory("").UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
		response.ErrorResponse(c, global.CurdUpdateFailMsg+", "+userName+" 已经被其他人使用", "")
		return
	}

	//注意：这里没有实现更加精细的权限控制逻辑，例如：超级管理管理员可以更新全部用户数据，普通用户只能修改自己的数据。目前只是验证了token有效、合法之后就可以进行后续操作
	// 实际使用请根据真是业务实现权限控制逻辑、再进行数据库操作
	if usersCrudObj.Update(int(userId), userName, pass, realName, phone, remark, userIp) {
		response.SuccessResponse(c, global.CurdStatusOkMsg, "")
	} else {
		response.ErrorResponse(c, global.CurdUpdateFailMsg, "")
	}

}

// Destroy 6.删除记录
func Destroy(c *gin.Context) {
	//表单参数验证中的int、int16、int32 、int64、float32、float64等数字键（字段），请统一使用 GetFloat64() 获取，其他函数无效
	userId := c.GetFloat64(global.ValidatorPrefix + "id")
	if auth.CreateUserFactory("").Destroy(int(userId)) {
		response.SuccessResponse(c, global.CurdStatusOkMsg, "")
	} else {
		response.ErrorResponse(c, global.CurdDeleteFailMsg, "")
	}
}
