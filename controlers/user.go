package controlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"fmt"
	"gin_api/common"
	"gin_api/models"

	"github.com/utrack/gin-csrf"
	//"log"

	"github.com/thedevsaddam/govalidator"
	"github.com/kr/beanstalk"
	"time"
)

var (
	TubeName1 string = "channel1"
	TubeName2 string = "channel2"
)

//队列任务添加
func Queue(c *gin.Context)  {

	fname := "P-1"
	tubeName:=TubeName1

	b, err := beanstalk.Dial("tcp", "192.168.1.168:11300")
	if err != nil {
		panic(err)
	}
	defer b.Close()

	b.Tube.Name = tubeName
	b.TubeSet.Name[tubeName] = true
	fmt.Println(fname, " [Producer] tubeName:", tubeName, " b.Tube.Name:", b.Tube.Name)

	for i := 0; i < 1; i++ {
		msg := `{"send_mail":"mail_body channel1"}`

		b.Put([]byte(msg), 30, 0, 120*time.Second)
		fmt.Println(fname, " [Producer] beanstalk put body:", msg)
	}

	fname = "P-2"
	tubeName =TubeName2

	b.Tube.Name = tubeName
	b.TubeSet.Name[tubeName] = true
	fmt.Println(fname, " [Producer] tubeName:", tubeName, " b.Tube.Name:", b.Tube.Name)

	for i := 0; i < 1; i++ {
		msg := `{"send_mail":"mail_body channel2"}`

		b.Put([]byte(msg), 30, 0, 120*time.Second)
		fmt.Println(fname, " [Producer] beanstalk put body:", msg)
	}

	b.Close()
}

func Validator(c *gin.Context) {

	rules := govalidator.MapData{
		"username": []string{"required", "between:3,8"},
		"email":    []string{"required", "min:4", "max:20", "email"},
		"web":      []string{"url"},
		"phone":    []string{"digits:11"},
		"agree":    []string{"bool"},
		"dob":      []string{"date"},
	}

	messages := govalidator.MapData{
		"username": []string{"required:用户名不能为空","between:3到8位"},
		"phone":    []string{"digits:手机号码为11位数字"},
		"web":    []string{"url:必须是URL格式"},
		"email":    []string{"required:Email不能为空", "min:Email4-20位", "max:Email4-20位", "email:必须为Email格式"},

	}

	opts := govalidator.Options{
		Request:         c.Request,        // request object
		Rules:           rules,    // rules map
		Messages:        messages, // custom message map (Optional)
		RequiredDefault: false,     // all the field to be pass the rules
	}
	v := govalidator.New(opts)
	e := v.Validate()
	c.JSON(200,e )
}

func Redis(c *gin.Context) {
	common.PutCache("youfa_name", "oqierqkj38lkj;lk.", 3)
	common.PutCache("youfa_name", "oqierqkj38lkj;lk.", 60)
	username := common.GetCache("youfa_name")
	fmt.Println("=====" + username)
}

func Model(c *gin.Context) {
	models.GetOrders()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
	})
}

func UserLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"csrf_token": csrf.GetToken(c),
	})
}

func UserLoginAuth(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "UserLoginAuth",
		"name":  c.PostForm("user_name"),
	})
}

func Session(c *gin.Context) {

	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	c.JSON(200, gin.H{"count": count})

}

func Cookie(c *gin.Context) {
	common.SetCookie(c, "name", "fuxiaojun00", 3600)
	v, _ := c.Cookie("name")
	c.JSON(200, gin.H{"count": v})
}
