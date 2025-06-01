package file

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type FolderName struct {
	FolderName string `json:"folderName"`
}

// AddFolder 新建文件夹的功能
// 将在对应的储存的文件内新建文件夹
func AddFolder(c *gin.Context) {
	// 从环境变量读取敏感配置
	userID := cosFile.GetID(c).ID
	//用户信息
	var folderName FolderName
	// 初始化 COS 客户端
	client := cosFile.Client()

	err := c.ShouldBind(&folderName)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    folderName})
		return
	}

	// 要创建的文件夹路径（必须以 / 结尾）
	folderKey := "users/" + strconv.Itoa(userID) + "/" + folderName.FolderName + "/"

	//exists, err := cosFile.CheckFolderExists(client, folderKey)
	//if err != nil {
	//	c.JSON(500, gin.H{
	//		"msg":  err.Error(),
	//		"code": 500,
	//		"data": "检查文件失败",
	//	})
	//	return
	//}

	//if exists {
	//	c.JSON(500, gin.H{
	//		"code": 500,
	//		"msg":  "文件夹已存在，跳过创建:" + folderKey,
	//	})
	//	return
	//}

	// 创建文件夹（上传空对象）
	_, err = client.Object.Put(context.Background(), folderKey, strings.NewReader(""), nil)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "创建文件夹失败: %v\n" + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "success", "data": "文件夹创建成功:" + folderKey})
}
