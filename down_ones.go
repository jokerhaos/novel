package main

import (
	"fmt"
	"net/http"
	"novel/config"
	"novel/models"
	"novel/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	Down()
}

func Down() {
	uri := "https://www.3zmmm.net"
	data := &models.Book{
		Href:   uri + "/files/article/html/71772/71772361/",
		Name:   "盛宠嫡女：医妃不好惹江初月萧景行",
		Author: "墨墨水田",
	}
	req := utils.NewSendRequest(http.Header{}, "")

	req.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Headers.Set("Cache-Control", "no-cache")
	req.Headers.Set("Pragma", "no-cache")
	req.Headers.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Headers.Set("Sec-Fetch-Dest", "document")
	req.Headers.Set("Sec-Fetch-Mode", "navigate")
	req.Headers.Set("Sec-Fetch-Site", "none")
	req.Headers.Set("Sec-Fetch-User", "?1")
	req.Headers.Set("Upgrade-Insecure-Requests", "1")

	body, _, err := req.Get(data.Href)

	if err != nil {
		return
	}

	// goquery 常规用法
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	start := false
	if data.Chapter == "" {
		utils.WriteToTxt(data.Name, data.Name, data.Author)
		utils.WriteToTxt("作者："+data.Author+"\r\n", data.Name, data.Author)
		start = true
	}
	// fmt.Println(string(body))
	dom.Find("dt").Eq(1).NextAllFiltered("dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Find("a").Text()
		url := uri + s.Find("a").AttrOr("href", "")
		defer func() {
			// 判断章节
			if data.Chapter == title {
				start = true
			}
		}()

		if !start {
			return true
		}
		// 开始爬虫
		body, _, err := req.RepeatSend("GET", url, nil)
		if err != nil {
			// 哦豁爬虫失败了
			config.Log.Error(fmt.Sprintf("[%s][%s][%s]爬虫异常：", data.Name, data.Author, title))
			return false
		}
		dom2, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		context := dom2.Find("#content").Text()
		// 写入小说
		utils.WriteToTxt(title+"\r\n", data.Name, data.Author)
		utils.WriteToTxt(context+"\r\n", data.Name, data.Author)
		fmt.Println(title, "写入成功")
		return true
	})
}
