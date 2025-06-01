package Tag

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/models/Tag"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func DeleteTag(c *gin.Context) {
	// 从环境变量读取敏感配置
	userID := cosFile.GetID(c).ID
	var tag Tag.Tag
	// 初始化 COS 客户端
	c.ShouldBind(&tag)
	client := cosFile.Client()
	name := "users" + "/" + strconv.Itoa(userID) + "/" + tag.Name
	_, err := client.Object.DeleteTagging(context.Background(), name)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
		})
		//ERROR
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
