package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/tecchu11/lambda-invoker-go/pkg/logger"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		in struct {
			msg string
			ctx context.Context
		}
		want struct {
			msg       string
			requestID string
			coldStart bool
		}
	}{
		"cold start": {
			in: struct {
				msg string
				ctx context.Context
			}{
				msg: "cold start",
				ctx: lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{AwsRequestID: "cold-start-id"}),
			},
			want: struct {
				msg       string
				requestID string
				coldStart bool
			}{
				msg:       "cold start",
				requestID: "cold-start-id",
				coldStart: true,
			},
		},
		"warm start": {
			in: struct {
				msg string
				ctx context.Context
			}{
				msg: "warm start",
				ctx: lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{AwsRequestID: "warm-start-id"}),
			},
			want: struct {
				msg       string
				requestID string
				coldStart bool
			}{
				msg:       "warm start",
				requestID: "warm-start-id",
				coldStart: false,
			},
		},
		"missing lambda context": {
			in: struct {
				msg string
				ctx context.Context
			}{
				msg: "missing",
				ctx: context.Background(),
			},
			want: struct {
				msg       string
				requestID string
				coldStart bool
			}{
				msg:       "missing",
				requestID: "",
				coldStart: false,
			},
		},
	}
	var testCases []string
	for k := range tests {
		testCases = append(testCases, k)
	}
	slices.Sort(testCases)

	for _, k := range testCases {
		t.Run(k, func(t *testing.T) {
			v := tests[k]
			buf := bytes.NewBuffer(nil)
			log := logger.New(v.in.ctx, slog.NewJSONHandler(buf, nil))

			log.Info(v.in.msg)

			type LogRecord struct {
				Msg       string `json:"msg"`
				ColdStart bool   `json:"coldStart"`
				RequestID string `json:"requestId"`
				Function  struct {
					Name    *string `json:"name"`
					Version *string `json:"version"`
				} `json:"function"`
			}
			var record LogRecord
			err := json.NewDecoder(buf).Decode(&record)
			if err != nil {
				t.Fatalf("decode record: %v", err)
			}
			if record.Msg != v.want.msg {
				t.Fatalf("record msg: %s", record.Msg)
			}
			if record.ColdStart != v.want.coldStart {
				t.Fatalf("record coldStart: %v", record.ColdStart)
			}
			if record.RequestID != v.want.requestID {
				t.Fatalf("record request id: %v", record.RequestID)
			}
			if record.Function.Name == nil {
				t.Fatal("record function.name must be not nil")
			}
			if record.Function.Version == nil {
				t.Fatal("record function.version must be not nil")
			}
		})
	}
}
