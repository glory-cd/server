/**
* @Author: xhzhang
* @Date: 2019-06-21 15:49
 */
package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
	"strconv"
)

type Task struct{}

func (t *Task) AddTask(ctx context.Context, in *pb.TaskAddRequest) (*pb.TaskAddReply, error) {
	var taskObj comm.Task
	taskObj = comm.Task{Name: in.Name}

	if err := comm.CreateRecord(&taskObj); err != nil {
		log.Slogger.Errorf("[Task] add [%s] failed. %s", in.Name, err)
		return &pb.TaskAddReply{}, err

	}
	log.Slogger.Infof("[Task] add [%s] successful.", in.Name)
	return &pb.TaskAddReply{Taskid: int32(taskObj.ID)}, nil
}

func (t *Task) DeleteTask(ctx context.Context, in *pb.TaskNameRequest) (*pb.EmptyReply, error) {
	taskObj := comm.Task{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &taskObj) {
		log.Slogger.Errorf("[Task] delete [%s] failed. not-exist", in.Name)
		return &pb.EmptyReply{}, errors.New("not-exist")
	}
	err := comm.DeleteRecord(&taskObj)
	if err != nil {
		log.Slogger.Errorf("[Task] delete [%s] failed. %s", taskObj.Name, err)
		return &pb.EmptyReply{}, err

	}

	log.Slogger.Infof("[Task] delete [%s] successful.", taskObj.Name)
	return &pb.EmptyReply{}, nil
}

func (t *Task) GetTasks(ctx context.Context, in *pb.GetTaskRequest) (*pb.TaskList, error) {
	var tasks []comm.Task
	var rtasks pb.TaskList
	queryCmd := comm.DB
	if in.Id != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Id)
	}

	if in.Name != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Name)
	}

	if err := queryCmd.Find(&tasks).Error; err != nil {
		return &rtasks, err
	}
	for _, task := range tasks {
		ti := &pb.TaskList_TaskInfo{Id: int32(task.ID),
			Name:      task.Name,
			Status:    int32(task.Status),
			Ctime:     task.CreatedAt.String(),
			Starttime: task.StartTime.String(),
			Endtime:   task.EndTime.String()}
		rtasks.Tasks = append(rtasks.Tasks, ti)
	}
	return &rtasks, nil
}

func (t *Task) SetTaskDetails(ctx context.Context, in *pb.TaskDetailsRequst) (*pb.EmptyReply, error) {
	// 校验: 任务是否存在
	var task comm.Task
	if err := comm.DB.First(&task, in.TaskID).Error; err != nil {
		errInfo := fmt.Sprintf("[SetTaskDetails]  check task [%d] error. %s", in.TaskID, err)
		log.Slogger.Error(errInfo)
		return &pb.EmptyReply{}, errors.New(errInfo)
	}

	var releaseCodes []comm.ReleaseCode
	rcMap := map[string]int{}
	if err := comm.DB.Where("release_id = ?", in.ReleaseID).Find(&releaseCodes).Error; err != nil {
		return &pb.EmptyReply{}, err
	}
	for _, rc := range releaseCodes {
		rcMap[rc.Name] = rc.ID
	}
	// 校验: 服务是否存在,
	//      服务是否属于同一group,
	//      服务是deploy,upgrade,校验代码是否存在
	serviceModuleMap := map[string]string{}
	sameGroupId := 0
	for _, ss := range in.Sslist {
		var service comm.Service
		if notOk := comm.CheckRecordWithStringID(ss.ServiceID, &service); notOk {
			errInfo := fmt.Sprintf("[SetTaskDetails]  check service [%s] not found", ss.ServiceID)
			log.Slogger.Error(errInfo)
			return &pb.EmptyReply{}, errors.New(errInfo)
		}
		serviceModuleMap[service.ID] = service.ModuleName

		if sameGroupId == 0 {
			sameGroupId = service.GroupID
		} else {
			if sameGroupId != service.GroupID {
				errInfo := fmt.Sprintf("[SetTaskDetails] check group is different.")
				log.Slogger.Error(errInfo)
				return &pb.EmptyReply{}, errors.New(errInfo)
			}
		}
		// 根据服务模块名称校验
		if OpMode(ss.Operation) == OperateDeploy || OpMode(ss.Operation) == OperateUpgrade {
			if in.ReleaseID == 0 {
				err := errors.New("[SetTaskDetails] operate is deploy or upgrade ,task's release must appoint.")
				log.Slogger.Errorf("%s", err)
				return &pb.EmptyReply{}, err
			}
			if _, ok := rcMap[service.ModuleName]; !ok {
				log.Slogger.Errorf("[SetTaskDetails] check module failed.")
				return &pb.EmptyReply{}, errors.New("check module failed.")
			}
		}
	}

	tx := comm.DB.Begin()
	for _, ss := range in.Sslist {
		moduleName := serviceModuleMap[ss.ServiceID]
		// 加入到切片
		if tx.Where("task_id = ?", in.TaskID).Where("service_id = ?", ss.ServiceID).Where("operation = ?", ss.Operation).First(&comm.Execution{}).RecordNotFound() {
			execution := comm.Execution{TaskID: int(in.TaskID), ServiceID: ss.ServiceID, Operation: int(ss.Operation), CustomUpgradePattern: ss.CustomUpgradePattern, ReleaseCodeID: rcMap[moduleName]}
			if err := tx.Create(&execution).Error; err != nil {
				tx.Rollback()
				return &pb.EmptyReply{}, err
			}
		}

	}
	tx.Commit()
	return &pb.EmptyReply{}, nil
}

// 根据task-id获取所有该任务的所有切片
func (t *Task) GetTaskExecutions(ctx context.Context, in *pb.TaskIdRequest) (*pb.ExecutionList, error) {
	var rEList pb.ExecutionList
	var eList []comm.Execution
	if err := comm.DB.Preload("Task").Preload("Service").Where("task_id = ?", in.Id).Find(&eList).Error; err != nil {
		return &rEList, err
	}
	for _, e := range eList {
		_tmpE := &pb.ExecutionList_ExecutionInfo{Id: int32(e.ID),
			Operation:            int32(e.Operation),
			RCode:                int32(e.ResultCode),
			RMsg:                 e.ResultMsg,
			ServiceName:          e.Service.Name,
			TaskName:             e.Task.Name,
			TaskID:               int32(e.TaskID),
			CustomUpgradePattern: e.CustomUpgradePattern}
		rEList.Executions = append(rEList.Executions, _tmpE)
	}
	return &rEList, nil
}

// 根据execution-id 获取该切片的步骤详情
func (t *Task) GetExecutionDetail(ctx context.Context, in *pb.GetExecutionDetailRequest) (*pb.ExecutionDetailsList, error) {
	var rEDetailList pb.ExecutionDetailsList
	var eDetails []comm.Execution_Detail
	if err := comm.DB.Preload("Execution").Where("execution_id = ?", in.ExecutionID).Find(&eDetails).Error; err != nil {
		return &rEDetailList, err
	}

	for _, ed := range eDetails {
		_tmp := &pb.ExecutionDetailsList_ExecutionDetail{StepNum: int32(ed.StepNum), StepName: ed.StepName, StepMsg: ed.StepMsg, StepState: int32(ed.StepState), StepTime: ed.StepTime.String()}
		rEDetailList.EDetails = append(rEDetailList.EDetails, _tmp)
	}
	return &rEDetailList, nil
}


func (t *Task) PublishTask(ctx context.Context, in *pb.TaskIdRequest) (*pb.ExecutionList, error) {
	log.Slogger.Infof("[PublishTask] exec task[%d]", in.Id)
	var el pb.ExecutionList
	//检查任务是否是可执行状态. 当任务状态是已经执行成功(1)或者正在执行(4)则无法再次执行此任务.
	isExecute, taskStatus, err := checkTaskIsExecute(int(in.Id))
	if !isExecute {
		return &el, err
	}
	//设置任务开始时间和running状态. 设置失败则设置其结束时间和失败状态
	task := comm.Task{ID: int(in.Id)}
	if err := task.SetTaskStartTimeAndRuningStatus(); err != nil {
		errInfo := fmt.Sprintf("[PublishTask] set task[%d] start-time and running-status failed. %v", in.Id, err)
		log.Slogger.Error(errInfo)
		_ = task.SetTaskEndTimeAndStatus(TaskStatus_Failed)
		return &el, errors.New(errInfo) // 此时返回，el为空[]
	}

	// 获取任务切片字符串
	eObjectList, err := getExecutions(&task, taskStatus, &el)
	if err != nil {
		return &el, err // 此时返回，el为空[]
	}
	//订阅任务结果通道
	resultChannel := make(chan string)
	subscribeChannel := "result." + strconv.Itoa(int(in.Id)) //channel:result.taskId
	go comm.SubscribeCMDResult(subscribeChannel, resultChannel)
	//发布任务
	publishErrResult := publishExecutionCMD(eObjectList, &el)
	//成功publish的命令个数
	publishSuccessLen := len(eObjectList) - len(publishErrResult)
	//收集任务结果
	cmdResult := collectResult(resultChannel, publishSuccessLen)
	//入库
	err = result2DB(task, cmdResult, &el)
	// 返回结果
	if err != nil{
		return &el,err
	}
	return &el,nil
}

/*
	设置任务为定时任务
    1. 检查Cron_Task中是否已经有记录.如果任务已经调度，则返回
    2. 检查Task中是否存在任务记录. 如果不存在，则返回，存在则将任务状态给为TaskStatus_Timed
    3.
*/
func (t *Task) SetTimedTask(ctx context.Context, in *pb.CronTaskAddRequest) (*pb.CronTaskAddReply, error) {
	taskID := int(in.TaskId)
	// check cron-task
	cronTaskObj := comm.Cron_Task{TaskID: taskID}
	if ! cronTaskObj.CheckRecord() {
		return &pb.CronTaskAddReply{}, errors.New("The task is already scheduled.")
	}
	// check task
	var taskObj comm.Task
	if comm.CheckRecordWithID(int(in.TaskId), &taskObj) {
		return &pb.CronTaskAddReply{}, errors.New("This task does not exist.")
	}
	err := taskObj.SetTaskStatus(TaskStatus_Cron)
	if err != nil {
		return &pb.CronTaskAddReply{}, err
	}
	// add-cron
	cronid, err := comm.AddCronTask(taskID, in.TimedSpec)
	// add-cron failed，rollback task.status to TaskStatus_NoExecute
	if err != nil {
		_ = taskObj.SetTaskStatus(TaskStatus_NoExecute)
		return &pb.CronTaskAddReply{}, err
	}
	// add-cron success,add cron-task record.
	taskTimedObj := comm.Cron_Task{TaskID: taskID, TimeSpec: in.TimedSpec, EntryID: cronid}
	err = comm.CreateRecord(&taskTimedObj)
	if err != nil {
		log.Slogger.Errorf(" add cron-task record to db failed. %s", err.Error())

	}
	return &pb.CronTaskAddReply{CronTaskID: int32(cronid)}, nil
}

func (t *Task) RemoveTimedTask(ctx context.Context, in *pb.RemoveCronTaskRequest) (*pb.EmptyReply, error) {
	comm.RemoveEntryFromCron(in.EntryID)

	taskID := int(in.TaskID)
	cronTaskObj := comm.Cron_Task{TaskID: taskID}
	err := comm.DeleteRecord(&cronTaskObj)
	if err != nil {
		return &pb.EmptyReply{}, errors.New("[cron-task] delete task faied. " + err.Error())
	}

	return &pb.EmptyReply{}, nil
}

func (t *Task) GetTimedTasks(ctx context.Context, in *pb.GetCronTaskRequest) (*pb.CronTaskList, error) {
	var timedTasks []comm.Cron_Task
	var timedTaskList pb.CronTaskList
	queryCMD := comm.DB.Preload("Task")
	if len(in.EntryIDs) > 0 {
		queryCMD = queryCMD.Where("entry_id in (?)", in.EntryIDs)
	}

	if len(in.TaskNames) > 0 {
		queryCMD = queryCMD.Joins("JOIN cdp_tasks ON cdp_tasks.id = cdp_task_timeds.task_id AND cdp_task.name in (?)", in.TaskNames)
	}

	if err := queryCMD.Find(&timedTasks).Error; err != nil {
		return &timedTaskList, err
	}

	for _, tt := range timedTasks {
		timedTaskList.TTasks = append(timedTaskList.TTasks, &pb.CronTaskList_CronTask{EntryId: int32(tt.EntryID), TaskId: int32(tt.TaskID), TaskName: tt.Task.Name, TaskExecTIme: tt.TimeSpec, CTime: tt.CreatedAt.String()})
	}
	return &timedTaskList, nil
}
