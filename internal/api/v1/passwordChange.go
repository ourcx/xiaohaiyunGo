package v1

//更改密码的

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/models/pwd"
	"xiaohaiyun/internal/utils"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// PasswordChange 这个是一个密码改变的操作
func PasswordChange(c *gin.Context) {
	userTable := cosFile.GetID(c)
	//得到身份信息

	var change pwd.Change

	err := c.ShouldBind(&change)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":  err,
			"code": 400,
		})
		return
	}
	//拿到相关的参数

	hashNew, err := utils.HashPassword(change.New)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error(), "msg": "新密码hash化错误"})
		return
	}

	//获得旧的hash

	var table models.UserReq
	_, err = app.Engine.Table("user_req").Where("email = ?", userTable.Email).Get(&table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error(), "msg": "提取旧密码错误"})
		return
	}

	//拿到数据库储存的值
	fmt.Println(utils.CheckPassword(change.Old, table.Password))

	if utils.CheckPassword(change.Old, table.Password) {
		table.Password = hashNew
		_, err := app.Engine.Table("user_req").Where("email = ?", userTable.Email).Update(&table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error(), "msg": "更新密码错误"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"err": "旧输入密码错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "success", "code": http.StatusOK, "msg": "密码成功更改"})

}
