package userFound

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/models/userData"
	"xiaohaiyun/internal/utils/profiles"

	cosFile "xiaohaiyun/internal/utils/cos"
)

// UpUserReqName UpUserReq 这个东西的功能就是更新用户user_req表的名字
func UpUserReqName(c *gin.Context) {
	userReq := cosFile.GetID(c)

	var userName userData.UserName
	err := c.ShouldBind(&userName)
	if err != nil {
		c.JSON(http.StatusHTTPVersionNotSupported, gin.H{"message": "err,错误的传参"})
		return
	}

	err = profiles.UpUserName(userReq.Email, userName)
	if err != nil {
		c.JSON(http.StatusHTTPVersionNotSupported, gin.H{"message": "err,更新发生错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "newName": userName})

}
