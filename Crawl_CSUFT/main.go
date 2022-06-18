package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"log"
	"os"
)

var (
	loginurl = "http://authserver.csuft.edu.cn/authserver/login?service=http://ehall.csuft.edu.cn/login"
	indexurl = "http://ehall.csuft.edu.cn/new/index.html"
)

func main() {

	writer, err := os.OpenFile("crawl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	HandleErr("file", err)
	loginpage := colly.NewCollector(colly.Debugger(&debug.LogDebugger{Output: writer}), colly.MaxDepth(2))
	errlogin := loginpage.Post(loginurl, map[string]string{"username": "20202524", "password": "LKD2020.0921."})
	HandleErr("login", errlogin)

	indexpage := loginpage.Clone()
	indexpage.OnHTML("div[class='amp-widget amp-widget-only-drag PC-CARD-HTML-4786696181714491-01 undefined undefined ampLoadingCompleteFlag']", func(e *colly.HTMLElement) {
		log.Println("welcome!")
		//id := e.ChildText("div[id='Top1_divLoginName']")
		//fmt.Println(id)
	})

	//设置随机useragent
	extensions.RandomUserAgent(loginpage)
	extensions.Referer(loginpage)

	//loginpage.SetCookies(loginurl, []*http.Cookie{
	//	{
	//		Name:     "SERVERID",
	//		Value:    "122",
	//		Path:     "/",
	//		Domain:   "jwgl.csuft.edu.cn",
	//		Secure:   false,
	//		HttpOnly: false,
	//	},
	//})

	// 在提出请求之前打印 "访问…"
	loginpage.OnRequest(func(r *colly.Request) {
		log.Println("访问登录页面:", r.URL)
	})

	loginpage.OnError(func(r *colly.Response, err error) {
		log.Println("正在访问登录页面:", r.Request.URL, "访问失败:", err.Error())
	})

	indexpage.OnRequest(func(r *colly.Request) {
		log.Println("访问首页:", r.URL)
	})

	indexpage.OnError(func(r *colly.Response, e error) {
		log.Println("访问首页:", r.Request.URL, "访问失败：", err.Error())
	})

	err = loginpage.Visit(loginurl)
	HandleErr("visit login", err)

	indexpage.Visit(indexurl)

}

func HandleErr(str string, err error) {
	if err != nil {
		fmt.Printf("%s failed:%v", str, err.Error())
	}
}
