package rBook

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/userData"
)

// UserProfileWithReq 定义联合查询结构体
type UserProfileWithReq struct {
	UserProfile  userData.UserProfile `xorm:"extends"` // 继承用户资料表字段
	UserReqEmail string               `xorm:"email"`   // 来自user_req表的邮箱
	// 可以添加其他需要联合查询的字段
}

func ReqProfiles(friendIDs []int, c *gin.Context) ([]UserProfileWithReq, error) {
	var results []UserProfileWithReq
	err := app.Engine.Table("user_profiles").
		Select("user_profiles.*, user_req.email"). // 明确选择字段
		Join("INNER", "user_req",
							"user_profiles.user_id = user_req.ID"). // 关联条件
		In("user_profiles.user_id", friendIDs). // 过滤条件
		Find(&results)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"data": err, "code": 502})
		return nil, fmt.Errorf("联合查询失败: %w", err)
	}

	if len(results) == 0 {
		c.JSON(http.StatusBadGateway, gin.H{"data": 0, "code": 502})
		return nil, errors.New("未找到匹配记录")
	}
	return results, nil
}
