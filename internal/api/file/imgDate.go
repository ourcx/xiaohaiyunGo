package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	cosFile "xiaohaiyun/internal/utils/cos"
	redis1 "xiaohaiyun/internal/utils/redis-1"
)

func ImgDate(c *gin.Context) {
	var fileList []string
	imgExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
		".bmp":  true,
	}
	userID := cosFile.GetID(c).ID
	var folderName cosFile.Folder
	fileList = TreeFileTree(userID, folderName, fileList)
	var date []string
	file := ClassifyFiles(DeduplicateUnordered(fileList), imgExts)
	for _, i := range file {
		base, err := imgBaseDate(i)
		if err != nil {
			return
		}
		date = append(date, base.LastModified)
	}
	c.JSON(200, gin.H{
		"code": 200, "data": date,
	})
}

// 提取图片数据的元数据
func imgBaseDate(name string) (ImgCOSObjectHeaders, error) {
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
			LastModified: dataJSON.LastModified.Format("2006-01-02"),
		}
		return imgData, nil
	}

	//拿到相关的响应头，然后进行相关的处理，完成修复

	base := ImgCOSObjectHeaders{
		LastModified: time.Now().Format("2006-01-02"),
	}
	jsonData, _ := json.Marshal(base) // 注意 := 和逗号
	err = clientRedis.HSet(ctx, name, name, jsonData).Err()
	return base, nil
}
