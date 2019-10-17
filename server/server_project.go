/**
* @Author: xhzhang
* @Date: 2019/7/17 13:29
 */
package server

import (
	"context"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Pro struct{}

func (p *Pro) AddProject(ctx context.Context, in *pb.ProjectNameRequest) (*pb.ProjectAddReply, error) {
	proObj := comm.Project{Name: in.Name}
	if err := comm.CreateRecord(&proObj); err != nil {
		log.Slogger.Errorf("[Project] add [%s] failed. %s", proObj.Name,err.Error())
		return &pb.ProjectAddReply{}, err
	} else {
		log.Slogger.Infof("[Project] add [%s] successful.", proObj.Name)
		return &pb.ProjectAddReply{Proid: int32(proObj.ID)}, err
	}
}

func (p *Pro) DeleteProject(ctx context.Context, in *pb.ProjectNameRequest) (*pb.EmptyReply, error) {
	proObj := comm.Project{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &proObj) {
		log.Slogger.Errorf("[Project] delete [%s] failed. not-exist", in.Name)
		return &pb.EmptyReply{}, errors.New("project not exist")
	}
	err := comm.DeleteRecord(&proObj)
	if err != nil {
		log.Slogger.Errorf("[Project] delete [%s] failed.: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}

	log.Slogger.Infof("[Project] delete [%s] successful.", in.Name)
	return &pb.EmptyReply{}, nil
}

func (p *Pro) GetProjects(ctx context.Context, in *pb.GetProRequest) (*pb.ProjectList, error) {
	var pros []comm.Project
	queryCmd := comm.DB

	if in.Names != nil  {
		queryCmd = queryCmd.Where("name in (?)", in.Names)
	}

	if in.Ids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Ids)
	}

	if err := queryCmd.Find(&pros).Error; err != nil {
		log.Slogger.Errorf("[Project] query error: %v", err)
		return nil, err
	}

	var rpros pb.ProjectList
	for _, pro := range pros {
		rpros.Pros = append(rpros.Pros, &pb.ProjectList_ProjectInfo{Id: int32(pro.ID), Name: pro.Name, Ctime: pro.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return &rpros, nil
}
