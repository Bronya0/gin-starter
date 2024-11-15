package service

import (
	"fmt"
	"gin-starter/internal/util/logger"
	"github.com/robfig/cron/v3"
)

func InitCronJob() {
	c := cron.New()
	//依次是 分 时 日 月 周。@every 1s、@every 1h、@every 1m、@every 1m2s、@every 1h30m10s
	_, err := c.AddFunc("@every 1m", PrintJob)
	if err != nil {
		panic(err)
	}
	c.Start()

	logger.Logger.Info("定时任务加载成功...")

}

func PrintJob() {
	fmt.Println("主人你好")
}
