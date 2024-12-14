package global

const (
	ErrorsParseTokenFail    string = "解析token失败"
	ErrorsTokenInvalid      string = "无效的token"
	ErrorsTokenNotActiveYet string = "token 尚未激活"
	ErrorsTokenMalFormed    string = "token 格式不正确"
)

const (
	//  登录认证相关
	ErrorsLoginFail           string = "登录失败"
	ErrorsLoginExpired        string = "登录已过期"
	ErrorsLoginInvalid        string = "无效的登录"
	ErrorsLoginTokenExpired   string = "token 已过期"
	ErrorsLoginTokenMalFormed string = "token 格式不正确"
)
