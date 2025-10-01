package File

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func AllFile(c *gin.Context) {
	client := cosFile.Client()
	opt := &cos.BucketGetOptions{
		Prefix:  "",   // prefix 表示要查询的文件夹 		// deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys: 1000, // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}

	v, _, err := client.Bucket.Get(context.Background(), opt)

	if err != nil {
		fmt.Println(err)
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  "发生了不知名的错误",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": v,
		"msg":  "success",
	})
}
