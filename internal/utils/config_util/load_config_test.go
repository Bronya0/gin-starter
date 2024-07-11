package config_util

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path"
	"runtime"
	"testing"
	"time"
)

type Key struct {
	Comcode struct {
		Mobile  string `json:"mobile"`
		Telecom string `json:"telecom"`
		Unicom  string `json:"unicom"`
	} `json:"comcode"`
}

var Keys Key

func TestConfig(t *testing.T) {
	_, filepath, _, _ := runtime.Caller(0)
	fmt.Println(filepath)
	config_path := path.Join(path.Dir(path.Dir(filepath)), "config")
	v := viper.New()
	v.SetConfigFile(path.Join(config_path, "key.json"))
	v.SetConfigType("json")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("无效的配置文件: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件变更:", e.Name)
		if err = v.Unmarshal(&Keys); err != nil {
			panic(err)
		}
	})
	if err = v.Unmarshal(&Keys); err != nil {
		panic(err)
	}
	fmt.Println("配置加载成功...")
	for {
		data_bytes, _ := json.Marshal(Keys)
		fmt.Println(string(data_bytes))
		time.Sleep(3 * time.Second)
	}

}
