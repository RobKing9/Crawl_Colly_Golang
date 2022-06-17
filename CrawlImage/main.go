package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/http"
)

func main() {
	/*
		爬取需要登录的网站
	*/
	//新建一个采集器

	// Instantiate default collector
	//c := colly.NewCollector(colly.AllowURLRevisit())

	//// Rotate two socks5 proxies
	//rp, err := proxy.RoundRobinProxySwitcher("socks5://us5.sgateway.link:706", "socks5://hk1.sgateway.link:706")
	//if err != nil {
	//	log.Println("err:", err.Error())
	//}
	//c.SetProxyFunc(rp)
	c := colly.NewCollector()
	url := "https://skyeysnow.com"
	cookie := "rkvl_2132_saltkey=iwJB5Bip"
	//设置随机useragent
	extensions.RandomUserAgent(c)
	//设置登录cookie
	c.SetCookies(url, []*http.Cookie{
		&http.Cookie{
			Name:     "remember_user_token",
			Value:    cookie,
			Path:     "/",
			Domain:   "skyeysnow.com",
			Secure:   true,
			HttpOnly: true,
		},
	})

	//err1 := c.Post("https://skyeysnow.com/login.php", map[string]string{"username": "Lzhjk", "password": "Porthack123"})
	//if err1 != nil {
	//	log.Println("err1:", err1.Error())
	//	return
	//} else {
	//	log.Println("login successfully!")
	//}
	c.OnHTML("*", func(e *colly.HTMLElement) {
		fmt.Println(e)
	})
	/*
		c.OnHTML("div[id='um']", func(e *colly.HTMLElement) {
			content := e.ChildText("p > a[id='ratio']")
			log.Println("welcome!")
			log.Printf("内容:  %v\n", content)
			log.Println()
			//c.Visit(e.Request.AbsoluteURL(link))
		})

	*/

	c.OnRequest(func(r *colly.Request) {
		log.Println("c正在访问:", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("c正在访问:", r.Request.URL, "访问失败:", err.Error())
	})

	c.OnResponse(func(r *colly.Response) {
		filename := "./log/log.txt"
		err := r.Save(filename)
		if err != nil {
			fmt.Println("保存失败！", "原因:", err.Error())
		} else {
			fmt.Println("保存成功！")
		}
		log.Println("response received", r.StatusCode)
	})

	err2 := c.Visit(url)
	if err2 != nil {
		log.Println("err2:", err2.Error())
	}
}

/*
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
		filename := fmt.Sprintf("./pageimagesplus/img%d.jpg", atomic.AddUint32(&count, 1))
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

*/

//err := cl.Limit(&colly.LimitRule{
//	DomainRegexp: `pic.netbian.com`,
//	RandomDelay:  500 * time.Millisecond,
//	Parallelism:  12,
//})
//if err != nil {
//	log.Fatal(err)
//}
//cl.WithTransport( &http.Transport{
//	Proxy: http.ProxyFromEnvironment,
//	DialContext: ( &net.Dialer{
//		Timeout: 30 *time.Second,
//		KeepAlive: 30 *time.Second,
//	}).DialContext,
//	MaxIdleConns: 100,
//	IdleConnTimeout: 90 *time.Second,
//	TLSHandshakeTimeout: 10 *time.Second,
//	ExpectContinueTimeout: 1 *time.Second,
//
//},
//	)

//cl.Wait()
//c2.Wait()
//c3.Wait()

/*
		爬取图片(无需登录）
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"), colly.MaxDepth(1), colly.Debugger(&debug.LogDebugger{}))
	//文章列表
	c.OnHTML("div[class='TypeList'] > ul", func(e *colly.HTMLElement) {
		//列表中每一项
		e.ForEach("li", func(i int, item *colly.HTMLElement) {
			//图片标题
			title := item.ChildText("a > span")
			//图片链接
			href := item.ChildAttr("a", "href")
			//图片jpg
			image := item.ChildAttr("a > img", "src")
			////文章摘要
			//summary := item.ChildText("div[class='content'] > p[class='abstract']")
			//fmt.Printf("标题:  %v\n文章链接:  %v\n", title, href)
			fmt.Printf("图片标题:  %v\n", title)
			fmt.Printf("图片链接:  %v\n", href)
			fmt.Printf("图片jpg:  %v\n", image)
			fmt.Println()
		})
	})

	err := c.Visit("https://umei.cc/meinvtupian/xingganmeinv/")
	if err != nil {
		fmt.Println(err.Error())
	}

*/
//}

//检查错误
//func HandleError (err error, why string) {
//	if err != nil {
//		fmt.Println(err, why)
//	}
//}

//func download(url string, i int) {
//	res, err := http.Get(url)
//	if err != nil {
//		panic(err)
//	}
//	defer res.Body.Close()
//
//	file, _ := os.Create(fmt.Sprintf("/Golang/Crawl/Colly/image/1_%d.jpg", i))
//
//	defer file.Close()
//	io.Copy(file, res.Body)
//
//}
