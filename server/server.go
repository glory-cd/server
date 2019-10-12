/**
* @Author: xhzhang
* @Date: 2019-06-19 14:55
 */
package server

import (
	"github.com/glory-cd/server/comm"
	pb "github.com/glory-cd/server/idlentity"
	"github.com/glory-cd/utils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
)

const rpcPort = "localhost:50051"

func InitRpcServer() {
	config := comm.Config().RPC
	rpcPort := config.HostPort
	lis, err := net.Listen("tcp", rpcPort)
	if err != nil {
		log.Slogger.Fatalf("[RPC] Listen failed. %v", err)
	}

	log.Slogger.Infof("[RPC] Listen sucessful. %s", rpcPort)

	creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
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
		log.Slogger.Fatalf("[RPC] register rpc service failed. %v", err)
	}
}
