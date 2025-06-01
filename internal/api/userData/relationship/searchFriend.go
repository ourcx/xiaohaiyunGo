package relationship

import (
	"github.com/gin-gonic/gin"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/RelationShip"
	"xiaohaiyun/internal/repositories"
)

// 详细的个人信息查找
func SearchFriend(c *gin.Context) {
	var email RelationShip.FriendEmail

	c.ShouldBind(&email)
	s := repositories.NewUserRepository(app.Engine)
	res, err := s.GetTableByEmail(email.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": res})
}
