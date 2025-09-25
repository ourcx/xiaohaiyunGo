package Login

import (
	"github.com/gin-gonic/gin"
	"net/http"
	utilsBack "xiaohaiyun/internal/api/Backend/utils"
	"xiaohaiyun/internal/utils"
)

func Req(c *gin.Context) {
	//好像不是很需要这个！
	//改成添加账户吧
	var user ListFolder
	_ = c.ShouldBind(&user)

	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userId int64
	userId = utilsBack.GenerateSnowflakeID()
	//使用雪花算法生成id
	//添加到数据库
	err = utilsBack.InsetUser(user.Email, hashPassword, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
	//简单的登录注册服务
}
