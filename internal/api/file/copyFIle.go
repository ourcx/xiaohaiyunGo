package file

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"path/filepath"
	"strconv"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// CopyFile 拷贝对象
func CopyFile(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()
	var renameFileName cos.MoveFile

	err := c.ShouldBind(&renameFileName)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
			"data": renameFileName,
		})
		return
	}

	for _, Key := range renameFileName.OldName {
		// 提取 Key 的文件名部分（去掉前面的目录结构）
		fileName := filepath.Base(Key) // 如果 Key = "sub/file.txt"，返回 "file.txt"
		// 旧路径构造（只保留文件名）
		oldKey := path.Join("users", strconv.Itoa(userID), Key)
		// 或者更简洁的写法：
		parts := []string{"users", strconv.Itoa(userID)}
		if renameFileName.NewName != "" {
			parts = append(parts, renameFileName.NewName)
		}
		newKey := path.Join(append(parts, fileName)...)

		fmt.Println(newKey, oldKey)

		source := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, oldKey)
		_, _, err = client.Object.Copy(context.Background(), newKey, source, nil)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": "文件已经移动了",
	})
}
