package share

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"regexp"
	"strings"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func Extract(c *gin.Context) {
	var Ext share.ExtractShareFiles
	userID := cosFile.GetID(c).ID
	client := cosFile.Client()

	err := c.ShouldBind(&Ext)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	var getUrlData share.UrlData
	oneUUID, err := uuid.Parse(Ext.OneId)
	_, err = app.Engine.Table("url_data").Where("share_id=?", oneUUID[:]).Get(&getUrlData)
	if err != nil {
		return
	}

	for _, key := range getUrlData.Files {
		// 改进正则：匹配完整用户路径（支持多位数ID）
		re := regexp.MustCompile(`^users?/\d+/`)
		cleanedKey := re.ReplaceAllString(key, "")

		// 安全构造新路径（处理可能的开头斜杠）
		newPath := strings.TrimPrefix(cleanedKey, "/")
		newKey := fmt.Sprintf("users/%d/%s", userID, newPath)

		// 统一使用正斜杠构造旧路径
		oldKey := fmt.Sprintf("share/%s/%s", Ext.OneId, newPath)
		oldKey = strings.ReplaceAll(oldKey, "\\", "/") // 确保统一分隔符

		// 构造源路径（根据SDK要求调整）
		source := fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.Host, oldKey)

		// 执行复制操作
		_, _, err = client.Object.Copy(context.Background(), newKey, source, nil)
		if err != nil {
			log.Printf("文件复制失败: %s -> %s, 错误: %v", source, newKey, err)
			c.JSON(500, gin.H{
				"code":    500,
				"message": err.Error(),
				"data":    fmt.Sprintf("在%s/%s这个文件发生系统性记录错误,请清理文件", Ext.OneId, newPath),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"one_id": oneUUID,
		},
	})

}
