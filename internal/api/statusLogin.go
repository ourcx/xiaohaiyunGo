package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"xiaohaiyun/internal/utils"
)

// JWT 结构体定义（确保字段标签匹配）
type JWT struct {
	Jwt string `json:"jwt"` // 确保与前端传参字段一致
}

// JwtStatus 检查JWT有效性
func JwtStatus(c *gin.Context) {
	var jwts JWT

	// 绑定请求参数（传递指针）
	if err := c.ShouldBindJSON(&jwts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 检查JWT是否为空
	if jwts.Jwt == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "JWT不能为空",
		})
		return
	}

	// 解析JWT
	claims, err := utils.ParseUserJWT(jwts.Jwt)
	if err != nil {
		log.Printf("JWT解析失败: %v", err) // 记录详细错误
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "无效的JWT: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "有效的JWT",
		"data":    claims,
	})
}
