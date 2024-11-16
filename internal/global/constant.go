package global

import (
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"log"
	"os"
	"path"
)

var (
	DB *gorm.DB

	RootPath   = path.Dir(getWorkDir())
	HttpClient = resty.New()
)

// 这里定义的常量，一般是具有错误代码+错误说明组成，一般用于接口返回
const ()

func getWorkDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return wd
}
