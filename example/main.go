package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var handler = func(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	for _, record := range event.Records {
		slog.InfoContext(ctx, "event received", slog.String("messageId", record.MessageId))
	}
	return events.SQSEventResponse{}, nil
}

func main() {
	lambda.Start(handler)
}
