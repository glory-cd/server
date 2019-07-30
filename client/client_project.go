/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "server/idlentity"
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
