package job

import (
	"gin-starter/internal/util/glog"
	"github.com/robfig/cron/v3"
)

func InitCronJob() {
	c := cron.New()
	// 依次是 分 时 日 月 周。@every 1s、@every 1h、@every 1m、@every 1m2s、@every 1h30m10s
	_, err := c.AddFunc("@every 1m", PrintJob)
	if err != nil {
		panic(err)
	}
	c.Start()

	glog.Log.Info("定时任务加载成功...")

}

func PrintJob() {
	glog.Log.Info("主人你好,定时任务运行...")
}
