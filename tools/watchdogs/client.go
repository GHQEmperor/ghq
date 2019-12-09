package watchdogs

import (
	"net/rpc"
)

const (
	container        = "ContainerMap."
	feedDogMethod    = "RpcFeedDog"
	deadDogMethod    = "RpcDeadDog"
	deadReturnMethod = "RpcDeadReturn"
	destroyMethod    = "RpcDestroy"
)

type Client struct {
	client    *rpc.Client
	container string
}

func NewRpcClient(serverIP, container string) (*Client, error) {
	client, err := rpc.DialHTTP("tcp", serverIP)
	if err != nil {
		return nil, err
	}

	return &Client{client: client, container: container}, nil
}

func (c *Client) FeedDog(key, context string) (result Message, err error) {
	return c.call(key, context, feedDogMethod)
}

func (c *Client) DeadDog(key string) (result Message, err error) {
	return c.call(key, "", deadDogMethod)
}

func (c *Client) DeadReturn(key string) (result Message, err error) {
	return c.call(key, "", deadReturnMethod)
}

func (c *Client) Destroy(key string) (result Message, err error) {
	return c.call(key, "", destroyMethod)
}

func (c *Client) call(key, context string, method string) (result Message, err error) {
	if err := c.client.Call(container+method, DogReq{
		Container: c.container,
		Key:       key,
		Context:   context,
	}, &result); err != nil {
		return result, err
	}
	return result, nil
}
