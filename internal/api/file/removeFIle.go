package file

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strconv"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// RemoveFile 删除对象，其实做到移除到回收站比较好
func RemoveFile(c *gin.Context) {
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

	oldKey := "users/" + strconv.Itoa(userID) + "/" + filename.FileName
	parts := []string{"users", strconv.Itoa(userID) + "trash/"}
	newKey := path.Join(parts...) + "/" + filename.FileName

	fmt.Println(newKey, oldKey)

	source := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, oldKey)
	_, _, err = client.Object.Copy(context.Background(), newKey, source, nil)
	if err == nil {
		_, err = client.Object.Delete(context.Background(), oldKey)

		if err != nil {
			c.JSON(500, gin.H{
				"code":    500,
				"message": err.Error(),
				"data":    "删除操作出现错误",
			})
			return
		}
	} else {
		return
	}
	
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    oldKey,
	})
}
