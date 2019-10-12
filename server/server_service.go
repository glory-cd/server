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
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
	"math/rand"
)

func GetMd5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateRandSalt(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func HashOsUser(ospass string) (string, error) {
	c := sha512_crypt.New()
	salt := "$6$" + GenerateRandSalt(12) + "$"
	hash, err := c.Generate([]byte(ospass), []byte(salt))
	return hash, err
}

type Service struct{}

func (s *Service) AddService(ctx context.Context, in *pb.ServiceAddRequest) (*pb.ServiceAddReply, error) {
	// 校验部分参数
	if in.Name == "" || in.Dir == "" || in.Osuser == "" || in.Ospass == "" || in.Agentid == "" || in.Startcmd == "" || in.Modulename == "" {
		return &pb.ServiceAddReply{}, errors.New("[Service] Parameter error, field that cannot be empty is empty")
	}

	serviceObj := comm.Service{Name: in.Name,
		Dir:          in.Dir,
		OsUser:       in.Osuser,
		ModuleName:   in.Modulename,
		CodePatterns: in.Codepattern,
		Pidfile:      in.Pidfile,
		StartCMD:     in.Startcmd,
		StopCMD:      in.Stopcmd,
		AgentID:      in.Agentid,
		GroupID:      int(in.Groupid)}
	serviceObj.ID = GetMd5String(in.Agentid + in.Dir)
	hashpass, err := HashOsUser(in.Ospass)
	if err != nil {
		return &pb.ServiceAddReply{}, err
	}

	serviceObj.OsPass = hashpass

	if err := comm.CreateRecord(&serviceObj); err != nil {
		log.Slogger.Errorf("[Service] add [%s] failed. %s", in.Name, err.Error())
		return &pb.ServiceAddReply{}, err
	}

	log.Slogger.Infof("[Service] add [%s] successful.", in.Name)
	return &pb.ServiceAddReply{Serviceid: serviceObj.ID}, nil
}

func (s *Service) DeleteService(ctx context.Context, in *pb.ServiceDeleteRequest) (*pb.EmptyReply, error) {
	serviceObj := comm.Service{ID: in.Id}
	if comm.CheckRecordWithStringID(in.Id, &serviceObj) {
		log.Slogger.Errorf("[Service] delete [%s] failed. not-exist", in.Id)
		return &pb.EmptyReply{}, errors.New("service not exist.")

	}
	err := comm.DeleteRecord(&serviceObj)
	if err != nil {
		log.Slogger.Errorf("[Service] delete [%s] failed. %s", serviceObj.Name, err)
		return &pb.EmptyReply{}, err
	}
	log.Slogger.Infof("[Release] delete [%s] successful.", serviceObj.Name)
	return &pb.EmptyReply{}, nil
}

func (s *Service) GetServices(ctx context.Context, in *pb.ServiceRequest) (*pb.ServiceList, error) {
	var services []comm.Service
	var rservices pb.ServiceList
	queryCmd := comm.DB.Preload("Agent").Preload("Group")

	if in.Serviceids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Serviceids)
	}

	if in.Servicenames != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Servicenames)
	}

	if in.Moudlenames != nil {
		queryCmd = queryCmd.Where("moudle_name in (?)", in.Moudlenames)
	}

	if in.Groupnames != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_groups ON cdp_groups.id = cdp_services.group_id AND cdp_groups.name in (?)", in.Groupnames)
	}

	if in.Agentids != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_agents ON cdp_agents.id = cdp_services.agent_id AND cdp_agents.id in (?)", in.Agentids)
	}

	if err := queryCmd.Find(&services).Error; err != nil {
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
			Agentid:     service.AgentID,
			Agentname:   service.Agent.Alias,
			Groupname:   service.Group.Name,
			Hostip:      service.Agent.HostIp,
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
		return &pb.EmptyReply{}, errors.New("service not-exist.")
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
		return &pb.EmptyReply{}, errors.New("parameter error")
	}

	if err != nil {
		return &pb.EmptyReply{}, err
	}

	return &pb.EmptyReply{}, nil
}
