package Describe

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func ForDescribe(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()
	var filename cos.FileNames
	// 生成预览 URL
	previewURLs := make([]string, 0, len(filename.FileName))

	err := c.ShouldBind(&filename)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
		})
		return
	}

	for _, res := range filename.FileName {
		objectKey := "users/" + strconv.Itoa(userID) + "/" + res
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
			log.Printf("Failed to generate presigned URL for %s: %v", objectKey, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Failed to generate download URL",
			})
			return
		}

		// 拼接预览参数
		previewURLs = append(previewURLs, presignedURL.String())
	}
	
	c.JSON(200, gin.H{
		"previewURLs": previewURLs,
	})
}
