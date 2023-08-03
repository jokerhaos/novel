package main

import (
	"fmt"
	"net/http"
	"novel/config"
	"novel/models"
	"novel/utils"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TextObj struct {
	Url    string
	Author string
	Name   string
}

func main() {
	// Down()
	// Down458t("32980")
	// var text string
	// fmt.Printf("请输入小说名字:")
	// fmt.Scanf("%d", &text)
	// 通过浏览器去爬
	// DownBiququ("4234")
	DownGashuw("/biquge_110771/")
}

func Search(textName string) TextObj {

	return TextObj{}
}

func Down458t(book string) {
	uri := "https://www.458t.com"
	data := &models.Book{
		Href: uri + fmt.Sprintf("/book/%s.html", book),
		// Name:   "过度反应",
		// Author: "阿司匹林",
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
	data.Author = dom.Find(".booktag .red").Eq(0).Text()
	data.Name = dom.Find(".booktitle").Text()

	utils.WriteToTxt(data.Name, data.Name, data.Author)
	utils.WriteToTxt("作者："+data.Author+"\r\n", data.Name, data.Author)
	fmt.Println("作者："+data.Author+"\r\n", data.Name, data.Author)
	// fmt.Println(string(body))
	dom.Find("#list-chapterAll dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Find("a").Text()
		url := s.Find("a").AttrOr("href", "")
		defer func() {
			// 判断章节
			// if data.Chapter == title {
			// 	start = true
			// }
		}()
		if title == "<<---展开全部章节--->>" {
			// 循环另一个
			// dom.Find(".dd_hide dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
			// 	return true
			// })
			return true
		}
		// if !start {
		// 	return true
		// }
		// 开始爬虫
		body, resp, err := req.RepeatSend("GET", url, nil)
		if resp.StatusCode != http.StatusServiceUnavailable && err != nil {
			// 哦豁爬虫失败了
			config.Log.Error(fmt.Sprintf("[%s][%s][%s]爬虫异常：", data.Name, data.Author, title))
			return false
		}
		dom2, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		context := dom2.Find(".readcontent").Text()
		// 写入小说
		re := regexp.MustCompile(`(\r\n|\n|<br\s*/?>)`)
		context = re.ReplaceAllString(context, "\r\n\r\n")
		// 写入小说
		utils.WriteToTxt(title+"\r\n\r\n", data.Name, data.Author)
		utils.WriteToTxt(context+"\r\n\r\n", data.Name, data.Author)
		fmt.Println(title, "写入成功")
		return true
	})
}

func Down3zmmm(bookId string) {
	biquge("https://www.3zmmm.net", "/files/article/html/71772/71772361/")
}

func DownGashuw(bookId string) {
	biquge("http://www.gashuw.com", bookId)
}

func biquge(uri string, bookId string) {
	data := &models.Book{
		Href: uri + bookId,
	}
	req := utils.NewSendRequest(http.Header{}, "")
	req.SetBiqugeHeaders()

	body, _, err := req.Get(data.Href)

	if err != nil {
		return
	}

	// goquery 常规用法
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	data.Name = dom.Find("#info h1").Text()
	data.Author = (strings.Split(dom.Find("#info p").Eq(0).Text(), "："))[1]

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

func DownBiququ(bookId string) {
	uri := "https://www.biququ.info"
	data := &models.Book{
		Href: uri + "/html/" + bookId,
	}

	req := utils.NewSendRequest(http.Header{}, "")
	req.SetBiqugeHeaders()
	fmt.Println(data.Href)

	body, _, err := req.Get(data.Href)
	fmt.Println(err)
	if err != nil {
		return
	}

	// goquery 常规用法
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	data.Name = dom.Find("#info h1").Text()
	data.Author = (strings.Split(dom.Find("#info p").Eq(0).Text(), "："))[1]

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
