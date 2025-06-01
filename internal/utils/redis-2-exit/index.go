package redis2exit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	config "xiaohaiyun/configs"
)

func Redis2() *redis.Client {
	// 配置Redis客户端参数，指定DB为1
	client := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host,     // Redis地址
		Password: config.Conf.Redis.Password, // 密码
		DB:       2,                          // 使用数据库1（db1）
		PoolSize: 10,                         // 连接池大小
	})

	// 测试连接
	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("连接Redis失败: %v", err))
	}
	fmt.Println("连接成功:", pong)
	return client
}
