# ghq
## 简介
~~~
这是一个golang菜鸟的极简web框架。（基于内置http包）
只有基本的读取配置文件，添加路由，获取表单参数，静态文件功能。
写代码一定记得接收错误 error，不要忽略！！！
写代码一定记得接收错误 error，不要忽略！！！
写代码一定记得接收错误 error，不要忽略！！！
代码不规范，找错两行泪！
~~~

## 使用说明
~~~
目录结构
|--Demo
	|--controller
	|--model
	|--static
		|--js
		|--css
		|--image
	|--views
		|--index.html
	|--config.json
	|--main.go
~~~

### 配置文件读取

 配置文件  **config.json**  总是在项目下的第一层

~~~
func (r *Router) GetConfig(configName string) (config string,ok bool)

Demo: 获取端口，如果配置文件中有port字段则ok为true
config,ok := r.GetConfig("port")
~~~

### 获取Form和PostForm参数

 获取Form和PostForm参数使用相同方法，有两种返回值分别为string和int

~~~
func (rw *RW) GetString(key string) (value string)
Demo: 如key为username
username := rw.GetString("username")

func (rw *RW) GetInt(key string) (value int, err error)
Demo: 如key为age
age ,err := rw.GetInt("age")
if err != nil {
    // 参数中没有age或age中值不为数字
}
~~~

### 添加路由

 Method有9种，分别为：

~~~
MethodGet
MethodHead
MethodPost
MethodPut
MethodPatch
MethodDelete
MethodConnect
MethodOptions
MethodTrace
~~~

 对应方法:

~~~
func (r *Router) Get(uri string, function FuncRW)
func (r *Router) Head(uri string, function FuncRW) 
func (r *Router) Post(uri string, function FuncRW)
等...依次对应

Demo: 
	uri为: /index ，
	function为:	func(rw ghq.RW) {
					rw.WriteHTML("index.html")
				 })
g.Get("/index", func(rw ghq.RW) {
		rw.WriteHTML("index.html")
})
即访问 localhost:****/index
~~~

### 静态文件

 添加静态文件路径

~~~
func (r *Router) SetStaticFile(uri, dirPath string)

Demo: 如uri 为 /static/， dirPath 为 static
r.SetStaticFile("/static/","static")
即访问 localhost:****/static
~~~

### 返回数据

#### html

~~~
func (rw *RW) WriteHTML(fileName string) (err error)

Demo: fileName为views目录下文件的文件名。如 index.html
if err := rw.WriteHTML("index.html");err != nil {
    // 处理错误
}
~~~

#### json

~~~
func (rw *RW) WriteJson(data interface{}) (err error)

Demo: data 类型可为 struct 和 map，如
type FileSaver struct {
	FileName string
	FilePath string
	FileMD5 string
}
var file = FileSaver{
	FileName:"*",
	FilePath:"**",
	FileMD5:"***",
}
if err := rw.WriteJson(file);err != nil {
    // 处理错误
}

~~~

#### xml

~~~
------
~~~

### 创建Router对象

~~~
func New() *Router

Demo:
func main() {
    g := ghq.New() // 创建Router对象
	g.Get("/index", func(rw ghq.RW) {
		rw.WriteHTML("index.html")
	})
	g.SetStaticFile("/static/","static")
	if err := g.Run(); err != nil {
		panic(err)
	}
}
~~~

### 服务启动

~~~
func (r *Router) Run() error
将所需的配置信息填入config.json，加入自己所需的路由和静态文件，
即可启动服务。错误信息均从error传出。
~~~

