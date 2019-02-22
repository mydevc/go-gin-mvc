go-gin-mvc 项目地址：[https://github.com/mydevc/go-gin-mvc](https://github.com/mydevc/go-gin-mvc)

基于golang语言的gin-gonic/gin 框架搭建的MVC架构的基础项目空架子供初学者学习参考，如果你是从PHP语言转过来的，一定会非常喜欢这个架构。

# 此项目集成了小型网站开发常用的功能:

## 1、基于redis连接池存储的cache操作;
utils/helper.go
```
// 建立连接池
RedisPool = &redis.Pool{
	MaxIdle:     max_idle,
	MaxActive:   max_active,
	IdleTimeout: idle_timeout,
	Wait:        true,
	Dial: func() (redis.Conn, error) {
		con, err := redis.Dial("tcp", host,
			//redis.DialPassword(conf["Password"].(string)),
			redis.DialDatabase(database),
			redis.DialConnectTimeout(timeout),
			redis.DialReadTimeout(timeout),
			redis.DialWriteTimeout(timeout))
		if err != nil {
			return nil, err
		}
		return con, nil
	},
}
```
ctrls/test_controller.go
```
//Redis测试
func RedisSetAction(ctx *gin.Context) {
	rds := utils.RedisPool.Get();
	count, _ := redis.Int(rds.Do("GET", "count"))
	count++
	rds.Do("SET", "count", count)
	ctx.JSON(200, gin.H{
		"message": count,
	})
}
```
## 2、基于redis连接池存储的session操作;
#####注意这里的连接池是独立于cache操作redis的连接池，需单独配置参数。
```
//不用连接池
//store,err:=redis.NewStore(10,"tcp","rs1.baidu.com:6379","",[]byte("asfajfa;lskjr2"))
//使用连接池
store, err := redis.NewStoreWithPool(utils.RedisPool, []byte("as&8(0fajfa;lskjr2"))
```
session配置
```
store, err := redis.NewStoreWithPool(utils.RedisPool, []byte("as&8(0fajfa;lskjr2"))
store.Options(sessions.Options{
	"/",
	domain,
	maxage,
	false, //https 时使用
	true,  //true:JS脚本无法获取cookie信息
})
if err != nil {
	// Handle the error. Probably bail out if we can't connect.
	fmt.Println("redis.NewStore error")
}
Router.Use(sessions.Sessions(session_name, store))
```
Sesssion测试
```
//Sesssion测试
func SessionAction(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count += 1
	}
	session.Set("count", count)
	session.Save()
	ctx.JSON(200, gin.H{"count": count})
}
```

## 3、基于xorm的数据库操作,主从分离配置，支持多从库配置,鲜活连接定时PING操作，集成xorm/cmd；
部分代码展示，完整代码见：models/orm.go
```
//从库添加
slaves := utils.Config.Section("mysql_slave").Keys()
for _, s_dsn := range slaves {
	_dbs, err := xorm.NewEngine("mysql", s_dsn.String())
	_dbs.SetMaxIdleConns(10)
	_dbs.SetMaxOpenConns(200)
	_dbs.ShowSQL(true)
	_dbs.ShowExecTime(true)

	if err!=nil {
		fmt.Println(err)
	}else{
		dbs = append(dbs, _dbs)
	}
}
```
## 4、基于rabbitmq的队列应用，注意生产者与消费者队列名称的一致性

####多个任务可发送到一个队列，也可以灵活应用一个队列一个任务; 生产者与消费者消息传递的是序列化的结构体，结构体由生产者提供，并自行反序列化操作；
- 消费者： console/queue_daemon.go （队列需要单独控制台命令启动，与http服务独立[避免相互影响]；）
```
subscriber := new(jobs.Subscribe)
forever := make(chan bool)
q := queue.NewQueue()
//队列执行的任务需要注册方可执行
q.PushJob("Dosome",jobs.HandlerFunc(subscriber.Dosome))
//q.PushJob("Fusome",jobs.HandlerFunc(subscriber.Fusome))
//提前规划好队列，可按延时时间来划分。可多个任务由一个队列来执行，也可以一个任务一个队列，一个队列可启动多个消费者
go q.NewShareQueue("SomeQueue")
//go q.NewShareQueue("SomeQueue")
//go q.NewShareQueue("SomeQueue")
defer q.Close()
<-forever
```
- 控制台生产者（这里仅测试使用，正式应用一般在web代码中） console/send_single.go
```
forever := make(chan bool)
go func() {
	for i := 0; i < 1000000; i++ {
		queue.NewSender("SomeQueue", "Dosome", jobs.Subscribe{Name: "We are doing..." + strconv.Itoa(i)}).Send()
	}
}()
defer queue.SendConn.Close()
<-forever
```

- WEB测试生产者： http://localhost:8080/queue
```
//队列生产者测试
func QueueAction(ctx *gin.Context)  {
	queue.NewSender("SomeQueue", "Dosome", jobs.Subscribe{Name: "We are doing..."}).Send()
}
```
## 5、csrf防跨站攻击

此功能集成此中间件完成[点这里](https://github.com/utrack/gin-csrf),更多[中间件](https://github.com/gin-gonic/contrib)。 这里要重点说一下，utrack/gin-csrf这个中间件没有加白名单机制排除一些例外，这在实际应用中是很常见的，尤其是对外合作接口中。 我把此中间件代码集成到我自己的代码中来，把白名单功能补上了。这里直接用包名+函数名来定位，在配置文件conf/csrf_except.ini中配置， key值随意，不空，不重复即可，因不是实时读取，修改后需要重启web服务才生效。
- 应用代码route/router.go
```
Router.Use(csrf.Middleware(csrf.Options{
	Secret: csrfscret,
	ErrorFunc: func(c *gin.Context) {
		c.String(400, "CSRF token mismatch")
		c.Abort()
	},
}))
```
- 补丁代码：middleware/csrf/csrf.go
```
fn := c.HandlerName();
fn = fn[strings.LastIndex(fn, "/"):]

for _, action := range IgnoreAction {
	if (strings.Contains(fn, action)) {
		fmt.Println(action)
		c.Next()
		return
	}
}
```

## 6、数据验证，可自定义友好错误提示，[更多实例参考](https://github.com/thedevsaddam/govalidator)；
ctrls/users/user_controller.go 中 UserAddAction()
```
//参数检验
rules := govalidator.MapData{
	"name": []string{"required", "between:3,8"},
	"age":  []string{"digits:11"},
}

messages := govalidator.MapData{
	"name": []string{"required:用户名不能为空", "between:3到8位"},
	"age":  []string{"digits:手机号码为11位数字"},
}

opts := govalidator.Options{
	Request:         ctx.Request, // request object
	Rules:           rules,       // rules map
	Messages:        messages,    // custom message map (Optional)
	RequiredDefault: false,       // all the field to be pass the rules
}
v := govalidator.New(opts)
e := v.Validate()

//校验结果判断
if len(e)>0 {
	ctx.JSON(200, e)
	return
}
```
## 7、INI配置文件读取操作，可分别加载多个配置文件;
utils/config.go
```
RootPath="/Users/baidu/data/golang/gopath/src/go-gin-mvc"
var err error
Config, err = ini.Load(RootPath+"/conf/config.ini");
if err != nil {
	fmt.Printf("Fail to read file: %v", err)
	os.Exit(1)
}
```
应用实例
models/orm.go
```
//从库添加
slaves := utils.Config.Section("mysql_slave").Keys()
for _, s_dsn := range slaves {
	_dbs, err := xorm.NewEngine("mysql", s_dsn.String())
    ...
}
```

## 8、定时任务;
main.go
```
//定时程序启动
c := cron.New()
//数据库状态检查
c.AddFunc("*/600 * * * * *", models.DbCheck)
c.Start()
```
## cmd/xorm安装注意事项
##### 正常安装命令： 
```
go get github.com/go-xorm/cmd/xorm
```
但会报错，有两个包无法安装，
```
cloud.google.com/go/civil
golang.org/x/crypto/md4
```
移步到https://github.com/GoogleCloudPlatform/google-cloud-go下载相应的包
GOPATH目录下新建cloud.google.com 文件夹（与github.com同级）
```
cloud.google.com/go/civil
golang.org/x/crypto/md4
```
进入cmd/xorm 运行命令
```
go build
```
查看帮助 xorm help reverse

## xorm生成struct
```
xorm reverse mysql "root:12345678@tcp(dbm1.baidu.com:3306)/test?charset=utf8" .
```
项目地址：[https://github.com/mydevc/go-gin-mvc](https://github.com/mydevc/go-gin-mvc)

## EMAIL：mydev@126.com