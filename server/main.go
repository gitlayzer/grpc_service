package main

import (
	"demo/pb"
	"demo/serviceImpl"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// new一个grpc服务
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterMessageServiceServer(grpcServer, &serviceImpl.MessageSenderServerImpl{})
	listener, err := net.Listen("tcp", ":8002") // 监听端口
	if err != nil {
		panic("TCP Listen Error: " + err.Error())
	}

	// 启动RPC服务
	err = grpcServer.Serve(listener)
	if err != nil {
		panic("gRPC Server Error: " + err.Error())
	}
}
