/**
* @Author: xhzhang
* @Date: 2019/7/19 9:26
 */
package client

const (
	Key_OrganizationClient string = "orgClient"
	Key_EnvironmentClient  string = "envClient"
	Key_ProjectClient      string = "proClient"
	Key_GroupClient        string = "groupClient"
	Key_ReleaseClient      string = "releaseClient"
	Key_ServiceClient      string = "serviceClient"
	Key_AgentClient        string = "AgentClient"
	Key_TaskClient         string = "TaskClient"
)

type OpMode int

const (
	Operate_Default OpMode = 0
	Operate_Deploy  OpMode = 1
	Operate_Upgrade OpMode = 2
	Operate_Start   OpMode = 3
	Operate_Stop    OpMode = 4
	Operate_Restart OpMode = 5
	Operate_Check   OpMode = 6
)

type Organization struct {
	ID        int32
	Name      string
	CreatTime string
}

type Environment struct {
	ID        int32
	Name      string
	CreatTime string
}

type Project struct {
	ID        int32
	Name      string
	CreatTime string
}

type Group struct {
	ID           int32
	Name         string
	Organization string
	Environment  string
	Project      string
}

type ReleaseCode struct {
	CodeName string
	CodePath string
}

type Service struct {
	ID          string
	Name        string
	Dir         string
	MoudleName  string
	OsUser      string
	CodePattern string
	PidFile     string
	StartCmd    string
	StopCmd     string
	AgentName   string
	GroupName   string
}

type Agent struct {
	ID        string
	Alias     string
	Host      string
	Ip        string
	Status    string
	CreatTime string
}

type TaskResult struct {
	TaskName    string
	ExecutionID int
	ServiceName string
	Operation   int32
	Resultcode  int
	Resultmsg   string
}
