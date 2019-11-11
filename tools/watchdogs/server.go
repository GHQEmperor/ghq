package watchdogs

import (
	"errors"
	"github.com/GHQEmperor/ghq"
	"github.com/gogf/gf/container/gmap"
	"net/http"
	"net/rpc"
	"time"
)

var ContMap ContainerMap

type ContainerMap map[string]ContainerFunc
type ContainerFunc struct {
	Container   Container
	TimeoutFunc func(context interface{})
	DeadFunc    func(context interface{})
	DeadReturn  func(context interface{})
	After       time.Duration
}

// 本地调用,创建 Container
// 并添加处理函数
func AddContainer(key string, timeoutFunc, deadFunc, deadReturn func(context interface{}), after time.Duration) error {
	if ContMap == nil {
		ContMap = make(map[string]ContainerFunc)
	}
	_, ok := ContMap[key]
	if ok {
		return errors.New("this container is exist")
	}
	ContMap[key] = ContainerFunc{
		Container: Container{
			_gmap: gmap.New(),
		},
		TimeoutFunc: timeoutFunc,
		DeadFunc:    deadFunc,
		DeadReturn:  deadReturn,
		After:       after,
	}
	return nil
}

func Run() {
	if err := rpc.Register(&ContMap); err != nil {
		panic(err)
	}
	//if err := rpc.Register(WatchDog{}); err != nil {
	//	panic(err)
	//}
	rpc.HandleHTTP()
	port, ok := ghq.GetConfig("rpc_dog_port")
	if !ok {
		panic("rpc_dog_port is not found")
	}
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

type Message struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type FeedDogReq struct {
	Container string `json:"container"`
	Key       string `json:"key"`
	Context   string `json:"context"`
}

func (c *ContainerMap) RpcFeedDog(req FeedDogReq, res *Message) error {
	if req.Container == "" || req.Key == "" || req.Context == "" {
		res = &Message{
			Status: 10999,
			Error:  "参数不完全",
		}
		return errors.New("参数不完全")
	}

	cont, ok := ContMap[req.Container]
	if !ok {
		res = &Message{
			Status: 10999,
			Error:  "无此 Container",
		}
		return errors.New("无此 Container")
	}
	if err := cont.Container.FeedDog(req); err != nil {
		if err != ErrNoThisDog {
			res = &Message{
				Status: 10001,
				Error:  err.Error(),
			}
			return err
		}
		if err := cont.Container.NewDog(req.Key, req.Context, cont.After,
			cont.TimeoutFunc,
			cont.DeadFunc,
			cont.DeadReturn); err != nil {
			res = &Message{
				Status: 10001,
				Error:  err.Error(),
			}
			return err
		}
	}
	return nil
}
