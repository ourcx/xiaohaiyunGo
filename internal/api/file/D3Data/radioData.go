package D3

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
	"xiaohaiyun/internal/api/file"
	Change "xiaohaiyun/internal/utils/D3"
	cosFile "xiaohaiyun/internal/utils/cos"
)

// Total TreeFile 这是一个可以拿到所有文件列表的东西
func Total(c *gin.Context) {
	file.List = []string{}
	userID := cosFile.GetID(c).ID
	var folderName cosFile.Folder
	var fileList []string
	fileList = file.TreeFileTree(userID, folderName, fileList)
	file.DeduplicateUnordered(file.List)

	resultCh := make(chan string, len(file.List)) // 结果通道
	var wg sync.WaitGroup
	var data int
	for _, i := range file.List {
		wg.Add(1)
		go func() {
			defer wg.Done()
			base, err := file.ImgBase(i)
			if err != nil {
				return
			}
			resultCh <- base.ContentLength
		}()
	}
	go func() {
		wg.Wait()
		close(resultCh) // 关闭通道以便后续遍历
	}()
	for res := range resultCh {
		resInt, _ := strconv.Atoi(res)
		data = data + resInt
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": Change.FormatBytes(uint64(data)),
	})

}
