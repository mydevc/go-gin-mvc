package users

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"go-gin-mvc/entitys"
	"go-gin-mvc/models/users"
	"net/http"
	"strconv"
)

func UserIndexAction(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "user_add.html", gin.H{
		"title": "Main website",
	})
}

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

func UserShowAction(ctx *gin.Context) {

	type Item struct {
		Ono string `json:ono`
		Oid int    `json:oid`
	}
	type Refund struct {
		Ono     string `json:ono`
		Item    int    `json:item`
		Content string `json:content`
		Imgs    string `json:imgs`
		Status  string `json:status`
	}

	type AlipayRemoteReqStruct struct {
		Ono         string   `json:ono`
		OrderItem   []Item   `json:item`
		OrderRefund []Refund `json:refund`
	}

	var m AlipayRemoteReqStruct
	m.Ono = "12345"
	m.OrderItem = append(m.OrderItem, Item{Ono: "Shanghai_VPN", Oid: 1})
	m.OrderItem = append(m.OrderItem, Item{Ono: "Beijing_VPN", Oid: 2})
	for i := 1; i < 6; i++ {
		str := []byte("物品")
		str = strconv.AppendInt(str, int64(i), 10)
		orderi := Item{Ono: string(str), Oid: i}
		m.OrderItem = append(m.OrderItem, orderi)
	}
	bytes, _ := json.Marshal(m)

	//ctx.String(200,string(bytes))

	fmt.Println(string(bytes))

	var js AlipayRemoteReqStruct
	err := json.Unmarshal(bytes, &js)
	if err != nil {
		fmt.Printf("format err:%s\n", err.Error())
		return
	}

	for _, v := range js.OrderItem {

		ctx.String(200, string(v.Ono)+"\n")
	}

	//users, has := model_users.UserOne()
	//
	//if has {
	//	ctx.JSON(200, gin.H{
	//		"message": users.Name,
	//	})
	//} else {
	//	ctx.JSON(201, gin.H{
	//		"error": "数据为空",
	//	})
	//}
}

func UserListAction(ctx *gin.Context) {

	users, err := model_users.UserList()

	if err == nil {
		for _, user := range users {

			fmt.Println(user.Name)
		}

	}
}

func UserEditAction(c *gin.Context) {

}
