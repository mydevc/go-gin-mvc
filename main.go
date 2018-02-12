package main

import (
	"gin_api/routers"
	"github.com/gin-gonic/gin"
	"gin_api/common"
	"github.com/robfig/cron"
	"gin_api/database"
)


func main() {
	gin.SetMode(gin.DebugMode)

	//gin工程实例 *gin.Engine
	router := routers.Router

	//session设置，注意顺序,路由配置前
	common.SetSession(router)

	//路由初始化
	routers.SetupRouter()

	//模板路径设置
	common.SetTemplate(router)

	c := cron.New()
	c.AddFunc("*/30 * * * * *", database.CheckDb)
	c.Start()

	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")
}


