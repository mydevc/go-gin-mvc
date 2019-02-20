package route

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go-gin-mvc/ctrls"
	"go-gin-mvc/ctrls/users"
	"go-gin-mvc/middleware/csrf"
	"go-gin-mvc/utils"
	"net/http"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()

	session_name := utils.Config.Section("session").Key("sessionname").String()
	domain := utils.Config.Section("session").Key("sessiondomain").String()
	maxage, _ := utils.Config.Section("session").Key("sessiongcmaxlifetime").Int()
	csrfscret := utils.Config.Section("session").Key("csrfscret").String()
	//Session配置，Redis存储
	//store,err:=redis.NewStore(10,"tcp","rs1.rs.youfa365.com:6379","",[]byte("asfajfa;lskjr2"))
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

	Router.Use(csrf.Middleware(csrf.Options{
		Secret: csrfscret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

}

func SetupRouter() *gin.Engine {

	//静态目录配置
	public := utils.Config.Section("router").Key("public").String()
	Router.Static("/public", public)

	//模板
	view_path := utils.Config.Section("router").Key("view_path").String()
	Router.LoadHTMLGlob(view_path)

	//Session测试-Redis存储
	Router.GET("/session", ctrls.SessionAction)
	//Cookie测试
	Router.GET("/cookie", ctrls.CookieAction)
	//Redis测试
	Router.GET("/redis", ctrls.RedisSetAction)

	//单个用户,显示数据提交混合，偷懒写法
	Router.GET("/user/add", users.UserAddAction)
	Router.POST("/user/add", users.UserAddAction)

	//单个用户
	Router.GET("/user/show", users.UserShowAction)
	//用户列表
	Router.GET("/user/index", users.UserListAction)

	//用户编辑
	Router.GET("/user/edit/:id", users.UserEditAction)

	Router.GET("/protected", ctrls.TokenAction)

	Router.POST("/protected", func(c *gin.Context) {
		c.String(200, "CSRF token is valid")
	})

	Router.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "404.html", "")
	})

	return Router
}
