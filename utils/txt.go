package utils

import (
	"bufio"
	"fmt"
	"os"
)

const (
	filePath = "./text/" //项目根目录前缀
)

func init() {
	Mkdir(filePath)
}

func Mkdir(filePath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		dir := filePath
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleDir(filePath string) error {
	err := os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return err
}

// 将字符串添加写入文本文件
func WriteToTxt(content, bookname, dirname string) {
	Mkdir(filePath + dirname)
	filepath := filePath + dirname + "/" + bookname + ".txt" // 存放小说的TXT文件路径
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(content)
	write.Flush()
	// fmt.Println("写入", bookname, "完成")
}
