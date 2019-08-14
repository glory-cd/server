/**
* @Author: xhzhang
* @Date: 2019/7/19 9:24
 */
package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "github.com/glory-cd/server/idlentity"
)

type CDPCClientAttr struct {
	CertFile   string
	ServerName string
	Address    string
}

type CDPClient struct {
	Conn   *grpc.ClientConn
	SNodes map[string]interface{}
}

func NewClient(cdpcattr CDPCClientAttr) (*CDPClient, error) {
	creds, err := credentials.NewClientTLSFromFile(cdpcattr.CertFile, cdpcattr.ServerName) // create the client TLS credentials
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(cdpcattr.Address, grpc.WithTransportCredentials(creds)) // initiate a connection with the server using creds
	if err != nil {
		return nil, err
	}
	return &CDPClient{Conn: conn, SNodes: make(map[string]interface{})}, nil
}

func (c *CDPClient) newServClient(stype string) (nsc interface{}) {
	if oclient, ok := c.SNodes[stype]; ok {
		nsc = oclient
	} else {
		switch stype {
		case KeyOrganizationClient:
			nsc = pb.NewOrganizationClient(c.Conn)
		case KeyEnvironmentClient:
			nsc = pb.NewEnvironmentClient(c.Conn)
		case KeyProjectClient:
			nsc = pb.NewProjectClient(c.Conn)
		case KeyGroupClient:
			nsc = pb.NewGroupClient(c.Conn)
		case KeyReleaseClient:
			nsc = pb.NewReleaseClient(c.Conn)
		case KeyServiceClient:
			nsc = pb.NewServiceClient(c.Conn)
		case KeyAgentClient:
			nsc = pb.NewAgentClient(c.Conn)
		case KeyTaskClient:
			nsc = pb.NewTaskClient(c.Conn)
		}
		c.SNodes[stype] = nsc
	}
	return
}

func (c *CDPClient) newOrganizationClient() (oc pb.OrganizationClient) {
	oci := c.newServClient(KeyOrganizationClient)
	oc = oci.(pb.OrganizationClient)
	return
}

func (c *CDPClient) newProjectClient() (pc pb.ProjectClient) {
	pci := c.newServClient(KeyProjectClient)
	pc = pci.(pb.ProjectClient)
	return
}

func (c *CDPClient) newEnvironmentClient() (ec pb.EnvironmentClient) {
	eci := c.newServClient(KeyEnvironmentClient)
	ec = eci.(pb.EnvironmentClient)
	return
}

func (c *CDPClient) newGroupClient() (gc pb.GroupClient) {
	gci := c.newServClient(KeyGroupClient)
	gc = gci.(pb.GroupClient)
	return
}

func (c *CDPClient) newServiceClient() (sc pb.ServiceClient) {
	sci := c.newServClient(KeyServiceClient)
	sc = sci.(pb.ServiceClient)
	return
}

func (c *CDPClient) newTaskClient() (sc pb.TaskClient) {
	sci := c.newServClient(KeyTaskClient)
	sc = sci.(pb.TaskClient)
	return
}

func (c *CDPClient) newAgentClient() (sc pb.AgentClient) {
	sci := c.newServClient(KeyAgentClient)
	sc = sci.(pb.AgentClient)
	return
}

func (c *CDPClient) newReleaseClient() (rc pb.ReleaseClient) {
	rci := c.newServClient(KeyReleaseClient)
	rc = rci.(pb.ReleaseClient)
	return
}
