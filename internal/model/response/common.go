package response

import (
	"gin-starter/internal/util/validator_zh"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

// SuccessResponse 直接返回成功
func SuccessResponse(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, 2000, msg, data)
}
func ErrorResponse(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, 5000, msg, nil)
	c.Abort()
}

func ReturnJson(c *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {

	//c.Header("key2020","value2020")  	//可以根据实际情况在头部添加额外的其他信息
	c.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

func ValidatorError(c *gin.Context, err error) {
	// 1、validator.ValidationErrors类型错误
	if errs, ok := err.(validator.ValidationErrors); ok {
		wrongParam := validator_zh.RemoveTopStruct(errs.Translate(validator_zh.Trans))
		ReturnJson(c, http.StatusBadRequest, 4000, "参数校验失败", wrongParam)
	} else {
		// 2、非validator.ValidationErrors类型错误
		errStr := err.Error()
		// 2.1、multipart:nextpart:eof 错误表示验证器需要一些参数，但是调用者没有提交任何参数
		if strings.ReplaceAll(strings.ToLower(errStr), " ", "") == "multipart:nextpart:eof" {
			ReturnJson(c, http.StatusBadRequest, 4000, "参数校验失败", gin.H{"tips": "该接口不允许所有参数都为空,请按照接口要求提交必填参数"})
		} else {
			// 2.2、正常错误返回
			ReturnJson(c, http.StatusBadRequest, 4000, "参数校验失败", gin.H{"tips": errStr})
		}
	}
	c.Abort()
}
