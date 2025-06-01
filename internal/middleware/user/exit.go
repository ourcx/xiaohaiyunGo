package userAuth

import (
	"context"
	"github.com/gin-gonic/gin"
	"xiaohaiyun/internal/models/exit"
	cosFile "xiaohaiyun/internal/utils/cos"
	redis2exit "xiaohaiyun/internal/utils/redis-2-exit"
)

func SetExitJwt(c *gin.Context) {

	user := cosFile.GetID(c)
	token := c.GetHeader("Authorization")
	clientRedis := redis2exit.Redis2()
	ctx := context.Background()

	var invalidated exit.Invalidated
	invalidated = exit.Invalidated{
		Jwt:   token,
		Email: user.Email,
	}

	result, err := clientRedis.Exists(ctx, user.Email).Result()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "redis数据库出错",
		})
		return
	}

	if result > 0 {
		return
	} else {
		clientRedis.MSet(ctx, "jwt_blacklist:"+invalidated.Jwt, "invalidated", 0)
		//如果后端redis没有黑名单就设置黑名单
		c.JSON(200, gin.H{
			"code": 200, "message": "success,已经退出登录，上次登录的jwt计入黑名单",
		})
	}
}

func GetExitJwt(c *gin.Context) {
	token := c.GetHeader("Authorization")
	clientRedis := redis2exit.Redis2()
	ctx := context.Background()

	result, err := clientRedis.Exists(ctx, "jwt_blacklist:"+token).Result()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "redis数据库出错",
		})
		return
	}

	if result > 0 {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "失效的jwt",
		})
		return
	}

	c.Next()
}
