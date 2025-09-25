package Login

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/api/Backend/models"
	utilsBack "xiaohaiyun/internal/api/Backend/utils"
	"xiaohaiyun/internal/utils"
)

type ListFolder struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login GetUser godoc
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.BackendLogin
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func Login(c *gin.Context) {
	//这个就是一个简单的后台用户验证
	//解析token
	var user ListFolder
	var getVal models.BackendLogin
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	//	调用gorm库，请求比对
	getVal, err, exists := utilsBack.GetUser(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "系统错误，请稍后再试"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "没有此账户，请先注册"})
		return
	}

	enter := utils.CheckPassword(getVal.Password, user.Password)
	if !enter {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "密码错误"})
		return
	}

	// 3. 密码比对成功，生成 JWT Token
	JWT, err2 := utils.GenerateJWTHS256(user.Email)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error(), "msg": "JWT生成错误"})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"token": JWT,
		},
		"msg": "success",
	})

}
