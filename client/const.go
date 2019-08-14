/**
* @Author: xhzhang
* @Date: 2019/7/19 9:26
 */
package client

const (
	KeyOrganizationClient string = "orgClient"
	KeyEnvironmentClient  string = "envClient"
	KeyProjectClient      string = "proClient"
	KeyGroupClient        string = "groupClient"
	KeyReleaseClient      string = "releaseClient"
	KeyServiceClient      string = "serviceClient"
	KeyAgentClient        string = "AgentClient"
	KeyTaskClient         string = "TaskClient"
)

type OpMode int

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
