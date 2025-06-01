package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/repositories"
)

// JWTConfig JWT 配置项（建议通过环境变量注入）
type JWTConfig struct {
	SecretKey  string        // 密钥（建议使用环境变量）
	Expiration time.Duration // 有效期（如 24*time.Hour）
	Issuer     string        // 签发者标识（如 "xiaohaiyun-auth"）
}

// GenerateJWTHS256V2 增强版JWT生成函数
// 返回值：(token, error)
func GenerateJWTHS256V2(email string, config JWTConfig) (string, error) {
	// ----------------------------
	// 1. 参数校验
	// ----------------------------
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	if config.SecretKey == "" {
		return "", fmt.Errorf("JWT secret key not configured")
	}

	// ----------------------------
	// 2. 获取用户信息
	// ----------------------------
	repo := repositories.NewUserRepository(app.Engine)
	user, err := repo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("database query failed: %w", err)
	}

	if user == nil {
		return "", fmt.Errorf("user with email %s not found", email)
	}

	// ----------------------------
	// 3. 构建标准化声明
	// ----------------------------
	now := time.Now()
	standardClaims := jwt.RegisteredClaims{
		Issuer:    config.Issuer,                                  // 签发者
		Subject:   strconv.Itoa(user.ID),                          // 用户ID作为主题
		IssuedAt:  jwt.NewNumericDate(now),                        // 签发时间
		ExpiresAt: jwt.NewNumericDate(now.Add(config.Expiration)), // 过期时间
	}

	// ----------------------------
	// 4. 构建自定义声明
	// ----------------------------
	customClaims := jwt.MapClaims{
		// 保留原始结构
		"user": map[string]string{
			"email": user.Email,
			"name":  user.Name,
			"ID":    strconv.Itoa(user.ID),
		},
		// 合并标准声明
		"iss": standardClaims.Issuer,
		"sub": standardClaims.Subject,
		"iat": standardClaims.IssuedAt.Unix(),
		"exp": standardClaims.ExpiresAt.Unix(),
	}

	// ----------------------------
	// 5. 生成令牌
	// ----------------------------
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// 使用环境变量密钥
	signingKey := []byte(os.Getenv(config.SecretKey))

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("token signing failed: %w", err)
	}

	return signedToken, nil
}
