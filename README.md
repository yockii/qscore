# qscore
QS通用库

lowcode分支，大部分简单的增删改查可以无需再写，仅建立相关的model、domain、service、controller即可，同时，亦可在service、controller中自行实现其他逻辑

若有事务操作，可在使用
```go
database.DB.Transaction(func(tx *gorm.DB) error {
	// 此中调用不同service的方法，并在其中传入tx，即可在同一事务中操作
})
```

# 使用方法
## 1. 引入库
```shell
go get -u github.com/yockii/qscore@lowcode
```

## 2. 使用库
### 建立model
作为与数据库沟通的桥梁

BaseModel中已经定义了ID主键
```go
type SomeModel struct {
	common.BaseModel
	// ...自定义字段，注意本框架使用gorm，故需要使用gorm的tag
	Name        string         `json:"name" gorm:"size:100;comment:名称"`
	CreateTime  int64          `json:"createTime" gorm:"autoCreateTime:milli"`
    UpdateTime  int64          `json:"updateTime" gorm:"autoUpdateTime:milli"`
    DeleteTime  gorm.DeletedAt `json:"deleteTime,omitempty" gorm:"index"`
}

// 按照需要实现common.Model的接口，可参考common.BaseModel的实现
```
### 建立domain
作为接口接收的对象
```go
type SomeDomain struct {
	model.SomeModel
    common.BaseDomain
    // ...自定义字段
	CreateTimeCondition *server.TimeCondition `json:"createTimeCondition,omitempty"`
}
```

### 建立service
作为业务逻辑的实现，需要继承BaseService，基础的增删改查可以无需再写，若有其他逻辑，可自行实现
```go
var SomeService = newSomeService()
type someService struct {
    common.BaseService[*model.SomeModel]
}
func newSomeService() *someService {
    s := new(someService)
	s.BaseService = common.BaseService[*model.SomeModel] {
		Service: s,
    }
	return s
}
```

### 建立controller
作为对外接口实现，需要继承BaseController，基础的增删改查可以无需再写，若有其他逻辑，可自行实现
```go
var SomeController = &someController{}
type someController struct {
    common.BaseController[*model.SomeModel, *domain.SomeDomain]
}
func init() {
    // 初始化路由
    r := server.Group("/api/v1/some")
    r.Get("/list", SomeController.List)
    r.Get("/detail", SomeController.Detail)
    r.Post("/add", SomeController.Add)
    ...
}
func (*someController) GetService() common.Service[*model.SomeModel] {
    return service.SomeService
}
```

### 处理入口程序
不需要的部分可以删除
```go
func main() {
	// 初始化日志
    config.InitialLogger()
	
	// 初始化数据库
    database.Initial()
    defer database.Close()

	// 雪花算法初始节点信息
    _ = util.InitNode(0)
	
	// 初始化数据，自行编码
	
	// 启动定时任务
	task.Start()
	defer task.Stop()
	
	// 初始化路由
	controller.InitRouter() // 自行编码，仅为了引用出controller包，也可尝试直接导入 _ "github.com/xxx/xxx/controller"
	
    // 启动服务
	for {
        err := server.Start()
        if err != nil {
            logger.Errorln(err)
        }
    }
}
```

### 配置文件
放置于程序同级目录的conf/config.toml中
```toml
userTokenExpire = 86400

[server]
port = 8080

[database]
driver = "mysql"
host = "localhost"
user = "root"
password = "root"
db = "xxxx"
port = 3306
prefix = "t_"
showSql = true

[logger]
level = "debug"

[redis]
host = "localhost"
port = 6379
password = ""
db = 0
app = "xxxx"

```


# 下一步计划
- [ ] 优化代码
- [ ] 利用go generate生成代码，进一步简化开发