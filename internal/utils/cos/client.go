package cosFileUtils

import (
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
)

func Client() *cos.Client {
	//获取用户的JWT信息
	secretID := os.Getenv("COS_MAIN_SECRET_ID")
	secretKey := os.Getenv("COS_MAIN_SECRET_KEY")

	// 检查环境变量是否为空
	if secretID == "" || secretKey == "" {
		fmt.Println("错误：请先设置环境变量 COS_MAIN_SECRET_ID, COS_MAIN_SECRET_KEY, COS_BUCKET_NAME, COS_BUCKET_REGION")
		os.Exit(1)
	}
	//检查环境变量

	u, _ := url.Parse("https://xiaohaiyun-1331891188.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，
			SecretID: os.Getenv("COS_MAIN_SECRET_ID"), // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，
			SecretKey: os.Getenv("COS_MAIN_SECRET_KEY"),
			// 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见
		},
	})

	return client
}
