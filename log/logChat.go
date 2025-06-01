package log

import (
	"bufio"
	"fmt"
	"os"
)

// 监听群聊里出现的违规词
func open(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	return file, err
}

// Writer 写入日志
func Writer(log string, c string) error {
	file, err := open(log)
	//错误在于没找到路径
	writer := bufio.NewWriter(file)

	_, err = writer.Write([]byte(c))
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return err
}

// Reader 读取日志
func Reader(log string) error {
	file, err := open(log)
	reader := bufio.NewReader(file)
	read, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Print(read)
	//打印日志
	return err
}
