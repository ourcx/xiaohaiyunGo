package cosFileUtils

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"os"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	modelsCOS "xiaohaiyun/internal/models/cos"
	"xiaohaiyun/internal/repositories"
	"xiaohaiyun/internal/utils"
)

// GenerateSecureUploadURL 生成临时密钥和预签名URL给前端
func GenerateSecureUploadURL(c *gin.Context) {
	filename := modelsCOS.FileName{}
	//获取用户的JWT信息
	userID := GetID(c).ID
	c.ShouldBind(&filename)
	//阻止多次请求？
	
	// ===================== 2. 配置COS客户端 =====================
	secretID := os.Getenv("COS_MAIN_SECRET_ID")   // 从环境变量读取
	secretKey := os.Getenv("COS_MAIN_SECRET_KEY") // 从环境变量读取
	bucketURL := "http://cdn.norubias.top"

	u, err := url.Parse(bucketURL)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "解析存储桶URL失败: " + err.Error()})
		return
	}

	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
			Transport: &debug.DebugRequestTransport{}, // 启用调试
		},
	})

	objectKey := fmt.Sprintf("users/%d/%s", userID, filename.FileName)
	fmt.Println(objectKey)
	presignedURL, err := client.Object.GetPresignedURL(
		context.Background(),
		http.MethodPut,
		objectKey,
		secretID,
		secretKey,
		5*time.Minute,
		nil, // 可设置Content-Type等参数
	)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "生成预签名URL失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"url":    presignedURL.String(),
			"method": "PUT",
			"path":   objectKey,
		},
	})
}

func GetID(c *gin.Context) *models.UserReq {
	claimsValue, exists := c.Get("processed_data")
	if !exists {
		c.JSON(500, gin.H{"code": 500, "msg": "用户信息丢失"})
		return nil
	}

	claims, ok := claimsValue.(*utils.UserClaims)
	if !ok {
		c.JSON(500, gin.H{"code": 500, "msg": "用户信息解析失败"})
		return nil
	}

	// ===================== 1. 获取用户数据 =====================
	s := repositories.NewUserRepository(app.Engine)
	userData, err := s.GetUserByEmail(claims.Email)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "用户查询失败: " + err.Error()})
		return nil
	}
	return userData

}
