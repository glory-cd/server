/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

func (c *CDPClient) GetAgents(opts ...Option) (AgentSlice, error) {
	agentQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&agentQueryOption)
	}

	ac := c.newAgentClient()
	ctx := context.TODO()
	var agents AgentSlice
	var aqs int32

	if agentQueryOption.AgentIsOnLine == true {
		aqs = 1
	} else {
		aqs = 2
	}

	agentList, err := ac.GetAgents(ctx, &pb.GetAgentRequest{Agentstatus: aqs, Id: agentQueryOption.AgentIDs, Name: agentQueryOption.Names})
	if err != nil {
		return agents, err
	}

	for _, a := range agentList.Agents {
		agents = append(agents, Agent{ID: a.Id, Alias: a.Alias, Host: a.Hostname, IP: a.Hostip, Status: a.Status, CreatTime: a.Ctime, LatestTime: a.Utime})
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

// 操作agent
func (c *CDPClient) OperateAgent(op string, agentIds ...string) error {
	ac := c.newAgentClient()
	ctx := context.TODO()
	for _, aid := range agentIds {
		_, err := ac.OperateAgent(ctx, &pb.AgentRestartRequest{Id: aid, Op: op})
		if err != nil {
			return err
		}
	}
	return nil
}
