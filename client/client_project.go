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

func (c *CDPClient) GetProjects(opts ...Option) (ProjectSlice, error) {
	proQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&proQueryOption)
	}
	pc := c.newProjectClient()
	ctx := context.TODO()
	var projects []Project
	prolist, err := pc.GetProjects(ctx, &pb.GetProRequest{Ids:proQueryOption.Ids,Names:proQueryOption.Names})
	if err != nil {
		return projects, err
	}

	for _, pro := range prolist.Pros {
		projects = append(projects, Project{ID: pro.Id, Name: pro.Name, CreatTime: pro.Ctime})
	}
	return projects, nil
}

func (c *CDPClient) GetAllProjects() (ProjectSlice, error) {
	pc := c.newProjectClient()
	ctx := context.TODO()
	var projects []Project
	prolist, err := pc.GetProjects(ctx, &pb.GetProRequest{})
	if err != nil {
		return projects, err
	}

	for _, pro := range prolist.Pros {
		projects = append(projects, Project{ID: pro.Id, Name: pro.Name, CreatTime: pro.Ctime})
	}
	return projects, nil
}

func (c *CDPClient) GetProjectsFromNames(names []string) (ProjectSlice, error) {
	pc := c.newProjectClient()
	ctx := context.TODO()
	var projects []Project
	prolist, err := pc.GetProjects(ctx, &pb.GetProRequest{Names: names})
	if err != nil {
		return projects, err
	}

	for _, pro := range prolist.Pros {
		projects = append(projects, Project{ID: pro.Id, Name: pro.Name, CreatTime: pro.Ctime})
	}
	return projects, nil
}

func (c *CDPClient) GetProjectsFromIDs(ids []int32) (ProjectSlice, error) {
	pc := c.newProjectClient()
	ctx := context.TODO()
	var projects []Project
	prolist, err := pc.GetProjects(ctx, &pb.GetProRequest{Ids: ids})
	if err != nil {
		return projects, err
	}

	for _, pro := range prolist.Pros {
		projects = append(projects, Project{ID: pro.Id, Name: pro.Name, CreatTime: pro.Ctime})
	}
	return projects, nil
}
