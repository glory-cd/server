/**
* @Author: xhzhang
* @Date: 2019-05-05 13:37
 */
package main

import (
	"flag"
	"github.com/glory-cd/server/comm"
	"github.com/glory-cd/server/server"
)

var etcdEndpoints string
var version string

func init() {
	flag.StringVar(&etcdEndpoints, "etcd", "", "etcd endpoints")
	flag.StringVar(&version, "version", "", "etcd endpoints")
}

func main() {
	flag.Parse()
	comm.ConnEtcd(etcdEndpoints)
	comm.InitLogLevel()
	comm.ConnDB()
	comm.ConnRedis()
	comm.WatchAgent()
	comm.WatchService()
	comm.StartCron()
	go server.InitRpcServer()
	select {}
}
