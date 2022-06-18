package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"net/http"
	"os"
	"time"
)

var (
	url = "https://skyeysnow.com//"
)

func main() {

	fmt.Println("罗子豪是一个大傻逼......")
	fmt.Println("接下来开始爬虫工作......")
	time.Sleep(time.Second * 1)

	writer, err := os.OpenFile("crawl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	HandleErr("file", err)

	cl := colly.NewCollector(colly.Debugger(&debug.LogDebugger{Output: writer}), colly.MaxDepth(2))

	cl.OnHTML("div[id='um']", func(e *colly.HTMLElement) {
		content := e.ChildText("p>a[id='ratio']")
		fmt.Printf("%v\n", content)
	})

	//设置随机useragent
	extensions.RandomUserAgent(cl)
	extensions.Referer(cl)
	//设置cookie
	cl.SetCookies(url, []*http.Cookie{
		{
			Name:     "rkvl_2132_auth",
			Value:    "5848v%2BsBF1fQi6WfN%2FWPVqWlwGkKr6hZ8Q1pXSzyrqfW2d1m%2Bls86cSu3tJv8LTMQJcd5dhfqHQs%2F33pQICGFiZXGA",
			Path:     "/",
			Domain:   "skyeysnow.com",
			Secure:   true,
			HttpOnly: true,
		},
		{
			Name:     "rkvl_2132_saltkey",
			Value:    "e0NQvZzC",
			Path:     "/",
			Domain:   "skyeysnow.com",
			Secure:   true,
			HttpOnly: true,
		},
	})

	// 在提出请求之前打印 "访问…"
	cl.OnRequest(func(r *colly.Request) {
		fmt.Println("蜘蛛正在爬取:", r.URL)
	})

	cl.OnError(func(r *colly.Response, err error) {
		fmt.Println("蜘蛛正在爬取:", r.Request.URL, "访问失败:", err.Error())
	})

	err = cl.Visit(url)
	HandleErr("visit", err)

	time.Sleep(time.Second * 3)
	fmt.Println("罗子豪大傻逼你已经成功获得你的积分了......")
	fmt.Println("Ctrl+C关闭窗口......")

	for {
	}
}

func HandleErr(str string, err error) {
	if err != nil {
		fmt.Printf("%s failed:%v", str, err.Error())
	}
}
