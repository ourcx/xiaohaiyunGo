package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	cosFile "xiaohaiyun/internal/utils/cos"
	redis1 "xiaohaiyun/internal/utils/redis-1"
)

type COSObjectHeaders struct {
	ContentType   string    `json:"Content-Type" header:"Content-Type"`
	ContentLength string    `json:"Content-Length" header:"Content-Length"` // 通常长度会用数字类型
	ETag          string    `json:"ETag" header:"ETag"`
	LastModified  time.Time `json:"Last-Modified" header:"Last-Modified"` // 时间类型
	XCosRequestID string    `json:"X-Cos-Request-Id" header:"X-Cos-Request-Id"`
}

type ImgCOSObjectHeaders struct {
	ContentLength string `json:"Content-Length" header:"Content-Length"` // 通常长度会用数字类型
	LastModified  string `json:"Last-Modified" header:"Last-Modified"`   // 时间类型
	Name          string `json:"name"`
	BasePath      string `json:"BasePath"`
}

type COSKey struct {
	Name string `json:"name"`
}

// BaseData 这个方法是列出对象的元信息
// 根据前端给我的数据来返回对应的元信息
func BaseData(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	var base COSObjectHeaders
	var CKey COSKey

	err := c.ShouldBind(&CKey)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	client := cosFile.Client()
	key := "users/" + strconv.Itoa(userID) + "/" + CKey.Name

	clientRedis := redis1.Redis1()
	ctx := context.Background()

	exists, err := clientRedis.Exists(ctx, key).Result()
	if err != nil {
		panic(fmt.Sprintf("检查Hash存在性失败: %v", err))
	}

	if exists > 0 {
		result, err := clientRedis.HGet(ctx, key, key).Result()
		if err != nil {
			return
		}
		var dataJSON COSObjectHeaders
		err = json.Unmarshal([]byte(result), &dataJSON)
		if err != nil {
			return
		}
		c.JSON(200, gin.H{
			"code":  200,
			"msg":   "ok",
			"data":  dataJSON,
			"title": "这个是缓存的数据",
		})
		return
	}

	resp, err := client.Object.Head(context.Background(), key, nil)
	if err != nil {
		c.JSON(403, gin.H{
			"code":    403,
			"message": err.Error(),
			"data":    "文件名称输入错误，文件并不存在，请检查",
		})
		return
	}
	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.Header.Get("Content-Length")
	etag := resp.Header.Get("ETag")
	reqid := resp.Header.Get("X-Cos-Request-Id")

	//拿到相关的响应头，然后进行相关的处理，完成修复

	base = COSObjectHeaders{
		ContentType:   contentType,
		ContentLength: contentLength,
		ETag:          etag,
		LastModified:  time.Now(),
		XCosRequestID: reqid,
	}
	jsonData, _ := json.Marshal(base) // 注意 := 和逗号
	err = clientRedis.HSet(ctx, key, key, jsonData).Err()
	if err != nil {
		panic(fmt.Sprintf("写入Hash失败: %v", err))
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": base,
	})
}
