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
		log.Slogger.Errorf("[Release] add [%s] failed: %s", in.Name, err)
		return &pb.ReleaseAddReply{}, err
	}
	// 则分析代码字符串，初始化release_code
	for _, ci := range in.Releasecodes {
		name := ci.Name
		relativePath := ci.Relativepath
		releasecodeObj := comm.ReleaseCode{Name: name, RelativePath: relativePath, ReleaseID: releaseObj.ID}
		if err := tx.Create(&releasecodeObj).Error; err != nil {
			tx.Rollback()
			log.Slogger.Errorf("[Release] add [%s] successful,but parse releasecode failed. err: %s", in.Name,err)
			return &pb.ReleaseAddReply{}, err
		}
	}
	tx.Commit()
	log.Slogger.Infof("[Release] add [%s] successful.",releaseObj.Name)
	return &pb.ReleaseAddReply{Releaseid: int32(releaseObj.ID)}, nil
}

func (r *Release) DeleteRelease(ctx context.Context, in *pb.ReleaseNameRequest) (*pb.EmptyReply, error) {
	releaseObj := comm.Release{Name: in.Name}
	if comm.CheckRecordWithName(in.Name, &releaseObj) {
		log.Slogger.Errorf("[Release] delete [%s] failed. err: do not exist and cannot be deleted", in.Name)
		return &pb.EmptyReply{}, errors.New("release do not exist and cannot be deleted")
	}
	if err := comm.DeleteRecord(&releaseObj); err != nil {
		log.Slogger.Errorf("[Release] delete [%s] failed. err: %s", in.Name, err)
		return &pb.EmptyReply{}, err
	}
	log.Slogger.Infof("[Release] delete [%s] successful.", in.Name)
	return &pb.EmptyReply{}, nil
}

func (r *Release) GetReleases(ctx context.Context, in *pb.GetReleaseRequest) (*pb.ReleaseList, error) {
	var rels []comm.Release
	queryCmd := comm.DB.Preload("Organization").Preload("Project")
	if in.Ids != nil {
		queryCmd = queryCmd.Where("id in (?)", in.Ids)
	}

	if in.Names != nil {
		queryCmd = queryCmd.Where("name in (?)", in.Names)
	}


	if in.Pros != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_projects ON cdp_projects.id = cdp_releases.project_id AND cdp_projects.name in (?)", in.Pros)
	}

	if in.Orgs != nil {
		queryCmd = queryCmd.Joins("JOIN cdp_organizations on cdp_organizations.id = cdp_releases.organization_id AND cdp_organizations.name in (?) ", in.Orgs)
	}

	// query cdp_releases
	err := queryCmd.Find(&rels).Error
	if err != nil {
		return nil, err
	}
	var rrels pb.ReleaseList
	for _, rel := range rels {
		// query cdp_releasecodes
		var rcs []comm.ReleaseCode
		err := comm.DB.Where("release_id = ?", rel.ID).Find(&rcs).Error
		if err != nil {
			return nil, err
		}
		var rcsmeta []*pb.QueryReleaseCode
		for _, rc := range rcs {
			rcsmeta = append(rcsmeta, &pb.QueryReleaseCode{Id:int32(rc.ID),Releaseid:int32(rc.ID),Rc:&pb.ReleaseCode{Name: rc.Name, Relativepath: rc.RelativePath}})
		}
		rrels.Releases = append(rrels.Releases, &pb.ReleaseList_ReleaseInfo{Id: int32(rel.ID), Name: rel.Name, Version: rel.Version, Orgname: rel.Organization.Name, Proname: rel.Project.Name, Rcs: rcsmeta})
	}
	return &rrels, nil
}


// 查询发布代码
func (r *Release) GetReleaseCodes(ctx context.Context, in *pb.GetReleaseCodeRequest) (*pb.ReleaseCodeList, error) {
	var rcs []comm.ReleaseCode

	if err := comm.DB.Where("release_id in (?)", in.Releaseids).Find(&rcs).Error; err != nil {
		return nil, err
	}

	var rcslist pb.ReleaseCodeList
	for _, rc := range rcs {
		rcslist.Rcs = append(rcslist.Rcs, &pb.QueryReleaseCode{Id:int32(rc.ID),Releaseid:int32(rc.ReleaseID),Rc:&pb.ReleaseCode{Name:rc.Name,Relativepath:rc.RelativePath}})
	}
	return &rcslist, nil
}