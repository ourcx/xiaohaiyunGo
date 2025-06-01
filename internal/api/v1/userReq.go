package v1

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/models/userData"
	"xiaohaiyun/internal/repositories"
	"xiaohaiyun/internal/utils"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func UserReq(c *gin.Context) {
	var userReq models.UserReq
	var userReqByEmail models.UserReqByEmail
	var profiles userData.UserProfile
	if err := c.ShouldBind(&userReqByEmail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userReq = models.UserReq{
		ID:       userReqByEmail.ID,
		Name:     userReqByEmail.Name,
		Email:    userReqByEmail.Email,
		Password: userReqByEmail.Password,
	}

	hashedPassword, err := utils.HashPassword(userReq.Password) // 实现你的哈希函数
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//哈希加密用户的密码

	userReq.Password = hashedPassword

	s := repositories.NewUserRepository(app.Engine)
	email, err := s.GetUserByEmail(userReq.Email) //等于空就说明数据库里面有这个邮箱了

	if email != nil {
		c.JSON(http.StatusConflict, gin.H{
			"code": http.StatusConflict,
			"msk":  "邮箱已经注册了,点击这里跳转登录页面",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(userReqByEmail)
	err = utils.ValidateCodeFromReq(userReqByEmail.Code, userReqByEmail.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "msg": "验证码出现问题"})
		return
	}

	//先验证是否重复再进行注册信息的提交
	_, err = app.Engine.Insert(userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = app.Engine.Table("user_req").Where("email=?", userReq.Email).Get(&userReq)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "msg": "数据库问题",
		})
		return
	}

	client := cosFile.Client()
	folderKey := "users/" + strconv.Itoa(userReq.ID) + "/"
	_, err = client.Object.Put(context.Background(), folderKey, strings.NewReader(""), nil)

	profiles = userData.UserProfile{
		UserId:    userReq.ID,
		Signature: "",
		AvatarUrl: "",
	}
	_, err = app.Engine.Table("user_profiles").Insert(&profiles)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "msg": "数据库问题",
		})
		return
	}

	userHash := models.SafeUserReq{
		Email:    userReq.Email,
		Name:     userReq.Name,
		Password: userReq.Password,
		//加密后的密码
	}

	c.JSON(http.StatusOK, gin.H{"data": userHash, "code": http.StatusOK})
	//插入数据到数据库中
	return
}
