/**
* @Author: xhzhang
* @Date: 2019/7/17 13:32
 */
package server

import (
	"context"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Group struct{}

func (g *Group) AddGroup(ctx context.Context, in *pb.GroupAddRequest) (*pb.GroupAddReply, error) {
	groupObj := comm.Group{Name: in.Name, EnvironmentID: int(in.Envid), OrganizationID: int(in.Orgid), ProjectID: int(in.Proid)}
	if err := comm.CreateRecord(&groupObj); err != nil {
		log.Slogger.Errorf("[Group] add group [%s] failed: %s", in.Name, err.Error())
		return &pb.GroupAddReply{}, err
	}
	log.Slogger.Infof("[Group] add group [%s] successful.", groupObj.Name)
	return &pb.GroupAddReply{Groupid: int32(groupObj.ID)}, nil
}

func (g *Group) DeleteGroup(ctx context.Context, in *pb.GroupNameRequest) (*pb.EmptyReply, error) {
	groupObj := comm.Group{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &groupObj) {
		log.Slogger.Errorf("[Group] delte group [%s] failed. not-exist", in.Name)
		return &pb.EmptyReply{}, errors.New("Group do not exist and cannot be deleted.")
	}
	err := comm.DeleteRecord(&groupObj)
	if err != nil {
		log.Slogger.Errorf("[Group] delte group [%s] failed.: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}

	log.Slogger.Infof("[Group] delte group [%s] successful.", in.Name)
	return &pb.EmptyReply{}, nil
}

func (g *Group) GetGroups(ctx context.Context, in *pb.GetGroupRequest) (*pb.GroupList, error) {
	var groups []comm.Group
	var rGrous pb.GroupList
	queryCmd := comm.DB.Preload("Project").Preload("Environment").Preload("Organization")
	if in.Ids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Ids)
	}

	if in.Names != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Names)
	}

	if in.Pros != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_projects ON cdp_projects.id = cdp_groups.project_id AND cdp_projects.name in (?)", in.Pros)
	}

	if in.Envs != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_environments on cdp_environments.id = cdp_groups.environment_id AND cdp_environments.name in (?)", in.Envs)
	}

	if in.Orgs != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_organizations on cdp_organizations.id = cdp_groups.organization_id AND cdp_organizations.name in (?) ", in.Orgs)
	}

	err := queryCmd.Find(&groups).Error
	if err != nil {
		log.Slogger.Errorf("[Group] query err: %s", err)
		return &rGrous, err
	}
	for _, group := range groups {
		rGrous.Groups = append(rGrous.Groups, &pb.GroupList_GroupInfo{Id: int32(group.ID), Name: group.Name, Orgname: group.Organization.Name, Envname: group.Environment.Name, Proname: group.Project.Name})
	}
	return &rGrous, nil
}

func (g *Group) GetAgentIdFromGroup(ctx context.Context, in *pb.GetAgentFromGroupRequest) (*pb.GroupAgentIds, error) {
	var services []comm.Service
	var rAgentIDs pb.GroupAgentIds
	queryCmd := comm.DB

	if in.Gids != nil {
		queryCmd = queryCmd.Where("group_id in (?)", in.Gids)
	}

	if err := queryCmd.Find(&services).Error; err != nil {
		return &rAgentIDs, err
	}

	for _, s := range services {
		rAgentIDs.Agentid = append(rAgentIDs.Agentid, s.AgentID)
	}

	return &rAgentIDs, nil
}
