package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/services"
	"xiaohaiyun/internal/utils"
)

func userLongin(c *gin.Context) {
	var user models.UserloginPost
	var log models.UserLogin
	clientIP := c.ClientIP()
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userHash := models.SafeUserReq{
		Email:    user.Email,
		Password: user.Password,
		//加密后的密码
	}

	log = models.UserLogin{
		Email:     user.Email,
		LoginTime: time.Now(),
		LoginIp:   clientIP,
	}

	JWT, err2 := utils.GenerateJWTHS256(userHash.Email)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error(), "msg": "JWT生成错误,请注册后登录"})
	}

	// 3. 解析验证Token
	claims, _ := utils.ParseUserJWT(JWT)

	// 4. 提取结构化用户数据
	userClaims, _ := utils.ExtractUserClaims(claims)
	fmt.Println(userClaims)

	//先验证是否重复再进行注册信息的提交
	bool, err := services.CheckPassword(userHash.Email, userHash.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "data": "系统错误"})
	}
	if bool {
		fmt.Println(log)
		_, err = app.Engine.Table("user_logins").Insert(&log)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": userHash.Email + "登录成功了", "code": http.StatusOK, "jwt": JWT})
	} else {
		c.JSON(http.StatusConflict, gin.H{"data": "请注册后再进行登录", "code": http.StatusConflict})
	}
	//返回的时候带上code是一个好习惯
	//插入数据到数据库中
	return
}
