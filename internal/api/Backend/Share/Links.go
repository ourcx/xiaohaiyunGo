package Share

import (
	"context"
	"github.com/gin-gonic/gin"
	"strconv"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
	cosFile "xiaohaiyun/internal/utils/cos"
)

func Links(c *gin.Context) {
	var dataList []share.UrlData
	err := app.Engine.Table("url_data").Find(&dataList)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"data": nil,
			"msg":  err.Error(),
		})
		return
	}

	// 使用管道处理数据
	processedData := processDataWithPipeline(dataList)

	c.JSON(200, gin.H{
		"code": 200,
		"data": processedData,
		"msg":  "success",
	})
}

// 使用多个管道处理数据
// 使用管道模式处理数据
func processDataWithPipeline(dataList []share.UrlData) []any {
	// 创建管道
	inputChan := make(chan share.UrlData, len(dataList))
	outputChan := make(chan []any, len(dataList))

	// 步骤1: 发送数据到输入管道
	go func() {
		defer close(inputChan)
		for _, data := range dataList {
			inputChan <- data
		}
	}()

	// 步骤2: 处理数据的工作器
	workerCount := 3 // 可以根据需要调整工作器数量
	for i := 0; i < workerCount; i++ {
		go processWorker(inputChan, outputChan)
	}

	// 步骤3: 收集处理后的数据
	var result []any
	go func() {
		for processedData := range outputChan {
			result = append(result, processedData)
		}
	}()

	return result
}

// 处理数据的工作器函数
func processWorker(inputChan <-chan share.UrlData, outputChan chan<- []any) {
	client := cosFile.Client()
	for data := range inputChan {
		var response []any
		for shareItem := range data.ShareID {
			res, err := client.Object.Head(context.Background(), strconv.Itoa(shareItem), nil)
			if err != nil {
				panic(err)
			}
			response = append(response, res)
		}
		outputChan <- response
	}
}
