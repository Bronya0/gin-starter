package utils

import (
	"fmt"
	"gin-starter/internal/global"
	"github.com/robfig/cron/v3"
)

func InitCronJob() {
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", PrintJob)
	if err != nil {
		panic(err)
	}

	global.Logger.Info("定时任务加载成功...")

}

func PrintJob() {
	fmt.Println("主人你好")
}
