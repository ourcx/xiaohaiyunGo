package File

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	cosFile "xiaohaiyun/internal/utils/cos"

	"github.com/gin-gonic/gin"
)

type SelfListener struct {
}

func (s SelfListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	//TODO implement me
	panic("implement me")
}

// AddFile 添加文件
func AddFile(c *gin.Context) {
	// 接收文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"msg":   "无法获取文件",
		})
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			// 记录日志，但不影响主要流程
			fmt.Printf("关闭文件流失败: %v\n", err)
		}
	}(file)
	if header.Size > 50*1024*1024 { // 50MB限制
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"msg":   "文件过大，请上传小于50MB的文件",
		})
		return
	}
	ext := filepath.Ext(header.Filename)
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
		".txt":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
	}
	if !allowedTypes[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"msg":   "不支持的文件类型",
		})
		return
	}

	// COS云存储
	client := cosFile.Client()

	// 生成唯一文件名
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	key := "./back/" + filename

	// 上传到COS
	_, err = client.Object.Put(context.Background(), key, file, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"msg":   "文件上传到云存储失败",
		})
		return
	}
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
			// 设置默认的进度回调函数
			Listener: &cos.DefaultProgressListener{},
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			// 如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
			XCosACL: "private",
		},
	}
	opt.Listener = &SelfListener{}
	filepath := "./back"
	_, err = client.Object.PutFromFile(context.Background(), filename, filepath, opt)
	if err != nil {
		panic(err)
	}

	fmt.Println(SelfListener{})
	//这样后台可以看到上传的进度

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"msg":      "文件上传成功",
		"filename": filename,
		"key":      key,
		"size":     header.Size,
		"ext":      ext,
	})
}

//先上传到服务器进行校验
