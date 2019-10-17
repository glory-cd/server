/**
* @Author: xhzhang
* @Date: 2019/10/11 10:19
 */
package client

type DeployServiceDetail struct {
	ServiceID     string
	ReleaseCodeID int32
}

type UpgradeServiceDetail struct {
	ServiceID            string
	CustomUpgradePattern []string
}

type StaticServiceDetail struct {
	ServiceID string
	Op        OpMode
}

type option struct {
	//id
	GroupID       int32
	ReleaseID     int32
	AgentID       string
	//ids
	Ids           []int32
	GroupIDs      []int32
	AgentIDs      []string
	ServiceIDs    []string
	CronEntryIDs  []int32
	//name
	OrgName       string
	ProName       string
	EnvName       string
	GroupName     string
	//names
	Names         []string
	OrgNames      []string
	ProNames      []string
	EnvNames      []string
	GroupNames    []string
	TaskNames     []string
	ReleaseNames  []string
	ModuleNames   []string
	//other
	AgentIsOnLine bool
	CodePattern   string
	PidFile       string
	StopCmd       string
	Op            OpMode
	Deploys       []DeployServiceDetail
	Upgrades      []UpgradeServiceDetail
	Statics       []StaticServiceDetail
	TaskIsShow    bool
}

type Option interface {
	apply(*option)
}

type funcOptionA struct {
	f func(*option)
}

func (fdo *funcOptionA) apply(do *option) {
	fdo.f(do)
}

func newFuncOptionA(f func(*option)) *funcOptionA {
	return &funcOptionA{f: f}
}

func defaultOption() option {
	return option{}
}

//id
func WithAgentId(id string) Option {
	return newFuncOptionA(func(o *option) { o.AgentID = id })
}

func WithGroupId(id int32) Option {
	return newFuncOptionA(func(o *option) { o.GroupID = id })
}

func WithReleaseId(id int32) Option {
	return newFuncOptionA(func(o *option) { o.ReleaseID = id })
}

//ids
func WithInt32Ids(ids []int32) Option {
	return newFuncOptionA(func(o *option) { o.Ids = ids })
}

func WithGroupIds(ids []int32) Option {
	return newFuncOptionA(func(o *option) { o.GroupIDs = ids })
}

func WithAgentIds(ids []string) Option {
	return newFuncOptionA(func(o *option) { o.AgentIDs = ids })
}

func WithServiceIds(ids []string) Option {
	return newFuncOptionA(func(o *option) { o.ServiceIDs = ids })
}

//name
func WithOrgName(name string) Option {
	return newFuncOptionA(func(o *option) { o.OrgName = name })
}

func WithProName(name string) Option {
	return newFuncOptionA(func(o *option) { o.ProName = name })
}

func WithEnvName(name string) Option {
	return newFuncOptionA(func(o *option) { o.EnvName = name })
}

func WithGroupName(name string) Option {
	return newFuncOptionA(func(o *option) { o.GroupName = name })
}

//names
func WithNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.Names = names })
}

func WithOrgNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.OrgNames = names })
}

func WithProNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.ProNames = names })
}

func WithEnvNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.EnvNames = names })
}

func WithGroupNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.GroupNames = names })
}

func WithReleaseNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.ReleaseNames = names })
}

func WithModuleNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.ModuleNames = names })
}

func WithAgentStatus(status bool) Option {
	return newFuncOptionA(func(o *option) { o.AgentIsOnLine = status })
}

func WithCronEntryIds(ids []int32) Option {
	return newFuncOptionA(func(o *option) { o.CronEntryIDs = ids })
}

func WithTaskNames(names []string) Option {
	return newFuncOptionA(func(o *option) { o.TaskNames = names })
}

//other
func WithCodePattern(cPattern string) Option {
	return newFuncOptionA(func(o *option) { o.CodePattern = cPattern })
}

func WithPidFile(pFile string) Option {
	return newFuncOptionA(func(o *option) { o.PidFile = pFile })
}

func WithStopCmd(scmd string) Option {
	return newFuncOptionA(func(o *option) { o.StopCmd = scmd })
}

func WithTaskOp(opn OpMode) Option {
	return newFuncOptionA(func(o *option) { o.Op = opn })
}

func WithTaskDeploy(d []DeployServiceDetail) Option {
	return newFuncOptionA(func(o *option) { o.Deploys = d })
}

func WithTaskUpgrade(u []UpgradeServiceDetail) Option {
	return newFuncOptionA(func(o *option) { o.Upgrades = u })
}

func WithTaskStatic(s []StaticServiceDetail) Option {
	return newFuncOptionA(func(o *option) { o.Statics = s })
}

func WithTaskShow(show bool) Option {
	return newFuncOptionA(func(o *option) { o.TaskIsShow = show })
}
