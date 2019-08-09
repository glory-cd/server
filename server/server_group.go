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
		log.Slogger.Errorf("[Group] 添加分组失败: %s", err.Error())
		return &pb.GroupAddReply{}, err
	}
	log.Slogger.Infof("[Group] 添加分组成功")
	return &pb.GroupAddReply{Groupid: int32(groupObj.ID)}, nil
}

func (g *Group) DeleteGroup(ctx context.Context, in *pb.GroupNameRequest) (*pb.EmptyReply, error) {
	groupObj := comm.Group{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &groupObj) {
		log.Slogger.Errorf("[Group] 删除分组[%s]失败: 不存在,无法删除", in.Name)
		return &pb.EmptyReply{}, errors.New("分组不存在，无法删除")
	}
	err := comm.DeleteRecord(&groupObj)
	if err != nil {
		log.Slogger.Errorf("[Group] 删除分组[%s]失败: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}

	log.Slogger.Infof("[Group] 删除分组[%s]成功", in.Name)
	return &pb.EmptyReply{}, nil
}

func (g *Group) GetGroups(ctx context.Context, in *pb.EmptyRequest) (*pb.GroupList, error) {
	var groups []comm.Group
	var rgrous pb.GroupList
	err := comm.DB.Preload("Organization").Preload("Environment").Preload("Project").Find(&groups)
	//err := db.DB.Find(&groups)
	if err.Error != nil {
		return &rgrous, err.Error
	}
	for _, group := range groups {
		rgrous.Groups = append(rgrous.Groups, &pb.GroupList_GroupInfo{Id: int32(group.ID), Name: group.Name, Orgname: group.Organization.Name, Envname: group.Environment.Name, Proname: group.Project.Name})
	}
	return &rgrous, nil
}

// 根据分组名称获取ID
func (o *Group) GetGroupID(ctx context.Context, in *pb.GroupNameRequest) (*pb.GroupAddReply, error) {
	var group comm.Group
	if err := comm.DB.Where("name=?", in.Name).Find(&group).Error; err != nil {
		return nil, err
	}
	return &pb.GroupAddReply{Groupid: int32(group.ID)}, nil
}
