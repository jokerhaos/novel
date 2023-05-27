package main

import (
	"log"
	"novel/config"
	"novel/getter"
	"novel/models"
	"novel/utils"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	models.ConnectDB()
	config.InitLog()
	Find("邪神")
	time.Sleep(time.Hour)
}

type pcfun func(string) ([]models.Book, error)

func Find(name string) {
	var wg sync.WaitGroup
	funs := []pcfun{
		// getter.Biqudd,
		getter.Ibiquges,
	}

	for _, f := range funs {
		wg.Add(1)
		go func(f pcfun, name string) {
			defer wg.Done()
			// defer func() {
			// 	if r := recover(); r != nil {
			// 		// 在这里处理panic异常
			// 		fmt.Println("捕获到panic异常:", r)
			// 	}
			// }()
			temp, err := f(name)
			if err != nil {
				return
			}
			// 判断数据库是否存在
			for _, v := range temp {
				v.CreateTime = time.Now()
				// 判断数据库是否存在
				data := models.Book{}
				models.DB.Table("book").Where(map[string]interface{}{
					"name":   v.Name,
					"author": v.Author,
				}).Last(&data)
				// .Where("new_chapter != " + v.NewChapter).Where("chapter != new_chapter") 程序判断，不要影响索引
				if utils.IsStructEmpty(data) {
					models.DB.Table("book").Create(&v)
					data = v
				} else {
					// 判断章节是否需要更新
					if (data.NewChapter == v.NewChapter && data.Chapter == v.NewChapter) || v.Lock == 1 {
						continue
					}
				}
				// 更新小说
				// func BqgCrawl(startUrl, bookname string, sign int)
				models.DB.Table("book").Where("id = ?", data.ID).Update("lock", 1)
				callback := func(id uint, name string, updateData models.Book) {
					// 因为结构体更新是非0属性，又不想用map那就改值叭
					config.Log.Info(name + " 爬取成功")
					updateData.Lock = 2
					models.DB.Table("book").Where("id = ?", id).Updates(updateData)
				}
				go func(data models.Book) {
					switch data.F {
					case "BqgCrawl":
						getter.BqgCrawl(data, callback)
					}
				}(data)
			}
		}(f, name)
	}
	wg.Wait()
	log.Println("All getters finished.")
}
