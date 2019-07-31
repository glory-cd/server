/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

type groupOption struct {
	OrgID int
	ProID int
	EnvID int
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

func WithOrg(id int) GroupOption {
	return newFuncOptionGroup(func(o *groupOption) { o.OrgID = id })
}

func WithPro(id int) GroupOption {
	return newFuncOptionGroup(func(o *groupOption) { o.ProID = id })
}

func WithEnv(id int) GroupOption {
	return newFuncOptionGroup(func(o *groupOption) { o.EnvID = id })
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
	res, err := gc.AddGroup(ctx, &pb.GroupAddRequest{Name: name, Orgid: int32(groupOption.OrgID), Envid: int32(groupOption.EnvID), Proid: int32(groupOption.ProID)})
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
