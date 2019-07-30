package comm

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"os"
	"sync"
	"time"
)

type RedisConfig struct {
	Host      string        `json:"host"`
	MaxIdle   int           `json:"maxidele"`
	MaxActive int           `json:"maxactive"`
	Timeout   time.Duration `json:"timeout"`
}

type EtcdConfig struct {
	Endpoint    []string      `json:"endpoint"`
	DialTimeout time.Duration `json:"dialtimeout"`
}

type MysqlConfig struct {
	Dsn          string `json:"dsn"`
	MaxOpenConns int    `json:"maxopenconns"`
	MaxIdleConns int    `json:"maxidleconns"`
}

type LogConfig struct {
	LogLevel   string `json:"loglevel"`   // 日志级别
	Filename   string `json:"filename"`   // 日志文件路径
	MaxSize    int    `json:"maxsize"`    // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int    `json:"maxbackups"` // 日志文件最多保存多少个备份
	MaxAge     int    `json:"maxage"`     // 文件最多保存多少天
	Compress   bool   `json:"compress"`   // 是否压缩
}

type GlobalConfig struct {
	Debug bool         `json:"debug"`
	Mysql *MysqlConfig `json:"mysql"`
	Redis *RedisConfig `json:"redis"`
	Etcd  *EtcdConfig  `json:"etcd"`
	Log   *LogConfig   `json:"log"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if !file.IsExist(cfg) {
		log.Fatal("config file:", cfg, "is not existent.")
		os.Exit(0)
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatal("read config file:", cfg, "fail:", err)
		os.Exit(0)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatal("parse config file:", cfg, "fail:", err)
		os.Exit(0)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Print("read config file:" + cfg + " successfully")
}
