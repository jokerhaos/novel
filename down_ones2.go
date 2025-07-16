package main

import (
	"bufio"
	"fmt"
	"net/http"
	"novel/config"
	"novel/models"
	"novel/utils"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

type TextObj struct {
	Url    string
	Author string
	Name   string
	Loeva  string // 最新章节
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 读取小说名
	fmt.Print("请输入小说名字: ")
	textLine, _ := reader.ReadString('\n')
	text := strings.TrimSpace(textLine)

	list, err := Search(text)
	if err != nil {
		panic(err)
	}
	PrintTable(list)

	// 读取序号
	var sel int
	for {
		fmt.Print("请根据输入序号: ")
		selLine, _ := reader.ReadString('\n')
		selLine = strings.TrimSpace(selLine)
		sel, err = strconv.Atoi(selLine)
		if err != nil || sel < 0 || sel >= len(list) {
			fmt.Println("序号不对，请重新输入。")
			continue
		}
		break
	}

	DownTxt(&list[sel])
}

func Search(textName string) ([]models.Book, error) {
	list := make([]models.Book, 0)
	uri := "https://www.banxia.la/modules/article/search_t.php"
	param := map[string]string{
		"searchkey": textName,
		"Submit":    "",
	}
	req := utils.NewSendRequest(http.Header{}, "")
	SetHeaders(req)
	body, resp, err := req.Post(uri, utils.MapToUrlValue(param))
	if err != nil {
		return nil, err
	}
	// fmt.Println(resp)
	// fmt.Println(string(body))
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	// ✅ 判断是否存在 .pop-books2 元素
	if dom.Find(".pop-books2").Length() == 0 {
		textObj := models.Book{}
		textObj.Name = dom.Find(".book-describe>h1").Text()
		textObj.Author = dom.Find(".book-describe>p").Eq(0).Find("a").Text()
		textObj.Loeva = dom.Find(".book-describe>p").Eq(4).Find("a").Text()
		textObj.Href = resp.Request.URL.String()
		list = append(list, textObj)
		return list, nil
	}
	dom.Find(".pop-books2 .pop-book2").EachWithBreak(func(i int, s *goquery.Selection) bool {
		textObj := models.Book{}
		textObj.Name = s.Find("a").Eq(1).Text()
		textObj.Author = s.Find(".pop-intro").Text()
		textObj.Loeva = ""
		textObj.Href = s.Find("a").Eq(1).AttrOr("href", "")
		list = append(list, textObj)
		return true
	})
	return list, nil
}

func DownTxt(data *models.Book) {
	req := utils.NewSendRequest(http.Header{}, "")
	SetHeaders(req)
	body, _, err := req.Get(data.Href)
	if err != nil {
		return
	}
	// goquery 常规用法
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	data.Author = dom.Find(".book-describe>p").Eq(0).Find("a").Text()
	data.Name = dom.Find(".book-describe>h1").Text()

	utils.WriteToTxt(data.Name, data.Name, data.Author)
	utils.WriteToTxt("作者："+data.Author+"\r\n", data.Name, data.Author)
	fmt.Println("作者："+data.Author+"\r\n", data.Name, data.Author)
	// fmt.Println(string(body))
	dom.Find(".book-list li").EachWithBreak(func(i int, s *goquery.Selection) bool {
		title := s.Find("a").Text()
		url := "https://www.banxia.la" + s.Find("a").AttrOr("href", "")
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
		context := dom2.Find("#nr1").Text()
		// 写入小说
		re := regexp.MustCompile(`(\r\n|\n|<br\s*/?>)`)
		context = re.ReplaceAllString(context, "\r\n\r\n")
		// 写入小说
		// utils.WriteToTxt(title+"\r\n\r\n", data.Name, data.Author)
		// utils.WriteToTxt(context+"\r\n\r\n", data.Name, data.Author)
		fmt.Println(title, "写入成功")
		return true
	})
}

func SetHeaders(req *utils.SendRequest) {
	// 设置请求头
	headers := map[string]string{
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"accept-language":           "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"cache-control":             "max-age=0",
		"content-type":              "application/x-www-form-urlencoded",
		"cookie":                    "HstCfa4825374=1719545648108; HstCmu4825374=1719545648108; HstCnv4825374=2; HstCns4825374=2; jieqiVisitTime=jieqiArticlesearchTime%3D1719555839; HstCla4825374=1719556037965; HstPn4825374=3; HstPt4825374=8",
		"origin":                    "https://www.458t.com",
		"priority":                  "u=0, i",
		"referer":                   "https://www.458t.com/modules/article/search.php",
		"sec-ch-ua":                 `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        `"Windows"`,
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	}
	// 设置请求头
	for key, value := range headers {
		req.Headers.Set(key, value)
	}
}

// 计算字符串的显示宽度（考虑多字节字符）
func GetDisplayWidth(str string) int {
	width := 0
	for _, r := range str {
		if utf8.RuneLen(r) > 1 {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

// 打印表格的函数
func PrintTable(list []models.Book) {
	const (
		idWidth      = 4
		nameWidth    = 40
		authorWidth  = 20
		chapterWidth = 30
		herfWidth    = 30
	)

	// 打印标题行
	fmt.Printf("%-*s %-*s %-*s %-*s %-*s\n", idWidth, "序号", nameWidth, "名字", authorWidth, "作者", chapterWidth, "最新章节", herfWidth, "下载地址")

	// 打印分割线
	fmt.Printf("%s %s %s %s %s\n",
		strings.Repeat("-", idWidth),
		strings.Repeat("-", nameWidth),
		strings.Repeat("-", authorWidth),
		strings.Repeat("-", chapterWidth),
		strings.Repeat("-", herfWidth),
	)

	// 打印数据行
	for k, v := range list {
		fmt.Printf("%-*d %-*s %-*s %-*s %-*s\n",
			idWidth, k,
			nameWidth+GetDisplayWidth(v.Name)-len(v.Name), v.Name,
			authorWidth+GetDisplayWidth(v.Author)-len(v.Author), v.Author,
			chapterWidth+GetDisplayWidth(v.Loeva)-len(v.Loeva), v.Loeva,
			chapterWidth+GetDisplayWidth(v.Href)-len(v.Href), v.Href,
		)
	}
}
