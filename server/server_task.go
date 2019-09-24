/**
* @Author: xhzhang
* @Date: 2019-06-21 15:49
 */
package server

import (
	"context"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
	"strconv"
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
		log.Slogger.Errorf("[Task] add [%s] failed. %s", in.Name, err)
		return &pb.TaskAddReply{}, err

	}
	log.Slogger.Infof("[Task] add [%s] successful.", in.Name)
	return &pb.TaskAddReply{Taskid: int32(taskObj.ID)}, nil

}

func (t *Task) DeleteTask(ctx context.Context, in *pb.TaksNameRequest) (*pb.EmptyReply, error) {
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
	queryCmd := comm.DB.Preload("Group").Preload("Release")
	if in.Id != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Id)
	}

	if in.Name != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Name)
	}

	if in.Release != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_releases ON cdp_release.id = cdp_tasks.release_id AND cdp_release.name in (?)", in.Release)
	}

	if in.Group != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_groups ON cdp_groups.id = cdp_tasks.group_id AND cdp_groups.name in (?)", in.Group)
	}

	if err := queryCmd.Find(&tasks).Error; err != nil {
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
				log.Slogger.Errorf("[SetTaskDetails] can not find releasecode [%d]. %s", ss.Releasecodeid,err)
				return &pb.EmptyReply{}, err
			}
			service := comm.Service{ID: ss.Serviceid}
			if err := tx.Find(&service).Update("module_name", releasecode.Name).Error; err != nil {
				tx.Rollback()
				log.Slogger.Errorf("[SetTaskDetails] deploy task set service moudle failed. %s", err)
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
	log.Slogger.Infof("[PublishTask] exec task [%d]", in.Id)
	var el pb.ExecutionList
	//检查任务是否是可执行状态
	isExecute, err := checkTaskIsExecute(int(in.Id))
	if !isExecute {
		if err != nil {
			return &el, errors.New("the task has executed，can not repeated execution. " + err.Error())
		} else {
			return &el, errors.New("the task has executed，can not repeated execution. ")
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
