package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"hr-app-diploma-backend/pkg/protobuf/hr-app-diploma-backend/auth-service/proto"
)

type AuthServiceClient struct {
	proto.AuthServiceClient
}

func NewAuthServiceClient(authServiceAddr string) (proto.AuthServiceClient, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth-service: %w", err)
	}

	client := proto.NewAuthServiceClient(conn)
	return &AuthServiceClient{client}, nil
}
