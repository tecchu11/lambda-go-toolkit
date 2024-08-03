package localclient_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tecchu11/lambda-invoker-go/internal/localclient"
)

func TestIntegration(t *testing.T) {
	type Event struct {
		Msg string
	}
	handler := func(event Event) (string, error) {
		return event.Msg, nil
	}
	go func() {
		t.Setenv("_LAMBDA_SERVER_PORT", "9000")
		lambda.Start(handler)
	}()
	time.Sleep(100 * time.Millisecond)

	client, err := localclient.New(9000)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := json.Marshal(Event{Msg: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(res) != `"ok"` {
		t.Fatal(string(res))
	}

}
