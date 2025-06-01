package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	cosFile "xiaohaiyun/internal/utils/cos"
)

var List []string
var mutex sync.Mutex

// TreeFile 这是一个可以拿到所有文件列表的东西
func TreeFile(c *gin.Context) {
	var fileList []string
	userID := cosFile.GetID(c).ID
	var folderName cosFile.Folder
	fileList = TreeFileTree(userID, folderName, fileList)

	c.JSON(200, gin.H{
		"code": 200,
		"data": DeduplicateUnordered(fileList),
	})
	fileList = []string{}
}

func TreeFileTree(userID int, folderName cosFile.Folder, fileList []string) []string {
	fileListT := cosFile.FList(userID, folderName)
	if len(fileListT) == 0 {
		List = append([]string{}, fileList...)
		return fileList
	}

	fileList = append(fileList, fileListT...)
	var folders cosFile.Folder

	for _, file := range fileListT {
		dir, end := filepath.Split(file)
		if dir != "users"+"/"+strconv.Itoa(userID)+"/"+folderName.Name {
			fmt.Println(dir)
			if end == "" {
				e, _ := trimBasePath(file, "users"+"/"+strconv.Itoa(userID)+"/")
				folders.Name = e + "/"
				fmt.Println(folders.Name)
				TreeFileTree(userID, folders, fileList)
			}
		}
	}
	List = append([]string{}, fileList...)
	//这里出现了一个拷贝上的深浅问题
	return fileList
}

func trimBasePath(rawPath, base string) (string, error) {
	// 统一路径格式（兼容不同操作系统）
	cleanPath := filepath.ToSlash(filepath.Clean(rawPath))
	cleanBase := filepath.ToSlash(filepath.Clean(base))

	// 确保基础路径以 / 结尾
	if !strings.HasSuffix(cleanBase, "/") {
		cleanBase += "/"
	}

	// 检查路径是否以基础路径开头
	if !strings.HasPrefix(cleanPath, cleanBase) {
		return "", fmt.Errorf("路径 '%s' 不以 '%s' 开头", rawPath, base)
	}

	// 截取剩余部分
	result := strings.TrimPrefix(cleanPath, cleanBase)

	// 处理空结果情况（当路径完全匹配基础路径时）
	if result == "" {
		return "", nil // 返回当前目录标识
	}

	return result, nil
}

func DeduplicateUnordered(slice []string) []string {
	seen := make(map[string]struct{}, len(slice))
	for _, item := range slice {
		seen[item] = struct{}{} // 空结构体节省内存
	}
	result := make([]string, 0, len(seen))
	dirList := make(map[string]struct{}, len(seen))
	for key := range seen {
		dir, end := filepath.Split(key)
		if end != "" {
			dirList[dir] = struct{}{}
		}
	}

	for key := range seen {
		dir, end := filepath.Split(key)
		if end != "" {
			result = append(result, key)
			continue
		}
		if _, exists := dirList[dir]; !exists {
			result = append(result, key)
			continue
		}
	}
	return result
}
