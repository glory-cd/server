/**
* @Author: xhzhang
* @Date: 2019/7/15 15:11
 */
package comm

import (
	"github.com/glory-cd/utils/cron"
	"github.com/glory-cd/utils/etcd"
	"github.com/glory-cd/utils/log"
	redis "github.com/glory-cd/utils/myredis"
	"github.com/jinzhu/gorm"
	"os"
	"strings"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var EtcdClient *etcd.BaseClient
var MyConfig *Config

func init() {
	fileName := "logs/server.log"
	log.InitLog(fileName, 128, 30, 7, true)
}

// 连接etcd
func ConnEtcd(endpoints string) {
	var Etcdhost string
	if endpoints != ""{
		Etcdhost = endpoints
	}else{
		Etcdhost = os.Getenv("ETCD_HOST")
	}

	if Etcdhost == ""{
		log.Slogger.Fatalf("[Etcd] Etcd host is empty.")
	}


	endpoint := strings.Split(Etcdhost, ";")
	client, err := etcd.NewBaseClient(endpoint, 10)
	if err != nil {
		log.Slogger.Fatalf("[Etcd] conn etcd failed. %s", err)
	}
	EtcdClient = client
	MyConfig = New()

	if (Config{}) == *MyConfig{
		log.Slogger.Fatal("config info is empty.")
	}else{
		log.Slogger.Infof("[Config] log level is %s", MyConfig.LogLevel)
		log.Slogger.Infof("[Config] db dsn is %s", MyConfig.DBDsn)
		log.Slogger.Infof("[Config] redis host is %s", MyConfig.RedisHost)
		log.Slogger.Infof("[Config] rpc host is %s", MyConfig.RPCHost)
		log.Slogger.Infof("[Config] rpc cert file is %s", MyConfig.RPCCertFile)
		log.Slogger.Infof("[Config] rpc key file is %s", MyConfig.RPCKeyFile)
	}

}

func InitLogLevel() {
	logLevel := MyConfig.LogLevel
	log.SetLevel(logLevel)
}

var DB *gorm.DB

func ConnDB() {
	/*
		sqlite conn string is file:dbName?cache=shared&mode=rwc
	*/
	dsn := MyConfig.DBDsn
	level := MyConfig.LogLevel
	var err error
	DB, err = gorm.Open("sqlite3", dsn)
	if err != nil {
		log.Slogger.Fatalf("[DB] Open DB Error:[%s]", err)
	}
	log.Slogger.Infof("[DB] Open DB successful.")
	// 数据库操作日志初始化
	if level == "debug" {
		DB.SetLogger(log.GetDBLogger())
		DB = DB.Debug()
	} else {
		DB.LogMode(false)
	}

	if err = DB.DB().Ping(); err != nil {
		log.Slogger.Fatalf("[DB] DB Ping Error:[%s]", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "cdp_" + defaultTableName
	}
	// 初始化表
	log.Slogger.Info("[DB] create tables if table is not exist...")
	inittable()
}

var RedisConn redis.RedisConn

func ConnRedis() {
	host := MyConfig.RedisHost
	redisPool := redis.NewRedisPool(host, 3, 0, 300)
	RedisConn = redis.NewRedisConn(redisPool)
}

var CronClient cron.CronClient

// init cron
func StartCron() {
	CronClient = cron.NewCronClient()
	CronClient.StartCron()
	log.Slogger.Infof("[Cron] start timed-task server....")

	cronMap, cronIdMap, err := GetCurrentCronTask()
	if err != nil {
		log.Slogger.Fatalf("[Corn] Get Cron init value failed. %v", err)
	}
	AddCronTasks(cronMap, cronIdMap)
}
