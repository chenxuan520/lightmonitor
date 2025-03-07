package main

import (
	"fmt"

	"github.com/chenxuan520/lightmonitor/internal/config"
	"github.com/chenxuan520/lightmonitor/internal/middlerware"
	"github.com/chenxuan520/lightmonitor/internal/service"
	"github.com/chxuan520/lightmonitor/internal/monitor"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Init()

	g := gin.New()
	w := service.NewWeb()

	api := g.Group("/api")
	api.Use(middlerware.PasswdAuth())
	{
		api.GET("/list", w.ListTasks)
		api.POST("/confirm", w.ConfirmTask)
		// 消息实时推送中台
		api.POST("/notify", w.NotifyMsg)
	}
	g.StaticFile("/", "./asserts/index.html")

	g.Run(fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))
}
