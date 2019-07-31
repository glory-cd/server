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
		log.Slogger.Errorf("[Project] 添加项目失败: %s", err.Error())
		return &pb.ProjectAddReply{}, err
	} else {
		log.Slogger.Infof("[Project] 添加项目成功")
		return &pb.ProjectAddReply{Proid: int32(proObj.ID)}, err
	}
}

func (p *Pro) DeleteProject(ctx context.Context, in *pb.ProjectNameRequest) (*pb.EmptyReply, error) {
	proObj := comm.Project{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &proObj) {
		log.Slogger.Errorf("[Project] 删除项目[%s]失败: 不存在,无法删除", in.Name)
		return &pb.EmptyReply{}, errors.New("项目不存在，无法删除")
	}
	err := comm.DeleteRecord(&proObj)
	if err != nil {
		log.Slogger.Errorf("[Project] 删除项目[%s]失败: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}

	log.Slogger.Infof("[Project] 删除项目[%s]成功", in.Name)
	return &pb.EmptyReply{}, nil
}

func (p *Pro) GetProjects(ctx context.Context, in *pb.EmptyRequest) (*pb.ProjectList, error) {
	var pros []comm.Project
	if err := comm.DB.Find(&pros).Error; err != nil {
		return nil, err
	}
	var rpros pb.ProjectList
	for _, pro := range pros {
		rpros.Pros = append(rpros.Pros, &pb.ProjectList_ProjectInfo{Id: int32(pro.ID), Name: pro.Name, Ctime: pro.CreatedAt.String()})
	}
	return &rpros, nil
}
