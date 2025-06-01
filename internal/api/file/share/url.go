package share

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
	"xiaohaiyun/internal/models/userData"
	"xiaohaiyun/internal/repositories"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func SetShare(c *gin.Context) {
	userTable := cosFile.GetID(c)

	id := uuid.New()                  // 生成UUIDv4
	binaryID, _ := id.MarshalBinary() // 转换为16字节

	var shareData share.UrlShare
	var bing share.UrlShareJSON
	err := c.ShouldBind(&bing)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err,
			"code": 502,
		})
		return
	}

	s := repositories.NewUserRepository(app.Engine)
	table, err := s.GetTableByEmail(userTable.Email)
	user, _ := table.(*userData.UserProfile)

	shareData = share.UrlShare{
		UserReqID: userTable.ID,
		ID:        binaryID,
		Url:       bing.Url,
		Username:  userTable.Name,
		Email:     userTable.Email,
		Signature: user.Signature,
		Avatar:    user.AvatarUrl,
	}

	update, err := app.Engine.Table("url_share").Update(shareData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if update == 0 {
		insert, err := app.Engine.Table("url_share").Insert(shareData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println(insert)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": shareData, "code": 200, "update": update})
}

func GetShare(c *gin.Context) {
	userTable := cosFile.GetID(c)
	var shareData share.UrlShare

	get, err := app.Engine.Table("url_share").Where("email=?", userTable.Email).Get(shareData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": get,
		"code": 200,
	})
}
