package userFound

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/userData"
	"xiaohaiyun/internal/repositories"
	"xiaohaiyun/internal/utils"
)

func GetProfiles(c *gin.Context) {

	var data userData.GetProfile

	claimsValue, exists := c.Get("processed_data")
	fmt.Println(claimsValue)

	//从中间件拿到自己搞定的上下文
	if !exists {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "用户信息丢失",
		})
		return // 必须终止执行
	}

	// 正确的类型断言方式
	claims, ok := claimsValue.(*utils.UserClaims)
	if !ok {
		c.JSON(500, gin.H{
			"code": 500,
		})
		return
	}

	//初始化那个表
	s := repositories.NewUserRepository(app.Engine)
	table, err := s.GetTableByEmail(claims.Email)
	userReq, err := s.GetUserByEmail(claims.Email)
	user, ok := table.(*userData.UserProfile)

	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": "数据库错误",
		})
		return
	}

	//包装我传回去的值
	data = userData.GetProfile{
		Email: claims.Email, Signature: user.Signature, AvatarUrl: user.AvatarUrl, UserName: userReq.Name,
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
	})
	//正确就返回那个
}

func PostProfiles(c *gin.Context) {
	var profile userData.PostProfile
	claimsValue, exists := c.Get("processed_data")
	//还是从中间件拿到上下文
	//感觉可以独立出去封装成一个函数

	if !exists {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "用户信息丢失",
		})
		return // 必须终止执行
	}

	// 正确的类型断言方式
	claims, ok := claimsValue.(*utils.UserClaims)
	if !ok {
		c.JSON(500, gin.H{
			"code": 500,
		})
		return
	}

	//从gin里面拿到我传的数据
	if err := c.ShouldBind(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//初始化那个表
	s := repositories.NewUserRepository(app.Engine)
	table, err := s.GetUserByEmail(claims.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	//包装我要插入到数据库里面的数据
	profiles := userData.UserProfile{
		UserId:    table.ID,
		Signature: profile.Signature,
		AvatarUrl: profile.AvatarUrl,
	}

	has, err := app.Engine.Table("user_profiles").Where("user_id=?", table.ID).Exist()
	if err != nil {
		return
	}

	//执行插入操作
	if has {
		_, err = app.Engine.Table("user_profiles").Where("user_id=?", table.ID).Update(&profiles)
		_, err = app.Engine.Table("url_share").Where("user_id=?", table.ID).Update(&profiles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "提示": "更新失败"})
			return
		}

	} else {
		_, err = app.Engine.Table("user_profiles").Where("user_id=?", table.ID).Insert(&profiles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "提示": "更新失败"})
			return
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"data": profile,
	})

	//返回结果和修改的签名hh
}
