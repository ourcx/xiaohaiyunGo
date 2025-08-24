package collect

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func CollectFile(c *gin.Context) {
	client := cosFile.Client()
	opt := &cos.BucketGetOptions{
		Prefix:    "users/",
		Delimiter: "/",  // prefix 表示要查询的文件夹 		// deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys:   1000, // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}
	v, _, err := client.Bucket.Get(context.Background(), opt)
	if err != nil {
	}

	fmt.Println(v)
	//现在这里是一个采集系统，传入文件，然后我把文件发送到cos里面去

}

func IsExist(key string) bool {
	client := cosFile.Client()
	ok, err := client.Object.IsExist(context.Background(), key)
	if err == nil && ok {
		return true
	} else if err != nil {
		return false
	} else {
		return false
	}
}
