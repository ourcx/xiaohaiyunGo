package file

import (
	"github.com/gin-gonic/gin"
	"path"
	"strings"
	cosFile "xiaohaiyun/internal/utils/cos"
)

type special struct {
	Exts string `json:"exts"`
}

// SpecialTreeFile TreeFile 这是一个可以拿到所有文件列表的东西
func SpecialTreeFile(c *gin.Context) {
	var fileList []string
	userID := cosFile.GetID(c).ID
	var folderName cosFile.Folder
	fileList = TreeFileTree(userID, folderName, fileList)
	var file []string
	var name special
	err := c.ShouldBind(&name)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500, "message": err.Error(), "data": folderName,
		})
		return
	}
	//文件容器
	// 定义支持的格式（统一转为小写判断）
	docExts := map[string]bool{
		".doc":  true,
		".docx": true,
		".pdf":  true,
		".md":   true,
		".txt":  true,
		".ppt":  true,
		".pptx": true,
		".xls":  true,
		".xlsx": true,
	}

	videoExts := map[string]bool{
		".mp4": true,
		".mov": true,
		".avi": true,
		".mkv": true,
		".flv": true,
		".wmv": true,
	}

	mp3Exts := map[string]bool{
		".mp3":  true,
		".wav":  true,
		".flac": true,
		".aac":  true,
		".ogg":  true,
	}

	if name.Exts == "doc" {
		file = ClassifyFiles(DeduplicateUnordered(fileList), docExts)
	} else if name.Exts == "video" {
		file = ClassifyFiles(DeduplicateUnordered(fileList), videoExts)
	} else if name.Exts == "mp3" {
		file = ClassifyFiles(DeduplicateUnordered(fileList), mp3Exts)
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": file,
	})
	fileList = []string{}
}

// ClassifyFiles 文件分类函数
func ClassifyFiles(files []string, Exts map[string]bool) (Files []string) {

	for _, file := range files {
		// 排除目录（以 / 结尾的路径）
		if strings.HasSuffix(file, "/") {
			continue
		}
		// 获取文件扩展名并转为小写
		ext := strings.ToLower(path.Ext(file))
		// 分类判断
		switch {
		case Exts[ext]:
			Files = append(Files, file)
		}
	}
	return
}
