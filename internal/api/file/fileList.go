package file

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"strconv"
	"time"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/locationFIleName"
	cosFile "xiaohaiyun/internal/utils/cos"
	sortW "xiaohaiyun/internal/utils/sort"
)

var (
	folder locationFIleName.Folder
)

type Folder struct {
	Name string `json:"name"`
}

// ReplayList 在查询的时候顺路更新一下数据库的表名
func ReplayList(c *gin.Context) {
	var fileList []string
	var folderList []string
	var folderName Folder
	userID := cosFile.GetID(c).ID

	client := cosFile.Client()
	//初始化客户端

	err := c.ShouldBind(&folderName)
	if err != nil {
		return
	}

	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:  "users/" + strconv.Itoa(userID) + "/" + folderName.Name, // prefix 表示要查询的文件夹 		// deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys: 1000,                                                    // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := client.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			c.JSON(403, gin.H{
				"code": 403,
				"msg":  err.Error(),
			})
			return
		}
		// 在循环内部收集文件 Key
		for _, content := range v.Contents {
			fileList = append(fileList, content.Key) // 将文件名存入切片
			fmt.Printf("Object: %v\n", content.Key)
		}

		// common prefix 表示表示被 delimiter 截断的路径, 如 delimter 设置为/, common prefix 则表示所有子目录的路径
		for _, commonPrefix := range v.CommonPrefixes {
			folderList = append(folderList, commonPrefix)
			fmt.Printf("CommonPrefixes: %v\n", commonPrefix)
		}
		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}

	fileList = append(folderList, fileList...)
	err = MysqlFile(fileList, c)
	if err != nil {
		c.JSON(403, gin.H{
			"code": 403,
			"msg":  err.Error(),
			"data": "保存列表数据的时候出错",
		})
		return
	}

	sortW.SortPaths(fileList, "dirsFirst")

	c.JSON(200, gin.H{
		"marker": marker,
		"data":   removeDuplicates(fileList),
		"code":   200,
		"msg":    "success",
	})
	//向前端返回相关的列表
}

// MysqlFile 把文件名进行本地化储存,先进行创建
func MysqlFile(list []string, c *gin.Context) error {

	userID := cosFile.GetID(c).ID

	//构建传进数据库的字段
	folder = locationFIleName.Folder{
		UserId:    userID,
		FileName:  list,
		FileType:  "?",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	has, err := app.Engine.Table("folder").Where("user_id = ?", userID).Exist(&folder.UserId)
	if err != nil {
		return err
	}

	if !has {
		fmt.Println("执行了插入操作")
		//在这里要用一个有去重功能的方法，区utils包装一个
		_, err = app.Engine.Insert(folder)
		if err != nil {
			c.JSON(403, gin.H{
				"msg": err.Error(),
			})
			return err
		}
	} else {
		fmt.Printf("执行了更新操作")
		has, err := app.Engine.Table("folder").Where("user_id=?", userID).Cols("file_name").Update(&folder)
		if err != nil {
			return err
		}
		fmt.Println(has)
	}

	return nil

}

func removeDuplicates(paths []string) []string {
	if len(paths) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(paths))
	index := 0

	for _, path := range paths {
		if _, exists := seen[path]; !exists {
			seen[path] = struct{}{}
			paths[index] = path // 复用原切片内存
			index++
		}
	}

	return paths[:index]
}
