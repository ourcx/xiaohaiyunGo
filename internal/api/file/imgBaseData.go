package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
	cosFile "xiaohaiyun/internal/utils/cos"
	redis1 "xiaohaiyun/internal/utils/redis-1"
)

type ImgRequest struct {
	Date string `json:"date"`
}

// ImgBaseData FilterByAge 根据传进来的日期查找hh
func ImgBaseData(c *gin.Context) {
	var fileList []string
	// 定义支持的格式（统一转为小写判断）
	imgExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
		".bmp":  true,
	}
	var date ImgRequest
	err := c.ShouldBindJSON(&date)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(),
		})
		return
	}
	userID := cosFile.GetID(c).ID
	var folderName cosFile.Folder
	fileList = TreeFileTree(userID, folderName, fileList)

	file := ClassifyFiles(DeduplicateUnordered(fileList), imgExts)
	resultCh := make(chan ImgCOSObjectHeaders, len(file)) // 结果通道
	var wg sync.WaitGroup
	var img []ImgCOSObjectHeaders
	for _, i := range file {
		wg.Add(1)
		go func() {
			defer wg.Done()
			base, err := ImgBase(i)
			if err != nil {
				return
			}
			if base.LastModified == date.Date {
				resultCh <- base
			}
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh) // 关闭通道以便后续遍历
	}()
	for res := range resultCh {
		img = append(img, res)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": img,
	})
	return
}

// 提取图片数据的元数据
func ImgBase(name string) (ImgCOSObjectHeaders, error) {
	client := cosFile.Client()

	clientRedis := redis1.Redis1()
	ctx := context.Background()

	exists, err := clientRedis.Exists(ctx, name).Result()
	if err != nil {
		panic(fmt.Sprintf("检查Hash存在性失败: %v", err))
	}

	if exists > 0 {
		result, _ := clientRedis.HGet(ctx, name, name).Result()
		var dataJSON COSObjectHeaders
		_ = json.Unmarshal([]byte(result), &dataJSON)
		var imgData = ImgCOSObjectHeaders{
			ContentLength: dataJSON.ContentLength,
			LastModified:  dataJSON.LastModified.Format("2006-01-02"),
			Name:          name,
		}
		return imgData, nil
	}

	resp, err := client.Object.Head(context.Background(), name, nil)
	contentLength := resp.Header.Get("Content-Length")

	//拿到相关的响应头，然后进行相关的处理，完成修复

	base := ImgCOSObjectHeaders{
		ContentLength: contentLength,
		LastModified:  time.Now().Format("2006-01-02"),
		Name:          name,
	}
	jsonData, _ := json.Marshal(base) // 注意 := 和逗号
	err = clientRedis.HSet(ctx, name, name, jsonData).Err()
	return base, nil
}
