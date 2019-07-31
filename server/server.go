/**
* @Author: xhzhang
* @Date: 2019-06-19 14:55
 */
package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
)

const rpcPort = ":50051"

func InitRpcServer() {
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Slogger.Fatalf("[RPC] 监听失败: %v", err)
	}

	log.Slogger.Infof("[RPC] 监听成功: %s", rpcPort)

	creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Slogger.Fatalf("[RPC] could not load TLS keys: %s", err)
	}
	// create a gRPC option array with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	s := grpc.NewServer(opts...)
	pb.RegisterOrganizationServer(s, &Org{})
	pb.RegisterProjectServer(s, &Pro{})
	pb.RegisterEnvironmentServer(s, &Env{})
	pb.RegisterGroupServer(s, &Group{})
	pb.RegisterReleaseServer(s, &Release{})
	pb.RegisterAgentServer(s, &Agent{})
	pb.RegisterServiceServer(s, &Service{})
	pb.RegisterTaskServer(s, &Task{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Slogger.Fatalf("[RPC] 监听连接rpc服务失败: %v", err)
	}
}
