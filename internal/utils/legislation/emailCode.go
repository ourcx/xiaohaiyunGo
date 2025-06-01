package legislations

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"xiaohaiyun/internal/app"
	cosFile "xiaohaiyun/internal/utils/cos"
)

//这里都是一些验证码处理函数

// CodeStoreFromEmail 储存你的验证码
func CodeStoreFromEmail(code string, email string) error {
	// 存储结构
	data := map[string]interface{}{
		"code":       code,
		"attempts":   3, // 初始尝试次数
		"expires_at": time.Now().Add(10 * time.Minute).UnixNano(),
	}

	// 使用管道操作保证原子性
	pipe := app.Rdb.Pipeline()
	pipe.HSet(app.Ctx, "vcode:"+email, data)
	pipe.Expire(app.Ctx, "vcode:"+email, 10*time.Minute)
	_, err := pipe.Exec(app.Ctx)

	return err
}

// StartCleanupScheduler 暂时废弃的清理的策略，等以后写应该全表删除吧，现在对go-redis和原子性不熟悉
func StartCleanupScheduler(c *gin.Context) {
	go func() {
		userTable := cosFile.GetID(c)
		ticker := time.NewTicker(10 * time.Minute)
		key := "vcode:" + userTable.Email
		defer ticker.Stop()

		err := app.Rdb.Watch(app.Ctx, func(tx *redis.Tx) error {
			data, _ := tx.HGetAll(app.Ctx, key).Result()
			// 检查过期
			expiresAt, _ := strconv.ParseInt(data["expires_at"], 10, 64)
			//转换函数
			if time.Now().UnixNano() > expiresAt {
				tx.Del(app.Ctx, key) // 自动清理过期验证码
				return fmt.Errorf("验证码已过期")
			}
			return nil
		}, key)

		if err != nil {
			fmt.Println("redis查找失败")
			return
		}
	}()
}

//验证重置密码的邮箱1
