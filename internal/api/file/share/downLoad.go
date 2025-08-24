package share

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"time"
	"xiaohaiyun/internal/models/share"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// DownLoadShare Describe 检测文件预览url
func DownLoadShare(c *gin.Context) {
	client := cosFile.Client()
	filename := share.DownLoadUrlShare{}
	// 生成预览 URL

	err := c.ShouldBind(&filename)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
		})
		return
	}

	fmt.Println(filename)
	opt := &cos.PresignedURLOptions{
		Query: &url.Values{
			"response-content-disposition": []string{
				"attachment; filename=\"" + filename.FileName + "\"",
			},
			// 添加这行强制二进制流类型
			"response-content-type": []string{"application/octet-stream"},
		},
	}

	objectKey := "share/" + filename.UUID + "/" + filename.FileName
	presignedURL, err := client.Object.GetPresignedURL(
		context.Background(),
		http.MethodGet,
		objectKey,
		os.Getenv("COS_MAIN_SECRET_ID"),
		os.Getenv("COS_MAIN_SECRET_KEY"),
		10*time.Minute,
		opt,
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
