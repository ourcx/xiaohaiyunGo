package sortW

import (
	"regexp"
	"sort"
	"strings"
)

func SortPaths(paths []string, dirsOrder string) {
	dirsFirst := dirsOrder == "dirsFirst"
	wordRegex := regexp.MustCompile(`[-_]`)

	sort.Slice(paths, func(i, j int) bool {
		p1, p2 := paths[i], paths[j]
		isDir1 := strings.HasSuffix(p1, "/")
		isDir2 := strings.HasSuffix(p2, "/")

		// 处理目录和文件的相对顺序
		if isDir1 != isDir2 {
			if dirsFirst {
				return isDir1
			}
			return isDir2
		}

		// 同类比较（目录 vs 目录 或 文件 vs 文件）
		if isDir1 {
			// 目录按完整路径升序
			return p1 < p2
		}

		// 文件比较逻辑
		getExt := func(p string) (string, int) {
			if dot := strings.LastIndex(p, "."); dot != -1 && dot != len(p)-1 {
				ext := p[dot+1:]
				return ext, len(wordRegex.Split(ext, -1))
			}
			return "", 0
		}

		ext1, cnt1 := getExt(p1)
		ext2, cnt2 := getExt(p2)

		// 先比较扩展名单词数量（降序）
		if cnt1 != cnt2 {
			return cnt1 > cnt2
		}

		// 再比较扩展名字母顺序（升序）
		if ext1 != ext2 {
			return ext1 < ext2
		}

		// 最后按完整路径升序
		return p1 < p2
	})
}
