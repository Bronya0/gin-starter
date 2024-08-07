package main

import (
	"gin-starter/internal/router"
	"gin-starter/internal/service"
	"gin-starter/internal/utils/logger"
	"gin-starter/internal/utils/validator_zh"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy

func main() {
	// 初始化logger
	logger.InitLogger()
	// 连接数据库
	//utils.InitDB()
	// 初始化定时任务
	service.InitCronJob()
	// 初始化校验器，并本地化，zh/en
	validator_zh.InitValidator("zh")
	// 注册路由，启动gin服务
	router.InitServer()

}
