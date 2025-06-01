package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"xiaohaiyun/internal/models"
)

func ParseUserJWT(tokenString string) (jwt.MapClaims, error) {
	// 1. 获取签名密钥
	signingKey := []byte(models.JwtKey)

	// 2. 解析Token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signingKey, nil
	})

	// 3. 处理解析错误
	if err != nil {
		return nil, fmt.Errorf("token parse failed: %w", err)
	}

	// 4. 验证Token有效性
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// 5. 类型断言获取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	// 6. 验证签发者
	if iss, _ := claims["iss"].(string); iss != "xiaohai" {
		return nil, fmt.Errorf("invalid issuer")
	}

	return claims, nil
}

type UserClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
}

// ExtractUserClaims 从原始声明中提取结构化数据
func ExtractUserClaims(claims jwt.MapClaims) (*UserClaims, error) {
	// 1. 验证用户数据存在
	userData, ok := claims["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing user claims")
	}

	// 2. 提取邮箱
	email, ok := userData["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("invalid email claim")
	}

	// 3. 提取用户名
	name, ok := userData["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid name claim")
	}

	// 4. 转换用户ID
	idStr, ok := userData["ID"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid ID claim")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}

	return &UserClaims{
		Email: email,
		Name:  name,
		ID:    id,
	}, nil
}

//这两个函数根据密钥对JWT提取了相关的数据
