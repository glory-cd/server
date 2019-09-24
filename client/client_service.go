/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

type addServiceOption struct {
	MoudleName  string
	CodePattern string
	PidFile     string
	StopCmd     string
	GroupID     int
}

type AddServiceOption interface {
	apply(*addServiceOption)
}

type funcAddServiceOption struct {
	f func(*addServiceOption)
}

func (fdo *funcAddServiceOption) apply(do *addServiceOption) {
	fdo.f(do)
}

func newFuncAddServiceOption(f func(*addServiceOption)) *funcAddServiceOption {
	return &funcAddServiceOption{f: f}
}

func WithMoudleName(mName string) AddServiceOption {
	return newFuncAddServiceOption(func(o *addServiceOption) { o.MoudleName = mName })
}

func WithCodePattern(cPattern string) AddServiceOption {
	return newFuncAddServiceOption(func(o *addServiceOption) { o.CodePattern = cPattern })
}

func WithPidFile(pfile string) AddServiceOption {
	return newFuncAddServiceOption(func(o *addServiceOption) { o.PidFile = pfile })
}

func WithStopCmd(scmd string) AddServiceOption {
	return newFuncAddServiceOption(func(o *addServiceOption) { o.StopCmd = scmd })
}

func WithGroup(id int) AddServiceOption {
	return newFuncAddServiceOption(func(o *addServiceOption) { o.GroupID = id })
}

func defaultAddServiceOption() addServiceOption {
	return addServiceOption{}
}

// 添加服务，返回组织ID和错误信息
func (c *CDPClient) AddService(name, dir, osUser, startCmd, agentID string, opts ...AddServiceOption) (string, error) {
	serviceOption := defaultAddServiceOption()
	for _, opt := range opts {
		opt.apply(&serviceOption)
	}
	addService := pb.ServiceAddRequest{Name: name,
		Dir:         dir,
		Moudlename:  serviceOption.MoudleName,
		Osuser:      osUser,
		Codepattern: serviceOption.CodePattern,
		Pidfile:     serviceOption.PidFile,
		Startcmd:    startCmd,
		Stopcmd:     serviceOption.StopCmd,
		Agentid:     agentID,
		Groupid:     int32(serviceOption.GroupID)}

	sc := c.newServiceClient()
	ctx := context.TODO()
	res, err := sc.AddService(ctx, &addService)
	if err != nil {
		return "", err
	}
	return res.Serviceid, err
}

func (c *CDPClient) DeleteService(id string) error {
	sc := c.newServiceClient()
	ctx := context.TODO()
	_, err := sc.DeleteService(ctx, &pb.ServiceDeleteRequest{Id: id})
	if err != nil {
		return err
	}
	return nil
}

// 查询服务
func (c *CDPClient) GetServices(opts ...QueryOption) (ServiceSlice, error) {
	serviceQueryOption := defaultQueryOption()
	for _, opt := range opts {
		opt.apply(&serviceQueryOption)
	}

	sc := c.newServiceClient()
	ctx := context.TODO()
	var services ServiceSlice
	servicelist, err := sc.GetServices(ctx, &pb.ServiceRequest{Agentids: serviceQueryOption.AgentIDs,
		Groupnames:   serviceQueryOption.GroupNames,
		Moudlenames:  serviceQueryOption.MoudleNames,
		Serviceids:   serviceQueryOption.ServiceIDs,
		Servicenames: serviceQueryOption.Names})
	if err != nil {
		return services, err
	}

	for _, s := range servicelist.Services {
		tmpService := Service{ID: s.Id, Name: s.Name, MoudleName: s.Moudlename, OsUser: s.Osuser, CodePattern: s.Codepattern, PidFile: s.Pidfile, StartCmd: s.Startcmd, StopCmd: s.Stopcmd, AgentName: s.Agentname, GroupName: s.Groupname}
		services = append(services, tmpService)
	}
	return services, nil
}

// 修改服务
type changeServiceOption struct {
	AgentID string
	GroupID int
}

type ChangeServiceOption interface {
	apply(*changeServiceOption)
}

type funcChangeServiceOption struct {
	f func(*changeServiceOption)
}

func (fdo *funcChangeServiceOption) apply(do *changeServiceOption) {
	fdo.f(do)
}

func newFuncChangeServiceOption(f func(*changeServiceOption)) *funcChangeServiceOption {
	return &funcChangeServiceOption{f: f}
}

func ChaneAgent(id string) ChangeServiceOption {
	return newFuncChangeServiceOption(func(o *changeServiceOption) { o.AgentID = id })
}

func ChangeGroup(id int) ChangeServiceOption {
	return newFuncChangeServiceOption(func(o *changeServiceOption) { o.GroupID = id })
}

func defaultChangeServiceOption() changeServiceOption {
	return changeServiceOption{}
}

/*
	修改service的agent属组,group属组
*/
func (c *CDPClient) ChangeServiceAgent(serviceid string, opts ...ChangeServiceOption) error {
	changeOption := defaultChangeServiceOption()
	for _, opt := range opts {
		opt.apply(&changeOption)
	}

	sc := c.newServiceClient()
	ctx := context.TODO()

	_, err := sc.ChangeServiceOwn(ctx, &pb.ServiceChangeOwnRequest{Id: serviceid, Agentid: changeOption.AgentID, Groupid: int32(changeOption.GroupID)})
	if err != nil {
		return err
	}
	return nil
}
