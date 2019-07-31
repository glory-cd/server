/**
* @Author: xhzhang
* @Date: 2019/7/10 13:50
 */
package server

import (
	"encoding/json"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"strings"
	"sync"
	"time"
	"github.com/glory-cd/utils/log"
)

const (
	Status_Fail    int = 0 //失败
	Status_Success int = 1 //成功
)

//struct-----------------------------------------------------------------
type CMD struct {
	TaskId               int      `json:"taskid"`
	ExecutionId          int      `json:"executionid"`
	ServiceId            string   `json:"serviceid"`
	ServiceOp            int      `json:"serviceop"`
	ServiceName          string   `json:"servicename"`
	ServiceOsUser        string   `json:"serviceosuser"`
	ServiceModuleName    string   `json:"servicemodulename"`
	ServiceDir           string   `json:"servicedir"`
	ServiceRemoteCode    string   `json:"serviceremotecode"`
	ServiceCodePattern   []string `json:"servicecodepattern"`
	ServiceCustomPattern []string `json:"servicecustompattern"`
	ServicePidfile       string   `json:"servicepidfile"`
	ServiceStartCmd      string   `json:"servicestartcmd"`
	ServiceStopCmd       string   `json:"servicestopcmd"`
}

type Result struct {
	TaskId      int    `json:"taskid"`
	ExecutionId int    `json:"executionid"`
	ResultCode  int    `json:"rcode"`
	ResultMsg   string `json:"rmsg"`
	ResultSteps []struct {
		StepNum   int    `json:"stepnum"`
		StepName  string `json:"stepname"`
		StepState int    `json:"stepstate"`
		StepMsg   string `json:"stepmsg"`
		StepTime  int64  `json:"steptime"`
	} `json:"rsteps"`
}

type ControlAgent struct {
	TaskId  int    `json:"taskid"`
	AgentID string `json:"agentid"`
	OPMode  string `json:"gracecmd"` //SIGHUP
}

//func----------------------------------------------------------------
func checkTaskIsExecute(taskid int) (bool, error) {
	task := comm.Task{}
	if err := comm.DB.Select("status").Where("id = ?", taskid).Find(&task).Error; err != nil {
		return false, err
	}
	if task.Status != 2 {
		return false, nil
	}
	return true, nil
}

func getExecutions(task *comm.Task, rrei interface{}) ([]map[string]interface{}, error) {
	var executionMapList []map[string]interface{}
	// 获取切片列表
	var elist []comm.Execution
	if err := comm.DB.Preload("Task").Preload("Service").Where("task_id = ?", task.ID).Find(&elist).Error; err != nil {
		_ = task.SetTaskEndTime()
		return executionMapList, err
	}
	// 解析
	for _, e := range elist {
		// 获取代码路径
		var releasecode comm.ReleaseCode
		comm.DB.Preload("Release").Where("release_id = ?", e.Task.ReleaseID).Where("name = ?", e.Service.ModuleName).Find(&releasecode)
		tmp := CMD{TaskId: e.TaskID,
			ExecutionId:        e.ID,
			ServiceId:          e.ServiceID,
			ServiceOp:          e.Operation,
			ServiceName:        e.Service.Name,
			ServiceOsUser:      e.Service.OsUser,
			ServiceModuleName:  e.Service.ModuleName,
			ServiceDir:         e.Service.Dir,
			ServiceRemoteCode:  releasecode.RelativePath,
			ServiceCodePattern: strings.Split(e.Service.CodePatterns, ";"),
			//ServiceCustomPattern: strings.Split(e.CustomUpgradePattern, ";"),
			ServiceCustomPattern: nil,
			ServicePidfile:       e.Service.Pidfile,
			ServiceStartCmd:      e.Service.StartCMD,
			ServiceStopCmd:       e.Service.StopCMD,
		}

		publishChannel := "cmd." + e.Service.AgentID
		executionMapList = append(executionMapList, map[string]interface{}{"pchannel": publishChannel, "eobject": tmp})
		rre, ok := rrei.(*pb.ExecutionList)
		if ok {
			rre.Executions = append(rre.Executions, &pb.ExecutionList_ExecutionInfo{Id: int32(e.ID), Taskname: e.Task.Name, Servicename: e.Service.Name, Operation: int32(e.Operation)})
		}

	}
	return executionMapList, nil
}

func publishExecutionCMD(ecmdList []map[string]interface{}, rre *pb.ExecutionList) (prl []map[int]error) {
	for _, e := range ecmdList {
		publishChannel := e["pchannel"]
		executionobj := e["eobject"]
		c, _ := publishChannel.(string)
		e, _ := executionobj.(CMD)

		eid := e.ExecutionId

		ebyte, err := json.Marshal(e)
		if err != nil {
			return
		}
		err = comm.PublishCMD(c, string(ebyte))
		if err != nil {
			log.Slogger.Errorf("[PublishTask] publish切片失败:[%s]=>[%s]", e, err)
			prl = append(prl, map[int]error{eid: err})
			for _, r := range rre.Executions {
				if int32(eid) == r.Id {
					r.Rcode = 0
					r.Rmsg = "publish任务切片失败.Error:" + err.Error()
					continue
				}
			}
		} else {
			log.Slogger.Infof("[PublishTask] publish切片成功:[%d]=>[%s]", e.ExecutionId, string(ebyte))
		}
	}
	return
}

func collectResult(resultChannel chan string, publishSuccessLen int) []string {
	//
	log.Slogger.Infof("[PublishTask] 启动收集任务结果goroutine")
	// 收集任务结果
	cmdResult := make([]string, publishSuccessLen)
	wg := sync.WaitGroup{}
	for i := 0; i < publishSuccessLen; i++ {
		wg.Add(1) // 计数加 1
		go func(i int) {
			cmdResult[i] = <-resultChannel
			defer wg.Done() // 计数减 1
		}(i)
	}
	wg.Wait()
	log.Slogger.Infof("[PublishTask] 结束收集任务结果goroutine")
	return cmdResult
}

func result2DB(task comm.Task, results []string, rre *pb.ExecutionList) error {
	// 入库
	var taskstatus = Status_Success

	tx := comm.DB.Begin()
	for _, r := range results {
		var rs Result
		err := json.Unmarshal([]byte(r), &rs)
		if err != nil {
			return err
		}
		executionID := rs.ExecutionId
		resultCode := rs.ResultCode
		resultMsg := rs.ResultMsg
		resultSteps := rs.ResultSteps

		// 设置切片结果
		e := comm.Execution{ID: executionID}
		if err := tx.Model(&e).Updates(map[string]interface{}{"result_code": resultCode, "result_msg": resultMsg}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 插入切片详情
		for _, rsp := range resultSteps {
			edetail := comm.Execution_Detail{StepNum: rsp.StepNum, StepName: rsp.StepName, StepState: rsp.StepState, StepMsg: rsp.StepMsg, StepTime: time.Unix(rsp.StepTime/1e9, 0), ExecutionID: executionID}
			if err = tx.Create(&edetail).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if resultCode == Status_Fail {
			taskstatus = Status_Fail
		}

		//设置返回数组中的结果部分
		for _, re := range rre.Executions {
			if int32(executionID) == re.Id {
				re.Rcode = int32(resultCode)
				re.Rmsg = resultMsg
			}
		}

	}
	// 设置任务结束时间
	if err := task.SetTaskEndTimeAndStatus(taskstatus); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
