package file

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// RemoveFolder 删除文件夹
func RemoveFolder(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()
	filename := cos.FileName{}

	err := c.ShouldBind(&filename)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    filename,
		})
		return
	}

	key := "users/" + strconv.Itoa(userID) + "/" + filename.FileName + "/"
	_, err = client.Object.Delete(context.Background(), key)

	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    "删除操作出现错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    key,
	})
}
