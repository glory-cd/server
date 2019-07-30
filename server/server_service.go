/**
* @Author: xhzhang
* @Date: 2019-06-21 15:49
 */
package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"server/comm"
	pb "server/idlentity"
	"utils/log"
)

func GetMd5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type Service struct{}

func (s *Service) AddService(ctx context.Context, in *pb.ServiceAddRequest) (*pb.ServiceAddReply, error) {
	// 校验部分参数
	if in.Name == "" || in.Dir == "" || in.Osuser == "" || in.Agentid == "" || in.Startcmd == "" {
		return &pb.ServiceAddReply{}, errors.New("参数错误: 不能为空的字段为空")
	}
	serviceObj := comm.Service{Name: in.Name,
		Dir:          in.Dir,
		OsUser:       in.Osuser,
		ModuleName:   in.Moudlename,
		CodePatterns: in.Codepattern,
		Pidfile:      in.Pidfile,
		StartCMD:     in.Startcmd,
		StopCMD:      in.Stopcmd,
		AgentID:      in.Agentid,
		GroupID:      int(in.Groupid)}
	serviceObj.ID = GetMd5String(in.Agentid + in.Dir)

	if err := comm.CreateRecord(&serviceObj); err != nil {
		log.Slogger.Errorf("[Service] 添加服务[%s]失败: %s", in.Name, err.Error())
		return &pb.ServiceAddReply{}, err
	}

	log.Slogger.Infof("[Service] 添加服务[%s]成功", in.Name)
	return &pb.ServiceAddReply{Serviceid: serviceObj.ID}, nil
}

func (s *Service) DeleteService(ctx context.Context, in *pb.ServiceDeleteRequest) (*pb.EmptyReply, error) {
	serviceObj := comm.Service{ID: in.Id}
	if comm.CheckRecordWithStringID(in.Id, &serviceObj) {
		log.Slogger.Errorf("[Service] 删除服务[%s]失败: 不存在,无法删除", in.Id)
		return &pb.EmptyReply{}, errors.New("服务不存在，无法删除")

	}
	err := comm.DeleteRecord(&serviceObj)
	if err != nil {
		log.Slogger.Errorf("[Service] 删除服务[%s]失败: %s", serviceObj.Name, err)
		return &pb.EmptyReply{}, err
	}
	log.Slogger.Infof("[Release] 删除服务[%s]成功", serviceObj.Name)
	return &pb.EmptyReply{}, nil
}

func (s *Service) GetServices(ctx context.Context, in *pb.ServiceRequest) (*pb.ServiceList, error) {
	var services []comm.Service
	var rservices pb.ServiceList
	queryCmd := comm.DB

	if in.Serviceids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Groupids)
	}

	if in.Groupids != nil {
		queryCmd = queryCmd.Where("group_id in (?)", in.Groupids)
	}

	if in.Agentids != nil {
		queryCmd = queryCmd.Where("agent_id in (?)", in.Agentids)
	}

	if in.Moudlenames != nil {
		queryCmd = queryCmd.Where("moudle_name in (?)", in.Moudlenames)
	}

	if err := queryCmd.Preload("Agent").Preload("Group").Find(&services).Error; err != nil {
		return &rservices, err
	}

	for _, service := range services {
		si := &pb.ServiceList_ServiceInfo{Id: service.ID,
			Name:        service.Name,
			Dir:         service.Dir,
			Moudlename:  service.ModuleName,
			Osuser:      service.OsUser,
			Codepattern: service.CodePatterns,
			Pidfile:     service.Pidfile,
			Startcmd:    service.StartCMD,
			Stopcmd:     service.StopCMD,
			Agentname:   service.Agent.Alias,
			Groupname:   service.Group.Name,
		}
		rservices.Services = append(rservices.Services, si)
	}
	return &rservices, nil
}

//修改服务的agentid,groupid
func (s *Service) ChangeServiceOwn(ctx context.Context, in *pb.ServiceChangeOwnRequest) (*pb.EmptyReply, error) {
	service := comm.Service{ID: in.Id}
	agentid := in.Agentid
	groupid := in.Groupid

	if comm.DB.First(&service).RecordNotFound() {
		return &pb.EmptyReply{}, errors.New("服务不存在")
	}

	var err error
	updateCMD := comm.DB.Model(&service)
	if agentid != "" && groupid != 0 {
		err = updateCMD.Updates(map[string]interface{}{"agent_id": agentid, "group_id": groupid}).Error
	} else if agentid != "" && groupid == 0 {
		err = updateCMD.Update("agent_id", agentid).Error
	} else if agentid == "" && groupid != 0 {
		err = updateCMD.Update("group_id", groupid).Error
	} else {
		return &pb.EmptyReply{}, errors.New("参数错误")
	}

	if err != nil {
		return &pb.EmptyReply{}, err
	}

	return &pb.EmptyReply{}, nil
}
