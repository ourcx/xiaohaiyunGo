package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/models/cos"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func Ourl(c *gin.Context) {
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

	key := "users/" + strconv.Itoa(userID) + "/" + filename.FileName
	ourl := client.Object.GetObjectURL(key)
	fmt.Println(ourl)

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"url":     ourl,
	})
}
