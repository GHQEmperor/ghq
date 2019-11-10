package watchdogs

import (
	"errors"
	"github.com/gogf/gf/container/gmap"
	"sync"
	"time"
)

var (
	ErrNoThisDog      = errors.New("the dog is not exist")
	ErrAnotherRunning = errors.New("another dog in this key")
)

type Container struct {
	_gmap *gmap.Map
}

func NewContainer() *Container {
	return &Container{gmap.New()}
}

func (c *Container) NewDog(key, context interface{}, after time.Duration, timeoutFunc, deadFunc, deadReturn func(context interface{})) error {
	newDog := WatchDog{
		key:            key,
		context:        context,
		foodChan:       make(chan struct{}, 1),
		deadChan:       make(chan struct{}, 1),
		deadReturnChan: make(chan struct{}, 1),
		destroyChan:    make(chan struct{}, 1),
		isClose:        false,
		mutex:          sync.Mutex{},
		container:      c,
	}
	if err := c.setNotRun(key, &newDog); err != nil {
		//
		return err
	}
	go newDog.dogRunning(timeoutFunc, deadFunc, deadReturn, after)
	_ = c.FeedDog(key)
	return nil
}

func (c *Container) FeedDog(key interface{}) error {
	dogInf := c._gmap.Get(key)
	if dogInf != nil {
		dog := dogInf.(*WatchDog)
		if key != dog.GetKey() {
			return ErrAnotherRunning
		}
		dog.feed()
		return nil
	}
	return ErrNoThisDog
}

func (c *Container) DeadDog(key interface{}) error {
	dogInf := c._gmap.Get(key)
	if dogInf != nil {
		dog := dogInf.(*WatchDog)
		dog.dead()
		return nil
	}
	return nil
}

func (c *Container) DeadDogReturn(key interface{}) error {
	dogInf := c._gmap.Get(key)
	if dogInf != nil {
		dog := dogInf.(*WatchDog)
		dog.deadReturn()
	}
	return nil
}

func (c *Container) Destroy(key interface{}) {
	dogInf := c._gmap.Get(key)
	if dogInf != nil {
		dog := dogInf.(*WatchDog)
		dog.destroy()
	}
	return
}

func (c *Container) removeDog(key interface{}) {
	dogInf := c._gmap.Get(key)
	if dogInf != nil {
		c._gmap.Remove(key)
	}
	return
}

func (c *Container) getNotRun(key interface{}) *WatchDog {
	return c._gmap.Get(key).(*WatchDog)
}

func (c *Container) setNotRun(key interface{}, run *WatchDog) error {
	dogInf := c._gmap.Get(key)
	if dogInf != nil {
		run.close()
		return ErrAnotherRunning
	}
	newDog := run
	c._gmap.Set(key, newDog)
	return nil
}
