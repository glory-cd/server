/**
* @Author: xhzhang
* @Date: 2019-06-21 15:49
 */
package server

import (
	"context"
	"errors"
	"server/comm"
	pb "server/idlentity"
	"strconv"
	"utils/log"
)

type Task struct{}

func (t *Task) AddTask(ctx context.Context, in *pb.TaskAddRequest) (*pb.TaskAddReply, error) {
	var taskObj comm.Task
	if in.Releaseid == 0 {
		taskObj = comm.Task{Name: in.Name, GroupID: int(in.Groupid)}
	} else {
		taskObj = comm.Task{Name: in.Name, GroupID: int(in.Groupid), ReleaseID: int(in.Releaseid)}
	}

	if err := comm.CreateRecord(&taskObj); err != nil {
		log.Slogger.Errorf("[Task] 添加任务[%s]失败: %s", in.Name, err)
		return &pb.TaskAddReply{}, err

	}
	log.Slogger.Infof("[Task] 添加任务[%s]成功", in.Name)
	return &pb.TaskAddReply{Taskid: int32(taskObj.ID)}, nil

}

func (t *Task) DeleteTask(ctx context.Context, in *pb.TaksNameRequest) (*pb.EmptyReply, error) {
	taskObj := comm.Task{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &taskObj) {
		log.Slogger.Errorf("[Task] 删除任务[%s]失败: 不存在,无法删除", in.Name)
		return &pb.EmptyReply{}, errors.New("服务不存在，无法删除")
	}
	err := comm.DeleteRecord(&taskObj)
	if err != nil {
		log.Slogger.Errorf("[Task] 删除任务[%s]失败: %s", taskObj.Name, err)
		return &pb.EmptyReply{}, err

	}

	log.Slogger.Infof("[Task] 删除任务[%s]成功", taskObj.Name)
	return &pb.EmptyReply{}, nil
}

func (t *Task) GetTasks(ctx context.Context, in *pb.EmptyRequest) (*pb.TaskList, error) {
	var tasks []comm.Task
	var rtasks pb.TaskList
	if err := comm.DB.Preload("Group").Preload("Release").Find(&tasks).Error; err != nil {
		return &rtasks, err
	}
	for _, task := range tasks {
		ti := &pb.TaskList_TaskInfo{Id: int32(task.ID),
			Name:        task.Name,
			Status:      int32(task.Status),
			Ctime:       task.CreatedAt.String(),
			Starttime:   task.StartTime.String(),
			Endtime:     task.EndTime.String(),
			Groupname:   task.Group.Name,
			Releasename: task.Release.Name}
		rtasks.Tasks = append(rtasks.Tasks, ti)
	}
	return &rtasks, nil
}

func (t *Task) SetTaskDetails(ctx context.Context, in *pb.TaskDetailsRequst) (*pb.EmptyReply, error) {
	tx := comm.DB.Begin()
	for _, ss := range in.Sslist {
		// 部署时，设置service的moudle_name
		if ss.Operation == 1 {
			var releasecode comm.ReleaseCode
			if err := tx.First(&releasecode, ss.Releasecodeid).Error; err != nil {
				tx.Rollback()
				log.Slogger.Errorf("[SetTaskDetails] 部署任务获取发布代码错误: %s", err)
				return &pb.EmptyReply{}, err
			}
			service := comm.Service{ID: ss.Serviceid}
			if err := tx.Find(&service).Update("module_name", releasecode.Name).Error; err != nil {
				tx.Rollback()
				log.Slogger.Errorf("[SetTaskDetails] 部署任务获设置服务MoudleName失败: %s", err)
				return &pb.EmptyReply{}, err
			}
		}
		execution := comm.Execution{TaskID: int(in.Taskid), ServiceID: ss.Serviceid, Operation: int(ss.Operation), CustomUpgradePattern: ss.Customupgradepattern}

		if tx.Where("task_id = ?", in.Taskid).Where("service_id = ?", ss.Serviceid).Where("operation = ?", ss.Operation).First(&comm.Execution{}).RecordNotFound() {
			if err := tx.Create(&execution).Error; err != nil {
				tx.Rollback()
				return &pb.EmptyReply{}, err
			}
		}

	}
	tx.Commit()
	return &pb.EmptyReply{}, nil
}

// 根据taskid获取所有该任务的所有切片
func (t *Task) GetTaskExecutions(ctx context.Context, in *pb.TaskIdRequest) (*pb.ExecutionList, error) {
	// + 判断任务是否存在
	var relist pb.ExecutionList
	var elist []comm.Execution
	if err := comm.DB.Preload("Task").Preload("Service").Where("task_id = ?", in.Id).Find(&elist).Error; err != nil {
		return &relist, err
	}
	for _, e := range elist {
		_tmpe := &pb.ExecutionList_ExecutionInfo{Id: int32(e.ID),
			Operation:            int32(e.Operation),
			Rcode:                int32(e.ResultCode),
			Rmsg:                 e.ResultMsg,
			Servicename:          e.Service.Name,
			Taskname:             e.Task.Name,
			Customupgradepattern: e.CustomUpgradePattern}
		relist.Executions = append(relist.Executions, _tmpe)
	}
	return &relist, nil

}

func (t *Task) PublishTask(ctx context.Context, in *pb.TaskIdRequest) (*pb.ExecutionList, error) {
	log.Slogger.Infof("[PublishTask] 执行任务[%d]", in.Id)
	var el pb.ExecutionList
	//检查任务是否是可执行状态
	isExecute, err := checkTaskIsExecute(int(in.Id))
	if !isExecute {
		if err != nil {
			return &el, errors.New("任务已经执行过，不能重复执行" + err.Error())
		} else {
			return &el, errors.New("任务已经执行过，不能重复执行")
		}
	}
	//设置任务开始时间
	task := comm.Task{ID: int(in.Id)}
	if err := task.SetTaskStartTime(); err != nil {
		_ = task.SetTaskEndTimeAndStatus(Status_Fail)
		return &el, err // 此时返回，el为空[]
	}
	// 获取任务切片字符串
	eobjectList, err := getExecutions(&task, &el)
	if err != nil {
		return &el, err // 此时返回，el为空[]
	}
	//订阅任务结果通道
	resultChannel := make(chan string)
	subscribChannel := "result." + strconv.Itoa(int(in.Id)) //channel:result.taskid
	go comm.SubscribeCMDResult(subscribChannel, resultChannel)
	//发布任务
	publistErrResult := publishExecutionCMD(eobjectList, &el)
	//成功publish的命令个数
	publishSuccessLen := len(eobjectList) - len(publistErrResult)
	//收集任务结果
	cmdResult := collectResult(resultChannel, publishSuccessLen)

	//入库
	err = result2DB(task, cmdResult, &el)
	if err != nil {
		return &el, err
	}
	return &el, nil
}
