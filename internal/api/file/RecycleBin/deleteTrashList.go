package RecycleBin

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func DeleteTrashList(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	var deleteFile cos.DeleteFile
	client := cosFile.Client()

	err := c.ShouldBind(&deleteFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	key := "users/" + strconv.Itoa(userID) + "trash/" + deleteFile.DeleteName
	_, err = client.Object.Delete(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": "服务器错误",
		})
		panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "已删除对象：" + key,
	})
}
