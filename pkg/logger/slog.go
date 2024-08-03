package logger

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

var coldStart = true

// New creates [slog.Logger] with lambda context attributes.
func New(ctx context.Context, h slog.Handler) *slog.Logger {
	var isColdStart bool
	if coldStart {
		isColdStart = true
		coldStart = false
	}
	lambdaCtx, ok := lambdacontext.FromContext(ctx)
	if !ok {
		h = h.WithAttrs([]slog.Attr{
			slog.Bool("coldStart", isColdStart),
			slog.Group("function",
				slog.String("name", lambdacontext.FunctionName),
				slog.String("version", lambdacontext.FunctionVersion)),
		})
		return slog.New(h)
	}
	h = h.WithAttrs([]slog.Attr{
		slog.Bool("coldStart", isColdStart),
		slog.Group("function",
			slog.String("name", lambdacontext.FunctionName),
			slog.String("version", lambdacontext.FunctionVersion)),
		slog.String("requestId", lambdaCtx.AwsRequestID),
	})
	return slog.New(h)
}
