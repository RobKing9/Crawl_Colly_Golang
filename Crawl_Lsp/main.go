package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

/*
	1、解决代理问题
	2、生成日志
*/
//var URL = "https://asiantolick.com/post-489/%E6%82%A0%E5%AE%9D-haruka-very-cute-ktv-uniform-girl"
//https://asiantolick.com/post-1240/%E7%A6%8F%E5%88%A9%E5%A7%AC-%E8%90%8C%E7%99%BD%E9%86%AC-%E7%B4%85%E8%89%B2%E5%A4%A7%E8%9D%B4%E8%9D%B6%E7%B5%90-66p
func main() {

	URL := os.Args[1]

	os.Mkdir("img", 0666)

	writer, err := os.OpenFile("crawl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file failed:", err.Error())
	}

	cl := colly.NewCollector(colly.Async(true), colly.Debugger(&debug.LogDebugger{Output: writer}), colly.MaxDepth(2))
	c2 := cl.Clone()

	var title string
	cl.OnHTML("article>div[class='spotlight-group']", func(e *colly.HTMLElement) {
		//列表中的每一项
		e.ForEach("div", func(i int, el *colly.HTMLElement) {
			//图片标题
			title = el.ChildAttr("img", "alt")
			//图片
			img := el.ChildAttr("img", "src")
			//img := "https://www.imn5.net" + image
			//c2.Visit(img)
			c2.Visit(img)
			//将内容输出
			fmt.Printf("图片标题:%v\n", title)
			fmt.Printf("图片jpg:%v\n", img)

		})
	})

	extensions.RandomUserAgent(cl)
	extensions.Referer(cl)
	//配置代理
	rp, err := proxy.RoundRobinProxySwitcher("http://192.168.123.220:7890", "http://127.0.0.1:7890")
	if err != nil {
		fmt.Println("proxy failed:", err.Error())
	}
	cl.SetProxyFunc(rp)

	cl.OnRequest(func(r *colly.Request) {
		fmt.Println("c1正在访问:", r.URL)
	})

	cl.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   300 * time.Second,
			KeepAlive: 300 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	cl.OnError(func(r *colly.Response, err error) {
		fmt.Println("c1正在访问:", r.Request.URL, "访问失败:", err.Error())
	})
	c2.OnError(func(r *colly.Response, err error) {
		fmt.Println("c2正在访问:", r.Request.URL, "访问失败:", err.Error())
	})

	//下载图片到本地
	var count uint32
	c2.OnResponse(func(r *colly.Response) {
		filename := fmt.Sprintf("./img/img%d_%s.jpg", atomic.AddUint32(&count, 1), title)
		err := r.Save(filename)
		if err != nil {
			fmt.Println("保存失败！", "原因:", err.Error())
		} else {
			fmt.Println("保存成功！")
		}
	})

	cl.Visit(URL)

	cl.Wait()
	c2.Wait()

}
