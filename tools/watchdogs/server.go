package watchdogs

import (
	"errors"
	"fmt"
	"github.com/GHQEmperor/ghq"
	"github.com/gogf/gf/container/gmap"
	"net/http"
	"net/rpc"
	"time"
)

var ContMap ContainerMap

type ContainerMap map[string]*ContainerFunc
type ContainerFunc struct {
	Container   Container
	TimeoutFunc func(wg *WatchDog)
	DeadFunc    func(wg *WatchDog)
	DeadReturn  func(wg *WatchDog)
	After       time.Duration
}

// 本地调用,创建 Container
// 并添加处理函数
func AddContainer(key string, timeoutFunc, deadFunc, deadReturn func(wg *WatchDog), after time.Duration) error {
	if ContMap == nil {
		ContMap = make(map[string]*ContainerFunc)
	}
	_, ok := ContMap[key]
	if ok {
		return errors.New("this container is exist")
	}
	ContMap[key] = &ContainerFunc{
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
	router := ghq.New()
	if err := router.LoadConfig(); err != nil {
		panic(err)
	}
	if err := rpc.Register(&ContMap); err != nil {
		panic(err)
	}
	rpc.HandleHTTP()
	port, ok := ghq.GetConfig("rpc_dog_port")
	if !ok {
		panic("rpc_dog_port is not found")
	}
	fmt.Printf("[ watchdogs running in port: %v ]\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

type Message struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

type DogReq struct {
	Container string `json:"container"`
	Key       string `json:"key"`
	Context   string `json:"context"`
}

func (c *ContainerMap) RpcFeedDog(req DogReq, res *Message) error {
	if req.Container == "" || req.Key == "" || req.Context == "" {
		res.Status = 10999
		res.Error = "参数不完全"
		return nil
	}

	cont, ok := (*c)[req.Container]
	if !ok {
		res.Status = 10999
		res.Error = "无此 Container"
		return nil
	}
	if err := cont.Container.FeedDog(req.Key); err != nil {
		//fmt.Println("cont.Container.FeedDog(req) error:", err)
		if err != ErrNoThisDog {
			res.Status = 10001
			res.Error = err.Error()
			return nil
		}
		if err := cont.Container.NewDog(req.Key, req.Context, cont.After, cont.TimeoutFunc, cont.DeadFunc, cont.DeadReturn); err != nil {
			res.Status = 10001
			res.Error = err.Error()
			//fmt.Println("cont.Container.NewDog error:", err)
			return nil
		}
	}
	res.Status = 10000
	return nil
}

func (c *ContainerMap) RpcDeadDog(req DogReq, res *Message) error {
	if req.Container == "" || req.Key == "" {
		res.Status = 10999
		res.Error = "参数不完全"
		return nil
	}
	cont, ok := (*c)[req.Container]
	if !ok {
		res.Status = 10999
		res.Error = "无此 Container"
		return nil
	}
	if err := cont.Container.DeadDog(req.Key); err != nil {
		res.Status = 10999
		res.Error = err.Error()
		return nil
	}

	res.Status = 10000
	return nil
}

func (c *ContainerMap) RpcDeadReturn(req DogReq, res *Message) error {
	if req.Container == "" || req.Key == "" {
		res.Status = 10999
		res.Error = "参数不完全"
		return nil
	}
	cont, ok := (*c)[req.Container]
	if !ok {
		res.Status = 10999
		res.Error = "无此 Container"
		return nil
	}
	if err := cont.Container.DeadDogReturn(req.Key); err != nil {
		res.Status = 10999
		res.Error = err.Error()
		return nil
	}

	res.Status = 10000
	return nil
}

func (c *ContainerMap) RpcDestroy(req DogReq, res *Message) error {
	if req.Container == "" || req.Key == "" {
		res.Status = 10999
		res.Error = "参数不完全"
		return nil
	}
	cont, ok := (*c)[req.Container]
	if !ok {
		res.Status = 10999
		res.Error = "无此 Container"
		return nil
	}
	if err := cont.Container.Destroy(req.Key); err != nil {
		res.Status = 10999
		res.Error = err.Error()
		return nil
	}

	res.Status = 10000
	return nil
}
