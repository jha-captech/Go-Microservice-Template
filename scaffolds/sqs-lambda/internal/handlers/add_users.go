package handlers

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/captechconsulting/go-microservice-templates/sqs-lambda/internal/models"
)

type FailedItems struct {
	ItemIdentifier string `form:"itemIdentifier"`
}

type ReturnFailures struct {
	BatchItemFailures []FailedItems `json:"batchItemFailures"`
}
type userCreator interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
}

// HandleCreateUsers adds users from an SQS event,
func HandleCreateUsers(logger *slog.Logger, service userCreator) HandlerFunc {
	return func(ctx context.Context, sqsEvent events.SQSEvent) (ReturnFailures, error) {
		var batchItemFailures []FailedItems

		for _, record := range sqsEvent.Records {
			// unmarshal and validate
			user, problems, err := decodeValidateBody[inputUser, models.User](record.Body)
			if err != nil {
				logger.Error("Failed to decode validate body", "error", err, "problems", problems)
				batchItemFailures = append(batchItemFailures, FailedItems{
					ItemIdentifier: record.MessageId,
				})
				continue
			}

			// process
			if _, err = service.CreateUser(ctx, user); err != nil {
				logger.Error("Failed to create user", "error", err)
				batchItemFailures = append(batchItemFailures, FailedItems{
					ItemIdentifier: record.MessageId,
				})
			}
		}

		// return response
		return ReturnFailures{BatchItemFailures: batchItemFailures}, nil
	}
}
