package main

import (
	"context"
	"demo/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 建立GRPC连接
	conn, err := grpc.Dial("localhost:8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("gRPC Connection Error: " + err.Error())
	}
	// 关闭连接
	defer conn.Close()

	// 创建客户端
	client := pb.NewMessageServiceClient(conn)
	resp, err := client.Send(context.Background(), &pb.MessageRequest{
		RequestSomething: "Hello GRPC!",
	})
	if err != nil {
		panic("gRPC Client Error: " + err.Error())
	}
	fmt.Println("Response: ", resp.GetResponseSomething())
}
