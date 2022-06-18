package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

/*
	1、提高下载速度
	2、异步爬取
	3、flag解析
*/

/*
	url
	https://www.imn5.net/XiuRen/XiuRen/10614.html
	https://www.imn5.net/XiuRen/XiuRen/10611.html
	https://www.imn5.net/XiuRen/XiuRen/10610.html
	https://www.imn5.net/XiuRen/XiuRen/10527.html
	https://www.imn5.net/XiuRen/XiuRen/10585.html
*/

var a int

func main() {
	URL := os.Args[1]

	os.Mkdir("img", 0666)

	cl := colly.NewCollector(colly.Async(true))
	//下载图片

	//c2 := cl.Clone()
	c3 := cl.Clone()

	//爬取网站
	var title string

	cl.OnHTML("div[class = 'imgwebp']>p", func(e *colly.HTMLElement) {
		//列表中的每一项
		e.ForEach("img", func(i int, el *colly.HTMLElement) {
			//图片标题
			title = el.Attr("alt")
			//图片
			image := el.Attr("src")
			img := "https://www.imn5.net" + image
			//c2.Visit(img)
			c3.Visit(img)
			//将内容输出
			fmt.Printf("图片标题:%v\n", title)
			fmt.Printf("图片jpg:https://www.imn5.net%v\n", image)

		})
	})
	extensions.RandomUserAgent(cl)
	extensions.Referer(cl)
	cl.OnRequest(func(r *colly.Request) {
		fmt.Println("c1正在访问:", r.URL)
	})

	cl.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	cl.OnError(func(r *colly.Response, err error) {
		fmt.Println("c1正在访问:", r.Request.URL, "访问失败:", err.Error())
	})

	//下载图片到本地
	var count uint32
	c3.OnResponse(func(r *colly.Response) {
		filename := fmt.Sprintf("./img/%d%s.jpg", atomic.AddUint32(&count, 1), title)
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
	url := strings.Replace(URL, ".html", "", -1)
	//开始爬取
	for page := 1; page <= 10; page++ {
		cl.Visit(fmt.Sprintf("%s_%d.html", url, page))
	}

	cl.Wait()
	//c2.Wait()
	c3.Wait()
}
