package localclient

import (
	"fmt"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

// Client interacts aws lambda run on local.
type Client struct {
	c *rpc.Client
}

// New creates [Client].
func New(port int) (*Client, error) {
	c, err := rpc.Dial("tcp", strings.Join([]string{"localhost", strconv.Itoa(port)}, ":"))
	if err != nil {
		return nil, fmt.Errorf("rpc dial\n%w", err)
	}
	return &Client{c: c}, nil
}

func (c *Client) Do(payload []byte) ([]byte, error) {
	var (
		reqPing messages.PingRequest
		resPing messages.InvokeResponse
	)
	err := c.c.Call("Function.Ping", reqPing, &resPing)
	if err != nil {
		return nil, fmt.Errorf("ping\n%w", err)
	}
	if resPing.Error != nil {
		return nil, fmt.Errorf("ping response\n%v", resPing.Error)
	}
	var (
		reqInvoke = messages.InvokeRequest{
			Payload: payload,
		}
		resInvoke messages.InvokeResponse
	)
	err = c.c.Call("Function.Invoke", reqInvoke, &resInvoke)
	if err != nil {
		return nil, fmt.Errorf("invoke\n%w", err)
	}
	if resInvoke.Error != nil {
		return nil, fmt.Errorf("invoke response\n%w", resInvoke.Error)
	}
	return resInvoke.Payload, nil
}
