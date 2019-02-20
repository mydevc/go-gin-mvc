package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"go-gin-mvc/utils"
	"go-gin-mvc/models"
	"go-gin-mvc/route"
)

func main() {

	gin.SetMode(gin.DebugMode)

	models.GetMaster()
	models.GetSlave()

	r := route.SetupRouter()



	//定时程序启动
	c := cron.New()
	//数据库状态检查
	c.AddFunc("*/600 * * * * *", models.DbCheck)
	c.Start()

	port := utils.Config.Section("system").Key("http_port").String()
	r.Run(":" + port)
}
