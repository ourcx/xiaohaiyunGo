package Tag

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"strconv"
	"xiaohaiyun/internal/models/Tag"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func SetTag(c *gin.Context) {
	// 从环境变量读取敏感配置
	userID := cosFile.GetID(c).ID
	// 初始化 COS 客户端
	client := cosFile.Client()

	opt := &cos.ObjectPutTaggingOptions{
		TagSet: []cos.ObjectTaggingTag{
			{
				Key:   "test_k2",
				Value: "test_v2",
			},
		},
	}
	var tag Tag.Tag
	// 初始化 COS 客户端
	c.ShouldBind(&tag)
	name := "users" + "/" + strconv.Itoa(userID) + "/" + tag.Name
	_, err := client.Object.PutTagging(context.Background(), name, opt)
	if err != nil {
		//ERROR
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
