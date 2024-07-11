package jwt

import (
	"errors"
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	"gin-starter/internal/service/auth"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// CreateUserFactory 创建 userToken 工厂
func CreateUserFactory() *userToken {
	return &userToken{
		userJwt: CreateMyJWT(config.GloConfig.Jwt.JwtTokenSignKey),
	}
}

// 自定义jwt的声明字段信息+标准字段，参考地址：https://blog.csdn.net/codeSquare/article/details/99288718
type CustomClaims struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"user_name"`
	Phone  string `json:"phone"`
	jwt.StandardClaims
}

type userToken struct {
	userJwt *JwtSign
}

// GenerateToken 生成token
func (u *userToken) GenerateToken(userid int64, username string, phone string, expireAt int64) (tokens string, err error) {

	// 根据实际业务自定义token需要包含的参数，生成token，注意：用户密码请勿包含在token
	customClaims := CustomClaims{
		UserId: userid,
		Name:   username,
		Phone:  phone,
		// 特别注意，针对前文的匿名结构体，初始化的时候必须指定键名，并且不带 jwt. 否则报错：Mixture of field: value and value initializers
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10,       // 生效开始时间
			ExpiresAt: time.Now().Unix() + expireAt, // 失效截止时间
		},
	}
	return u.userJwt.CreateToken(customClaims)
}

// RecordLoginToken 用户login成功，记录用户token
func (u *userToken) RecordLoginToken(userToken, clientIp string) bool {
	if customClaims, err := u.userJwt.ParseToken(userToken); err == nil {
		userId := customClaims.UserId
		expiresAt := customClaims.ExpiresAt
		return auth.CreateUserFactory("").OauthLoginToken(userId, userToken, expiresAt, clientIp)
	} else {
		return false
	}
}

// TokenIsMeetRefreshCondition 检查token是否满足刷新条件
func (u *userToken) TokenIsMeetRefreshCondition(token string) bool {
	// token基本信息是否有效：1.过期时间在允许的过期范围内;2.基本格式正确
	customClaims, code := u.isNotExpired(token, config.GloConfig.Jwt.JwtTokenRefreshAllowSec)
	switch code {
	case global.JwtTokenOK, global.JwtTokenExpired:
		//在数据库的存储信息是否也符合过期刷新刷新条件
		if auth.CreateUserFactory("").OauthRefreshConditionCheck(customClaims.UserId, token) {
			return true
		}
	}
	return false
}

// RefreshToken 刷新token的有效期（默认+3600秒，参见常量配置项）
func (u *userToken) RefreshToken(oldToken, clientIp string) (newToken string, res bool) {
	var err error
	//如果token是有效的、或者在过期时间内，那么执行更新，换取新token
	if newToken, err = u.userJwt.RefreshToken(oldToken, config.GloConfig.Jwt.JwtTokenRefreshAllowSec); err == nil {
		if customClaims, err := u.userJwt.ParseToken(newToken); err == nil {
			userId := customClaims.UserId
			expiresAt := customClaims.ExpiresAt
			if auth.CreateUserFactory("").OauthRefreshToken(userId, expiresAt, oldToken, newToken, clientIp) {
				return newToken, true
			}
		}
	}

	return "", false
}

// 判断token本身是否未过期
// 参数解释：
// token： 待处理的token值
// expireAtSec： 过期时间延长的秒数，主要用于用户刷新token时，判断是否在延长的时间范围内，非刷新逻辑默认为0
func (u *userToken) isNotExpired(token string, expireAtSec int64) (*CustomClaims, int) {
	if customClaims, err := u.userJwt.ParseToken(token); err == nil {

		if time.Now().Unix()-(customClaims.ExpiresAt+expireAtSec) < 0 {
			// token有效
			return customClaims, global.JwtTokenOK
		} else {
			// 过期的token
			return customClaims, global.JwtTokenExpired
		}
	} else {
		// 无效的token
		return nil, global.JwtTokenInvalid
	}
}

// IsEffective 判断token是否有效（未过期+数据库用户信息正常）
func (u *userToken) IsEffective(token string) bool {
	customClaims, code := u.isNotExpired(token, 0)
	if global.JwtTokenOK == code {
		//1.首先在redis检测是否存在某个用户对应的有效token，如果存在就直接返回，不再继续查询mysql，否则最后查询mysql逻辑，确保万无一失
		//if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
		//	tokenRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(customClaims.UserId)
		//	if tokenRedisFact != nil {
		//		defer tokenRedisFact.ReleaseRedisConn()
		//		if tokenRedisFact.TokenCacheIsExists(token) {
		//			return true
		//		}
		//	}
		//}
		//2.token符合token本身的规则以后，继续在数据库校验是不是符合本系统其他设置，例如：一个用户默认只允许10个账号同时在线（10个token同时有效）
		if auth.CreateUserFactory("").OauthCheckTokenIsOk(customClaims.UserId, token) {
			return true
		}
	}
	return false
}

// ParseToken 将 token 解析为绑定时传递的参数
func (u *userToken) ParseToken(tokenStr string) (CustomClaims CustomClaims, err error) {
	if customClaims, err := u.userJwt.ParseToken(tokenStr); err == nil {
		return *customClaims, nil
	} else {
		return CustomClaims, errors.New(global.ErrorsParseTokenFail)
	}
}

// DestroyToken 销毁token，基本用不到，因为一个网站的用户退出都是直接关闭浏览器窗口，极少有户会点击“注销、退出”等按钮，销毁token其实无多大意义
func (u *userToken) DestroyToken() {

}
