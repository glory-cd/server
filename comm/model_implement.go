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

func SyncAgentFromEtcdToDBOnline(agentid, hostname, hostip string) {
	//(1) etcd中存在，数据库中存在，但是状态下线，则将数据库中该agent状态更新为上线
	//(2) etcd中存在，数据库中不存在，则将该agent信息插入数据库，状态为上线
	a := Agent{ID: agentid, HostName: hostname, HostIp: hostip}
	if err := a.Online(); err != nil {
		log.Slogger.Errorf("[SyncAgent] 同步agent[%s]失败[%s]。", hostip, err)
	} else {
		log.Slogger.Infof("[SyncAgent] 同步agent[%s]成功。", hostip)
	}

}

func SyncAgentFromEtcdToDBOffline(agentkeys []string) {
	//(3) etcd中不存在，数据库中存在，但是状态是上线，则将数据库中中该agent状态更新为下线
	var needOffline []Agent
	if err := DB.Not(agentkeys).Find(&needOffline).UpdateColumn("status", "0").Error; err != nil {
		log.Slogger.Errorf("[SyncAgent] 同步下线agent失败[%s]。", err)
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
		return UpdateRecord(&s)
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

// josn字符串转换成Service对象
func NewService(sjson string, agentid string) (s Service, err error) {
	err = json.Unmarshal([]byte(sjson), &s)
	s.AgentID = agentid
	s.CodePatterns = strings.Join(s.CodePattern, ";")
	return
}

func SyncServiceFromEtcdToDB(agentid, service string) {
	s, err := NewService(service, agentid)
	if err != nil {
		log.Slogger.Errorf("[SyncService] 转换服务json字符串为Service对象出错。%s", err)
		return
	}

	if s.CheckRecord() {
		if err = CreateRecord(&s); err != nil {
			log.Slogger.Errorf("[SyncService] 同步服务[%s]失败。%s", s.ID, err)
		} else {
			log.Slogger.Infof("[SyncService] 同步服务[%s]成功。", s.Name)
		}
	} else {
		if err = UpdateRecord(&s); err != nil {
			log.Slogger.Errorf("[SyncService] 同步更新服务[%s]失败。%s", s.ID, err)
		} else {
			log.Slogger.Infof("[SyncService] 同步更新服务[%s]成功。", s.ID)
		}
	}
}

//--------------------------------------------------------------
//task
// 设置该任务开始时间
/*func (t *Task) SetTaskStartTime() error {
	return DB.Model(&t).UpdateColumn("start_time", time.Now()).Error
}*/

func (t *Task) SetTaskStartTimeAndRuningStatus() error{
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
func (t *Cron_Task) CheckRecord() bool{
	return DB.Find(&t).RecordNotFound()
}

func (t *Cron_Task) SetEntryID(newID int) error {
	return DB.Model(&t).UpdateColumn("entry_id", newID).Error
}

/*func (t *Cron_Task) SetEffective(e bool) error {
	return DB.Model(&t).UpdateColumn("effective", e).Error
}*/

