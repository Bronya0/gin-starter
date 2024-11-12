package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	DB *gorm.DB

	Enforcer   *casbin.SyncedEnforcer
	RootPath   = path.Dir(getCurrentAbPath())
	HttpClient = resty.New()
)

// 这里定义的常量，一般是具有错误代码+错误说明组成，一般用于接口返回
const (
	// ValidatorPrefix 表单验证器前缀
	ValidatorPrefix string = "Form_Validator_"

	// JwtTokenOK token相关
	JwtTokenOK        int    = 200100                                  //token有效
	JwtTokenInvalid   int    = -400100                                 //无效的token
	JwtTokenExpired   int    = -400101                                 //过期的token
	JwtTokenMustValid string = "token为必填项,请在请求header部分提交!" //提交的 token 格式错误

	// CurdStatusOkMsg CURD 常用业务状态码
	CurdStatusOkMsg          string = "Success"
	CurdCreatFailCode        int    = -400200
	CurdCreatFailMsg         string = "新增失败"
	CurdUpdateFailCode       int    = -400201
	CurdUpdateFailMsg        string = "更新失败"
	CurdDeleteFailCode       int    = -400202
	CurdDeleteFailMsg        string = "删除失败"
	CurdSelectFailCode       int    = -400203
	CurdSelectFailMsg        string = "查询无数据"
	CurdRegisterFailCode     int    = -400204
	CurdRegisterFailMsg      string = "注册失败"
	CurdLoginFailCode        int    = -400205
	CurdLoginFailMsg         string = "登录失败"
	CurdRefreshTokenFailCode int    = -400206
	CurdRefreshTokenFailMsg  string = "刷新Token失败"
)

const (
	CSTLayout = "2006-01-02 15:04:05"
)

// 最终方案-全兼容
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
