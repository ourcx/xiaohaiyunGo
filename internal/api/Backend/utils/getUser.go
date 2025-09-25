package utilsBack

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
	"xiaohaiyun/internal/api/Backend/models"
	"xiaohaiyun/internal/app"
)

// GetUser 根据email得到用户
func GetUser(email string) (models.BackendLogin, error, bool) {
	var getVal models.BackendLogin
	exists, err := app.Engine.Table("BackendLogin").Where("Email=?", email).Get(&getVal)

	if err != nil {
		return getVal, err, exists
	} else if exists {
		return getVal, err, exists
	}

	return getVal, nil, exists
}

// InsetUser 用户注册的插入函数
func InsetUser(email string, password string, userId int64) error {
	var Val models.BackendLogin
	Val.UserID = userId
	Val.Password = password
	Val.Email = email
	Val.Username = "未知用户"
	Val.CreatedAt = time.Now()
	Val.UpdatedAt = time.Now()
	insert, err := app.Engine.Table("BackendLogin").Insert(Val)
	if err != nil {
		return err
	}
	fmt.Print(insert, "插入用户表")
	return nil
}

// GenerateSnowflakeID 雪花算法
func GenerateSnowflakeID() int64 {
	node, err := snowflake.NewNode(1) // 节点ID，可以根据机器配置
	if err != nil {
		// 备用方案：使用时间戳
		return time.Now().UnixNano() / 1e6
	}
	return node.Generate().Int64()
}
