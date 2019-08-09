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

/*
根据组织名称获取组织ID，返回组织名称和ID的map.eg:map[cdporg:1 org2:9 org3:10]
当参数为空时，则获取所有的组织ID,当参数指定时，则获取指定的组织ID。指定参数可以为一个或者多个
example:
        1. cdpclient.GetOrganizationID()
        2. cdpclient.GetOrganizationID("org2")
        3. cdpclient.GetOrganizationID("org3","org2")
*/
func (c *CDPClient) GetEnvironmentID(envName ...string) (map[string]int32, error) {
	envNameId := make(map[string]int32)
	oc := c.newEnvironmentClient()
	ctx := context.TODO()
	if len(envName) == 0 {
		res, err := oc.GetEnvironments(ctx, &pb.EmptyRequest{})
		if err != nil {
			return envNameId, err
		}
		for _, r := range res.Envs {
			envNameId[r.Name] = r.Id
		}

	} else {
		for _, ename := range envName {
			res, err := oc.GetEnvironmentID(ctx, &pb.EnvNameRequest{Name: ename})
			if err != nil {
				return envNameId, err
			}
			envNameId[ename] = res.Envid
		}
	}
	return envNameId, nil
}
