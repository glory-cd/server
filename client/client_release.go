/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

// 添加发布，返回发布ID和错误信息
func (c *CDPClient) AddRelease(name string, version string, organizationid, projectid int, codes []ReleaseCode) (int32, error) {
	rc := c.newReleaseClient()
	ctx := context.TODO()
	var rcs []*pb.ReleaseCode
	for _, c := range codes {
		_tmp := pb.ReleaseCode{Name: c.CodeName, Relativepath: c.CodePath}
		rcs = append(rcs, &_tmp)
	}
	res, err := rc.AddRelease(ctx, &pb.AddReleaseRequest{Name: name, Orgid: int32(organizationid), Proid: int32(projectid), Releasecodes: rcs})
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

func (c *CDPClient) GetReleaseCode(releaseID int) (map[string]int, error) {
	rc := c.newReleaseClient()
	ctx := context.TODO()
	rcmap := map[string]int{}
	res, err := rc.GetReleaseCode(ctx, &pb.ReleaseIdRequest{Id: int32(releaseID)})
	if err != nil {
		return rcmap, err
	}

	for _, r := range res.Releasecodes {
		rcmap[r.Name] = int(r.Id)
	}
	return rcmap, nil
}
