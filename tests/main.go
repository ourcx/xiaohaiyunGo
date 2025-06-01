package main

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/qcloud-cos-sts-sdk/go"
	"net/http"
	"net/url"
	"os"
	"time"
	"xiaohaiyun/configs"
)

// GenerateSecureUploadURL 生成临时密钥和预签名URL
func GenerateSecureUploadURL(userID string) (string, error) {
	// ===================== 1. 生成临时密钥 (STS) =====================
	// 从环境变量读取主账号密钥（用于申请临时密钥）
	mainSecretID := os.Getenv(config.Conf.Cos.SecretID)   // 主账号 SecretId
	mainSecretKey := os.Getenv(config.Conf.Cos.SecretKey) // 主账号 SecretKey

	stsClient := sts.NewClient(mainSecretID, mainSecretKey, nil)

	// 动态生成用户隔离路径
	cosPath := fmt.Sprintf("users/%s/${filename}", userID) // 保留文件名灵活性

	// STS 策略：仅允许上传到用户专属目录
	policy := &sts.CredentialPolicy{
		Statement: []sts.CredentialPolicyStatement{
			{
				Effect: "allow",
				Action: []string{
					"cos:PutObject", // 只允许上传操作
				},
				Resource: []string{
					fmt.Sprintf("qcs::cos:ap-guangzhou:uid/1250000000:xiaohaiyun-1331891188/%s", cosPath),
				},
			},
		},
	}

	// 申请临时密钥（有效期15分钟）
	stsRes, err := stsClient.GetCredential(&sts.CredentialOptions{
		Policy:          policy,
		DurationSeconds: 900, // 15分钟
	})
	if err != nil {
		return "", fmt.Errorf("STS请求失败: %v", err)
	}

	// ===================== 2. 生成预签名URL =====================
	// 初始化COS客户端（使用临时密钥）
	u, _ := url.Parse("https://xiaohaiyun-1331891188.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  stsRes.Credentials.TmpSecretID,
			SecretKey: stsRes.Credentials.TmpSecretKey,
			Transport: &http.Transport{},
		},
	})

	// 生成预签名URL（允许客户端上传到动态路径）
	presignedURL, err := client.Object.GetPresignedURL(
		context.Background(),
		http.MethodPut,
		cosPath, // 示例：users/123/${filename}
		stsRes.Credentials.TmpSecretID,
		stsRes.Credentials.TmpSecretKey,
		15*time.Minute, // URL有效期 <= 临时密钥有效期
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %v", err)
	}

	return presignedURL.String(), nil
}

func main() {
	// 示例：用户ID为123的用户上传文件
	presignedURL, err := GenerateSecureUploadURL("123")
	if err != nil {
		panic(err)
	}
	fmt.Println("安全的上传URL:", presignedURL)
}
