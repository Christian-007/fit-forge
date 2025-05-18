package services

import (
	"context"
	"log/slog"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
)

type SubscriptionService struct {
	SubscriptionServiceOptions
}

type SubscriptionServiceOptions struct {
	PointsRepository repositories.PointsRepository
	Logger           applog.Logger
}

func NewSubscriptionService(options SubscriptionServiceOptions) SubscriptionService {
	return SubscriptionService{options}
}

func (s SubscriptionService) ProcessDueSubscriptions(ctx context.Context, dueDate string) error {
	usersDueForSubscription, err := s.PointsRepository.FindUsersForSubscriptionDeduction(ctx, dueDate)
	if err != nil {
		return err
	}

	for _, user := range usersDueForSubscription.EligibleForDeduction {
		_, err := s.PointsRepository.UpdateWithTransactionHistory(ctx, user.Id, domains.SubscriptionDeductionAmount)
		if err != nil {
			s.Logger.Error("failed to deduct points", slog.Any("user", user), slog.String("error", err.Error()))
		}
	}

	return nil
}
