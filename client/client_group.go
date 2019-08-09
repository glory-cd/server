/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	"errors"
	pb "github.com/glory-cd/server/idlentity"
)

type groupOption struct {
	OrgName string
	ProName string
	EnvName string
}

type GroupOption interface {
	apply(*groupOption)
}

type funcOptionGroup struct {
	f func(*groupOption)
}

func (fdo *funcOptionGroup) apply(do *groupOption) {
	fdo.f(do)
}

func newFuncOptionGroup(f func(*groupOption)) *funcOptionGroup {
	return &funcOptionGroup{f: f}
}

func WithOrg(name string) GroupOption {
	return newFuncOptionGroup(func(o *groupOption) { o.OrgName = name })
}

func WithPro(name string) GroupOption {
	return newFuncOptionGroup(func(o *groupOption) { o.ProName = name })
}

func WithEnv(name string) GroupOption {
	return newFuncOptionGroup(func(o *groupOption) { o.EnvName = name})
}

//默认参数
func defaultOptionGroup() groupOption {
	return groupOption{}
}

// 添加组织，返回组织ID和错误信息
func (c *CDPClient) AddGroup(name string, opts ...GroupOption) (int32, error) {
	groupOption := defaultOptionGroup()
	for _, opt := range opts {
		opt.apply(&groupOption)
	}
	gc := c.newGroupClient()
	ctx := context.TODO()
	attr := pb.GroupAddRequest{Name: name}
	if groupOption.OrgName != ""{
		gidMap,err := c.GetOrganizationID(groupOption.OrgName)
		if err != nil{
			return 0,errors.New("获取组织ID错误: " + err.Error())
		}
		attr.Orgid = gidMap[groupOption.OrgName]
	}

	if groupOption.EnvName != ""{
		eidMap,err := c.GetEnvironmentID(groupOption.EnvName)
		if err != nil{
			return 0,errors.New("获取环境ID错误: " + err.Error())
		}
		attr.Envid = eidMap[groupOption.EnvName]
	}

	if groupOption.ProName != ""{
		pidMap,err := c.GetProjectID(groupOption.ProName)
		if err != nil{
			return 0,errors.New("获取项目ID错误: " + err.Error())
		}
		attr.Proid = pidMap[groupOption.ProName]
	}

	res, err := gc.AddGroup(ctx, &attr)
	if err != nil {
		return 0, err
	}
	return res.Groupid, err
}

func (c *CDPClient) DeleteGroup(name string) error {
	gc := c.newGroupClient()
	ctx := context.TODO()
	_, err := gc.DeleteGroup(ctx, &pb.GroupNameRequest{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *CDPClient) GetGroups() ([]Group, error) {
	gc := c.newGroupClient()
	ctx := context.TODO()
	var groups []Group
	grouplist, err := gc.GetGroups(ctx, &pb.EmptyRequest{})
	if err != nil {
		return groups, err
	}
	for _, org := range grouplist.Groups {
		groups = append(groups, Group{ID: org.Id, Name: org.Name, Organization: org.Orgname, Environment: org.Envname, Project: org.Orgname})
	}
	return groups, nil
}

/*
根据组织名称获取组织ID，返回组织名称和ID的map.eg:map[cdporg:1 org2:9 org3:10]
当参数为空时，则获取所有的组织ID,当参数指定时，则获取指定的组织ID。指定参数可以为一个或者多个
example:
        1. cdpclient.GetOrganizationID()
        2. cdpclient.GetOrganizationID("org2")
        3. cdpclient.GetOrganizationID("org3","org2")
*/
func (c *CDPClient) GetGroupID(groupName ...string) (map[string]int32, error) {
	groupNameId := make(map[string]int32)
	oc := c.newGroupClient()
	ctx := context.TODO()
	if len(groupName) == 0 {
		res, err := oc.GetGroups(ctx, &pb.EmptyRequest{})
		if err != nil {
			return groupNameId, err
		}
		for _, r := range res.Groups {
			groupNameId[r.Name] = r.Id
		}

	} else {
		for _, ename := range groupName {
			res, err := oc.GetGroupID(ctx, &pb.GroupNameRequest{Name: ename})
			if err != nil {
				return groupNameId, err
			}
			groupNameId[ename] = res.Groupid
		}
	}
	return groupNameId, nil
}
