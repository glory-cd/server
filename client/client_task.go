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


type DeployServiceDeatl struct {
	ServiceID     string
	ReleaseCodeID int32
}

type UpgradeServiceDeatl struct {
	ServiceID            string
	CustomUpgradePattern []string
}

type StaticServiceDeatl struct {
	ServiceID string
	Op        OpMode
}

type taskOption struct {
	ReleaseID int32
	Op        OpMode
	Deploys   []DeployServiceDeatl
	Upgrades  []UpgradeServiceDeatl
	Statics   []StaticServiceDeatl
}

type TaskOption interface {
	apply(*taskOption)
}

type funcOption struct {
	f func(*taskOption)
}

func (fdo *funcOption) apply(do *taskOption) {
	fdo.f(do)
}

func newFuncOption(f func(*taskOption)) *funcOption {
	return &funcOption{f: f}
}

func WithRelease(id int32) TaskOption {
	return newFuncOption(func(o *taskOption) { o.ReleaseID = id })
}

func WithOp(opn OpMode) TaskOption {
	return newFuncOption(func(o *taskOption) { o.Op = opn })
}

func WithDeploy(d []DeployServiceDeatl) TaskOption {
	return newFuncOption(func(o *taskOption) { o.Deploys = d })
}

func WithUpgrade(u []UpgradeServiceDeatl) TaskOption {
	return newFuncOption(func(o *taskOption) { o.Upgrades = u })
}

func WithStatic(s []StaticServiceDeatl) TaskOption {
	return newFuncOption(func(o *taskOption) { o.Statics = s })
}

//默认参数
func defaultOptions() taskOption {
	return taskOption{}
}

/*添加任务
** 可选参数:
   1. OP: 操作类型. 是全局的，针对任务关联的group下的所有服务。只能为Operate_Start, Operate_Stop, Operate_Restart, Operate_Check，Operate_Upgrade
          如果设置了OP，将会忽略WithDeploy,WithUpgrade,WithStatic中设置的参数
   2. ReleaseID: 发布ID. 如果本次任务设置了WithDeploy,WithUpgrade,则该参数必须设置
   3. Deploys:  部署详细参数。其中定义了ServiceID和ReleasecodeID。ServiceID自然是服务ID，ReleasecodeID是ReleaseID发布中的codeID
   4. Upgrades: 升级详细参数。其中定义了ServiceID和CustomUpgradePattern。ServiceID自然是服务ID，CustomUpgradePattern是本次自定义的要更新的文件或者文件夹
   5. Statics:  其他操作参数详情。其中定义了ServiceID和Op。ServiceID自然是服务ID，Op则是具体操作类型。
*/

func (c *CDPClient) AddTask(taskName string, groupName string, opts ...TaskOption) (int32, error) {
	taskOption := defaultOptions()
	for _, opt := range opts {
		opt.apply(&taskOption)
	}
	releaseID := taskOption.ReleaseID
	op := taskOption.Op
	deploys := taskOption.Deploys
	upgrades := taskOption.Upgrades
	statics := taskOption.Statics

	if op == OperateDeploy {
		return 0, errors.New("Parameter error: OP is global, for all services in designated group，So it can't be [deploy], " +
			"its [1].If you want to achieve the deploy task, use the 'WithDeploy' parameter")
	}

	if op == OperateDefault && len(deploys) == 0 && len(upgrades) == 0 && len(statics) == 0 {
		return 0, errors.New("Parameter error: neither OP nor detailed parameters are set.")
	}

	if releaseID == 0 && (len(deploys) > 0 || len(upgrades) > 0 || op == OperateUpgrade) {
		return 0, errors.New("Parameter error: to deploy or upgrade tasks, you must use WithRelease to set the release's name.")
	}

	ctx := context.TODO()
	// 获取分组ID
	groups,err := c.GetGroups(WithNames([]string{groupName}))
	if err != nil{
		return 0,err
	}
	groupID := groups.GetID()

	// 获取所有服务信息
	services, err := c.GetServices(WithGroups([]string{groupName}))
	if err != nil {
		return 0, err
	}

	if len(services) == 0{
		return 0,errors.New("service list is empty.")
	}

	sc := c.newTaskClient()
	// 添加任务，在任务表中插入一条记录
	res, err := sc.AddTask(ctx, &pb.TaskAddRequest{Name: taskName, Groupid: groupID, Releaseid: int32(releaseID)})
	if err != nil {
		return 0, err
	}
	// 设置任务详情
	var ss []*pb.SpecificService
	switch op {
	case OperateDeploy:
		return 0, errors.New("Parameter error: OP is global, for all services in designated group，So it can't be [deploy],its [1]")
	case OperateUpgrade:
		for _, s := range services {
			ss = append(ss, &pb.SpecificService{Serviceid: s.ID, Operation: int32(OperateUpgrade)})
		}
	case OperateStart, OperateStop, OperateRestart, OperateCheck, OperateBackUp, OperateRollBack:
		for _, s := range services {
			ss = append(ss, &pb.SpecificService{Serviceid: s.ID, Operation: int32(op)})
		}
	// 当op没有设置时，根据deploys,upgrades,statics信息设置
	case OperateDefault:
		if deploys != nil {
			var serviceIDList []string
			var releaseCodeIDList []int32
			for _, d := range deploys {
				serviceIDList = append(serviceIDList, d.ServiceID)
				releaseCodeIDList = append(releaseCodeIDList, d.ReleaseCodeID)
			}
			// 验证服务ID
			err := CheckServiceOwnGroup(serviceIDList, services)
			if err != nil {
				return 0, err
			}
			// 验证发布代码
			rCodeParentMap, err := c.GetReleaseCodeMap(releaseID)
			if err != nil {
				return 0, err
			}
			err = CheckRCIDOwnRelease(releaseCodeIDList, rCodeParentMap)
			if err != nil {
				return 0, err
			}
		}
		// 校验upgrades
		if upgrades !=nil {
			var serviceIDList []string
			for _, u := range upgrades {
				serviceIDList = append(serviceIDList, u.ServiceID)
			}
			// 验证服务ID
			err := CheckServiceOwnGroup(serviceIDList, services)
			if err != nil {
				return 0, err
			}

			// 验证发布代码
			upgradeServices, err := c.GetServices(WithServices(serviceIDList))
			var upgradeServiceMoudleNames []string
			for _, us := range upgradeServices {
				upgradeServiceMoudleNames = append(upgradeServiceMoudleNames, us.MoudleName)
			}
			rCodeParentMap, err := c.GetReleaseCodeMap(releaseID)
			if err != nil {
				return 0, err
			}

			err = CheckRCNameOwnRelease(upgradeServiceMoudleNames, rCodeParentMap)
			if err != nil {
				return 0, err
			}

		}

		// 校验static
		if statics != nil {
			var serviceIDList []string
			for _, s := range statics {
				serviceIDList = append(serviceIDList, s.ServiceID)
				//验证OpMode
				if s.Op == OperateDeploy || s.Op == OperateUpgrade {
					return 0, errors.New("参数错误: 静态模式的OpMode不能是OperateDeploy和OperateUpgrade.")
				}

			}
			// 验证服务ID
			err := CheckServiceOwnGroup(serviceIDList, services)
			if err != nil {
				return 0, err
			}
		}

		for _, d := range deploys {
			ss = append(ss, &pb.SpecificService{Serviceid: d.ServiceID, Operation: int32(OperateDeploy), Releasecodeid: int32(d.ReleaseCodeID)})
		}

		for _, u := range upgrades {
			ss = append(ss, &pb.SpecificService{Serviceid: u.ServiceID, Operation: int32(OperateUpgrade), Customupgradepattern: strings.Join(u.CustomUpgradePattern, ";")})
		}

		for _, s := range statics {
			ss = append(ss, &pb.SpecificService{Serviceid: s.ServiceID, Operation: int32(s.Op)})
		}
	}

	_, err = sc.SetTaskDetails(ctx, &pb.TaskDetailsRequst{Taskid: res.Taskid, Sslist: ss})
	if err != nil {
		return res.Taskid, err
	}

	return res.Taskid, nil
}

func (c *CDPClient) DeleteTask(taskName string) error {
	sc := c.newTaskClient()
	ctx := context.TODO()
	_, err := sc.DeleteTask(ctx, &pb.TaksNameRequest{Name: taskName})
	if err != nil {
		return err
	}
	return nil
}

func (c *CDPClient) GetTasks(opts ...QueryOption) (TaskSlice, error) {
	taskQueryOption := defaultQueryOption()
	for _, opt := range opts {
		opt.apply(&taskQueryOption)
	}

	sc := c.newTaskClient()
	ctx := context.TODO()
	var tasks TaskSlice
	tasklist, err := sc.GetTasks(ctx, &pb.GetTaskRequest{Id: taskQueryOption.Ids,
		Name:    taskQueryOption.Names,
		Release: taskQueryOption.ReleaseNames,
		Group:   taskQueryOption.GroupNames})
	if err != nil {
		return tasks, err
	}
	for _, t := range tasklist.Tasks {
		tasks = append(tasks, Task{ID: t.Id, Name: t.Name, Status: t.Status, StartTime: t.Starttime, EndTime: t.Endtime, GroupName: t.Groupname, ReleaseName: t.Releasename, CreateTime: t.Ctime})
	}
	return tasks, nil
}

func (c *CDPClient) ExecuteTask(taskID int32) ([]TaskResult, error) {
	sc := c.newTaskClient()
	ctx := context.TODO()

	var tresult []TaskResult

	res, err := sc.PublishTask(ctx, &pb.TaskIdRequest{Id: int32(taskID)})
	if err != nil {
		return tresult, err
	}

	for _, r := range res.Executions {
		tmp := TaskResult{TaskName: r.Taskname, ExecutionID: int(r.Id), ServiceName: r.Servicename, Operation: r.Operation, Resultcode: int(r.Rcode), Resultmsg: r.Rmsg}
		tresult = append(tresult, tmp)
	}
	return tresult, nil
}

func (c *CDPClient) GetTaskExecutions(taskID int32) (ExecutionSlice, error) {
	sc := c.newTaskClient()
	ctx := context.TODO()
	es, err := sc.GetTaskExecutions(ctx, &pb.TaskIdRequest{Id: taskID})
	if err != nil {
		return nil, err
	}
	res := ExecutionSlice{}
	for _, e := range es.Executions {
		res = append(res, Execution{TaskName: e.Taskname, ServiceName: e.Servicename, ID: e.Id, Op: OpMap[OpMode(e.Operation)], ReturnCode: e.Rcode, ReturnMsg: e.Rmsg})
	}
	return res, nil
}
