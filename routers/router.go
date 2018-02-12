package routers

import (
	"github.com/gin-gonic/gin"
	"gin_api/controlers"
	"os"
	"time"
	"github.com/gin-contrib/cache/persistence"
	"gin_api/common"
	"github.com/gin-contrib/cache"
	"github.com/utrack/gin-csrf"
	"fmt"
	"github.com/gin-contrib/sessions"


	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
)

var DB = make(map[string]string)
var Router *gin.Engine


func init() {
	Router = gin.Default()
}

func SetupRouter() {
	//缓存存储初始化操作
	store := persistence.NewRedisCache(common.GetConfig("redis","address").String(), "", time.Second)

	//静态目录配置
	Router.Static("/static", os.Getenv("GINPATH")+"/static")

	//缓存例子
	Router.GET("/queue", controlers.Queue)

	//缓存例子
	Router.GET("/cache_model", cache.CachePage(store, time.Minute,controlers.Model))

	//数据库查询测试
	Router.GET("/model", controlers.Model)

	//redis测试
	Router.GET("/redis", controlers.Redis)

	//session测试
	Router.GET("/session", controlers.Session)

	//cookie测试
	Router.GET("/cookie", controlers.Cookie)

	// Ping test
	Router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	//表单提交 + csrf测试
	Router.GET("/login", controlers.UserLogin)
	Router.POST("/login", controlers.UserLoginAuth)

	//表单验证
	Router.GET("/validator", controlers.Validator)


	//csrf测试
	Router.GET("/protected", func(c *gin.Context){
		session:=sessions.Default(c)
		fmt.Println(session.Get("csrfSalt"))

		c.String(200, csrf.GetToken(c))
	})

	Router.POST("/protected", func(c *gin.Context){
		c.String(200, "CSRF token is valid")
	})


	//jwt test
	var mysupersecretpassword = "unicornsAreAwesome"

	public := Router.Group("/api")

	public.GET("/", func(c *gin.Context) {
		// Create the token
		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		// Set some claims
		token.Claims = jwt_lib.MapClaims{
			"Id":  "Christopher",
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		}
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(mysupersecretpassword))
		if err != nil {
			c.JSON(500, gin.H{"message": "Could not generate token"})
		}
		c.JSON(200, gin.H{"token": tokenString})
	})

	private := Router.Group("/api/private")
	private.Use(jwt.Auth(mysupersecretpassword))

	/*
		Set this header in your request to get here.
		Authorization: Bearer `token`
	*/

	private.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from private"})
	})

	private.GET("/go", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from private go"})
	})




	// Get user value
	Router.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := DB[user]
		if ok {
			c.JSON(200, gin.H{"user": user, "value": value})
		} else {
			c.JSON(200, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := Router.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	//authorized.GET("/login", controlers.UserLogin)

	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			DB[user] = json.Value
			c.JSON(200, gin.H{"status": "ok"})
		}
	})

}
