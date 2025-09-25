package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"xiaohaiyun/configs"
	"xiaohaiyun/internal/api/v1"
	"xiaohaiyun/internal/app"

	// swagger 依赖
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	//_ "xiaohaiyun/docs" // 重要：导入生成的docs文件夹
)

// @title           小海云 API文档
// @version         1.0
// @description     小海云后端API接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.xiaohaiyun.com/support
// @contact.email  support@xiaohaiyun.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 加载配置文件
	err := config.LoadConfig()
	if err != nil {
		log.Errorf("配置文件加载错误: %v", err)
		return
	}

	// 初始化所有模块
	err = app.InitializeAll()
	if err != nil {
		log.Errorf("模块初始化错误: %v", err)
		return
	}

	r := gin.Default()

	// 添加Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化路由模块
	v1.SetupRoutes(r, app.Engine)

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"port":   config.Conf.App.Port,
		})
	})

	log.Infof("服务启动在端口: %d", config.Conf.App.Port)
	err = r.Run(fmt.Sprintf(":%d", config.Conf.App.Port))
	if err != nil {
		log.Errorf("服务启动错误: %v", err)
		return
	}
}
