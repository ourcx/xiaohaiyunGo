package app

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"xiaohaiyun/configs"
)

var Engine *xorm.Engine
var (
	Rdb *redis.Client
	// Ctx RdbFirst *redis.Client
	Ctx = context.Background()
)

// InitializeMySQL 数据库初始化
func InitializeMySQL() error {
	var err error
	// 创建数据库引擎
	Engine, err = xorm.NewEngine(config.Conf.Database.Driver, config.Conf.Database.Source)
	if err != nil {
		log.Error("数据库初始化失败: %v", err)
		return err
	}

	// 测试数据库连接
	if err = Engine.Ping(); err != nil {
		log.Error("数据库连接失败: %v", err)
		return err
	}

	//repositories.NewUserRepository(Engine)

	return nil
}

// InitRedis 链接到了redis数据库的初始化
func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})

	// 测试连接
	if err := Rdb.Ping(Ctx).Err(); err != nil {
		//redis的链接测试，出现多次突然断联导致程序终止
		//panic(fmt.Sprintf("Redis连接失败: %v", err))
		fmt.Print(fmt.Sprintf("Redis连接失败: %v", err))
		return err
	}
	return nil
}

//func InitRedisFirst() error {
//	RdbFirst = redis.NewClient(&redis.Options{
//		Addr:     config.Conf.Redis.Host,
//		Password: config.Conf.Redis.Password,
//		DB:       config.Conf.Redis.Db1,
//	})
//
//	// 测试连接
//	if err := RdbFirst.Ping(Ctx).Err(); err != nil {
//		panic(fmt.Sprintf("Redis连接失败: %v", err))
//		return err
//	}
//	return nil
//}
