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

func (c *CDPClient) GetOrganizations(opts ...Option) (OrganizationSlice, error) {
	orgQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&orgQueryOption)
	}

	oc := c.newOrganizationClient()
	ctx := context.TODO()
	organizations := []Organization{}
	orglist, err := oc.GetOrganizations(ctx, &pb.GetOrgRequest{Ids:orgQueryOption.Ids,Names:orgQueryOption.Names})
	if err != nil {
		return organizations, err
	}

	for _, org := range orglist.Orgs {
		organizations = append(organizations, Organization{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return organizations, nil
}

/*
	access all organization records
*/
func (c *CDPClient) GetAllOrganizations() (OrganizationSlice, error) {
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	organizations := []Organization{}
	orglist, err := oc.GetOrganizations(ctx, &pb.GetOrgRequest{})
	if err != nil {
		return organizations, err
	}
	for _, org := range orglist.Orgs {
		organizations = append(organizations, Organization{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return organizations, nil
}

/*
	Obtain organization record information based on the organization name provided;
*/
func (c *CDPClient) GetOrganizationsFromNames(names []string) (OrganizationSlice, error) {
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	organizations := []Organization{}
	orglist, err := oc.GetOrganizations(ctx, &pb.GetOrgRequest{Names:names})
	if err != nil {
		return organizations, err
	}
	for _, org := range orglist.Orgs {
		organizations = append(organizations, Organization{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return organizations, nil
}

/*
	Obtain organization record information based on the organization id provided;
*/
func (c *CDPClient) GetOrganizationsFromIDs(ids []int32) (OrganizationSlice, error) {
	oc := c.newOrganizationClient()
	ctx := context.TODO()
	organizations := []Organization{}
	orglist, err := oc.GetOrganizations(ctx, &pb.GetOrgRequest{Ids: ids})
	if err != nil {
		return organizations, err
	}
	for _, org := range orglist.Orgs {
		organizations = append(organizations, Organization{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return organizations, nil
}