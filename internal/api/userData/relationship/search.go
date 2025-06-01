package relationship

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/repositories"
)

type user struct {
	Email string `json:"email"`
}

type formatTable struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// SearchUser 根据前端传过来的email返回这个人的信息,查找陌生人的方法
func SearchUser(c *gin.Context) {
	var userType user
	var formatTables formatTable

	err := c.ShouldBind(&userType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "错误的接收"})
		return
	}

	s := repositories.NewUserRepository(app.Engine)
	table, err := s.GetUserByEmail(userType.Email) //等于空就说明数据库里面有这个邮箱了
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误"})
		return
	}

	formatTables = formatTable{
		Email: table.Email,
		Name:  table.Name,
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": formatTables})
}
