package cosFileUtils

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"strconv"
)

type Folder struct {
	Name string `json:"name"`
}

// DeckFolderExists 判断是否有重复的子文件
func DeckFolderExists(Key string, folder Folder, userID int) (bool, error) {
	list := FList(userID, folder)
	for _, i := range list {
		if Key == i {
			return false, nil
		}
	}
	return true, nil
}

// FList 这是一个去重的函数，在于在访问某个文件夹的时候，拿到目录里面的内容，和你进行比对
func FList(userID int, folderName Folder) []string {
	var fileList []string
	var folderList []string

	client := Client()
	//初始化客户端

	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    "users/" + strconv.Itoa(userID) + "/" + folderName.Name, // prefix 表示要查询的文件夹
		Delimiter: "",                                                      // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys:   1000,                                                    // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := client.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		// 在循环内部收集文件 Key
		for _, content := range v.Contents {
			fileList = append(fileList, content.Key) // 将文件名存入切片
		}

		// common prefix 表示表示被 delimiter 截断的路径, 如 delimter 设置为/, common prefix 则表示所有子目录的路径
		for _, commonPrefix := range v.CommonPrefixes {
			folderList = append(folderList, commonPrefix)
		}
		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}

	fileList = append(folderList, fileList...)
	return fileList
}
