package D3

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
	"xiaohaiyun/internal/api/file"
	cosFile "xiaohaiyun/internal/utils/cos"
)

var (
	docExts = map[string]bool{
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
	videoExts = map[string]bool{
		".mp4": true,
		".mov": true,
		".avi": true,
		".mkv": true,
		".flv": true,
		".wmv": true,
	}
	mp3Exts = map[string]bool{
		".mp3":  true,
		".wav":  true,
		".flac": true,
		".aac":  true,
		".ogg":  true,
	}
	imgExts = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
		".bmp":  true,
	}
)

func Proportion(c *gin.Context) {
	userID := cosFile.GetID(c).ID
	var fileList []string
	var folderName cosFile.Folder
	fileList = file.TreeFileTree(userID, folderName, fileList)
	file.DeduplicateUnordered(file.List)
	fmt.Print(file.List)
	doc, video, mp3, img := ProcessAllCategories(file.List)
	
	c.JSON(200, gin.H{
		"doc":   doc,
		"video": video,
		"mp3":   mp3,
		"img":   img,
		"code":  200,
	})
}

func SingleDOC(files []string) []string {
	fileDOC := file.ClassifyFiles(file.DeduplicateUnordered(files), docExts)
	return fileDOC
}
func SingleVideo(files []string) []string {
	fileVideo := file.ClassifyFiles(file.DeduplicateUnordered(files), videoExts)
	return fileVideo
}

func SingleMp3(files []string) []string {
	fileMp3 := file.ClassifyFiles(file.DeduplicateUnordered(files), mp3Exts)
	return fileMp3
}

func SingleImg(files []string) []string {
	fileImg := file.ClassifyFiles(file.DeduplicateUnordered(files), imgExts)
	return fileImg
}

// 获取指定文件列表的总大小
func getTotalSize(files []string) int {
	total := 0
	resultCh := make(chan string, len(files))
	var wg sync.WaitGroup

	// 并发处理每个文件
	for _, f := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			base, err := file.ImgBase(filePath)
			if err != nil {
				return
			}
			resultCh <- base.ContentLength
		}(f)
	}

	// 等待所有任务完成后关闭通道
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// 累加所有结果
	for res := range resultCh {
		if resInt, err := strconv.Atoi(res); err == nil {
			total += resInt
		}
	}

	return total
}

// ProcessAllCategories 主处理函数：并发处理四类文件并返回各自总大小
func ProcessAllCategories(fileList []string) (docSize, videoSize, mp3Size, imgSize int) {
	var wg sync.WaitGroup
	type task struct {
		classify func([]string) []string
		result   *int
	}

	// 定义四个分类任务
	tasks := []task{
		{SingleDOC, &docSize},
		{SingleVideo, &videoSize},
		{SingleMp3, &mp3Size},
		{SingleImg, &imgSize},
	}

	// 并发执行每个分类的处理
	for _, t := range tasks {
		wg.Add(1)
		go func(t task) {
			defer wg.Done()
			classifiedFiles := t.classify(fileList)
			*t.result = getTotalSize(classifiedFiles)
		}(t)
	}

	wg.Wait()
	return
}
