package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"sync/atomic"
)

func main() {
	cl := colly.NewCollector()
	//访问这些图片链接
	c2 := cl.Clone()
	//下载图片
	c3 := cl.Clone()

	//爬取网站

	cl.OnHTML("ul[class = 'clearfix']", func(e *colly.HTMLElement) {
		//列表中的每一项
		e.ForEach("li", func(i int, el *colly.HTMLElement) {
			//图片标题
			title := el.ChildText("a > b")
			//图片链接
			href := el.ChildAttr("a", "href")
			//访问这些图片链接
			c2.Visit(e.Request.AbsoluteURL(href))
			//图片
			image := el.ChildAttr("a > img", "src")
			image1 := "https://pic.netbian.com" + image
			c3.Visit(image1)
			//将内容输出
			fmt.Printf("图片标题:  %v\n", title)
			//fmt.Printf("图片链接:  %v\n", href)
			fmt.Printf("图片jpg:  https://pic.netbian.com%v\n", image)
			fmt.Println()

		})
	})

	cl.OnRequest(func(r *colly.Request) {
		fmt.Println("c1正在访问:", r.URL)
	})

	cl.OnError(func(r *colly.Response, err error) {
		fmt.Println("c1正在访问:", r.Request.URL, "访问失败:", err.Error())
	})

	//下载图片到本地
	var count uint32
	c3.OnResponse(func(r *colly.Response) {
		filename := fmt.Sprintf("./img/img%d.jpg", atomic.AddUint32(&count, 1))
		err := r.Save(filename)
		if err != nil {
			fmt.Println("保存失败！", "原因:", err.Error())
		} else {
			fmt.Println("保存成功！")
		}

	})
	c3.OnRequest(func(r *colly.Request) {
		fmt.Println("c3正在访问：", r.URL)
	})

	//开始爬取

	for page := 2; page <= 10; page++ {
		cl.Visit(fmt.Sprintf("https://pic.netbian.com/4kmeinv/index_%d.html", page))
	}
}
