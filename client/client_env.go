/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "server/idlentity"
)

// 添加环境，返回组织ID和错误信息
func (c *CDPClient) AddEnvironment(name string) (int32, error) {
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	res, err := ec.AddEnvironment(ctx, &pb.EnvNameRequest{Name: name})
	if err != nil {
		return 0, err
	}
	return res.Envid, err
}

func (c *CDPClient) DeleteEnvironment(name string) error {
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	_, err := ec.DeleteEnvironment(ctx, &pb.EnvNameRequest{Name: name})
	if err != nil {
		return err
	}
	return nil
}

func (c *CDPClient) GetEnvironments() (*[]Environment, error) {
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	var environments []Environment
	envlist, err := ec.GetEnvironments(ctx, &pb.EmptyRequest{})
	if err != nil {
		return &environments, err
	}

	for _, org := range envlist.Envs {
		environments = append(environments, Environment{ID: org.Id, Name: org.Name, CreatTime: org.Ctime})
	}
	return &environments, nil
}
