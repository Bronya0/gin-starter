package auth

import (
	"fmt"
	"gin-starter/internal/config"
	"gin-starter/internal/global"
	"gin-starter/internal/utils"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// 针对数据库选型为 postgresql 的开发者使用

// CreateUserFactory 创建 userFactory
// 参数说明： 传递空值，默认使用 配置文件选项：UseDbType（mysql）
func CreateUserFactory(sqlType string) *UsersModel {
	return &UsersModel{BaseModel: BaseModel{DB: UseDbConn()}}
}

type UsersModel struct {
	BaseModel
	UserName    string `gorm:"column:user_name" json:"user_name"`
	Pass        string `json:"-"`
	Phone       string `json:"phone"`
	RealName    string `gorm:"column:real_name" json:"real_name"`
	Status      int    `json:"status"`
	Token       string `json:"token"`
	LastLoginIp string `gorm:"column:last_login_ip" json:"last_login_ip"`
}

// TableName 表名
func (u *UsersModel) TableName() string {
	return "web.tb_users"
}

// Register 用户注册（写一个最简单的使用账号、密码注册即可）
func (u *UsersModel) Register(userName, pass, userIp string) bool {
	sql := "INSERT  INTO web.tb_users(user_name,pass,last_login_ip) " +
		"(SELECT ?,?,?    WHERE NOT EXISTS (SELECT 1  FROM web.tb_users WHERE  user_name=?))"
	fmt.Println("==========", userName, pass, userIp)
	result := global.DB.Exec(sql, userName, pass, userIp, userName)
	//result := u.Exec(sql, userName, pass, userIp, userName)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// Login 用户登录,
func (u *UsersModel) Login(userName string, pass string) *UsersModel {
	sql := "select id, user_name,real_name,pass,phone  from web.tb_users where  user_name=? limit 1"
	result := global.DB.Raw(sql, userName).First(u)
	if result.Error == nil {
		// 账号密码验证成功
		if len(u.Pass) > 0 && (u.Pass == utils.Base64Md5(pass)) {
			return u
		}
	} else {
		global.Logger.Error("根据账号查询单条记录出错:", zap.Error(result.Error))
	}
	return nil
}

// OauthLoginToken 记录用户登陆（login）生成的token，每次登陆记录一次token
func (u *UsersModel) OauthLoginToken(userId int64, token string, expiresAt int64, clientIp string) bool {
	sql := `INSERT   INTO  web.tb_oauth_access_tokens(fr_user_id,action_name,token,expires_at,client_ip)
		SELECT  ?,'login',? ,?,?    WHERE   NOT   EXISTS(SELECT  1  FROM  web.tb_oauth_access_tokens a WHERE  a.fr_user_id=?  AND a.action_name='login' AND a.token=?)
		`
	//注意：token的精确度为秒，如果在一秒之内，一个账号多次调用接口生成的token其实是相同的，这样写入数据库，第二次的影响行数为0，知己实际上操作仍然是有效的。
	//所以这里只判断无错误即可，判断影响行数的话，>=0 都是ok的
	if global.DB.Exec(sql, userId, token, time.Unix(expiresAt, 0).Format(time.DateTime), clientIp, userId, token).Error == nil {
		// 异步缓存用户有效的token到redis
		//if config.GloConfig.Jwt.GetInt("Token.IsCacheToRedis") == 1 {
		//	go u.ValidTokenCacheToRedis(userId)
		//}
		return true
	}
	return false
}

// OauthRefreshConditionCheck 用户刷新token,条件检查: 相关token在过期的时间之内，就符合刷新条件
func (u *UsersModel) OauthRefreshConditionCheck(userId int64, oldToken string) bool {
	// 首先判断旧token在本系统自带的数据库已经存在，才允许继续执行刷新逻辑
	var oldTokenIsExists int
	sql := "SELECT count(*)  as  counts FROM  web.tb_oauth_access_tokens  WHERE fr_user_id =? and token=? and NOW() < (expires_at + cast(? as interval)) "
	refreshSec := config.GloConfig.Jwt.JwtTokenRefreshAllowSec
	if global.DB.Raw(sql, userId, oldToken, strconv.FormatInt(refreshSec, 10)+" second").First(&oldTokenIsExists).Error == nil && oldTokenIsExists == 1 {
		return true
	}
	return false
}

// OauthRefreshToken 用户刷新token
func (u *UsersModel) OauthRefreshToken(userId, expiresAt int64, oldToken, newToken, clientIp string) bool {
	sql := "UPDATE   web.tb_oauth_access_tokens   SET  token=? ,expires_at=?,client_ip=?,updated_at= now() ,action_name='refresh'  WHERE   fr_user_id=? AND token=?"
	if global.DB.Exec(sql, newToken, time.Unix(expiresAt, 0).Format(time.DateTime), clientIp, userId, oldToken).Error == nil {
		// 异步缓存用户有效的token到redis
		//if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
		//	go u.ValidTokenCacheToRedis(userId)
		//}
		go u.UpdateUserloginInfo(clientIp, userId)
		return true
	}
	return false
}

// UpdateUserloginInfo 更新用户登陆次数、最近一次登录ip、最近一次登录时间
func (u *UsersModel) UpdateUserloginInfo(last_login_ip string, userId int64) {
	sql := "UPDATE  web.tb_users   SET  login_times=COALESCE(login_times,0)+1,last_login_ip=?,last_login_time=?  WHERE   id=?  "
	_ = global.DB.Exec(sql, last_login_ip, time.Now().Format(time.DateTime), userId)
}

// OauthResetToken 当用户更改密码后，所有的token都失效，必须重新登录
func (u *UsersModel) OauthResetToken(userId int, newPass, clientIp string) bool {
	//如果用户新旧密码一致，直接返回true，不需要处理
	userItem, err := u.ShowOneItem(userId)
	if userItem != nil && err == nil && userItem.Pass == newPass {
		return true
	} else if userItem != nil {

		// 如果用户密码被修改，那么redis中的token值也清除
		//if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
		//	go u.DelTokenCacheFromRedis(int64(userId))
		//}

		sql := "UPDATE  web.tb_oauth_access_tokens  SET  revoked=1,updated_at= now() ,action_name='ResetPass',client_ip=?  WHERE  fr_user_id=?  "
		if global.DB.Exec(sql, clientIp, userId).Error == nil {
			return true
		}
	}
	return false
}

// OauthDestroyToken 当tb_users 删除数据，相关的token同步删除
func (u *UsersModel) OauthDestroyToken(userId int) bool {
	//如果用户新旧密码一致，直接返回true，不需要处理
	sql := "DELETE FROM  web.tb_oauth_access_tokens WHERE  fr_user_id=?  "
	//判断>=0, 有些没有登录过的用户没有相关token，此语句执行影响行数为0，但是仍然是执行成功
	if global.DB.Exec(sql, userId).Error == nil {
		return true
	}
	return false
}

// OauthCheckTokenIsOk 判断用户token是否在数据库存在+状态OK
func (u *UsersModel) OauthCheckTokenIsOk(userId int64, token string) bool {
	sql := `
		SELECT  token  FROM  web.tb_oauth_access_tokens
		WHERE   fr_user_id=?  AND  revoked=0  AND  expires_at> now()
		ORDER  BY   expires_at  DESC , updated_at  DESC limit ?
	`
	maxOnlineUsers := config.GloConfig.Jwt.JwtTokenOnlineUsers
	rows, err := global.DB.Raw(sql, userId, maxOnlineUsers).Rows()
	defer func() {
		//  凡是查询类记得释放记录集
		_ = rows.Close()
	}()
	if err == nil && rows != nil {
		for rows.Next() {
			var tempToken string
			err := rows.Scan(&tempToken)
			if err == nil {
				if tempToken == token {
					_ = rows.Close()
					return true
				}
			}
		}
	}
	return false
}

// SetTokenInvalid 禁用一个用户的: 1.tb_users表的 status 设置为 0，web.tb_oauth_access_tokens 表的所有token删除
// 禁用一个用户的token请求（本质上就是把tb_users表的 status 字段设置为 0 即可）
func (u *UsersModel) SetTokenInvalid(userId int) bool {
	sql := "delete from  web.tb_oauth_access_tokens  where  fr_user_id=?  "
	if global.DB.Exec(sql, userId).Error == nil {
		if global.DB.Exec("update  web.tb_users  set  status=0 where   id=?", userId).Error == nil {
			return true
		}
	}
	return false
}

// ShowOneItem 根据用户ID查询一条信息
func (u *UsersModel) ShowOneItem(userId int) (*UsersModel, error) {
	sql := "SELECT id, user_name,pass, real_name, phone, status,TO_CHAR(created_at,'yyyy-mm-dd hh24:mi:ss') as created_at, TO_CHAR(updated_at,'yyyy-mm-dd hh24:mi:ss') as updated_at FROM  web.tb_users  WHERE status=1 and   id=? limit 1"
	result := global.DB.Raw(sql, userId).First(u)
	if result.Error == nil {
		return u, nil
	} else {
		return nil, result.Error
	}
}

// counts 查询数据之前统计条数
func (u *UsersModel) counts(userName string) (counts int64) {
	sql := "SELECT  count(*) as counts  FROM  web.tb_users  WHERE status=1 and   user_name like ?"
	if res := global.DB.Raw(sql, "%"+userName+"%").First(&counts); res.Error != nil {
		global.Logger.Error("UsersModel - counts 查询数据条数出错", zap.Error(res.Error))
	}
	return counts
}

// Show 查询（根据关键词模糊查询）
func (u *UsersModel) Show(userName string, limitStart, limitItems int) (counts int64, temp []UsersModel) {
	if counts = u.counts(userName); counts > 0 {
		sql := `
		SELECT  id, user_name, real_name, phone, status, last_login_ip,phone,
		TO_CHAR(created_at,'yyyy-mm-dd hh24:mi:ss') as created_at, TO_CHAR(updated_at,'yyyy-mm-dd hh24:mi:ss') as updated_at
		FROM  web.tb_users  WHERE status=1 and   user_name like ? limit ? offset ?
		`
		if res := global.DB.Raw(sql, "%"+userName+"%", limitItems, limitStart).Find(&temp); res.RowsAffected > 0 {
			return counts, temp
		}
	}
	return 0, nil
}

// Store 新增
func (u *UsersModel) Store(userName string, pass string, realName string, phone string, remark string) bool {
	sql := "INSERT  INTO web.tb_users(user_name,pass,real_name,phone,remark) SELECT ?,?,?,?,?   WHERE NOT EXISTS (SELECT 1  FROM web.tb_users WHERE  user_name=?)"
	if global.DB.Exec(sql, userName, pass, realName, phone, remark, userName).RowsAffected > 0 {
		return true
	}
	return false
}

// UpdateDataCheckUserNameIsUsed 更新前检查新的用户名是否已经存在（避免和别的账号重名）
func (u *UsersModel) UpdateDataCheckUserNameIsUsed(userId int, userName string) (exists int64) {
	sql := "select count(*) as counts from web.tb_users where  id!=?  AND user_name=?"
	_ = global.DB.Raw(sql, userId, userName).First(&exists)
	return exists
}

// Update 更新
func (u *UsersModel) Update(id int, userName string, pass string, realName string, phone string, remark string, clientIp string) bool {
	sql := "update web.tb_users set user_name=?,pass=?,real_name=?,phone=?,remark=?  WHERE status=1 AND id=?"
	if global.DB.Exec(sql, userName, pass, realName, phone, remark, id).RowsAffected >= 0 {
		if u.OauthResetToken(id, pass, clientIp) {
			return true
		}
	}
	return false
}

// Destroy 删除用户以及关联的token记录
func (u *UsersModel) Destroy(id int) bool {
	// 删除用户时，清除用户缓存在redis的全部token
	//if variable.ConfigYml.GetInt("Token.IsCacheToRedis") == 1 {
	//	go u.DelTokenCacheFromRedis(int64(id))
	//}
	if u.Delete(u, id).Error == nil {
		if u.OauthDestroyToken(id) {
			return true
		}
	}
	return false
}

// 后续两个函数专门处理用户 token 缓存到 redis 逻辑
//
//func (u *UsersModel) ValidTokenCacheToRedis(userId int64) {
//	tokenCacheRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(userId)
//	if tokenCacheRedisFact == nil {
//		global.Logger.Error("redis连接失败，请检查配置")
//		return
//	}
//	defer tokenCacheRedisFact.ReleaseRedisConn()
//
//	sql := "SELECT   token,to_char(expires_at,'yyyy-mm-dd hh24:mi:ss') as expires_at  FROM  web.tb_oauth_access_tokens  WHERE   fr_user_id=?  AND  revoked=0  AND  expires_at>NOW() ORDER  BY  expires_at  DESC , updated_at  DESC  LIMIT ?"
//	maxOnlineUsers := config.GloConfig.Jwt.JwtTokenOnlineUsers
//	rows, err := u.Raw(sql, userId, maxOnlineUsers).Rows()
//	defer func() {
//		//  凡是获取原生结果集的查询，记得释放记录集
//		_ = rows.Close()
//	}()
//
//	var tempToken, expires string
//	if err == nil && rows != nil {
//		for i := 1; rows.Next(); i++ {
//			err = rows.Scan(&tempToken, &expires)
//			if err == nil {
//				if ts, err := time.ParseInLocation(variable.DateFormat, expires, time.Local); err == nil {
//					tokenCacheRedisFact.SetTokenCache(ts.Unix(), tempToken)
//					// 因为每个用户的token是按照过期时间倒叙排列的，第一个是有效期最长的，将该用户的总键设置一个最大过期时间，到期则自动清理，避免不必要的数据残留
//					if i == 1 {
//						tokenCacheRedisFact.SetUserTokenExpire(ts.Unix())
//					}
//				} else {
//					global.Logger.Error("expires_at 转换位时间戳出错", zap.Error(err))
//				}
//			}
//		}
//	}
//	// 缓存结束之后删除超过系统设置最大在线数量的token
//	tokenCacheRedisFact.DelOverMaxOnlineCache()
//}
//
//// DelTokenCacheFromRedis 用户密码修改后，删除redis所有的token
//func (u *UsersModel) DelTokenCacheFromRedis(userId int64) {
//	tokenCacheRedisFact := token_cache_redis.CreateUsersTokenCacheFactory(userId)
//	if tokenCacheRedisFact == nil {
//		global.Logger.Error("redis连接失败，请检查配置")
//		return
//	}
//	tokenCacheRedisFact.ClearUserToken()
//	tokenCacheRedisFact.ReleaseRedisConn()
//}
