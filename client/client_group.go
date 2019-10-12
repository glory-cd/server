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

// 添加组织，返回组织ID和错误信息
func (c *CDPClient) AddGroup(name string, opts ...Option) (int32, error) {
	groupOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&groupOption)
	}
	gc := c.newGroupClient()
	ctx := context.TODO()
	attr := pb.GroupAddRequest{Name: name}
	if groupOption.OrgName != "" {
		orgs, err := c.GetOrganizationsFromNames([]string{groupOption.OrgName})
		if err != nil {
			return 0, errors.New("get Org ID err: " + err.Error())
		}
		attr.Orgid = orgs.GetID()
	}

	if groupOption.EnvName != "" {
		envs, err := c.GetEnvironmentsFromNames([]string{groupOption.EnvName})
		if err != nil {
			return 0, errors.New("get env ID err: " + err.Error())
		}
		attr.Envid = envs.GetID()
	}

	if groupOption.ProName != "" {
		pros, err := c.GetProjectsFromNames([]string{groupOption.ProName})
		if err != nil {
			return 0, errors.New("get pro ID err: " + err.Error())
		}
		attr.Proid = pros.GetID()
	}

	res, err := gc.AddGroup(ctx, &attr)
	if err != nil {
		return 0, err
	}
	return res.Groupid, err
}

// delete group
func (c *CDPClient) DeleteGroup(name string) error {
	gc := c.newGroupClient()
	ctx := context.TODO()
	_, err := gc.DeleteGroup(ctx, &pb.GroupNameRequest{Name: name})
	if err != nil {
		return err
	}
	return nil
}

// query group
func (c *CDPClient) GetGroups(opts ...Option) (GroupSlice, error) {
	queryGroupOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&queryGroupOption)
	}

	gc := c.newGroupClient()
	ctx := context.TODO()
	var groups GroupSlice
	grouplist, err := gc.GetGroups(ctx, &pb.GetGroupRequest{Ids: queryGroupOption.Ids,
		Names: queryGroupOption.Names,
		Orgs:  queryGroupOption.OrgNames,
		Envs:  queryGroupOption.EnvNames,
		Pros:  queryGroupOption.ProNames})
	if err != nil {
		return groups, err
	}
	for _, gro := range grouplist.Groups {
		groups = append(groups, Group{ID: gro.Id, Name: gro.Name, Organization: gro.Orgname, Environment: gro.Envname, Project: gro.Proname})
	}
	return groups, nil
}

// get agentids from group id slice
func (c *CDPClient) GetAgentsFromGroup(groupIds []int32) ([]string, error) {
	gc := c.newGroupClient()
	ctx := context.TODO()

	aids, err := gc.GetAgentIdFromGroup(ctx, &pb.GetAgentFromGroupRequest{Gids: groupIds})
	if err != nil {
		return nil, err
	}
	return aids.Agentid, err
}
