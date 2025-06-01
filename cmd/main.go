package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"xiaohaiyun/configs"
	"xiaohaiyun/internal/api/v1"
	"xiaohaiyun/internal/app"
)

func main() {
	// 加载配置文件
	err := config.LoadConfig()
	if err != nil {
		log.Error("配置文件加载错误: %v", err)
		return
	}

	// 初始化所有模块
	err = app.InitializeAll()
	if err != nil {
		log.Error("模块初始化错误: %v", err)
		return
	}

	r := gin.Default()
	//初始化路由模块
	v1.SetupRoutes(r, app.Engine)

	err = r.Run(fmt.Sprintf(":%d", config.Conf.App.Port))
	if err != nil {
		log.Error("服务启动错误: %v", err)
		return
	}
}
