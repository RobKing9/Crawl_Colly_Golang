package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	url = "https://www.juzikong.com"
)

type JuZi struct {
	TopId   string
	Content string
	Link    string
}

func main() {

	fmt.Println("接下来开始爬虫工作......")
	time.Sleep(time.Second * 1)
	var juZis []JuZi

	cl := colly.NewCollector(colly.MaxDepth(2))

	//设置请求头部

	//爬取网站
	cl.OnHTML("div[style=\"display:;\"]", func(e *colly.HTMLElement) {
		topId := 0
		//列表中的每一项
		e.ForEach("article[class=\"container_3mvaj\"]", func(i int, el *colly.HTMLElement) {
			topId += 1
			//句子链接
			href := el.ChildAttr("div[class=\"body_2l9IL\"] > p[class=\"content_1Nfe9\"] > a", "href")
			href = "了解更多点击链接:" + url + href + "\n" + "\n"
			//句子内容
			juziContent := el.ChildText("div[class=\"body_2l9IL\"] > p[class=\"content_1Nfe9\"] > a > span > span")
			juziContent = "Top" + strconv.Itoa(topId) + "句子内容：\n" + juziContent + "\n"
			topIdStr := "句子控热榜Top" + strconv.Itoa(topId) + "\n"
			juZi := JuZi{
				TopId:   topIdStr,
				Content: juziContent,
				Link:    href,
			}
			juZis = append(juZis, juZi)
			//将内容输出
			log.Printf("%v", href)
			log.Printf("%v", juziContent)
			log.Println()
		})
	})

	//设置随机useragent
	extensions.RandomUserAgent(cl)
	extensions.Referer(cl)
	rp, err := proxy.RoundRobinProxySwitcher("http://127.0.0.1:7890", "https://127.0.0.1:7890")
	if err != nil {
		fmt.Println("proxy failed:", err.Error())
	}
	cl.SetProxyFunc(rp)

	// 在提出请求之前打印 "访问…"
	cl.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		fmt.Println("正在爬取:", r.URL)
	})

	flag := true
	cl.OnError(func(r *colly.Response, err error) {
		flag = false
		fmt.Println("正在爬取:", r.Request.URL, "访问失败:", err.Error())
	})

	err = cl.Visit(url)
	HandleErr("visit", err)

	//创建一个新文件，写入内容
	filePath := "./juzi.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	HandleErr("open file failed", err)
	//及时关闭
	defer file.Close()

	writer := bufio.NewWriter(file)
	//当前时间
	if flag {
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		head := fmt.Sprintf("句子控实时热榜（时间：%v）\n\n", timeNow)
		writer.WriteString(head)
	}
	for _, juzi := range juZis {
		writer.WriteString(juzi.TopId)
		writer.WriteString(juzi.Content)
		writer.WriteString(juzi.Link)
	}
	if flag {
		fgx := "----------------------------------------------我是分割线-------------------------------------------------" + "\n" + "\n" + "\n"
		writer.WriteString(fgx)
	}
	writer.Flush()

	for {

	}
}

func HandleErr(str string, err error) {
	if err != nil {
		fmt.Printf("%s failed:%v\n", str, err.Error())
	}
}
