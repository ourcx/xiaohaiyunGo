package userData

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/utils"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type UserData struct {
	NewPwd string `json:"newPwd"`
	Code   string `json:"code"`
}

func ForgetPwd(c *gin.Context) {
	userEmail := cosFile.GetID(c).Email
	var userData UserData

	if err := c.ShouldBind(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":  err,
			"code": 400,
		})
	}

	err := utils.ValidateCodeFromReq(userData.Code, userEmail+":forgetPwd")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":  err,
			"code": 400,
			"msg":  "验证码错误",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(userData.NewPwd) // 实现你的哈希函数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//哈希加密用户的密码
	userData.NewPwd = hashedPassword

	_, err = app.Engine.Table("user_req").Where("email=?", userEmail).Update(map[string]interface{}{"password": hashedPassword})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "密码完成更新",
	})

}
