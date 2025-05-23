package main

import (
	"net"

	pointrepositories "github.com/Christian-007/fit-forge/internal/app/points/repositories"
	subscriptiongrpc "github.com/Christian-007/fit-forge/internal/app/subscription/delivery/grpc"
	subscriptionservices "github.com/Christian-007/fit-forge/internal/app/subscription/services"
	userrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"google.golang.org/grpc"
)

func InitGrpcServices(appCtx appcontext.AppContext) func(*grpc.Server) {
	pointRepository := pointrepositories.NewPointsRepositoryPg(appCtx.Pool)
	userRepository := userrepositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepository,
	})
	subscriptionService := subscriptionservices.NewSubscriptionService(subscriptionservices.SubscriptionServiceOptions{
		PointsRepository: pointRepository,
		UsersService: userService,
		Logger: appCtx.Logger,
	})

	return func(s *grpc.Server) {
		subscriptiongrpc.RegisterSubscriptionGrpc(s, &subscriptionService)
	}
}

func StartGrpcServer(addr string, registerFn func(*grpc.Server)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	registerFn(grpcServer)

	return grpcServer.Serve(listener)
}
