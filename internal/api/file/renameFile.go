package file

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// RenameFile 重命名一个不大于5G的文件的名字
func RenameFile(c *gin.Context) {
	userID := cosFile.GetID(c).ID

	client := cosFile.Client()

	var renameFileName cos.RenameFile
	err := c.ShouldBind(&renameFileName)
	if err != nil {
		return
	}

	oldKey := "users" + "/" + strconv.Itoa(userID) + "/" + renameFileName.OldName
	newKey := "users" + "/" + strconv.Itoa(userID) + "/" + renameFileName.NewName

	ctx := context.Background()

	// 1. 复制对象到新名称
	//opt := &tencentcos.ObjectCopyOptions{
	//	XCosCopySourceIfNoneMatch: "*", // 仅当目标文件不存在时复制
	//}
	source := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, oldKey)
	_, _, err = client.Object.Copy(ctx, newKey, source, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "复制错误",
		})
	}

	// 2. 删除旧对象
	_, err = client.Object.Delete(ctx, oldKey)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "删除错误",
		})
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "已经重命名为" + newKey,
	})
}
