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

// 获取上线agent
func (c *CDPClient) GetOnLineAgents() ([]Agent, error) {
	var agents []Agent
	var err error
	var agentlist *pb.AgentList
	ac := c.newAgentClient()
	ctx := context.TODO()
	agentlist, err = ac.GetAgents(ctx, &pb.AgentGetRequest{Ac: 1})
	if err != nil {
		return agents, err
	}

	for _, a := range agentlist.Agents {
		agents = append(agents, Agent{ID: a.Id, Alias: a.Alias, Host: a.Hostname, Ip: a.Hostip, Status: a.Status, CreatTime: a.Ctime})
	}
	return agents, nil
}

// 获取下线agent
func (c *CDPClient) GetOffLineAgents() ([]Agent, error) {
	var agents []Agent
	var err error
	var agentlist *pb.AgentList
	ac := c.newAgentClient()
	ctx := context.TODO()
	agentlist, err = ac.GetAgents(ctx, &pb.AgentGetRequest{Ac: 2})
	if err != nil {
		return agents, err
	}

	for _, a := range agentlist.Agents {
		agents = append(agents, Agent{ID: a.Id, Alias: a.Alias, Host: a.Hostname, Ip: a.Hostip, Status: a.Status, CreatTime: a.Ctime})
	}
	return agents, nil
}

// 获取所有agent
func (c *CDPClient) GetAllAgents() ([]Agent, error) {
	var agents []Agent
	var err error
	var agentlist *pb.AgentList
	ac := c.newAgentClient()
	ctx := context.TODO()
	agentlist, err = ac.GetAgents(ctx, &pb.AgentGetRequest{Ac: 3})
	if err != nil {
		return agents, err
	}

	for _, a := range agentlist.Agents {
		agents = append(agents, Agent{ID: a.Id, Alias: a.Alias, Host: a.Hostname, Ip: a.Hostip, Status: a.Status, CreatTime: a.Ctime})
	}
	return agents, nil
}

// 根据groupID获取agent
func (c *CDPClient) GetGroupAgents(groupID int) ([]Agent, error) {
	var agents []Agent
	var err error
	var agentlist *pb.AgentList
	ac := c.newAgentClient()
	ctx := context.TODO()
	agentlist, err = ac.GetAgents(ctx, &pb.AgentGetRequest{Groupid: int32(groupID)})
	if err != nil {
		return agents, err
	}

	for _, a := range agentlist.Agents {
		agents = append(agents, Agent{ID: a.Id, Alias: a.Alias, Host: a.Hostname, Ip: a.Hostip, Status: a.Status, CreatTime: a.Ctime})
	}
	return agents, nil
}

// 设置agent别名
func (c *CDPClient) SetAgentAlias(agentID, agentAlias string) error {
	ac := c.newAgentClient()
	ctx := context.TODO()

	_, err := ac.SetAgentAlias(ctx, &pb.AgentAliasRequest{Id: agentID, Alias: agentAlias})

	if err != nil {
		return err
	}
	return nil
}

type agentOperateOption struct {
	AgentIDs []string
	GroupID  int
}

type AgentOperateOption interface {
	apply(*agentOperateOption)
}

type funcOptionAgentOperate struct {
	f func(*agentOperateOption)
}

func (fdo *funcOptionAgentOperate) apply(do *agentOperateOption) {
	fdo.f(do)
}

func newFuncOptionAgentOperate(f func(*agentOperateOption)) *funcOptionAgentOperate {
	return &funcOptionAgentOperate{f: f}
}

func WithAgentID(ids ...string) AgentOperateOption {
	return newFuncOptionAgentOperate(func(o *agentOperateOption) { o.AgentIDs = ids })
}

func WithGroupID(id int) AgentOperateOption {
	return newFuncOptionAgentOperate(func(o *agentOperateOption) { o.GroupID = id })
}

//默认参数
func defaultOptionAgentOperate() agentOperateOption {
	return agentOperateOption{}
}

// 操作agent
func (c *CDPClient) OperateAgent(op string, opts ...AgentOperateOption) error {
	agentOperateOption := defaultOptionAgentOperate()
	for _, opt := range opts {
		opt.apply(&agentOperateOption)
	}

	var agentIDList []string
	if len(agentOperateOption.AgentIDs) > 0 && agentOperateOption.GroupID != 0 {
		return errors.New("参数错误: 两个参数不能同时存在！")
	}

	for _, agentid := range agentOperateOption.AgentIDs {
		agentIDList = append(agentIDList, agentid)
	}

	agents, err := c.GetGroupAgents(agentOperateOption.GroupID)
	if err != nil {
		return err
	}

	for _, a := range agents {
		agentIDList = append(agentIDList, a.ID)
	}

	ac := c.newAgentClient()
	ctx := context.TODO()
	for _, aid := range agentIDList {
		_, err := ac.OperateAgent(ctx, &pb.AgentRestartRequest{Id: aid, Op: op})
		if err != nil {
			return err
		}
	}
	return nil
}
