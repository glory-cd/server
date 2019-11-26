/**
* @Author: xhzhang
* @Date: 2019/7/12 17:34
 */
package comm

import (
	"encoding/json"
	"github.com/glory-cd/utils/log"
	"strings"
	"time"
)

//------------------------------------------------------------
// agent
func (a *Agent) CheckRecord() bool {
	return DB.First(&Agent{}, "id = ?", a.ID).RecordNotFound()
}

func (a *Agent) Online() error {
	if a.CheckRecord() {
		return CreateRecord(&a)
	} else {
		a.Status = "1"
		return UpdateRecord(&a)
	}
}

func (a *Agent) Offline() error {
	return DB.Model(&a).UpdateColumn("status", '0').Error
}

func (a *Agent) SetAlias(name string) error {
	return DB.Model(&a).UpdateColumn("alias", name).Error
}

func SyncAgentFromEtcdToDBOnline(agentId, hostName, hostIp string) {
	//(1) etcd中存在，数据库中存在，但是状态下线，则将数据库中该agent状态更新为上线
	//(2) etcd中存在，数据库中不存在，则将该agent信息插入数据库，状态为上线
	a := Agent{ID: agentId, HostName: hostName, HostIp: hostIp}
	if err := a.Online(); err != nil {
		log.Slogger.Errorf("[SyncAgentToDBOnline] %s. agent id is [%s]", err, agentId)
	} else {
		log.Slogger.Infof("[SyncAgentToDBOnline] success. agent id is [%s]", hostIp)
	}

}

func SyncAgentFromEtcdToDBOffline(agentKeys []string) {
	//(3) etcd中不存在，数据库中存在，但是状态是上线，则将数据库中中该agent状态更新为下线
	var needOffline []Agent
	if err := DB.Not(agentKeys).Find(&needOffline).UpdateColumn("status", "0").Error; err != nil {
		log.Slogger.Errorf("[SyncAgentToDBOffline] %s.", err)
	}
}

//---------------------------------------------------------------------------
// service
func (s *Service) CheckRecord() bool {
	return DB.First(&Service{}, "id = ?", s.ID).RecordNotFound()
}

func (s *Service) OnLine() error {
	if s.CheckRecord() {
		return CreateRecord(&s)
	} else {
		return UpdatePartRecord(&s, Service{CodePatterns: s.CodePatterns, StartCMD: s.StartCMD, StopCMD: s.StopCMD, Pidfile: s.Pidfile})
	}
}

func (s *Service) OffLine() error {
	if !s.CheckRecord() {
		// 删除记录
		return DB.Delete(&s).Error
	}
	return nil
}

func (s *Service) ChangeGroup(newsgroup Group) error {
	return nil
}

// json字符串转换成Service对象
func NewService(sJson string, agentId string) (s Service, err error) {
	err = json.Unmarshal([]byte(sJson), &s)
	s.AgentID = agentId
	s.CodePatterns = strings.Join(s.CodePattern, ";")
	return
}

func SyncServiceFromEtcdToDB(agentId, service string) {
	s, err := NewService(service, agentId)
	if err != nil {
		log.Slogger.Errorf("[SyncServiceToDB] convert service string to service object failed. %s", err)
		return
	}

	if s.CheckRecord() {
		if err = CreateRecord(&s); err != nil {
			log.Slogger.Errorf("[SyncServiceToDB] %s. service id is %s", err, s.ID)
		} else {
			log.Slogger.Infof("[SyncServiceToDB] success. service id is [%s]", s.ID)
		}
	} else {
		if err = UpdatePartRecord(&s, Service{CodePatterns: s.CodePatterns, StartCMD: s.StartCMD, StopCMD: s.StopCMD, Pidfile: s.Pidfile}); err != nil {
			log.Slogger.Errorf("[SyncServiceToDB] %s. service id is %s", err, s.ID)
		} else {
			log.Slogger.Infof("[SyncServiceToDB] success. service id is [%s]", s.ID)
		}
	}
}

//--------------------------------------------------------------
//task
// 设置该任务开始时间
/*func (t *Task) SetTaskStartTime() error {
	return DB.Model(&t).UpdateColumn("start_time", time.Now()).Error
}*/

func (t *Task) SetTaskStartTimeAndRunningStatus() error {
	return DB.Model(&t).Updates(map[string]interface{}{"status": 4, "start_time": time.Now()}).Error
}

// 设置该任务结束时间
func (t *Task) SetTaskEndTime() error {
	return DB.Model(&t).UpdateColumn("end_time", time.Now()).Error
}

// 设置任务状态
func (t *Task) SetTaskStatus(status int) error {
	return DB.Model(&t).UpdateColumn("status", status).Error
}

// 设置该任务状态和结束时间
func (t *Task) SetTaskEndTimeAndStatus(status int) error {
	return DB.Model(&t).Updates(map[string]interface{}{"status": status, "end_time": time.Now()}).Error
}

//------------------------------------------------------------------
//CronTask
func (t *Cron_Task) CheckRecord() bool {
	return DB.Find(&t).RecordNotFound()
}

func (t *Cron_Task) SetEntryID(newID int) error {
	return DB.Model(&t).UpdateColumn("entry_id", newID).Error
}

/*func (t *Cron_Task) SetEffective(e bool) error {
	return DB.Model(&t).UpdateColumn("effective", e).Error
}*/
