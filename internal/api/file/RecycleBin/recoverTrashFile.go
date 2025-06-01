package RecycleBin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strconv"
	"strings"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func RecoverTrashFile(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()
	var renameFileName cos.RecoverTrashFile

	err := c.ShouldBind(&renameFileName)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
			"data": renameFileName,
		})
		return
	}

	for _, Key := range renameFileName.OldName {
		// 旧路径构造（只保留文件名）
		oldKey := path.Join("users", strconv.Itoa(userID)+"trash/", Key)
		// 或者更简洁的写法：
		parts := []string{"users", strconv.Itoa(userID)}
		prefix := "users/" + strconv.Itoa(userID) + "trash/"

		// 直接移除前缀
		result := strings.TrimPrefix(oldKey, prefix)
		parts = append(parts, result)
		newKey := path.Join(parts...)

		fmt.Println(newKey, oldKey)

		source := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, oldKey)
		_, _, err = client.Object.Copy(context.Background(), newKey, source, nil)
		if err == nil {
			_, err = client.Object.Delete(context.Background(), oldKey)
			if err != nil {
				c.JSON(500, gin.H{
					"code": 500, "message": err.Error(),
				})
				// Error
				return
			}
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "文件已经恢复了",
	})
}
