package auth

import (
	"context"
	"errors"
	"gin-starter/internal/global"
	"gin-starter/internal/model/sys"
	"gin-starter/internal/util/gredis"
	"github.com/duke-git/lancet/v2/cryptor"
	"gorm.io/gorm"
)

// checkBlacklist 检查用户是否在黑名单中
func CheckBlacklist(userId string) (bool, error) {
	redisClient := gredis.GetRedisClient()
	// 检查用户是否在黑名单集合中
	exists, err := redisClient.Client.SIsMember(context.Background(), "auth:user_blacklist", userId).Result()
	if err != nil {
		return false, err
	}
	return exists, nil
}

// 用户名密码检查
func CheckUser(username, password string) (*sys.AuthUser, error) {
	// 从数据库中查询用户信息
	pwdHash := cryptor.Sha256(password)
	user := &sys.AuthUser{}
	err := global.DB.Where("username = ? and password = ?", username, pwdHash).First(user).Error
	// 判断是否存在，查不到：检查 ErrRecordNotFound 错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New(global.ErrorsLoginInvalid)
	}
	return nil, nil
}
