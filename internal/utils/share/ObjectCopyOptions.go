package share

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// CopyOptions 把文件从用户去移动到共享桶,处理新的文件
func CopyOptions(objectKey string, destPrefix string) error {
	client := cosFile.Client()

	name := objectKey
	sourceURL := fmt.Sprintf("cdn.norubias.top/%s", name)
	//源路径
	dest := "share/" + strPath(objectKey, destPrefix)
	//目标路径
	_, _, err := client.Object.Copy(context.Background(), dest, sourceURL, nil)
	if err != nil {
		return err
	}
	return nil
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
