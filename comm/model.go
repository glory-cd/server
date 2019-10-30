/**
* @Author: xhzhang
* @Date: 2019/7/12 17:25
 */
package comm

import (
	"time"
)

type Organization struct {
	ID        int       `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Name      string    `gorm:"type:varchar(128);UNIQUE;NOT NULL"`
	CreatedAt time.Time `gorm:"column:ctime;NOT NULL"`
}

type Environment struct {
	ID        int       `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Name      string    `gorm:"type:varchar(128);UNIQUE;NOT NULL"`
	CreatedAt time.Time `gorm:"column:ctime;NOT NULL"`
}

type Project struct {
	ID        int       `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Name      string    `gorm:"type:varchar(128);UNIQUE;NOT NULL"`
	CreatedAt time.Time `gorm:"column:ctime;NOT NULL"`
}

type Group struct {
	ID             int          `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Name           string       `gorm:"type:varchar(128);UNIQUE;NOT NULL"`
	CreatedAt      time.Time    `gorm:"column:ctime;NOT NULL"`
	OrganizationID int          `gorm:"column:organization_id;type:int;NOT NULL;DEFAULT:1"`
	EnvironmentID  int          `gorm:"column:environment_id;type:int;NOT NULL;DEFAULT:1"`
	ProjectID      int          `gorm:"column:project_id;type:int;NOT NULL;DEFAULT:1"`
	Organization   Organization `gorm:"FOREIGNKEY:OrganizationID;ASSOCIATION_FOREIGNKEY:ID"` //one-to-one
	Environment    Environment  `gorm:"FOREIGNKEY:EnvironmentID;ASSOCIATION_FOREIGNKEY:ID"`
	Project        Project      `gorm:"FOREIGNKEY:ProjectID;ASSOCIATION_FOREIGNKEY:ID"`
}

type Release struct {
	ID             int          `gorm:"type:int AUTO_INCREMENT;PRIMARY_KEY"`
	Name           string       `gorm:"type:varchar(128);UNIQUE;NOT NULL"`
	Version        string       `gorm:"type:varchar(32);NOT NULL"`
	CreatedAt      time.Time    `gorm:"column:ctime;NOT NULL"`
	OrganizationID int          `gorm:"column:organization_id;type:int;NOT NULL;DEFAULT:1"`
	ProjectID      int          `gorm:"column:project_id;type:int;NOT NULL;DEFAULT:1"`
	Organization   Organization `gorm:"FOREIGNKEY:OrganizationID;ASSOCIATION_FOREIGNKEY:ID"`
	Project        Project      `gorm:"FOREIGNKEY:ProjectID;ASSOCIATION_FOREIGNKEY:ID"`
}

type ReleaseCode struct {
	ID           int     `gorm:"type:int AUTO_INCREMENT;PRIMARY_KEY"`
	Name         string  `gorm:"type:varchar(128);NOT NULL"`
	RelativePath string  `gorm:"column:relative_path;type:varchar(1024);NOT NULL"`
	ReleaseID    int     `gorm:"column:release_id;type:int;NOT NULL"`
	Release      Release `gorm:"FOREIGNKEY:ReleaseID;ASSOCIATION_FOREIGNKEY:ID"`
}

type Agent struct {
	ID        string    `gorm:"column:id;type:varchar(36);NOT NULL;PRIMARY_KEY"`
	Alias     string    `gorm:"column:alias;type:varchar(32)"`
	HostName  string    `gorm:"column:hostname;type:varchar(128);NOT NULL"`
	HostIp    string    `gorm:"column:hostip;type:varchar(128);NOT NULL"`
	Status    string    `gorm:"column:status;type:char(1);NOT NULL;default:'1'"`
	CreatedAt time.Time `gorm:"column:rtime;NOT NULL"` //agent第一次注册时间
	UpdatedAt time.Time `gorm:"column:ctime;NOT NULL"` //agent最近注册时间
}

type Service struct {
	ID           string   `gorm:"column:id;type:varchar(32);NOT NULL;PRIMARY_KEY" json:"serviceid"`
	Name         string   `gorm:"column:name;type:varchar(128);NOT NULL" json:"servicename"`
	Dir          string   `gorm:"column:dir;type:varchar(1024);NOT NULL" json:"servicedir"`
	ModuleName   string   `gorm:"column:module_name;type:varchar(128);NOT NULL" json:"servicemodulename"`
	OsUser       string   `gorm:"column:os_user;type:varchar(128);NOT NULL" json:"serviceosuser"`
	OsPass       string   `gorm:"column:os_pass;type:varchar(128);NOT NULL" json:"serviceospass"`
	CodePattern  []string `gorm:"-" json:"servicecodepattern"`
	CodePatterns string   `gorm:"column:codes;type:varchar(1000)"`
	Pidfile      string   `gorm:"column:pid_file" json:"servicepidfile"`
	StartCMD     string   `gorm:"column:start_cmd;type:varchar(128)" json:"servicestartcmd"`
	StopCMD      string   `gorm:"column:stop_cmd;type:varchar(128)" json:"servicestopcmd"`
	AgentID      string   `gorm:"column:agent_id;NOT NULL" json:"agentid"`
	GroupID      int      `gorm:"column:group_id;type:int;NOT NULL;DEFAULT:1" json:"groupid"`
	Agent        Agent    `gorm:"FOREIGNKEY:AgentID;ASSOCIATION_FOREIGNKEY:ID"`
	Group        Group    `gorm:"FOREIGNKEY:GroupID;ASSOCIATION_FOREIGNKEY:ID"`

}

type Task struct {
	ID        int       `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Name      string    `gorm:"column:name;type:varchar(128);UNIQUE;NOT NULL"`
	Status    int       `gorm:"column:status;type:int;NOT NULL;DEFAULT:2"` //0: 执行失败; 1:执行成功; 2:未执行(默认值); 3:定时任务; 4:正在执行
	CreatedAt time.Time `gorm:"column:ctime"`
	StartTime time.Time `gorm:"column:start_time"`
	EndTime   time.Time `gorm:"column:end_time"`
	IsShow    bool      `gorm:"type:bool;column:is_show;NOT NULL;DEFAULT:false"`
}

type Execution struct {
	ID                   int         `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	Operation            int         `gorm:"column:operation;type:int;NOT NULL"`
	ResultCode           int         `gorm:"column:result_code;type:int"`
	ResultMsg            string      `gorm:"column:result_msg;type:varchar(1024)"`
	TaskID               int         `gorm:"column:task_id;type:int;NOT NULL"`
	ServiceID            string      `gorm:"column:service_id;type:varchar(32);NOT NULL"`
	Task                 Task        `gorm:"FOREIGNKEY:TaskID;ASSOCIATION_FOREIGNKEY:ID"`
	Service              Service     `gorm:"FOREIGNKEY:ServiceID;ASSOCIATION_FOREIGNKEY:ID"`
	CustomUpgradePattern string      `gorm:"column:custom_upgradepattern;type:varchar(1024)"`
	ReleaseCodeID        int         `gorm:"column:releasecode_id;type:int;DEFAULT:NULL"`
	ReleaseCode          ReleaseCode `gorm:"FOREIGNKEY:ReleaseCodeID;ASSOCIATION_FOREIGNKEY:ID"`
}

type Execution_Detail struct {
	ID          int       `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	StepNum     int       `gorm:"column:step_num;NOT NULL"`
	StepName    string    `gorm:"column:step_name;type:varchar(128);NOT NULL"`
	StepMsg     string    `gorm:"column:step_msg;type:varchar(1024);NOT NULL"`
	StepState   int       `gorm:"column:step_state;type:int;NOT NULL"`
	StepTime    time.Time `gorm:"column:step_time;NOT NULL"`
	CreatedAt   time.Time `gorm:"column:ctime;NOT NULL"`
	ExecutionID int       `gorm:"column:execution_id;type:int;NOT NULL"`
	Execution   Execution `gorm:"FOREIGNKEY:ExecutionID;ASSOCIATION_FOREIGNKEY:ID"`
}

type Agent_Operation struct {
	ID        int       `gorm:"type:int AUTO_INCREMENT;NOT NULL;PRIMARY_KEY"`
	AgentID   string    `gorm:"column:agent_id;NOT NULL"`
	Agent     Agent     `gorm:"FOREIGNKEY:AgentID;ASSOCIATION_FOREIGNKEY:ID"`
	OpMode    string    `gorm:"column:opmode;type:varchar(20);NOT NULL"`
	CreatedAt time.Time `gorm:"column:ctime;NOT NULL"`
}

type Cron_Task struct {
	TaskID    int       `gorm:"column:task_id;type:int;NOT NULL;PRIMARY_KEY"`
	EntryID   int       `gorm:"column:entry_id;type:int;NOT NULL"`
	Task      Task      `gorm:"FOREIGNKEY:TaskID;ASSOCIATION_FOREIGNKEY:ID"`
	TimeSpec  string    `gorm:"column:time_spec;type:varchar(30);NOT NULL"`
	CreatedAt time.Time `gorm:"column:ctime;NOT NULL"`
}
