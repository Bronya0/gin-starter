package config_util

import (
	"fmt"
	"gin-starter/internal/util/file"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitConfig 将配置文件映射城结构体
func InitConfig[T any](configPath string, Config T) {
	// 判断文件是否存在
	if !file.IsExist(configPath) {
		panic(fmt.Errorf("配置文件不存在: %s \n", configPath))
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("无效的配置文件: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件变更:", e.Name)
		if err = v.Unmarshal(&Config); err != nil {
			panic(err)
		}
	})
	if err = v.Unmarshal(&Config); err != nil {
		panic(err)
	}
	fmt.Println("配置加载成功...", configPath)
}
