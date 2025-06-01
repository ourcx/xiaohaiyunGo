package userFound

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func SendLog(c *gin.Context) {
	userEmail := cosFile.GetID(c).Email
	var data []models.UserLogin
	err := app.Engine.Table("user_logins").Where("email=?", userEmail).Find(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请求失败",
			"err":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}
