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
func (c *CDPClient) AddProject(name string) (int32, error) {
	pc := c.newProjectClient()
	ctx := context.TODO()
	res, err := pc.AddProject(ctx, &pb.ProjectNameRequest{Name: name})
	if err != nil {
		return 0, err
	}
	return res.Proid, err
}

func (c *CDPClient) DeleteProject(name string) error {
	pc := c.newProjectClient()
	ctx := context.TODO()
	_, err := pc.DeleteProject(ctx, &pb.ProjectNameRequest{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *CDPClient) GetProjects() (*[]Project, error) {
	pc := c.newProjectClient()
	ctx := context.TODO()
	var projects []Project
	prolist, err := pc.GetProjects(ctx, &pb.EmptyRequest{})
	if err != nil {
		return &projects, err
	}

	for _, org := range prolist.Pros {
		projects = append(projects, Project{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return &projects, nil
}

/*
根据组织名称获取组织ID，返回组织名称和ID的map.eg:map[cdporg:1 org2:9 org3:10]
当参数为空时，则获取所有的组织ID,当参数指定时，则获取指定的组织ID。指定参数可以为一个或者多个
example:
        1. cdpclient.GetOrganizationID()
        2. cdpclient.GetOrganizationID("org2")
        3. cdpclient.GetOrganizationID("org3","org2")
*/
func (c *CDPClient) GetProjectID(proName ...string) (map[string]int32, error) {
	proNameId := make(map[string]int32)
	oc := c.newProjectClient()
	ctx := context.TODO()
	if len(proName) == 0 {
		res, err := oc.GetProjects(ctx, &pb.EmptyRequest{})
		if err != nil {
			return proNameId, err
		}
		for _, r := range res.Pros {
			proNameId[r.Name] = r.Id
		}

	} else {
		for _, ename := range proName {
			res, err := oc.GetProjectID(ctx, &pb.ProjectNameRequest{Name: ename})
			if err != nil {
				return proNameId, err
			}
			proNameId[ename] = res.Proid
		}
	}
	return proNameId, nil
}
