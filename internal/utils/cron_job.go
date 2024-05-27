package utils

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func InitCronJob() {
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", PrintJob)
	if err != nil {
		panic(err)
	}

	fmt.Println("定时任务加载成功...")

}

func PrintJob() {
	fmt.Println("主人你好")
}
