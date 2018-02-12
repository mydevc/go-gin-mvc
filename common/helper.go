package common

import (
	"github.com/go-ini/ini"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/utrack/gin-csrf"
)

var configPath string

func init() {
	//配置文件设置
	SetConfigPath( os.Getenv("GINPATH") + "/conf/config.ini")
}

func SetConfigPath(path string)  {
	configPath = path
}

func GetConfig(section string,key string) *ini.Key {
	cfg, _ := ini.InsensitiveLoad(configPath)

	v, _ := cfg.Section(section).GetKey(key)
	return  v
}

func SetCookie(c *gin.Context,name string, value string, maxAge int)  {
	domain := GetConfig("cookie","domain").String()
	c.SetCookie(name,value,maxAge,"/",domain,false,true)
}

func SetTemplate(engine *gin.Engine) {
	engine.LoadHTMLGlob(os.Getenv("GINPATH") + "/views/*/*")
}

func SetSession(engine *gin.Engine) {

	address:=GetConfig("redis","address").String()
	sessionsecret:=GetConfig("session","sessionsecret").String()
	sessionname:=GetConfig("session","sessionname").String()

	store, _ := sessions.NewRedisStore(10, "tcp", address, "", []byte(sessionsecret))
	engine.Use(sessions.Sessions(sessionname, store))

	//csrf支持 form表单：_csrf，url参数：_csrf，Heder参数：X-CSRF-TOKEN 或 X-XSRF-TOKEN
	//忽略的请求："GET", "HEAD", "OPTIONS"
	engine.Use(csrf.Middleware(csrf.Options{
		Secret: GetConfig("session","csrfscret").String(),
		ErrorFunc: func(c *gin.Context){
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
}