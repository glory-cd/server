/**
* @Author: xhzhang
* @Date: 2019/9/24 15:26
 */
package comm

import (
	"github.com/glory-cd/server/client"
	"github.com/glory-cd/utils/log"
	"github.com/robfig/cron/v3"
)

type ExecTaskJob struct {
	TaskID int32
}

func (ej ExecTaskJob) Run() {
	rpcPort := MyConfig.RPCHost
	certFile := MyConfig.RPCCertFile
	cdpAttr := client.CDPCClientAttr{CertFile: certFile, Address: rpcPort}
	conn, err := client.NewClient(cdpAttr)
	if err != nil {
		log.Slogger.Errorf("[ExecuteCronTask] Conn Server failed. [%v]", err)
	}

	_, err = conn.ExecuteTask(ej.TaskID)
	if err != nil {
		log.Slogger.Errorf("[ExecuteCronTask]: [%v]", err)
		return
	}
}

/*
	add task to Cron.
    paras taskID: 任务ID
    paras taskSpec: 定时时间字符串
*/
func AddCronTask(taskID int, taskSpec string) (int, error) {
	id, err := CronClient.AddJob(taskSpec, ExecTaskJob{TaskID: int32(taskID)})
	if err != nil {
		log.Slogger.Errorf("[AddTaskToCron] add task[%d] to cron-job failed. %s", taskID, err.Error())
		return 0, err
	}
	log.Slogger.Infof("[AddTaskToCron] add task[%d] to cron[%d] successful. ", taskID, id)
	return int(id), nil
}

/*
	When initializing cron, valid cron-tasks in the database are added to the scheduling queue.
*/

func AddCronTasks(cronTask map[int]string, cronIDMap map[int]int) {
	for taskID, taskSpec := range cronTask {
		newID, _ := AddCronTask(taskID, taskSpec)
		// 如果新生成的定时任务ID与数据库中不一致，则校正数据库ID.
		originalID := cronIDMap[taskID]
		if originalID != newID {
			cronTask := Cron_Task{TaskID: taskID}
			err := cronTask.SetEntryID(newID)
			if err != nil {
				log.Slogger.Errorf("[AddTaskToCron] update task[%d] entryid form[%d] to [%d] failed. %s", taskID, originalID, newID, err.Error())
			}
		}

	}
}

func RemoveEntryFromCron(entryID int32) {
	CronClient.RemoveJob(cron.EntryID(entryID))
}
