/**
* @Author: xhzhang
* @Date: 2019/7/10 13:50
 */
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
	"strings"
	"sync"
	"time"
)

const (
	TaskStatus_Failed     int = 0 //失败
	TaskStatus_Successful int = 1 //成功
	TaskStatus_NoExecute  int = 2 //未执行
	TaskStatus_Cron       int = 3 //定时任务
	TaskStatus_Running    int = 4 //正在执行
)

type OpMode int32

const (
	OperateDefault  OpMode = 0
	OperateDeploy   OpMode = 1
	OperateUpgrade  OpMode = 2
	OperateStart    OpMode = 3
	OperateStop     OpMode = 4
	OperateRestart  OpMode = 5
	OperateCheck    OpMode = 6
	OperateBackUp   OpMode = 7
	OperateRollBack OpMode = 8
)

//struct-----------------------------------------------------------------
type CMD struct {
	TaskId               int      `json:"taskid"`
	ExecutionId          int      `json:"executionid"`
	ServiceId            string   `json:"serviceid"`
	ServiceOp            OpMode   `json:"serviceop"`
	ServiceName          string   `json:"servicename"`
	ServiceOsUser        string   `json:"serviceosuser"`
	ServiceOsPass        string   `json:"serviceosuserpass"`
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
/*
	检查任务当前状态
    返回结果: 是否可执行;任务状态;错误信息
*/
func  checkTaskIsExecute(taskid int) (bool, int, error) {
	task := comm.Task{}
	if err := comm.DB.Select("status").Where("id = ?", taskid).Find(&task).Error; err != nil {
		return false, 0, err
	}
	if task.Status == TaskStatus_Successful {
		return false, task.Status, errors.New("the task has been executed successful.")
	} else if task.Status == TaskStatus_Running {
		return false, task.Status, errors.New("the task is running.")
	} else {
		return true, task.Status, nil
	}
}
/*
	获取需要执行的任务切片.
    1. 任务状态=0(失败)，则选择出失败的切片
    2. 任务状态=2(未执行)，则选择出所有相关切片
    3. 任务状态=3(定时任务)，则选择出所有相关切片
*/
func getExecutions(task *comm.Task, taskStatus int, rrei interface{}) ([]map[string]interface{}, error) {
	var executionMapList []map[string]interface{}
	// 获取切片列表
	queryCmd := comm.DB.Preload("Task").Preload("Service").Where("task_id = ?", task.ID)
	var eList []comm.Execution
	// 如果任务状态是执行失败(0),则只选出失败的切片去执行
	if taskStatus == TaskStatus_Failed{
		queryCmd = queryCmd.Where("result_code = ?",TaskStatus_Failed)
	}
	if err := queryCmd.Find(&eList).Error; err != nil {
		_ = task.SetTaskEndTime()
		return executionMapList, err
	}
	// 解析
	for _, e := range eList {
		// 获取代码路径
		var releaseCode comm.ReleaseCode
		comm.DB.First(&releaseCode, e.ReleaseCodeID)
		// 获取agent状态
		var agent comm.Agent
		if comm.CheckRecordWithStringID(e.Service.AgentID,&agent){
			errInfo := fmt.Sprintf("service[%s] ownership agent[%s] not-found",e.ServiceID,e.Service.AgentID)
			return executionMapList, errors.New(errInfo)
		}
		if agent.Status == "0" {
			errInfo := fmt.Sprintf("service[%s] ownership agent[%s] not-online",e.ServiceID,e.Service.AgentID)
			return executionMapList, errors.New(errInfo)
		}

		tmp := CMD{TaskId: e.TaskID,
			ExecutionId:          e.ID,
			ServiceId:            e.ServiceID,
			ServiceOp:            OpMode(e.Operation),
			ServiceName:          e.Service.Name,
			ServiceOsUser:        e.Service.OsUser,
			ServiceOsPass:        e.Service.OsPass,
			ServiceModuleName:    e.Service.ModuleName,
			ServiceDir:           e.Service.Dir,
			ServiceRemoteCode:    releaseCode.RelativePath,
			ServiceCodePattern:   strings.Split(e.Service.CodePatterns, ";"),
			ServiceCustomPattern: strings.Split(e.CustomUpgradePattern, ";"),
			ServicePidfile:       e.Service.Pidfile,
			ServiceStartCmd:      e.Service.StartCMD,
			ServiceStopCmd:       e.Service.StopCMD,
		}

		publishChannel := "cmd." + e.Service.AgentID
		executionMapList = append(executionMapList, map[string]interface{}{"pchannel": publishChannel, "eobject": tmp})
		rre, ok := rrei.(*pb.ExecutionList)
		if ok {
			rre.Executions = append(rre.Executions, &pb.ExecutionList_ExecutionInfo{Id: int32(e.ID), TaskName: e.Task.Name, ServiceName: e.Service.Name, Operation: int32(e.Operation)})
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
			log.Slogger.Errorf("[PublishTask] publish task execution failed. [%s]=>[%s]", e, err)
			prl = append(prl, map[int]error{eid: err})
			for _, r := range rre.Executions {
				if int32(eid) == r.Id {
					r.RCode = 0
					r.RMsg = "publish task execution failed. " + err.Error()
					continue
				}
			}
		} else {
			log.Slogger.Infof("[PublishTask] publish task execution successful. [%d]=>[%s]", e.ExecutionId, string(ebyte))
		}
	}
	return
}


func collectResult(resultChannel chan string, publishSuccessLen int) []string {
	log.Slogger.Infof("[PublishTask] start task results collection goroutine....")
	//收集任务结果
	cmdResult := make([]string, publishSuccessLen)
	// 状态
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < publishSuccessLen; i++ {
		wg.Add(1) // 计数加 1
		go func(i int) {
			cmdResult[i] = <-resultChannel
			defer wg.Done() // 计数减 1
		}(i)
	}
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Slogger.Info("[PublishTask] collect task results success")
	case <-time.After(60 * time.Second):
		log.Slogger.Error("[PublishTask] collect task results partial-timeout")
	}
	log.Slogger.Infof("[PublishTask] Close task results collection goroutine.")
	return cmdResult
}


func result2DB(task comm.Task, results []string, rre *pb.ExecutionList) error {
	// 入库
	var taskstatus = TaskStatus_Successful

	tx := comm.DB.Begin()
	for _, r := range results {
		var executionID int
		var resultMsg string
		rs := Result{}
		if r != ""{
			err := json.Unmarshal([]byte(r), &rs)
			if err != nil {
				return err
			}

		}
		executionID = rs.ExecutionId
		resultCode := rs.ResultCode
		resultMsg = rs.ResultMsg
		resultSteps := rs.ResultSteps

		//设置返回数组中的结果部分
		for _, re := range rre.Executions {
			if int32(executionID) == re.Id {
				re.RCode = int32(resultCode)
				re.RMsg = resultMsg
				break
			}
			// 超时未获取到结果的切片
			if executionID == 0 && re.RMsg == "" {
				executionID = int(re.Id) //切片超时，返回结果为空
				re.RMsg = "operate timeout"
				resultMsg = "timeout"
				break
			}
		}


		// 设置切片结果
		e := comm.Execution{ID: executionID}
		if err := tx.Model(&e).Updates(map[string]interface{}{"result_code": resultCode, "result_msg": resultMsg}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// 插入切片详情
		for _, rsp := range resultSteps {
			eDetail := comm.Execution_Detail{StepNum: rsp.StepNum, StepName: rsp.StepName, StepState: rsp.StepState, StepMsg: rsp.StepMsg, StepTime: time.Unix(rsp.StepTime/1e9, 0), ExecutionID: executionID}
			if err := tx.Create(&eDetail).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if resultCode == TaskStatus_Failed {
			taskstatus = TaskStatus_Failed
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
