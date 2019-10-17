/**
* @Author: xhzhang
* @Date: 2019-06-21 15:49
 */
package server

import (
	"context"
	"encoding/json"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Agent struct{}

func (a *Agent) GetAgents(ctx context.Context, in *pb.GetAgentRequest) (*pb.AgentList, error) {
	var agents []comm.Agent
	var ragents pb.AgentList
	queryCmd := comm.DB

	if in.Id != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Id)
	}

	if in.Name != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Name)
	}

	if in.Agentstatus == 1 {
		queryCmd = queryCmd.Where("status = ? ", 1)
	}  else {
		queryCmd = queryCmd
	}

	if err := queryCmd.Find(&agents).Error; err != nil {
		return &ragents, err
	}

	for _, agent := range agents {
		ragents.Agents = append(ragents.Agents, &pb.AgentList_AgentInfo{Id: agent.ID, Alias: agent.Alias, Hostname: agent.HostName, Hostip: agent.HostIp, Status: agent.Status, Ctime: agent.CreatedAt.String(), Utime: agent.UpdatedAt.Format("2006-01-02 15:04:05")})
	}

	return &ragents, nil
}

func (a *Agent) SetAgentAlias(ctx context.Context, in *pb.AgentAliasRequest) (*pb.EmptyReply, error) {
	agent := comm.Agent{ID: in.Id}
	err := agent.SetAlias(in.Alias)
	if err != nil {
		return &pb.EmptyReply{}, err
	}
	return &pb.EmptyReply{}, nil
}

// 重启agent
func (a *Agent) OperateAgent(ctx context.Context, in *pb.AgentRestartRequest) (*pb.EmptyReply, error) {
	log.Slogger.Infof("[OPAgent] [%s]->[%s]", in.Op, in.Id)
	operateAgent := comm.Agent_Operation{AgentID: in.Id, OpMode: in.Op}
	err := comm.CreateRecord(&operateAgent)
	if err != nil {
		return &pb.EmptyReply{}, err
	}
	opAgentChannel := "grace." + in.Id
	c := ControlAgent{AgentID: in.Id, OPMode: in.Op}
	by, err := json.Marshal(&c)
	err = comm.PublishCMD(opAgentChannel, string(by))
	if err != nil {
		return &pb.EmptyReply{}, err
	}
	return &pb.EmptyReply{}, nil
}
