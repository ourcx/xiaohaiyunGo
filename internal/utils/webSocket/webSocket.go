package webSocketLeg

import (
	legislations "xiaohaiyun/internal/utils/legislation"
)

var trieV1 *legislations.TrieV1

func init() {
	//初始化关键词
	trieV1 = legislations.NewTrieV1()
	for _, word := range legislations.SensitiveWords {
		trieV1.Insert(word)
	}
}

// Legislation 处理违规关键词的
func Legislation(text string) bool {
	//添加关键词
	return trieV1.Contains(text)
}
