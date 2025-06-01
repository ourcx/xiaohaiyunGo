package Tag

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/models/Tag"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func GetTag(c *gin.Context) {
	// 从环境变量读取敏感配置
	userID := cosFile.GetID(c).ID
	var tag Tag.Tag
	// 初始化 COS 客户端
	c.ShouldBind(&tag)
	client := cosFile.Client()
	name := "users" + "/" + strconv.Itoa(userID) + "/" + tag.Name
	res, _, err := client.Object.GetTagging(context.Background(), name)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
		})
		//ERROR
	}
	fmt.Println(res)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    res.TagSet,
	})
}
