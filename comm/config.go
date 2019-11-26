/**
* @Author: xhzhang
* @Date: 2019/11/25 15:01
 */
package comm

import (
	"github.com/glory-cd/utils/log"
)

type Config struct {
	LogLevel    string
	DBDsn       string
	RedisHost   string
	RPCHost     string
	RPCCertFile string
	RPCKeyFile  string
}

func New() *Config {
	var c Config
	c.GetLogLevel()
	c.GetDBDsn()
	c.GetRedisHost()
	c.GetRpcConfig()
	return &c
}

func (c *Config) GetLogLevel() {
	key := "/config/server/log/level"
	config, err := EtcdClient.Get(key, false)
	if err != nil {
		log.Slogger.Fatalf("[Etcd] get log config failed. %s", err)
	}
	c.LogLevel = config[key]
}

func (c *Config) GetDBDsn() {
	key := "/config/server/db/dsn"
	config, err := EtcdClient.Get(key, false)
	if err != nil {
		log.Slogger.Fatalf("[Etcd] get db config failed. %s", err)
	}
	c.DBDsn = config[key]
}

func (c *Config) GetRedisHost() {
	key := "/config/server/redis/host"
	config, err := EtcdClient.Get(key, false)
	if err != nil {
		log.Slogger.Fatalf("[Etcd] get redis config failed. %s", err)
	}
	c.RedisHost = config[key]
}

func (c *Config) GetRpcConfig() {
	key := "/config/server/rpc"
	config, err := EtcdClient.Get(key, true)
	if err != nil {
		log.Slogger.Fatalf("[Etcd] get rpc config failed. %s", err)
	}
	c.RPCHost = config[key+"/host"]
	c.RPCCertFile = config[key+"/certfile"]
	c.RPCKeyFile = config[key+"/keyfile"]
}
