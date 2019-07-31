/**
* @Author: xhzhang
* @Date: 2019/7/15 15:18
 */
package comm

import (
	"fmt"
	"strings"
	"github.com/glory-cd/utils/log"
)

//db--------------------------------------------------------
// 根据ID(int)查看记录是否存在
func CheckRecordWithID(id int, r interface{}) bool {
	return DB.Where("id = ?", id).First(r).RecordNotFound()
}

func CheckRecordWithStringID(id string, r interface{}) bool {
	return DB.Where("id = ?", id).First(r).RecordNotFound()
}

// 根据Name查看记录是否存在
func CheckRecordWithName(name string, r interface{}) bool {
	return DB.Where("name = ?", name).First(r).RecordNotFound()
}

// 创建记录
func CreateRecord(r interface{}) error {
	return DB.Create(r).Error
}

// 删除记录
func DeleteRecord(r interface{}) error {
	return DB.Delete(r).Error
}

// 更新全部字段
func UpdateRecord(r interface{}) error {
	return DB.Save(r).Error
}

//etcd----------------------------------------------------------------
/*watch agent
1. 从etcd中获取当前存在的agent信息
2. 同步获取的agent信息到数据库。
     存在以下情况：
     (1) etcd中存在，数据库中存在，但是状态下线，则将数据库中该agent状态更新为上线
     (2) etcd中存在，数据库中不存在，则将该agent信息插入数据库，状态为上线
	 (3) etcd中不存在，数据库中的状态是上线，则将数据库中中该agent状态更新为下线
*/
func WatchAgent() {
	agentmap, err := EtcdClient.GetAgents(HandleAgentMessage)
	if err != nil {
		log.Slogger.Errorf("[WatchAgent] 从etcd中获取Agent失败. %s", err)
		return
	}
	log.Slogger.Debugf("[WatchAgent] 从etcd中获取Agent成功.%s", agentmap)

	var agentKeys []string
	for k, v := range agentmap {
		agentid := SeparateAgentKey(k)
		hostname, hostip, err := SeparateAgentVal(v)
		if err != nil {
			log.Slogger.Errorf("[SyncAgent] %s", err)
			return
		}
		SyncAgentFromEtcdToDBOnline(agentid, hostname, hostip)
		agentKeys = append(agentKeys, agentid)
	}
	SyncAgentFromEtcdToDBOffline(agentKeys)
}

func HandleAgentMessage(messageType, messageKey, messageVal string) {
	switch messageType {
	case "PUT":
		agentOnline(messageKey, messageVal)
	case "DELETE":
		agentOffline(messageKey, messageVal)
	}

}

func agentOnline(key, val string) {
	agentid := SeparateAgentKey(key)
	hostname, hostip, err := SeparateAgentVal(val)
	if err != nil {
		log.Slogger.Errorf("[AgentOnline] [key:%s]-[val:%s]=>[%s]", key, val, err)
		return
	}
	agent := Agent{ID: agentid, HostName: hostname, HostIp: hostip}
	if err = agent.Online(); err == nil {
		log.Slogger.Infof("[AgentOnline] Agent上线入库成功：[key:%s]-[val:%s]", key, val)
	} else {
		log.Slogger.Errorf("[AgentOnline] Agent上线入库失败：[key:%s]-[val:%s]=>[%s]", key, val, err)
	}
}

func agentOffline(key, val string) {
	agentid := SeparateAgentKey(key)
	hostname, hostip, err := SeparateAgentVal(val)
	if err != nil {
		log.Slogger.Errorf("[AgentOffline] [key:%s]-[val:%s]=>[%s]", key, val, err)
		return
	}
	agent := Agent{ID: agentid, HostName: hostname, HostIp: hostip}
	if err := agent.Offline(); err == nil {
		log.Slogger.Infof("[AgentOffline] Agent下线入库成功.[key:%s]", key)
	} else {
		log.Slogger.Errorf("[AgentOffline] Agent下线入库失败.[key:%s]=>[%s]", key, err)
	}
}

// 从etcd中agent注册的key值拆分出agentid
func SeparateAgentKey(key string) string {
	return strings.Split(key, "/")[2]
}

// 从etcd中agent注册的val值拆分出agent属性值
func SeparateAgentVal(val string) (string, string, error) {
	if strings.ContainsAny(val, ":") {
		hostinfo := strings.Split(val, ":")
		return hostinfo[0], hostinfo[1], nil
	} else if val == "" { //删除key时，val是""
		return "", "", nil
	} else {
		return "", "", fmt.Errorf("[SeparateAgentAtt] Agent的注册信息的值格式错误[%s].", val)
	}
}

// watch service
func WatchService() {
	servicemap, err := EtcdClient.GetServices(HandleServiceMessage)
	if err != nil {
		log.Slogger.Errorf("[WatchService] 从etcd中获取service失败. %s", err)
		return
	}
	log.Slogger.Info("[WatchService] 从etcd中获取service成功.")
	log.Slogger.Debugf("[WatchService] etcd中的服务: %s", servicemap)

	for k, v := range servicemap {
		agentid := SeparateAgentKey(k)
		SyncServiceFromEtcdToDB(agentid, v)
	}
}

func HandleServiceMessage(messageType, messageKey, messageVal string) {
	switch messageType {
	case "PUT":
		serviceOnline(messageKey, messageVal)
	case "DELETE":
		serviceOffline(messageKey, messageVal)
	}
}

func serviceOnline(key, val string) {
	agentid := SeparateAgentKey(key)
	s, err := NewService(val, agentid)
	if err != nil {
		log.Slogger.Errorf("[ServiceOnline] %s", err)
		return
	}

	if err = s.OnLine(); err != nil {
		log.Slogger.Errorf("[ServiceOnline] 服务上线入库失败：[key:%s]-[val:%s]=>[%s]", key, val, err)
	}
	log.Slogger.Infof("[ServiceOnline] 服务上线入库成功：[key:%s]-[val:%s]", key, val)
}

func serviceOffline(key, val string) {
	agentid := SeparateAgentKey(key)
	s, err := NewService(val, agentid)
	if err != nil {
		log.Slogger.Errorf("[ServiceOffline] %s", err)
	}
	if err = s.OffLine(); err != nil {
		log.Slogger.Errorf("[ServiceOffline] 服务下线入库失败：[key:%s]-[val:%s]=>[%s]", key, val, err)
	}
	log.Slogger.Infof("[ServiceOffline] 服务下线入库成功：[key:%s]-[val:%s]=>[%s]", key, val)

}

//redis---------------------------------------------------------------
func PublishCMD(channel string, cmd string) (pr error) {
	publishResultChan := make(chan error)
	go func() {
		_, err := RedisConn.Publish(channel, cmd)
		publishResultChan <- err
	}()
	return <-publishResultChan
}

func SubscribeCMDResult(channel string, resultChan chan string) {
	psc, err := RedisConn.Subscribe(channel)
	if err != nil {
		return
	}
	RedisConn.HandleCMDResultMessage(psc, resultChan)
}
