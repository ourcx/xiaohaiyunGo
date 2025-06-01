package cosFileUtils

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"
)

// LocationUrl 本地获取预签名URL//传入用户id和要查询的文件，返回对应的预览url
func LocationUrl(userID int, name string) string {
	client := Client()
	// 生成预览 URL
	objectKey := "users/" + strconv.Itoa(userID) + "/" + name
	presignedURL, err := client.Object.GetPresignedURL(
		context.Background(),
		http.MethodGet,
		objectKey,
		os.Getenv("COS_MAIN_SECRET_ID"),
		os.Getenv("COS_MAIN_SECRET_KEY"),
		//从本地获取密钥和id
		10*time.Minute,
		nil,
	)
	if err != nil {
		panic(err)
	}

	// 拼接预览参数
	previewURL := presignedURL.String() +
		"?ci-process=doc-preview" +
		"&dstType=html" +
		"&htmlwaterword=Q09T5paH5qGj6aKE6KeI"

	return previewURL
}
