package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"time"
	"xiaohaiyun/internal/app"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type input struct {
	code string
}

// ValidateCode 校验验证码
func ValidateCode(c *gin.Context) {
	userTable := cosFile.GetID(c)
	var inputCode input

	err := c.ShouldBind(&inputCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err, "code": 400})
		return
	}

	key := "vcode:" + userTable.Email

	// 使用 Watch 实现原子操作
	err = app.Rdb.Watch(app.Ctx, func(tx *redis.Tx) error {
		// 获取存储数据
		data, err := tx.HGetAll(app.Ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "你未申请验证码或者验证码过期",
			})
			return err
		}

		// 检查是否存在
		if len(data) == 0 {
			return fmt.Errorf("请先获取验证码")
		}

		// 检查过期
		expiresAt, _ := strconv.ParseInt(data["expires_at"], 10, 64)
		if time.Now().UnixNano() > expiresAt {
			tx.Del(app.Ctx, key) // 自动清理过期验证码
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码过期",
			})
			return fmt.Errorf("验证码已过期")
		}

		// 检查尝试次数
		attempts, _ := strconv.Atoi(data["attempts"])
		if attempts <= 0 {
			return fmt.Errorf("尝试次数超限")
		}

		// 验证码比对
		if data["code"] != inputCode.code {
			// 原子减少尝试次数
			_, err = tx.HIncrBy(app.Ctx, key, "attempts", -1).Result()
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码输入错误",
			})
			return err
		}

		// 验证成功删除 Key
		_, err = tx.Del(app.Ctx, key).Result()
		return err
	}, key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "错误的操作",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码正确",
	})
	return
}
