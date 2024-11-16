package validator_zh

import (
	"fmt"
	"gin-starter/internal/util/logger"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

// Trans 定义一个全局翻译器T
var Trans ut.Translator

func InitValidator(locale string) {
	if err := InitTrans(locale); err != nil {
		panic(err)
	}
	logger.Log.Info("验证器初始化并本土化成功...")
}

// InitTrans 初始化表单参数验证器的翻译器
func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		//初始化翻译器
		zhT := zh.New()
		enT := en.New()
		// 第一个参数是备用（fallback）的语言环境
		// 后面的参数是应该支持的语言环境（支持多个）
		// uni := ut.New(zhT, zhT) 也是可以的
		uni := ut.New(enT, zhT, enT)
		// locale 通常取决于 http 请求头的 'Accept-Language'
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		//注册翻译器
		//默认注册中文，en 注册英文 zh 注册中文
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		default:
			err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		}
		return
	}
	return
}

// RemoveTopStruct 将返回的结构体名去除掉，只留下需要的字段名
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.LastIndex(field, ".")+1:]] = err
	}
	return res
}
