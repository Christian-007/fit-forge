package pubsub

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/topics"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill/message"
)

type PointRewardsPayload struct {
	ID     string
	Topic  string 
	Points int
	UserID string 
}

func Routes(router *message.Router, subscriber *googlecloud.Subscriber, appCtx appcontext.AppContext) {
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
				domains.CreatePointTransactions{
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

			messageId := watermill.NewUUID()
			pointRewardsPayload := PointRewardsPayload{
				ID: messageId,
				Topic: topics.PointsRewarded,
				Points: addedPoint,
				UserID: stringUserId,
			}
			payload, err := json.Marshal(pointRewardsPayload)
			if err != nil {
				appCtx.Logger.Error("failed to marshal pointRewardsPayload", slog.String("error", err.Error()))
				return err
			}

			pubsubMessage := message.NewMessage(messageId, payload)
			err = appCtx.Publisher.Publish(topics.PointsRewarded, pubsubMessage)
			if err != nil {
				appCtx.Logger.Error("failed to publish PointsRewarded", slog.String("error", err.Error()))
				return err
			}

			appCtx.Logger.Info("[update_points_and_transactions Handler] published an event successfully",
				slog.String("UUID", messageId),
				slog.Any("payload", pointRewardsPayload),
			)
			return nil
		},
	)
}
