/**
* @Author: xhzhang
* @Date: 2019/7/17 13:26
 */
package server

import (
	"context"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Org struct{}

func (o *Org) AddOrganization(ctx context.Context, in *pb.OrgNameRequest) (*pb.OrgAddReply, error) {
	orgObj := comm.Organization{Name: in.Name}
	if err := comm.CreateRecord(&orgObj); err != nil {
		log.Slogger.Errorf("[Organization] add [%s] failed. %s", in.Name, err.Error())
		return &pb.OrgAddReply{}, err
	} else {
		log.Slogger.Infof("[Organization]  add [%s] successful.", orgObj.Name)
		return &pb.OrgAddReply{Orgid: int32(orgObj.ID)}, nil
	}
}

func (o *Org) DeleteOrganization(ctx context.Context, in *pb.OrgNameRequest) (*pb.EmptyReply, error) {
	orgObj := comm.Organization{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &orgObj) {
		log.Slogger.Errorf("[Organization] delete [%s] failed. not-exist", in.Name)
		return &pb.EmptyReply{}, errors.New("not exist.")
	}
	err := comm.DeleteRecord(&orgObj)
	if err != nil {
		log.Slogger.Errorf("[Organization] delete [%s] failed. %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}
	log.Slogger.Infof("[Organization] delete [%s] successful.", in.Name)
	return &pb.EmptyReply{}, nil
}

func (o *Org) GetOrganizations(ctx context.Context, in *pb.GetOrgRequest) (*pb.OrganizationList, error) {
	var orgs []comm.Organization
	queryCmd := comm.DB
	if in.Names != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Names)
	}

	if in.Ids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Ids)
	}

	if err := queryCmd.Find(&orgs).Error; err != nil {
		log.Slogger.Errorf("[Organization] query error: %v", err)
		return nil, err
	}

	var rorgs pb.OrganizationList
	for _, org := range orgs {
		rorgs.Orgs = append(rorgs.Orgs, &pb.OrganizationList_OrganizationInfo{Id: int32(org.ID), Name: org.Name, Ctime: org.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return &rorgs, nil
}

