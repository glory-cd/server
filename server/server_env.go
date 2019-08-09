/**
* @Author: xhzhang
* @Date: 2019/7/17 13:31
 */
package server

import (
	"context"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Env struct{}

func (e *Env) AddEnvironment(ctx context.Context, in *pb.EnvNameRequest) (*pb.EnvAddReply, error) {
	envObj := comm.Environment{Name: in.Name}
	if err := comm.CreateRecord(&envObj); err != nil {
		log.Slogger.Errorf("[Environment] 添加环境失败: %s", err.Error())
		return &pb.EnvAddReply{}, err
	} else {
		log.Slogger.Infof("[Environment] 添加环境[%s]成功", envObj.Name)
		return &pb.EnvAddReply{Envid: int32(envObj.ID)}, err
	}
}

func (e *Env) DeleteEnvironment(ctx context.Context, in *pb.EnvNameRequest) (*pb.EmptyReply, error) {
	envObj := comm.Environment{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &envObj) {
		log.Slogger.Errorf("[Environment] 删除环境[%s]失败: 不存在,无法删除", in.Name)
		return &pb.EmptyReply{}, errors.New("环境不存在，无法删除")

	}
	err := comm.DeleteRecord(&envObj)
	if err != nil {
		log.Slogger.Errorf("[Environment] 删除环境[%s]失败: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}

	log.Slogger.Infof("[Environment] 删除环境[%s]成功", in.Name)
	return &pb.EmptyReply{}, nil
}

func (e *Env) GetEnvironments(ctx context.Context, in *pb.EmptyRequest) (*pb.EnvironmentList, error) {
	var envs []comm.Environment
	if err := comm.DB.Find(&envs).Error; err != nil {
		return nil, err
	}

	var renvs pb.EnvironmentList
	for _, env := range envs {
		renvs.Envs = append(renvs.Envs, &pb.EnvironmentList_EnvironmentInfo{Id: int32(env.ID), Name: env.Name, Ctime: env.CreatedAt.String()})
	}

	return &renvs, nil
}
