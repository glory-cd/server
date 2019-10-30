/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

// 添加服务，返回组织ID和错误信息
func (c *CDPClient) AddService(name, dir, osUser, osPass, agentID, moduleName string, opts ...Option) (string, error) {
	serviceOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&serviceOption)
	}

	/*var groupID int32 = 1
	if serviceOption.GroupName != "" {
		groups, err := c.GetGroups(WithGroupNames([]string{serviceOption.GroupName}))
		if err != nil {
			return "", errors.New("get group ID err: " + err.Error())
		}
		groupID = groups.GetID()
	}*/


	addService := pb.ServiceAddRequest{Name: name,
		Dir:         dir,
		Modulename:  moduleName,
		Osuser:      osUser,
		Ospass:      osPass,
		//Codepattern: serviceOption.CodePattern,
		//Pidfile:     serviceOption.PidFile,
		//Stopcmd:     serviceOption.StopCmd,
		Agentid:     agentID,
		Groupid:     serviceOption.GroupID}

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
func (c *CDPClient) GetServices(opts ...Option) (ServiceSlice, error) {
	serviceQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&serviceQueryOption)
	}

	sc := c.newServiceClient()
	ctx := context.TODO()
	var services ServiceSlice
	serviceList, err := sc.GetServices(ctx, &pb.ServiceRequest{Agentids: serviceQueryOption.AgentIDs,
		Groupnames:   serviceQueryOption.GroupNames,
		Moudlenames:  serviceQueryOption.ModuleNames,
		Serviceids:   serviceQueryOption.ServiceIDs,
		Servicenames: serviceQueryOption.Names})
	if err != nil {
		return services, err
	}

	for _, s := range serviceList.Services {
		tmpService := Service{ID: s.Id,
			Name:        s.Name,
			ModuleName:  s.Moudlename,
			OsUser:      s.Osuser,
			CodePattern: s.Codepattern,
			PidFile:     s.Pidfile,
			StartCmd:    s.Startcmd,
			StopCmd:     s.Stopcmd,
			HostIp:      s.Hostip,
			AgentName:   s.Agentname,
			AgentID:     s.Agentid,
			GroupName:   s.Groupname}
		services = append(services, tmpService)
	}
	return services, nil
}



/*
	修改service的agent属组,group属组
*/
func (c *CDPClient) ChangeServiceAgent(serviceId string, opts ...Option) error {
	changeOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&changeOption)
	}

	sc := c.newServiceClient()
	ctx := context.TODO()

	_, err := sc.ChangeServiceOwn(ctx, &pb.ServiceChangeOwnRequest{Id: serviceId, Agentid: changeOption.AgentID, Groupid: int32(changeOption.GroupID)})
	if err != nil {
		return err
	}
	return nil
}
