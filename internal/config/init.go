package config

import (
	"gin-starter/internal/global"
	"gin-starter/internal/util/gconf"
	"path/filepath"
)

var (
	GloConfig *Config
)

func init() {
	confFile := "server.yaml"
	conf := filepath.Join(global.RootPath, "conf", confFile)
	gconf.InitConfig(conf, &GloConfig)
}

type Config struct {
	Server    Server    `yaml:"server"`
	DB        DB        `yaml:"DB"`
	Redis     Redis     `yaml:"Redis"`
	Logs      Logs      `yaml:"Logs"`
	Jwt       Jwt       `yaml:"Jwt"`
	Websocket Websocket `yaml:"Websocket"`
}

type Server struct {
	Debug        bool   `yaml:"debug"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"ReadTimeout"`
	WriteTimeout int    `yaml:"WriteTimeout"`
	IdleTimeout  int    `yaml:"IdleTimeout"`
}

type DB struct {
	Enable       bool   `yaml:"Enable"`
	Type         string `yaml:"Type"`
	DSN          string `yaml:"DSN"`
	MaxLifetime  int    `yaml:"MaxLifetime"`
	MaxIdletime  int    `yaml:"MaxIdletime"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
}

type Logs struct {
	Level string `yaml:"Level"`
	Path  string `yaml:"Path"`
	DbLog string `yaml:"DbLog"`
}

type Redis struct {
	Port               int    `yaml:"Port"`
	Auth               string `yaml:"Auth"`
	MaxIdle            int    `yaml:"MaxIdle"`
	ReConnectInterval  int    `yaml:"ReConnectInterval"`
	Host               string `yaml:"Host"`
	MaxActive          int    `yaml:"MaxActive"`
	IdleTimeout        int    `yaml:"IdleTimeout"`
	IndexDb            int    `yaml:"IndexDb"`
	ConnFailRetryTimes int    `yaml:"ConnFailRetryTimes"`
}

type Jwt struct {
	JwtTokenSignKey string `yaml:"JwtTokenSignKey"`
	ExpiresTime     string `json:"ExpiresTime" yaml:"ExpiresTime"` // 过期时间
}

type Websocket struct {
	ReadDeadline          int `yaml:"ReadDeadline"`
	WriteDeadline         int `yaml:"WriteDeadline"`
	Start                 int `yaml:"Start"`
	WriteReadBufferSize   int `yaml:"WriteReadBufferSize"`
	MaxMessageSize        int `yaml:"MaxMessageSize"`
	PingPeriod            int `yaml:"PingPeriod"`
	HeartbeatFailMaxTimes int `yaml:"HeartbeatFailMaxTimes"`
}