package grpc

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/subscription/services"
	subscriptionpb "github.com/Christian-007/fit-forge/proto/subscription/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SubscriptionGrpcHandler struct {
	subscriptionpb.UnimplementedSubscriptionServiceServer
	subscriptionService *services.SubscriptionService
}

func RegisterSubscriptionGrpc(s *grpc.Server, svc *services.SubscriptionService) {
	subscriptionpb.RegisterSubscriptionServiceServer(s, &SubscriptionGrpcHandler{subscriptionService: svc})
}

func (s *SubscriptionGrpcHandler) ProcessDueSubscriptions(ctx context.Context, req *subscriptionpb.DueSubscriptionRequest) (*emptypb.Empty, error) {
	err := s.subscriptionService.ProcessDueSubscriptions(ctx, req.GetDueDate())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
