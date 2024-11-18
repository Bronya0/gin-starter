package middleware

import (
	"fmt"
	"gin-starter/internal/model/resp"
	"gin-starter/internal/util/glog"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// CustomRecovery 自定义错误(panic等)拦截中间件、对可能发生的错误进行拦截、统一记录
// 只能recover主线程的panic错误
func CustomRecovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 这里针对发生的panic等异常进行统一响应即可
		// 这里的 err 数据类型为 ：runtime.boundsError  ，需要转为普通数据类型才可以输出
		resp.Error(c, "", fmt.Sprintf("%s", err))
	})
}

// PanicExceptionRecord  panic等异常记录
type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	glog.Log.Error("【CustomRecovery】主协程内部错误：", err)
	return len(errStr), err
}
