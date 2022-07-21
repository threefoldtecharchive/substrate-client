package substrate

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
)

type RetryingClient struct {
	client.Client
}

func NewRetryingClient(cl client.Client) RetryingClient {
	return RetryingClient{cl}
}

func (c *RetryingClient) Call(result interface{}, method string, args ...interface{}) error {
	err := c.Client.Call(result, method, args...)
	if err == nil {
		return nil
	}
	if err.Error() == "use of closed network connection" {
		err = c.Client.Call(result, method, args...)
	}
	return err
}

func NewSubstrateAPI(url string) (*gsrpc.SubstrateAPI, error) {
	cl, err := client.Connect(url)
	if err != nil {
		return nil, err
	}
	rcl := NewRetryingClient(cl)
	newRPC, err := rpc.NewRPC(&rcl)
	if err != nil {
		return nil, err
	}

	return &gsrpc.SubstrateAPI{
		RPC:    newRPC,
		Client: &rcl,
	}, nil
}
