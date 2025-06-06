package Describe

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// Describe 检测文件预览url
func Describe(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()
	filename := cos.FileName{}
	// 生成预览 URL

	err := c.ShouldBind(&filename)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
		})
		return
	}

	fmt.Println(filename)

	objectKey := "users/" + strconv.Itoa(userID) + "/" + filename.FileName
	presignedURL, err := client.Object.GetPresignedURL(
		context.Background(),
		http.MethodGet,
		objectKey,
		os.Getenv("COS_MAIN_SECRET_ID"),
		os.Getenv("COS_MAIN_SECRET_KEY"),
		10*time.Minute,
		nil,
	)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		panic(err)
	}

	// 拼接预览参数
	previewURL := presignedURL.String()

	fmt.Println("文档预览链接:", previewURL)
	c.JSON(200, gin.H{
		"previewURL": previewURL,
	})
}

//问题应该在于无法通过默认域名获得文件预览这里，我其他的配置都是正常的
