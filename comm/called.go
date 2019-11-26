/**
* @Author: xhzhang
* @Date: 2019/7/15 15:18
 */
package comm

import (
	"fmt"
	"github.com/glory-cd/utils/log"
	"strings"
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

// 更新部分字段
func UpdatePartRecord(r interface{}, u interface{}) error {
	return DB.Model(r).UpdateColumns(u).Error
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

/*
|----------------|
| watch Agent  |
|----------------|
*/
func WatchAgent() {
	agentMap, err := EtcdClient.GetAgents(HandleAgentMessage)
	if err != nil {
		log.Slogger.Errorf("[WatchAgent] get agents from etcd failed. %s", err)
		return
	}
	log.Slogger.Infof("[WatchAgent] get agents from etcd success.")
	log.Slogger.Debugf("[WatchAgent] current agents is:  %+v.", agentMap)

	var agentKeys []string
	for k, v := range agentMap {
		agentId := SeparateAgentKey(k)
		hostname, hostIp, err := SeparateAgentVal(v)
		if err != nil {
			log.Slogger.Errorf("[SyncAgent] %s", err)
			return
		}
		SyncAgentFromEtcdToDBOnline(agentId, hostname, hostIp)
		agentKeys = append(agentKeys, agentId)
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
	log.Slogger.Debugf("[AgentOnline] receive agent online message. %s => %+v.", key, val)
	agentId := SeparateAgentKey(key)
	hostName, hostIp, err := SeparateAgentVal(val)
	if err != nil {
		log.Slogger.Errorf("[AgentOnline] %s. agent id is [%s], register info is [%s]", err, agentId, val)
		return
	}
	log.Slogger.Debugf("[AgentOnline] agent id: [%s]", agentId)
	log.Slogger.Debugf("[AgentOnline] agent hostname: [%s]", hostName)
	log.Slogger.Debugf("[AgentOnline] agent host ip: [%s]", hostIp)

	agent := Agent{ID: agentId, HostName: hostName, HostIp: hostIp}
	if err = agent.Online(); err == nil {
		log.Slogger.Infof("[AgentOnline] success. id is [%s]", agentId)
	} else {
		log.Slogger.Errorf("[AgentOnline] %s. agent id is [%s]", err, agentId)
	}
}

func agentOffline(key, val string) {
	log.Slogger.Debugf("[AgentOffline] receive agent offline message. %s => %+v.", key)
	agentId := SeparateAgentKey(key)

	agent := Agent{ID: agentId}
	if err := agent.Offline(); err == nil {
		log.Slogger.Infof("[AgentOffline] success. [%s]", agentId)
	} else {
		log.Slogger.Errorf("[AgentOffline] %s. agent id is [%s]", err, agentId)
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
		return "", "", fmt.Errorf("agent register info format error. [%s].", val)
	}
}

/*
|----------------|
| watch service  |
|----------------|
*/

func WatchService() {
	serviceMap, err := EtcdClient.GetServices(HandleServiceMessage)
	if err != nil {
		log.Slogger.Errorf("[WatchService] get services from etcd failed. %s", err)
		return
	}
	log.Slogger.Info("[WatchService] get services from etcd success. ")
	log.Slogger.Debugf("[WatchService] current services is: %+v", serviceMap)

	for k, v := range serviceMap {
		agentId := SeparateAgentKey(k)
		SyncServiceFromEtcdToDB(agentId, v)
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
	log.Slogger.Debugf("[ServiceOnline] receive service online message. %s => %+v.", key, val)
	agentId := SeparateAgentKey(key)
	s, err := NewService(val, agentId)
	if err != nil {
		log.Slogger.Errorf("[ServiceOnline] convert service string to service object failed. %s", err)
		return
	}

	if err = s.OnLine(); err != nil {
		log.Slogger.Errorf("[ServiceOnline] %s. service id is [%s]", err, s.ID)
	} else {
		log.Slogger.Infof("[ServiceOnline] success. service id is [%s].", s.ID)
	}

}

func serviceOffline(key, val string) {
	agentId, serviceId := getServiceKeyDetail(key)
	service := Service{ID: serviceId, AgentID: agentId}
	if err := service.OffLine(); err != nil {
		log.Slogger.Errorf("[ServiceOffline] %s. service id is [%s]", err, serviceId)
	} else {
		log.Slogger.Infof("[ServiceOffline] success. service id is [%s]", serviceId)
	}
}

// 获取etcd中服务的key中的agentID和serviceID
func getServiceKeyDetail(key string) (string, string) {
	detail := strings.Split(key, "/")
	return detail[2], detail[3]
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

//cron------------------------------------------------------------------
func GetCurrentCronTask() (map[int]string, map[int]int, error) {
	var cronTasks []Cron_Task
	cronTaskMap := map[int]string{}
	cronTaskIdMap := map[int]int{}
	if err := DB.Find(&cronTasks).Error; err != nil {
		log.Slogger.Errorf("[Cron] get current timed tasks failed. %v", err)
		return cronTaskMap, cronTaskIdMap, err
	}

	for _, ct := range cronTasks {
		cronTaskMap[ct.TaskID] = ct.TimeSpec
		cronTaskIdMap[ct.TaskID] = ct.EntryID
	}
	log.Slogger.Debugf("[Cron] current timed-task is: %+v", cronTaskMap)
	return cronTaskMap, cronTaskIdMap, nil
}
