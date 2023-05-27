package getter

import (
	"fmt"
	"net/http"
	"novel/config"
	"novel/models"
	"novel/utils"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Ibiquges(name string) (books []models.Book, err error) {
	// https://www.ibiquges.com/
	pollURL := "https://www.ibiquges.com/modules/article/waps.php"
	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Pragma", "no-cache")
	headers.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Set("Sec-Ch-Ua-Mobile", "?0")
	headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	headers.Set("Sec-Fetch-Dest", "document")
	headers.Set("Sec-Fetch-Mode", "navigate")
	headers.Set("Sec-Fetch-Site", "same-origin")
	headers.Set("Sec-Fetch-User", "?1")
	headers.Set("Upgrade-Insecure-Requests", "1")
	headers.Set("Referer", "https://www.ibiquges.com/modules/article/waps.php")
	req := utils.NewSendRequest(headers, "")

	// fmt.Println(name)
	param := map[string]string{
		"searchkey": name,
	}
	body, _, err := req.Post(pollURL, utils.MapToUrlValue(param))
	if err != nil {
		return nil, err
	}

	// goquery 常规用法
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	dom.Find("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		_title := s.Find("a").Eq(0).Text()
		_newChapter := s.Find("td").Eq(1).Find("a").Eq(0).Text()
		_author := s.Find("td").Eq(2).Text()
		layout := "2006-01-02 15:04:05"
		dateStr := s.Find("td").Eq(3).Text()
		year := time.Now().Year() // 获取当前年份
		dateTimeStr := fmt.Sprintf("%d-%s 00:00:00", year, dateStr)
		_updateTime, _ := time.Parse(layout, dateTimeStr)
		if err != nil {
			_updateTime = time.Now()
		}
		url, exists := s.Find("a").Eq(0).Attr("href")
		if exists {
			// config.Log.Info("书名：", _title, " 作者：", _author, " url:"+url)
			book := models.Book{
				Author:     _author,
				Name:       _title,
				Href:       url,
				NewChapter: _newChapter,
				UpdateTime: _updateTime,
				F:          "BqgCrawl",
			}
			books = append(books, book)
		}
		return true
	})

	return books, nil
}

func BqgCrawl(data models.Book, callback func(uint, string, models.Book)) {
	updateData := &models.Book{}
	defer callback(data.ID, data.Name, *updateData)
	headers := http.Header{}
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Pragma", "no-cache")
	headers.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	headers.Set("Sec-Ch-Ua-Mobile", "?0")
	headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	headers.Set("Sec-Fetch-Dest", "document")
	headers.Set("Sec-Fetch-Mode", "navigate")
	headers.Set("Sec-Fetch-Site", "same-origin")
	headers.Set("Sec-Fetch-User", "?1")
	headers.Set("Upgrade-Insecure-Requests", "1")
	headers.Set("Referer", "https://www.ibiquges.com/modules/article/waps.php")
	req := utils.NewSendRequest(headers, "")
	body, _, err := req.Get(data.Href)

	if err != nil {
		return
	}

	// goquery 常规用法
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	updateTime := dom.Find("#info p").Eq(2).Text()
	updateTime = (strings.Split(updateTime, "："))[1]
	layout := "2006-01-02 15:04:05"
	_updateTime, err := time.Parse(layout, updateTime)
	if err != nil {
		updateData.UpdateTime = _updateTime
	}
	updateData.Image = dom.Find("#fmimg img").AttrOr("src", "")
	updateData.Describe = dom.Find("#intro p").Eq(1).Text()
	start := false
	if data.Chapter == "" {
		utils.WriteToTxt(data.Name, data.Name, data.Author)
		utils.WriteToTxt("作者："+data.Author+"\r\n", data.Name, data.Author)
		start = true
	}

	dom.Find("#list dd").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Find("a").Text()
		url := "https://www.ibiquges.com" + s.Find("a").AttrOr("href", "")
		defer func() {
			// 判断章节
			if data.Chapter == title {
				start = true
				updateData.Chapter = title
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

		return true
	})
}
