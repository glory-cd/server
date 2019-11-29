/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	"errors"
	pb "github.com/glory-cd/server/idlentity"
	"strings"
)

/*添加任务
** 可选参数:
   1. OP: 操作类型. 是全局的，针对任务关联的group下的所有服务。只能为Operate_Start, Operate_Stop, Operate_Restart, Operate_Check，Operate_Upgrade
          如果设置了OP，将会忽略WithDeploy,WithUpgrade,WithStatic中设置的参数
   2. ReleaseID: 发布ID. 如果本次任务设置了WithDeploy,WithUpgrade,则该参数必须设置
   3. Deploys:  部署详细参数。其中定义了ServiceID。ServiceID自然是服务ID，ReleasecodeID是ReleaseID发布中的codeID
   4. Upgrades: 升级详细参数。其中定义了ServiceID和CustomUpgradePattern。ServiceID自然是服务ID，CustomUpgradePattern是本次自定义的要更新的文件或者文件夹
   5. Statics:  其他操作参数详情。其中定义了ServiceID和Op。ServiceID自然是服务ID，Op则是具体操作类型。
*/

func (c *CDPClient) AddTask(taskName string, opts ...Option) (int32, error) {
	taskOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&taskOption)
	}
	releaseID := taskOption.ReleaseID
	op := taskOption.Op
	group := taskOption.GroupName
	deploys := taskOption.Deploys
	upgrades := taskOption.Upgrades
	statics := taskOption.Statics
	taskIsShow := taskOption.TaskIsShow
	// 校验
	//if op == OperateDeploy {
	//	return 0, errors.New("Parameter error: OP is global, for all services in designated group，So it can't be deploy." +
	//		"If you want to achieve the deploy task, use the 'WithDeploy' parameter")
	//}

	if op == OperateDefault && len(deploys) == 0 && len(upgrades) == 0 && len(statics) == 0 {
		return 0, errors.New("Parameter error: neither OP nor detailed parameters are set.")
	}

	if releaseID == 0 && (len(deploys) > 0 || len(upgrades) > 0 || op == OperateUpgrade || op == OperateDeploy) {
		return 0, errors.New("Parameter error: to deploy or upgrade tasks, you must use WithRelease to set the release's name.")
	}

	// 格式化任务详情
	var ss []*pb.SpecificService
	if op == OperateDefault {
		for _, d := range deploys {
			ss = append(ss, &pb.SpecificService{ServiceID: d.ServiceID, Operation: int32(OperateDeploy)})
		}

		for _, u := range upgrades {
			ss = append(ss, &pb.SpecificService{ServiceID: u.ServiceID, Operation: int32(OperateUpgrade), CustomUpgradePattern: strings.Join(u.CustomUpgradePattern, ";")})
		}

		for _, s := range statics {
			if s.Op == OperateDeploy || s.Op == OperateUpgrade {
				return 0, errors.New("Parameter error: static-OpMode not be OperateDeploy and OperateUpgrade.")
			}
			ss = append(ss, &pb.SpecificService{ServiceID: s.ServiceID, Operation: int32(s.Op)})
		}

	} else {
		//获取所有服务信息
		services, err := c.GetServices(WithGroupNames([]string{group}))
		if err != nil {
			return 0, err
		}

		if len(services) == 0 {
			return 0, errors.New("service list is empty.")
		}
		for _, s := range services {
			ss = append(ss, &pb.SpecificService{ServiceID: s.ID, Operation: int32(op)})
		}
	}
	ctx := context.TODO()
	sc := c.newTaskClient()
	// 添加任务，在任务表中插入一条记录
	res, err := sc.AddTask(ctx, &pb.TaskAddRequest{Name: taskName, IsShow: taskIsShow})
	if err != nil {
		return 0, err
	}

	_, err = sc.SetTaskDetails(ctx, &pb.TaskDetailsRequst{TaskID: res.Taskid, ReleaseID: releaseID, Sslist: ss})
	if err != nil {
		return res.Taskid, err
	}

	return res.Taskid, nil
}

func (c *CDPClient) DeleteTask(taskName string) error {
	sc := c.newTaskClient()
	ctx := context.TODO()
	_, err := sc.DeleteTask(ctx, &pb.TaskNameRequest{Name: taskName})
	if err != nil {
		return err
	}
	return nil
}

func (c *CDPClient) GetTasks(opts ...Option) (TaskSlice, error) {
	taskQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&taskQueryOption)
	}

	sc := c.newTaskClient()
	ctx := context.TODO()
	var tasks TaskSlice
	taskList, err := sc.GetTasks(ctx, &pb.GetTaskRequest{Id: taskQueryOption.Ids, Name:taskQueryOption.Names})
	if err != nil {
		return tasks, err
	}
	for _, t := range taskList.Tasks {
		tasks = append(tasks, Task{ID: t.Id, Name: t.Name, Status: TaskStatus[t.Status], StartTime: t.Starttime, EndTime: t.Endtime, CreateTime: t.Ctime, IsShow: t.IsShow})
	}
	return tasks, nil
}

func (c *CDPClient) ExecuteTask(taskID int32) ([]TaskResult, error) {
	sc := c.newTaskClient()
	ctx := context.TODO()
	var tResult []TaskResult
	res, err := sc.PublishTask(ctx, &pb.TaskIdRequest{Id: int32(taskID)})

	if res != nil {
		for _, r := range res.Executions {
			tmp := TaskResult{TaskName: r.TaskName, WorkID: int(r.Id), ServiceName: r.ServiceName, Operation: OpMap[OpMode(r.Operation)], ResultCode: ResultMap[r.RCode], ResultMsg: r.RMsg}
			tResult = append(tResult, tmp)
		}

	}
	return tResult, err
}

/*
	获取任务切片
*/
func (c *CDPClient) GetTaskExecutions(taskID int32) (ExecutionSlice, error) {
	sc := c.newTaskClient()
	ctx := context.TODO()
	es, err := sc.GetTaskExecutions(ctx, &pb.TaskIdRequest{Id: taskID})
	if err != nil {
		return nil, err
	}
	res := ExecutionSlice{}
	for _, e := range es.Executions {
		res = append(res, Execution{TaskID: e.TaskID,
			TaskName:          e.TaskName,
			ServiceName:       e.ServiceName,
			WorkID:            e.Id,
			WorkOp:            OpMap[OpMode(e.Operation)],
			WorkReturnCode:    ResultMap[e.RCode],
			WorkReturnMsg:     e.RMsg,
			WorkCustomPattern: e.CustomUpgradePattern})
	}
	return res, nil
}

/*
	获取切片步骤
*/
func (c *CDPClient) GetTaskExecutionDetails(executionID int32) (ExecutionDetailSlice, error) {
	sc := c.newTaskClient()
	ctx := context.TODO()
	es, err := sc.GetExecutionDetail(ctx, &pb.GetExecutionDetailRequest{ExecutionID: executionID})
	if err != nil {
		return nil, err
	}
	res := ExecutionDetailSlice{}
	for _, e := range es.EDetails {
		res = append(res, ExecutionDetail{StepNum: e.StepNum, StepName: e.StepName, StepMsg: e.StepMsg, StepTime: e.StepTime, StepStatus: e.StepState})
	}
	return res, nil
}


/*
	设置任务为定时任务
*/
func (c *CDPClient) SetTaskToTimed(taksID int32, timedSpec string) (int32, error) {
	sc := c.newTaskClient()
	ctx := context.TODO()
	ts, err := sc.GetTasks(ctx, &pb.GetTaskRequest{Id: []int32{taksID}})
	if err != nil {
		return 0, err
	}
	if len(ts.Tasks) == 0 {
		return 0, errors.New("The task ID is empty.")
	}

	r, err := sc.SetTimedTask(ctx, &pb.CronTaskAddRequest{TaskId: taksID, TimedSpec: timedSpec})
	if err != nil {
		return 0, err
	}
	return r.CronTaskID, err
}

func (c *CDPClient) GetTimedTask(opts ...Option) (CronTaskSlice, error) {
	timedTaskQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&timedTaskQueryOption)
	}

	sc := c.newTaskClient()
	ctx := context.TODO()
	var tts CronTaskSlice

	ttaskList, err := sc.GetTimedTasks(ctx, &pb.GetCronTaskRequest{EntryIDs: timedTaskQueryOption.CronEntryIDs, TaskNames: timedTaskQueryOption.TaskNames})
	if err != nil {
		return tts, err
	}

	for _, t := range ttaskList.TTasks {
		tts = append(tts, CronTask{EntryID: t.EntryId, TaskID: t.TaskId, TaskName: t.TaskName, TaskExecTime: t.TaskExecTIme, CreateTime: t.CTime})
	}
	return tts, nil
}

/*
  删除定时任务
*/
func (c *CDPClient) RemoveTimedTask(entryID int32) error {
	tts, err := c.GetTimedTask(WithCronEntryIds([]int32{entryID}))
	if err != nil {
		return err
	}
	taskId := tts.GetTaskID()

	sc := c.newTaskClient()
	ctx := context.TODO()
	_, err = sc.RemoveTimedTask(ctx, &pb.RemoveCronTaskRequest{TaskID: taskId, EntryID: int32(entryID)})
	return err
}
