package utils

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"xiaohaiyun/internal/app"
)

// ValidateCodeFromReq ValidateCode 校验验证码
func ValidateCodeFromReq(inputCode string, email string) error {

	key := "vcode:" + email

	// 使用 Watch 实现原子操作
	err := app.Rdb.Watch(app.Ctx, func(tx *redis.Tx) error {
		// 获取存储数据
		data, err := tx.HGetAll(app.Ctx, key).Result()
		fmt.Println(data)
		if err != nil {
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
			return fmt.Errorf("验证码已过期")
		}

		// 检查尝试次数
		attempts, _ := strconv.Atoi(data["attempts"])
		if attempts <= 0 {
			return fmt.Errorf("尝试次数超限")
		}

		// 验证码比对
		if data["code"] != inputCode {
			// 原子减少尝试次数
			_, err = tx.HIncrBy(app.Ctx, key, "attempts", -1).Result()
			return err
		}

		// 验证成功删除 Key
		_, err = tx.Del(app.Ctx, key).Result()
		return err
	}, key)

	if err != nil {
		return err
	}
	return nil
}
