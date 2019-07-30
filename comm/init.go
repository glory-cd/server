/**
* @Author: xhzhang
* @Date: 2019/7/15 15:11
 */
package comm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"utils/etcd"
	"utils/log"
	"utils/myredis"
)

func InitLog() {
	config := Config().Log
	log.InitLog(config.Filename, config.MaxSize, config.MaxBackups, config.MaxAge, config.Compress)
	log.SetLevel(config.LogLevel)
}

var DB *gorm.DB

func ConnDB() {
	var err error
	config := Config().Mysql
	logLevel := Config().Log.LogLevel
	DB, err = gorm.Open("mysql", config.Dsn)
	if err != nil {
		log.Slogger.Errorf("CreateDBPool Error:[%s]", err)
		return
	}
	DB.DB().SetMaxOpenConns(config.MaxOpenConns)
	DB.DB().SetMaxIdleConns(config.MaxIdleConns)
	// 数据库操作日志初始化
	if logLevel == "debug" {
		//DB.LogMode(true)  默认
		DB.SetLogger(log.GetDBLogger())
	} else {
		DB.LogMode(false)
	}

	if err = DB.DB().Ping(); err != nil {
		log.Slogger.Errorf("CreateDBPoolPing Error:[%s]", err)
		return
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "cdp_" + defaultTableName
	}
	// 初始化表
	inittable()
}

var RedisConn redis.RedisConn

func ConnRedis() {
	config := Config().Redis
	redispool := redis.NewRedisPool(config.Host, config.MaxIdle, config.MaxActive, config.Timeout)
	RedisConn = redis.NewRedisConn(redispool)
}

var EtcdClient *etcd.BaseClient

// 初始化etcd
func ConnEtcd() {
	config := Config().Etcd
	endpoint := config.Endpoint
	dialtimeout := config.DialTimeout
	client, err := etcd.NewBaseClient(endpoint, dialtimeout)
	if err != nil {
		log.Slogger.Errorf("连接etcd失败. %s", err)
		return
	}
	EtcdClient = client
}
