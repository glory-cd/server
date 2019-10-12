/**
* @Author: xhzhang
* @Date: 2019/7/19 9:23
 */
package client

import (
	"context"
	pb "github.com/glory-cd/server/idlentity"
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

func (c *CDPClient) GetEnvironments(opts ...Option) (EnvironmentSlice, error) {
	envQueryOption := defaultOption()
	for _, opt := range opts {
		opt.apply(&envQueryOption)
	}
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	var environments []Environment
	envlist, err := ec.GetEnvironments(ctx, &pb.GetEnvRequest{Ids:envQueryOption.Ids,Names:envQueryOption.Names})
	if err != nil {
		return environments, err
	}

	for _, env := range envlist.Envs {
		environments = append(environments, Environment{ID: env.Id, Name: env.Name, CreatTime: env.Ctime})
	}
	return environments, nil
}

func (c *CDPClient) GetAllEnvironments() (EnvironmentSlice, error) {
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	var environments []Environment
	envlist, err := ec.GetEnvironments(ctx, &pb.GetEnvRequest{})
	if err != nil {
		return environments, err
	}

	for _, env := range envlist.Envs {
		environments = append(environments, Environment{ID: env.Id, Name: env.Name, CreatTime: env.Ctime})
	}
	return environments, nil
}

func (c *CDPClient) GetEnvironmentsFromNames(names []string) (EnvironmentSlice, error) {
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	var environments []Environment
	envlist, err := ec.GetEnvironments(ctx, &pb.GetEnvRequest{Names: names})
	if err != nil {
		return environments, err
	}

	for _, env := range envlist.Envs {
		environments = append(environments, Environment{ID: env.Id, Name: env.Name, CreatTime: env.Ctime})
	}
	return environments, nil
}

func (c *CDPClient) GetEnvironmentsFromIDs(ids []int32) (EnvironmentSlice, error) {
	ec := c.newEnvironmentClient()
	ctx := context.TODO()
	var environments []Environment
	envlist, err := ec.GetEnvironments(ctx, &pb.GetEnvRequest{Ids: ids})
	if err != nil {
		return environments, err
	}

	for _, env := range envlist.Envs {
		environments = append(environments, Environment{ID: env.Id, Name: env.Name, CreatTime: env.Ctime})
	}
	return environments, nil
}