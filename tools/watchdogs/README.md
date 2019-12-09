# 看门狗
## 单机使用
### 创建狗
- func NewContainer() *Container        创建 信号容器 
- func (c *Container) NewDog()          创建 看门狗 
                                        
|参数|说明|
|:-----         |-----                   |
| key           | 需要定时的内容           |
| context       | 需要定时的内容(备用)     |
| after         | 超时时长                |
| timeoutFunc   | 超时后执行的函数         |
| deadFunc      | dead信号执行的函数       |
| deadReturn    | deadReturn信号执行的函数 |

### 发送信号
- func (c *Container) FeedDog(key interface{}) error    
- func (c *Container) DeadDog(key interface{}) error
- func (c *Container) DeadDogReturn(key interface{}) error
- func (c *Container) Destroy(key interface{}) error 


## rpc多机使用
### server 创建 container
- func AddContainer(key string, timeoutFunc, deadFunc, deadReturn func(wg *WatchDog), after time.Duration) error

|参数|说明|
|:-----         |-----                   |
| key           | 需要定时的内容           |
| context       | 需要定时的内容(备用)     |
| after         | 超时时长                |
| timeoutFunc   | 超时后执行的函数         |
| deadFunc      | dead信号执行的函数       |
| deadReturn    | deadReturn信号执行的函数 |

### client 进行 feed,dead,deadReturn和destroy 操作
- func NewRpcClient() (*Client, error)      创建client
- func (c *Client) FeedDog(cont, key, context string) (result Message, err error)  
- func (c *Client) DeadDog(cont, key string) (result Message, err error) 
- func (c *Client) DeadReturn(cont, key string) (result Message, err error) 
- func (c *Client) Destroy(cont, key string) (result Message, err error)
