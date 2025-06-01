package shareUrlUtils

import (
	"bytes"
	"fmt"
	"xiaohaiyun/internal/app"
	"xiaohaiyun/internal/models/share"
)

//这是一些关于校验前端url的函数

// GetUrlShareByID 接受前端传过来的数值进行对比
// 根据二进制 ID 查询记录
func GetUrlShareByID(id []byte) (*share.UrlShare, error) {
	// 创建空结构体并设置查询条件
	urlShare := &share.UrlShare{ID: id}

	// 执行查询
	has, err := app.Engine.Get(urlShare)
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败: %w", err)
	}
	if !has {
		return nil, fmt.Errorf("未找到 ID 为 %x 的记录", id)
	}
	return urlShare, nil
}

// CompareData 完整比较函数
func CompareData(frontendData *share.UrlShare) (bool, error) {
	// 1. 获取数据库记录
	dbData, err := GetUrlShareByID(frontendData.ID)

	if err != nil {
		return false, err
	}

	// 2. 深度比较结构体
	return DeepCompare(dbData, frontendData), nil
}

// DeepCompare 深度比较结构体字段
func DeepCompare(a, b *share.UrlShare) bool {
	return bytes.Equal(a.ID, b.ID) &&
		a.Url == b.Url &&
		a.Username == b.Username &&
		a.Signature == b.Signature &&
		a.Email == b.Email &&
		a.Avatar == b.Avatar
}
