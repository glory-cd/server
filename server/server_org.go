/**
* @Author: xhzhang
* @Date: 2019/7/17 13:26
 */
package server

import (
	"context"
	"errors"
	"server/comm"
	pb "server/idlentity"
	"utils/log"
)

type Org struct{}

func (o *Org) AddOrganization(ctx context.Context, in *pb.OrgNameRequest) (*pb.OrgAddReply, error) {
	orgObj := comm.Organization{Name: in.Name}
	if err := comm.CreateRecord(&orgObj); err != nil {
		log.Slogger.Errorf("[Organization] 添加组织[%s]失败: %s", in.Name, err.Error())
		return &pb.OrgAddReply{}, err
	} else {
		log.Slogger.Infof("[Organization] 添加组织[%s]成功", orgObj.Name)
		return &pb.OrgAddReply{Orgid: int32(orgObj.ID)}, nil
	}
}

func (o *Org) DeleteOrganization(ctx context.Context, in *pb.OrgNameRequest) (*pb.EmptyReply, error) {
	orgObj := comm.Organization{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &orgObj) {
		log.Slogger.Errorf("[Organization] 删除组织[%s]失败: 不存在,无法删除", in.Name)
		return &pb.EmptyReply{}, errors.New("组织不存在，无法删除")
	}
	err := comm.DeleteRecord(&orgObj)
	if err != nil {
		log.Slogger.Errorf("[Organization] 删除组织[%s]失败: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}
	log.Slogger.Infof("[Organization] 删除组织[%s]成功", in.Name)
	return &pb.EmptyReply{}, nil
}

func (o *Org) GetOrganizations(ctx context.Context, in *pb.EmptyRequest) (*pb.OrganizationList, error) {
	var orgs []comm.Organization
	if err := comm.DB.Find(&orgs).Error; err != nil {
		log.Slogger.Errorf("[Organization] %s", err)
		return nil, err
	}
	var rorgs pb.OrganizationList
	for _, org := range orgs {
		rorgs.Orgs = append(rorgs.Orgs, &pb.OrganizationList_OrganizationInfo{Id: int32(org.ID), Name: org.Name, Ctime: org.CreatedAt.String()})
	}
	return &rorgs, nil
}

// 根据组织名称获取ID
func (o *Org) GetOrganizationID(ctx context.Context, in *pb.OrgNameRequest) (*pb.OrgAddReply, error) {
	var org comm.Organization
	if err := comm.DB.Where("name=?", in.Name).Find(&org).Error; err != nil {
		return nil, err
	}
	return &pb.OrgAddReply{Orgid: int32(org.ID)}, nil
}
