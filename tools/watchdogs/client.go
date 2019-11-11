package watchdogs

import (
	"github.com/GHQEmperor/ghq"
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
	client *rpc.Client
}

func NewRpcClient() (*Client, error) {
	router := ghq.New()
	if err := router.LoadConfig(); err != nil {
		return nil, err
	}
	serverIP, ok := ghq.GetConfig("server_ip")
	if !ok {
		panic("server_ip is not exist")
	}
	client, err := rpc.DialHTTP("tcp", serverIP)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) FeedDog(cont, key, context string) (result Message, err error) {
	return c.call(cont, key, context, feedDogMethod)
}

func (c *Client) DeadDog(cont, key string) (result Message, err error) {
	return c.call(cont, key, "", deadDogMethod)
}

func (c *Client) DeadReturn(cont, key string) (result Message, err error) {
	return c.call(cont, key, "", deadReturnMethod)
}

func (c *Client) Destroy(cont, key string) (result Message, err error) {
	return c.call(cont, key, "", destroyMethod)
}

func (c *Client) call(cont, key, context string, method string) (result Message, err error) {
	if err := c.client.Call(container+method, DogReq{
		Container: cont,
		Key:       key,
		Context:   context,
	}, &result); err != nil {
		return result, err
	}
	return result, nil
}
