package cosFileUtils

import (
	"context"
	"errors"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
)

// CheckFolderExists 检查空文件夹是否存在
func CheckFolderExists(client *cos.Client, folderKey string) (bool, error) {
	// COS 中文件夹存在 = 存在以 '/' 结尾的空对象
	resp, err := client.Object.Head(context.Background(), folderKey, nil)
	if err != nil {
		var cosErr *cos.ErrorResponse
		if errors.As(err, &cosErr) && cosErr.Code == "404" {
			return false, nil // 不存在
		}
		return false, err // 其他错误（如权限问题）
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// 检查对象内容长度是否为0（确认是空对象）
	if resp.ContentLength == 0 {
		return true, nil
	}
	return false, nil
}
