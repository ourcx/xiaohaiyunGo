package userAuth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"xiaohaiyun/internal/utils"
	logs "xiaohaiyun/log"
)

// LoggerMiddleware dd
// 示例中间件1：日志记录,记录相关的操作
func LoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next() // 继续执行后续中间件或路由处理函数
	latency := time.Since(start)
	_ = logs.Writer("log/loggerMiddleware.log", "\n有用户进行了获取数据的操作"+strconv.FormatInt(int64(latency), 10))
}

// AuthMiddleware 示例中间件2：身份验证这个东西是接受用户的jwt并进行认证
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")

	//这个请求头给的是JWT
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}
	// 验证Token逻辑（此处简化）
	claims, err := utils.ParseUserJWT(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	processedData, err := utils.ExtractUserClaims(claims)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "data": "解析jwt出现问题"})
		return
	}
	c.Set("processed_data", processedData)
	//从jwt得到用户的信息，用来提供给数据库查询和验证
	c.Next()
}
