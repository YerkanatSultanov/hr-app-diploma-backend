package grpc

import (
	"context"
	"hr-app-diploma-backend/auth-service/service"
	"hr-app-diploma-backend/pkg/protobuf/hr-app-diploma-backend/auth-service/proto"
	"net"

	"google.golang.org/grpc"
)

type AuthGRPCServer struct {
	proto.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthGRPCServer(authService *service.AuthService) *AuthGRPCServer {
	return &AuthGRPCServer{authService: authService}
}

func (s *AuthGRPCServer) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	userID, err := s.authService.VerifyToken(req.Token)
	if err != nil {
		return &proto.VerifyTokenResponse{Valid: false}, nil
	}

	return &proto.VerifyTokenResponse{
		UserId: int32(userID),
		Valid:  true,
	}, nil
}

func StartGRPCServer(authService *service.AuthService) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic("Failed to start gRPC server: " + err.Error())
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAuthServiceServer(grpcServer, NewAuthGRPCServer(authService))

	if err := grpcServer.Serve(listener); err != nil {
		panic("Failed to serve gRPC: " + err.Error())
	}
}
