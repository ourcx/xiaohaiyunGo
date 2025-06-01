package file

import (
	"github.com/gin-gonic/gin"
	"strconv"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// 一个简单的激活函数，返回一个序号
func Init(c *gin.Context) {
	userTable := cosFile.GetID(c)
	c.JSON(200, gin.H{
		"code": 200,
		"id":   strconv.Itoa(userTable.ID),
	})
}
