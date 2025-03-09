package domains

import (
	"time"

	"github.com/google/uuid"
)

type PointTransactionsModel struct {
	ID              uuid.UUID       `json:"id"`
	TransactionType TransactionType `json:"transaction_type"`
	Points          int             `json:"points"`
	Reason          *string         `json:"reason"`
	UserID          int             `json:"user_id"`
	CreatedAt       time.Time       `json:"created_at"`
}

type TransactionType string

const (
	EarnTransactionType   TransactionType = "earn"
	SpendTransactionType  TransactionType = "spend"
	ExpireTransactionType TransactionType = "expire"
)
