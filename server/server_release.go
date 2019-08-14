/**
* @Author: xhzhang
* @Date: 2019-05-27 15:32
 */
package server

import (
	"context"
	"errors"
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

type Release struct{}

func (r *Release) AddRelease(ctx context.Context, in *pb.AddReleaseRequest) (*pb.ReleaseAddReply, error) {
	tx := comm.DB.Begin()
	releaseObj := comm.Release{Name: in.Name, Version: in.Version, OrganizationID: int(in.Orgid), ProjectID: int(in.Proid)}
	if err := tx.Create(&releaseObj).Error; err != nil {
		tx.Rollback()
		log.Slogger.Errorf("[Release] 创建发布失败: %s", err)
		return &pb.ReleaseAddReply{}, err
	}
	// 则分析代码字符串，初始化release_code
	for _, ci := range in.Releasecodes {
		name := ci.Name
		relativePath := ci.Relativepath
		releasecodeObj := comm.ReleaseCode{Name: name, RelativePath: relativePath, ReleaseID: releaseObj.ID}
		if err := tx.Create(&releasecodeObj).Error; err != nil {
			tx.Rollback()
			log.Slogger.Errorf("[Release] 发布添加成功，发布代码解析失败: %s", err)
			return &pb.ReleaseAddReply{}, err
		}
	}
	tx.Commit()
	log.Slogger.Infof("[Release] 发布添加成功")
	return &pb.ReleaseAddReply{Releaseid: int32(releaseObj.ID)}, nil
}

func (r *Release) DeleteRelease(ctx context.Context, in *pb.ReleaseNameRequest) (*pb.EmptyReply, error) {
	releaseObj := comm.Release{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &releaseObj) {
		log.Slogger.Errorf("[Release] 删除发布[%s]失败: 不存在,无法删除", in.Name)
		return &pb.EmptyReply{}, errors.New("发布不存在，无法删除")
	}
	if err := comm.DeleteRecord(&releaseObj); err != nil {
		log.Slogger.Errorf("[Release] 删除发布[%s]失败: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}
	log.Slogger.Infof("[Release] 删除分组[%s]成功", in.Name)
	return &pb.EmptyReply{}, nil
}

func (r *Release) GetReleaseCode(ctx context.Context, in *pb.ReleaseIdRequest) (*pb.ReleaseCodeList, error) {
	var rcl pb.ReleaseCodeList
	var rcs []comm.ReleaseCode
	if err := comm.DB.Where("release_id = ?", in.Id).Find(&rcs).Error; err != nil {
		return &rcl, err
	}
	for _, rc := range rcs {
		rcl.Releasecodes = append(rcl.Releasecodes, &pb.ReleaseCodeList_ReleaseCodeInfo{Id: int32(rc.ID), Name: rc.Name})
	}
	return &rcl, nil
}

func (r *Release) GetReleases(ctx context.Context, in *pb.EmptyRequest) (*pb.ReleaseList, error) {
	var rels []comm.Release
	if err := comm.DB.Preload("Organization").Preload("Project").Find(&rels).Error; err != nil {
		return nil, err
	}
	var rrels pb.ReleaseList
	for _, rel := range rels {
		rrels.Releases = append(rrels.Releases, &pb.ReleaseList_ReleaseInfo{Id: int32(rel.ID), Name: rel.Name, Version: rel.Version, Orgname: rel.Organization.Name, Proname: rel.Project.Name})
	}
	return &rrels, nil
}

func (r *Release) GetReleaseID(ctx context.Context, in *pb.ReleaseNameRequest) (*pb.ReleaseAddReply, error) {
	var rl comm.Release
	var rlr pb.ReleaseAddReply
	if err := comm.DB.Where("name = ?", in.Name).Find(&rl).Error; err != nil {
		return &rlr, err
	}
	rlr.Releaseid = int32(rl.ID)
	return &rlr, nil
}
