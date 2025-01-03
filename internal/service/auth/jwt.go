package auth

import (
	"errors"
	"gin-starter/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// CustomClaims 自定义声明类型 并内嵌jwt.RegisteredClaims
// jwt包自带的jwt.RegisteredClaims只包含了官方字段
// 假设我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	// 可根据需要自行添加字段
	UserID               string `json:"user_id"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte(config.GloConfig.Jwt.JwtTokenSignKey)

// GenToken 生成JWT
func GenToken(userID string) (string, error) {
	// TokenExpireDuration jwt token 的过期时间
	tokenExpire, err := time.ParseDuration(config.GloConfig.Jwt.ExpiresTime)
	if err != nil {
		return "", err
	}
	// 创建一个我们自己的声明
	claims := CustomClaims{
		userID, // 自定义的用户名字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpire)),
			Issuer:    "gin-starter", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
