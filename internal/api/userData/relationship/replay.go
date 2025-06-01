package relationship

import (
	"github.com/gin-gonic/gin"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/repositories"
	"xiaohaiyun/internal/utils"
)

func ReplayToUser(c *gin.Context) {
	//用户返回jwt，我根据jwt解析出来的email查询数据库，返回符合的所有列表，，一个get方法
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
	s := repositories.NewUserRepository(app.Engine)
	masterTable, err := s.GetUserRelationSHipByEmail(masterEmail)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"masterEmail": masterEmail,
			"masterTable": masterTable,
		},
	})

}
