package serviceImpl

import (
	"context"
	"demo/pb"
	"fmt"
)

type MessageSenderServerImpl struct {
	*pb.UnimplementedMessageServiceServer
}

func (MessageSenderServerImpl) Send(ctx context.Context, request *pb.MessageRequest) (*pb.MessageResponse, error) {
	fmt.Println("Received Message: ", request.GetRequestSomething())
	resp := &pb.MessageResponse{}
	resp.ResponseSomething = "Roger That!"
	return resp, nil
}
