package relationship

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/RelationShip"
	"xiaohaiyun/internal/repositories"
	"xiaohaiyun/internal/utils"
)

var (
	friendEmail RelationShip.FriendEmail
)

func ApplyByEmail(c *gin.Context) {

	//申请添加好友，把好友的邮箱发过来
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

	masterEmail := claims.Email
	//先从中间件上下文拿到用户是谁吧

	//从gin里面拿到我传的数据
	if err := c.ShouldBind(&friendEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := repositories.NewUserRepository(app.Engine)
	masterTable, err := s.GetUserByEmail(masterEmail)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	friendTable, err := s.GetUserByEmail(friendEmail.Email)

	friendCreate := RelationShip.Friend{
		FriendId:  friendTable.ID,
		UserId:    masterTable.ID,
		Status:    "accepted",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = app.Engine.Insert(friendCreate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": "已存在"})
		return
	}

	friendCreate = RelationShip.Friend{
		FriendId:  masterTable.ID,
		UserId:    friendTable.ID,
		Status:    "accepted",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = app.Engine.Insert(friendCreate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"friend": friendCreate,
			"email":  friendEmail.Email,
			"status": "accepted",
		},
	})

}
