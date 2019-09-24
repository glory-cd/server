/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

func (c *CDPClient) AddRelease(name string, version string, orgName, projectName string, codes []ReleaseCode) (int32, error) {
	rc := c.newReleaseClient()
	ctx := context.TODO()
	var rcs []*pb.ReleaseCode
	for _, c := range codes {
		_tmp := pb.ReleaseCode{Name: c.CodeName, Relativepath: c.CodePath}
		rcs = append(rcs, &_tmp)
	}

	oc := c.newOrganizationClient()
	orgs, err := oc.GetOrganizations(ctx, &pb.GetOrgRequest{Names: []string{orgName}})
	if err != nil {
		return 0, err
	}
	orgid := orgs.Orgs[0].Id

	pc := c.newProjectClient()
	pros, err := pc.GetProjects(ctx, &pb.GetProRequest{Names: []string{projectName}})
	if err != nil {
		return 0, err
	}
	projectid := pros.Pros[0].Id

	res, err := rc.AddRelease(ctx, &pb.AddReleaseRequest{Name: name, Version: version, Orgid: orgid, Proid: projectid, Releasecodes: rcs})
	if err != nil {
		return 0, err
	}
	return res.Releaseid, err
}

func (c *CDPClient) DeleteRelease(name string) error {
	rc := c.newReleaseClient()
	ctx := context.TODO()
	_, err := rc.DeleteRelease(ctx, &pb.ReleaseNameRequest{Name: name})
	if err != nil {
		return err
	}
	return nil
}

// query release
func (c *CDPClient) GetReleases(opts ...QueryOption) (ReleaseSlice, error) {
	queryReleaseOption := defaultQueryOption()
	for _, opt := range opts {
		opt.apply(&queryReleaseOption)
	}

	rc := c.newReleaseClient()
	ctx := context.TODO()
	var releases ReleaseSlice
	releaselist, err := rc.GetReleases(ctx, &pb.GetReleaseRequest{Ids: queryReleaseOption.Ids,
		Names: queryReleaseOption.Names,
		Orgs:  queryReleaseOption.OrgNames,
		Pros:  queryReleaseOption.ProNames})
	if err != nil {
		return releases, err
	}
	for _, r := range releaselist.Releases {
		var rcmeta []int32
		for _,rc := range r.Rcs{
			rcmeta = append(rcmeta,rc.Id)
		}
		releases = append(releases, Release{ID: r.Id, Name: r.Name, Version: r.Version, ProName: r.Proname, OrgName: r.Orgname,ReleaseCodes:rcmeta})
	}
	return releases, nil
}
/*
	Getting the releasecodes  based on release ID
    para releaseId:
    return map[string]int32
*/
func (c *CDPClient) GetReleaseCodeMap(releaseId int32) (map[string]int32, error) {
	mapCodeNameId := map[string]int32{}
	ctx := context.TODO()
	rc := c.newReleaseClient()

	releasecodeList, err := rc.GetReleaseCodes(ctx, &pb.GetReleaseCodeRequest{Releaseids: []int32{releaseId}})
	if err != nil {
		return mapCodeNameId, err
	}

	for _, rc := range releasecodeList.Rcs {
		mapCodeNameId[rc.Rc.Name] = rc.Id
	}
	return mapCodeNameId, nil
}

func (c *CDPClient) GetReleaseCodes(releaseIds []int32) (ReleaseCodeSlice, error) {
	var rcs ReleaseCodeSlice
	ctx := context.TODO()
	rc := c.newReleaseClient()

	releasecodeList, err := rc.GetReleaseCodes(ctx, &pb.GetReleaseCodeRequest{Releaseids: releaseIds})
	if err != nil {
		return rcs, err
	}

	for _, rc := range releasecodeList.Rcs {
		rcs = append(rcs, ReleaseCode{ReleaseID:rc.Releaseid,Id:rc.Id,CodeName: rc.Rc.Name, CodePath: rc.Rc.Relativepath})
	}
	return rcs, nil
}

