// What we learned so far with Correlation ID is that it can be used for logging in a monolithic service
// We can use this unique ID to trace back where the error is
// Correlation ID can be implemented on the following actions:
// 1. Generate: before publishing events, when receiving incoming request from Frontend (the FE who is generating the UUID)
// 2. Logging: subscribing events, encountering an error

package decorator

import (
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/ThreeDotsLabs/watermill/message"
)

const CorrelationIdMetadataKey = "correlation_id"

type PublishWithCorrelationId struct {
	message.Publisher
}

func (p PublishWithCorrelationId) Publish(topic string, messages ...*message.Message) error {
	for _, message := range messages {
		if message.Metadata.Get(CorrelationIdMetadataKey) != "" {
			continue
		}

		correlationId, _ := requestctx.CorrelationId(message.Context())
		message.Metadata.Set(CorrelationIdMetadataKey, correlationId)
	}

	return p.Publisher.Publish(topic, messages...)
}
