package share

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// CopyOptions 把文件从用户去移动到共享桶,处理新的文件
func CopyOptions(objectKey string, destPrefix string) error {
	srcBucket := "xiaohaiyun-1331891188" // 源存储桶
	region := "ap-guangzhou"             // 地域代码（广州）
	client := Client()
	destKey := strPath(objectKey, destPrefix)

	// 构建复制请求（正确参数）
	sourceURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s",
		srcBucket,
		region,
		url.PathEscape(objectKey), // 关键编码步骤
	)

	fmt.Println(sourceURL)
	// 执行复制操作（关键调用方式）
	_, resp, err := client.Object.Copy(
		context.Background(),
		destKey,   // 目标路径（保持同名）
		sourceURL, // 源文件完整URL（此处是核心变化）
		nil,       // 不需要额外选项
	)

	if err != nil {
		fmt.Printf("复制失败: %v\n", err)
		os.Exit(1)
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("成功复制 %s 到目标桶\n", objectKey)
	} else {
		fmt.Printf("复制异常，状态码: %d\n", resp.StatusCode)
		return err
	}
	return nil
}

// Client 生成目标桶客户端，共享桶的访问客户端
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

	u, _ := url.Parse("https://share-1331891188.cos.ap-guangzhou.myqcloud.com")
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

func strPath(objectKey, destPrefix string) string {
	parts := strings.SplitN(objectKey, "/", 3) // 分割为 ["users", "14", "xxxx/sss"]
	if len(parts) < 3 {
		fmt.Println("源路径格式不符合要求")
		os.Exit(1)
	}
	strippedKey := parts[2]                       // 获取 xxxx/sss
	destKey := path.Join(destPrefix, strippedKey) // 生成 new_prefix/xxxx/sss
	return destKey
}
