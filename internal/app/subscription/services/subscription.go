package services

import (
	"context"
	"log/slog"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	usersdomain "github.com/Christian-007/fit-forge/internal/app/users/domains"
	usersdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
	usersservice "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
)

type SubscriptionService struct {
	SubscriptionServiceOptions
}

type SubscriptionServiceOptions struct {
	PointsRepository repositories.PointsRepository
	UsersService     usersservice.UserService
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
		_, err := s.PointsRepository.UpdateWithTransactionHistory(ctx, user.Id, domains.PointTransactionsModel{
			Points:          domains.SubscriptionDeductionAmount,
			TransactionType: domains.SubscriptionDeductionTransactionType,
			Reason:          domains.SubscriptionDeductionReason,
		})
		if err != nil {
			s.Logger.Error("failed to deduct points", slog.Any("user", user), slog.String("error", err.Error()))
		}
	}

	for _, user := range usersDueForSubscription.InsufficientPoints {
		inactiveSubscriptionStatus := usersdomain.InactiveSubscriptionStatus
		_, err := s.UsersService.UpdateOneByEmail(user.Email, usersdto.UpdateUserRequest{SubscriptionStatus: &inactiveSubscriptionStatus})
		if err != nil {
			s.Logger.Error("failed to update subscription status to 'INACTIVE'", slog.Any("user", user), slog.String("error", err.Error()))
		}
	}

	return nil
}
