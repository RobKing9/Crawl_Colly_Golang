package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"sync/atomic"
)

func main() {
	urlCollector := colly.NewCollector()
	//访问这些图片链接
	imageLink := urlCollector.Clone()
	//下载图片
	downloadImage := urlCollector.Clone()

	extensions.RandomUserAgent(urlCollector)
	extensions.Referer(urlCollector)
	rp, err := proxy.RoundRobinProxySwitcher("http://192.168.123.220:7890", "socks://127.0.0.1:1080")
	if err != nil {
		fmt.Println("proxy failed:", err.Error())
	}
	urlCollector.SetProxyFunc(rp)

	//爬取网站

	urlCollector.OnHTML("ul[class = 'clearfix']", func(e *colly.HTMLElement) {
		//列表中的每一项
		e.ForEach("li", func(i int, el *colly.HTMLElement) {
			//图片标题
			title := el.ChildText("a > b")
			//图片链接
			href := el.ChildAttr("a", "href")
			//访问这些图片链接
			imageLink.Visit(e.Request.AbsoluteURL(href))
			//图片
			image := el.ChildAttr("a > img", "src")
			image1 := "https://pic.netbian.com" + image
			downloadImage.Visit(image1)
			//将内容输出
			log.Printf("图片标题:  %v\n", title)
			log.Printf("图片jpg:  https://pic.netbian.com%v\n", image)
			log.Println()

		})
	})

	urlCollector.OnRequest(func(r *colly.Request) {
		log.Println("c1正在访问:", r.URL)
	})

	urlCollector.OnError(func(r *colly.Response, err error) {
		log.Println("c1正在访问:", r.Request.URL, "访问失败:", err.Error())
	})

	//下载图片到本地
	var count uint32
	downloadImage.OnResponse(func(r *colly.Response) {
		filename := fmt.Sprintf("./img/img%d.jpg", atomic.AddUint32(&count, 1))
		err := r.Save(filename)
		if err != nil {
			log.Println("保存失败！", "原因:", err.Error())
		} else {
			log.Println("保存成功！")
		}

	})
	downloadImage.OnRequest(func(r *colly.Request) {
		log.Println("downloadImage正在访问：", r.URL)
	})

	//开始爬取

	for page := 2; page <= 10; page++ {
		urlCollector.Visit(fmt.Sprintf("https://pic.netbian.com/4kmeinv/index_%d.html", page))
	}

}
