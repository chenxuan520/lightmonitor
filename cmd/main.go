package main

import (
	"fmt"

	"github.com/chenxuan520/lightmonitor/internal/config"
	"github.com/chenxuan520/lightmonitor/internal/middlerware"
	"github.com/chenxuan520/lightmonitor/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	g := gin.New()
	w := service.NewWeb()

	g.Use(middlerware.Cors())
	api := g.Group("/api")
	api.Use(middlerware.PasswdAuth())
	{
		// monitor 网页监控
		monitorApi := api.Group("/monitor")
		{
			monitorApi.GET("/list", w.ListMonitorTask)
			monitorApi.POST("/confirm", w.ConfirmMonitorTask)
		}
		// 消息实时推送中台
		api.POST("/notify", w.NotifyMsg)
		api.POST("/list", w.ListTask)
		api.POST("/cancel", w.CancelTask)
	}
	g.StaticFile("/", "./asserts/index.html")

	g.Run(fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))
}
