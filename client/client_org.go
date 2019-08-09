/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
)

// 添加组织，返回组织ID和错误信息
func (c *CDPClient) AddOrganization(name string) (int32, error) {
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	res, err := oc.AddOrganization(ctx, &pb.OrgNameRequest{Name: name})
	if err != nil {
		return 0, err
	}
	return res.Orgid, err
}

func (c *CDPClient) DeleteOrganization(name string) error {
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	_, err := oc.DeleteOrganization(ctx, &pb.OrgNameRequest{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *CDPClient) GetOrganizations() (*[]Organization, error) {
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	var organizations []Organization
	orglist, err := oc.GetOrganizations(ctx, &pb.EmptyRequest{})
	if err != nil {
		return &organizations, err
	}
	for _, org := range orglist.Orgs {
		organizations = append(organizations, Organization{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return &organizations, nil
}

/*
根据组织名称获取组织ID，返回组织名称和ID的map.eg:map[cdporg:1 org2:9 org3:10]
当参数为空时，则获取所有的组织ID,当参数指定时，则获取指定的组织ID。指定参数可以为一个或者多个
example:
        1. cdpclient.GetOrganizationID()
        2. cdpclient.GetOrganizationID("org2")
        3. cdpclient.GetOrganizationID("org3","org2")
*/
func (c *CDPClient) GetOrganizationID(orgName ...string) (map[string]int32, error) {
	orgNameId := make(map[string]int32)
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	if len(orgName) == 0 {
		res, err := oc.GetOrganizations(ctx, &pb.EmptyRequest{})
		if err != nil {
			return orgNameId, err
		}
		for _, r := range res.Orgs {
			orgNameId[r.Name] = r.Id
		}

	} else {
		for _, oname := range orgName {
			res, err := oc.GetOrganizationID(ctx, &pb.OrgNameRequest{Name: oname})
			if err != nil {
				return orgNameId, err
			}
			orgNameId[oname] = res.Orgid
		}
	}
	return orgNameId, nil
}
