package ctrls

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"go-gin-mvc/utils"
	"strconv"
)

func RedisSetAction(ctx *gin.Context) {


	rds := utils.RedisPool.Get();

	count, _ := redis.Int(rds.Do("GET", "count"))
	count++
	rds.Do("SET", "count", count)
	ctx.JSON(200, gin.H{
		"message": count,
	})

}


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


func CookieAction(ctx *gin.Context) {

	var count int
	v, _ := ctx.Cookie("count")

	count,_ = strconv.Atoi(v)
	count++
	utils.SetCookie(ctx, "count", strconv.Itoa(count), 3600*24)

	ctx.JSON(200, gin.H{"count": v})

}

func TokenAction(ctx *gin.Context)  {

	fmt.Println(ctx.HandlerName())

	//ctx.HTML(http.StatusOK, "token_test.html", gin.H{
	//	"title": "token测试",
	//	"token":csrf.GetToken(ctx),
	//})
}