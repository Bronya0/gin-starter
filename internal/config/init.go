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
	Enable        bool   `yaml:"Enable"`
	Type          string `yaml:"Type"`
	DSN           string `yaml:"DSN"`
	DbLog         string `yaml:"DbLog"`
	MaxLifetime   int    `yaml:"MaxLifetime"`
	MaxIdletime   int    `yaml:"MaxIdletime"`
	MaxOpenConns  int    `yaml:"MaxOpenConns"`
	MaxIdleConns  int    `yaml:"MaxIdleConns"`
	SlowThreshold int    `yaml:"SlowThreshold"`
}

type Logs struct {
	Level string `yaml:"Level"`
	Path  string `yaml:"Path"`
}

type Redis struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	MaxIdle      int    `yaml:"maxIdle"`
	MaxActive    int    `yaml:"maxActive"`
	IdleTimeout  int    `yaml:"idleTimeout"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
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
