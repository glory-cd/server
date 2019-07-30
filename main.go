/**
* @Author: xhzhang
* @Date: 2019-05-05 13:37
 */
package main

import (
	"server/comm"
	"server/server"
)

func main() {
	cfg := "conf/server.json"
	comm.ParseConfig(cfg)
	comm.InitLog()
	comm.ConnDB()
	comm.ConnRedis()
	comm.ConnEtcd()
	comm.WatchAgent()
	comm.WatchService()
	go server.InitRpcServer()
	select {}
}
