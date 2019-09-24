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

var OpMap = map[OpMode]string{OperateDefault: "",
	OperateDeploy:   "deploy",
	OperateUpgrade:  "upgrade",
	OperateStart:    "start",
	OperateStop:     "stop",
	OperateRestart:  "restart",
	OperateCheck:    "check",
	OperateBackUp:   "backup",
	OperateRollBack: "rollback"}



//-----------------------------------------------
type Organization struct {
	ID        int32
	Name      string
	CreatTime string
}

type OrganizationSlice []Organization

func (os OrganizationSlice) GetID() int32 {
	if len(os) == 1 {
		return os[0].ID
	} else {
		return 0
	}
}
//------------------------------------------------------
type Environment struct {
	ID        int32
	Name      string
	CreatTime string
}

type EnvironmentSlice []Environment

func (es EnvironmentSlice) GetID() int32 {
	if len(es) == 1 {
		return es[0].ID
	} else {
		return 0
	}
}

//---------------------------------------------------------
type Project struct {
	ID        int32
	Name      string
	CreatTime string
}

type ProjectSlice []Project

func (ps ProjectSlice) GetID() int32 {
	if len(ps) == 1 {
		return ps[0].ID
	} else {
		return 0
	}
}

//---------------------------------------------------------
type Group struct {
	ID           int32
	Name         string
	Organization string
	Environment  string
	Project      string
}

type GroupSlice []Group

func (gs GroupSlice) GetID() int32 {
	if len(gs) == 1 {
		return gs[0].ID
	} else {
		return 0
	}
}

//-----------------------------------------------------------
type Release struct {
	ID           int32
	Name         string
	Version      string
	OrgName      string
	ProName      string
	ReleaseCodes []int32   // releasecode id slice
}

type ReleaseSlice []Release

func (rs ReleaseSlice) GetID() int32 {
	if len(rs) == 1 {
		return rs[0].ID
	} else {
		return 0
	}
}
//--------------------------------------------------------
type ReleaseCode struct {
	ReleaseID int32
	Id        int32
	CodeName  string
	CodePath  string
}
type ReleaseCodeSlice []ReleaseCode

//------------------------------------------------------

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

type ServiceSlice []Service

func (ss ServiceSlice) GetID() string {
	if len(ss) == 1 {
		return ss[0].ID
	} else {
		return ""
	}
}
//-----------------------------------------------------
type Agent struct {
	ID        string
	Alias     string
	Host      string
	Ip        string
	Status    string
	CreatTime string
}

type AgentSlice []Agent

func (as AgentSlice) GetID() string {
	if len(as) == 1 {
		return as[0].ID
	} else {
		return ""
	}
}

//------------------------------------------------------
type Task struct {
	ID          int32
	Name        string
	Status      int32
	StartTime   string
	EndTime     string
	ReleaseName string
	GroupName   string
	CreateTime  string
}

type TaskSlice []Task

func (ts TaskSlice) GetID() int32 {
	if len(ts) == 1 {
		return ts[0].ID
	} else {
		return 0
	}
}

//--------------------------------------------------------

type TaskResult struct {
	TaskName    string
	ExecutionID int
	ServiceName string
	Operation   int32
	Resultcode  int
	Resultmsg   string
}

type Execution struct {
	TaskName       string
	ServiceName    string
	ID             int32
	Op             string
	ReturnCode     int32
	ReturnMsg      string
	CustomePattern string
}

type ExecutionSlice []Execution

//-------------------------------------------------