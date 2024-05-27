package models

type Config struct {
	Server    Server    `yaml:"server"`
	DB        DB        `yaml:"DB"`
	Redis     Redis     `yaml:"Redis"`
	Logs      Logs      `yaml:"Logs"`
	Jwt       Jwt       `yaml:"Jwt"`
	Websocket Websocket `yaml:"Websocket"`
}

type Server struct {
	Debug bool   `yaml:"debug"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
}

type DB struct {
	Type         string `yaml:"Type"`
	DSN          string `yaml:"DSN"`
	MaxLifetime  int    `yaml:"MaxLifetime"`
	MaxOpenConns int    `yaml:"MaxOpenConns"`
	MaxIdleConns int    `yaml:"MaxIdleConns"`
}

type Logs struct {
	Level         string `yaml:"Level"`
	TextFormat    string `yaml:"TextFormat"`
	TimePrecision string `yaml:"TimePrecision"`
	MaxSize       int    `yaml:"MaxSize"`
	MaxBackups    int    `yaml:"MaxBackups"`
	MaxAge        int    `yaml:"MaxAge"`
	Compress      bool   `yaml:"Compress"`
	AccessLog     string `yaml:"AccessLog"`
	ErrorLog      string `yaml:"ErrorLog"`
	Path          string `yaml:"Path"`
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
	JwtTokenSignKey         string `yaml:"JwtTokenSignKey"`
	JwtTokenOnlineUsers     int    `yaml:"JwtTokenOnlineUsers"`
	JwtTokenCreatedExpireAt int64  `yaml:"JwtTokenCreatedExpireAt"`
	JwtTokenRefreshAllowSec int64  `yaml:"JwtTokenRefreshAllowSec"`
	JwtTokenRefreshExpireAt int    `yaml:"JwtTokenRefreshExpireAt"`
	BindcKeyName            string `yaml:"BindcKeyName"`
	IsCacheToRedis          int    `yaml:"IsCacheToRedis"`
	ExpiresTime             string `mapstructure:"ExpiresTime" json:"ExpiresTime" yaml:"ExpiresTime"` // 过期时间
	BufferTime              string `mapstructure:"BufferTime" json:"BufferTime" yaml:"BufferTime"`    // 缓冲时间
	Issuer                  string `mapstructure:"Issuer" json:"Issuer" yaml:"Issuer"`                // 签发者
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
