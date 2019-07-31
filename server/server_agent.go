/**
* @Author: xhzhang
* @Date: 2019-06-21 15:49
 */
package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Agent struct{}

func (a *Agent) GetAgents(ctx context.Context, in *pb.AgentGetRequest) (*pb.AgentList, error) {
	var agents []comm.Agent
	var ragents pb.AgentList
	if in.Groupid == 0 && in.Ac == 0 {
		return &ragents, errors.New("参数错误")
	}

	if in.Ac == 1 {
		if err := comm.DB.Where("status = ? ", 1).Find(&agents).Error; err != nil {
			return &ragents, err
		}
	} else if in.Ac == 2 {
		if err := comm.DB.Where("status = ? ", 0).Find(&agents).Error; err != nil {
			return &ragents, err
		}
	} else if in.Ac == 3 {
		if err := comm.DB.Find(&agents).Error; err != nil {
			return &ragents, err
		}
	} else {
		if err := comm.DB.Where("group_id = ? ", in.Groupid).Find(&agents).Error; err != nil {
			return &ragents, err
		}
	}
	for _, agent := range agents {
		ragents.Agents = append(ragents.Agents, &pb.AgentList_AgentInfo{Id: agent.ID, Alias: agent.Alias, Hostname: agent.HostName, Hostip: agent.HostIp, Status: agent.Status, Ctime: agent.CreatedAt.String(), Utime: agent.UpdatedAt.String()})
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
	opagentchannel := "grace." + in.Id
	c := ControlAgent{AgentID: in.Id, OPMode: in.Op}
	by, err := json.Marshal(&c)
	err = comm.PublishCMD(opagentchannel, string(by))
	if err != nil {
		return &pb.EmptyReply{}, err
	}
	return &pb.EmptyReply{}, nil
}
