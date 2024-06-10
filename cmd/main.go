package main

import (
	"gin-starter/internal/global"
	"gin-starter/internal/router"
	"gin-starter/internal/utils"
	"gin-starter/internal/utils/logger"
	"gin-starter/internal/utils/validator_zh"
	"path"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy

func main() {
	// 加载配置文件
	utils.InitConfig(path.Join(global.RootPath, "config/config.yaml"))
	// 初始化logger
	logger.InitLogger()
	// 连接数据库
	utils.InitDB()
	// 初始化定时任务
	utils.InitCronJob()
	// 初始化校验器，并本地化，zh/en
	validator_zh.InitValidator("zh")
	// 注册路由，启动gin服务
	router.InitServer()

}
