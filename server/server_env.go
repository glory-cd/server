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
		log.Slogger.Errorf("[Environment] add [%s] failed. %s", in.Name,err.Error())
		return &pb.EnvAddReply{}, err
	} else {
		log.Slogger.Infof("[Environment] add [%s] failed.", envObj.Name)
		return &pb.EnvAddReply{Envid: int32(envObj.ID)}, err
	}
}

func (e *Env) DeleteEnvironment(ctx context.Context, in *pb.EnvNameRequest) (*pb.EmptyReply, error) {
	envObj := comm.Environment{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &envObj) {
		log.Slogger.Errorf("[Environment] delete [%s] failed. not-exist", in.Name)
		return &pb.EmptyReply{}, errors.New("env not-exist")

	}
	err := comm.DeleteRecord(&envObj)
	if err != nil {
		log.Slogger.Errorf("[Environment] delete [%s] failed. %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}

	log.Slogger.Infof("[Environment] delete [%s] successful.", in.Name)
	return &pb.EmptyReply{}, nil
}

func (o *Env) GetEnvironments(ctx context.Context, in *pb.GetEnvRequest) (*pb.EnvironmentList, error) {
	var envs []comm.Environment
	queryCmd := comm.DB

	if in.Names != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Names)
	}

	if in.Ids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Ids)
	}

	if err := queryCmd.Find(&envs).Error; err != nil {
		log.Slogger.Errorf("[Environment] query error: %v", err)
		return nil, err
	}

	var renvs pb.EnvironmentList
	for _, env := range envs {
		renvs.Envs = append(renvs.Envs, &pb.EnvironmentList_EnvironmentInfo{Id: int32(env.ID), Name: env.Name, Ctime: env.CreatedAt.Format("2006-01-02 15:04:05")})
	}
	return &renvs, nil
}