package file

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type FileInfo struct {
	FileName string `json:"file_name"`
}

func Exist(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	var fileInfo FileInfo
	err := c.ShouldBind(&fileInfo)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}
	client := cosFile.Client()
	key := "users/" + strconv.Itoa(userID) + "/" + fileInfo.FileName
	ok, err := client.Object.IsExist(context.Background(), key)
	if err == nil && ok {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  true,
		})
	} else if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  false,
		})
	}
}
