package pubsub

import (
	"log/slog"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/topics"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

func Routes(router *message.Router, subscriber *amqp.Subscriber, appCtx appcontext.AppContext) {
	// Instantiate dependencies
	pointsRepository := repositories.NewPointsRepositoryPg(appCtx.Pool)

	// Add handler into router
	router.AddNoPublisherHandler(
		"update_points_and_transactions",
		topics.TodoCompleted,
		subscriber,
		func(msg *message.Message) error {
			stringUserId := string(msg.Payload)
			userId, err := strconv.Atoi(stringUserId)
			if err != nil {
				appCtx.Logger.Info("[update_points_and_transactions Handler] Error getting userId value from update_points_and_transactions handler",
					slog.String("UUID", msg.UUID),
					slog.String("error", err.Error()),
				)
				return err
			}

			appCtx.Logger.Info("[update_points_and_transactions Handler] Handling a message",
				slog.String("UUID", msg.UUID),
				slog.Int("userId", userId),
			)

			addedPoint := 5
			pointModel, err := pointsRepository.UpdateWithTransactionHistory(
				msg.Context(),
				userId,
				domains.PointTransactionsModel{
					TransactionType: domains.EarnTransactionType,
					Reason:          domains.CompleteTodoReason,
					Points:          addedPoint,
				},
			)
			if err != nil {
				appCtx.Logger.Info("[update_points_and_transactions Handler] Error updating point with transaction history",
					slog.String("UUID", msg.UUID),
					slog.String("error", err.Error()),
				)
				return err
			}

			appCtx.Logger.Info("[update_points_and_transactions Handler] Successfully updating point with transaction history",
				slog.String("UUID", msg.UUID),
				slog.Any("payload", pointModel),
			)
			return nil
		},
	)
}
