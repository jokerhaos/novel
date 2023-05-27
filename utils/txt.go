package utils

import (
	"bufio"
	"fmt"
	"os"
)

const (
	msg = "======================================笔趣阁小说爬虫v1.0======================================\n" +
		"免责声明：本爬虫仅供资料学习，请勿滥用，造成的一切后果与作者无关，地址：https://github.com/jokerhaos/novel\n" +
		"测试阶段存在bug，不是所有小说都可以爬取。\n" +
		"==============================================================================================="
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
