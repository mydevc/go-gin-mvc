# go-gin-mvc
基于gin-gonic/gin 框架搭建的MVC架构的基础项目空架子，未经商用检验，请谨慎参考。

## 此项目集成了小型网站开发常用的功能:  
- 基于redis连接池存储的session操作;
```
详见代码：route/router.go
关键代码：
//不用连接池
//store,err:=redis.NewStore(10,"tcp","rs1.rs.youfa365.com:6379","",[]byte("asfajfa;lskjr2"))
//使用连接池
store, err := redis.NewStoreWithPool(utils.RedisPool, []byte("as&8(0fajfa;lskjr2"))
```
 
- 基于redis连接池存储的cache操作;  

- 基于xorm的数据库操作,主从分离配置，支持多从库配置,鲜活连接定时PING操作，集成xorm/cmd；

- 基于rabbitmq的队列应用，注意生产者与消费者队列名称的一致性

>>多个任务可发送到一个队列，也可以灵活应用一个队列一个任务;
生产者与消费者消息传递的是序列化的结构体，结构体由生产者提供，并自行反序列化操作；  

>>消费者：
console/queue_daemon.go
（队列需要单独控制台命令启动，与http服务独立[避免相互影响]；）

>>生产者（这里仅测试使用，正式应用一般在web代码中）
console/send_single.go

>>WEB测试生产者：
http://localhost:8080/queue



- csrf防跨站攻击,此功能集成此中间件完成[点这里](https://github.com/utrack/gin-csrf),更多[中间件](https://github.com/gin-gonic/contrib)。
  这里要重点说一下，utrack/gin-csrf这个中间件没有加白名单机制排除一些例外，这在实际应用中是很常见的，尤其是对外合作接口中。
  我把此中间件代码集成到我自己的代码中来，把白名单功能补上了。这里直接用包名+函数名来定位，在配置文件conf/csrf_except.ini中配置，
  key值随意，不空，不重复即可，因不是实时读取，修改后需要重启web服务才生效。

- 数据验证，可自定义友好错误提示，[更多实例参考](https://github.com/thedevsaddam/govalidator)；
```
user_controller.go
func UserAddAction(ctx *gin.Context) {

	if ctx.Request.Method == "POST" {

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

		name := ctx.PostForm("name")
		age_str := ctx.PostForm("age")
		age_int, _ := strconv.Atoi(age_str)

		model_users.UserAdd(&entitys.User{Name: name, Age: age_int})

		ctx.HTML(http.StatusOK, "user_add.html", gin.H{
			"name": name,
			"age":  age_int,
		})

	} else {
		ctx.HTML(http.StatusOK, "user_add.html", gin.H{
			"title": "用户添加",
		})
	}

}
```  

- INI配置文件读取操作，可分别加载多个配置文件;
```
utils/config.go
package utils

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
)

var Config *ini.File
var CsrfExcept *ini.File
var RootPath string

func init()  {
	RootPath="/Users/fuxiaojun/data/golang/gopath/src/go-gin-mvc"
	var err error
	Config, err = ini.Load(RootPath+"/conf/config.ini");
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	CsrfExcept, err = ini.Load(RootPath+"/conf/csrf_except.ini");
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
}

使用：
m_dsn := utils.Config.Section("mysql_master").Key("master").String()

slaves := utils.Config.Section("mysql_slave").Keys()
for _, s_dsn := range slaves {
	fmt.Println(s_dsn.String())
}

```

- 定时任务;
```
main.go
//定时程序启动
c := cron.New()
//数据库状态检查
c.AddFunc("*/600 * * * * *", models.DbCheck)
c.Start()
```

---

###cmd/xorm安装注意事项
```
正常安装命令：  
go get github.com/go-xorm/cmd/xorm
但会报错，有两个包无法安装，cloud.google.com/go/civil，golang.org/x/crypto/md4，移步到https://github.com/GoogleCloudPlatform/google-cloud-go下载相应的包
GOPATH目录下新建cloud.google.com 文件夹（与github.com同级）
cloud.google.com/go/civil
golang.org/x/crypto/md4

进入cmd/xorm 运行命令 go build
查看帮助 xorm help reverse
```
##xorm生成struct
```
xorm reverse mysql "root:12345678@tcp(dbm1.rs.youfa365.com:3306)/test?charset=utf8" .
```


## 技术支持
- Mail:mydev@126.com