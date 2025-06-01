package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models"
	"xiaohaiyun/internal/repositories"
)

// GenerateJWTHS256 生成JWT
func GenerateJWTHS256(email string) (string, error) {
	key := []byte(models.JwtKey)
	tokenDuration := 76 * time.Hour
	now := time.Now()
	//从数据库拿到name和id
	s := repositories.NewUserRepository(app.Engine)
	user, _ := s.GetUserByEmail(email) //等于空就说明数据库里面有这个邮箱了

	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]string{
			"email": email,
			"name":  user.Name,
			"ID":    strconv.Itoa(user.ID),
		},
		"iat": now.Unix(),
		"exp": now.Add(tokenDuration).Unix(),
		"iss": "xiaohai",
	})

	return t.SignedString(key)
}

//JWT是token的一种实现方式，token包含jwt

//更新jwt
//更新就是把他的数据再生成一下
