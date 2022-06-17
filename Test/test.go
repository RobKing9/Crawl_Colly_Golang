package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	reQQEmail = "(\\d+)@qq.com"
)

func CrawlEmail(url string) {
	//网站爬取数据
	resp, err := http.Get(url)
	HandleError(err, "http.Get url") //处理异常
	defer resp.Body.Close()          //关闭程序
	//读取页面内容
	PageByte, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll")
	//将字节转换为字符串
	PageString := string(PageByte)
	log.Println(PageString)
	//fmt.Println(PageString)
	//筛选出QQ邮箱
	//re := regexp.MustCompile(reQQEmail)
	//-1代表全部邮箱
	//results := re.FindAllStringSubmatch(PageString, -1)
	//fmt.Println(results)
	//遍历结果，得到自己想要的格式
	//for _, result := range results {
	//	fmt.Println("Email = ", result[0])
	//}
}

func HandleError(err error, why string) {
	if err != nil {
		fmt.Println(why, err)
	}
}

func main() {
	CrawlEmail("https://skyeysnow.com//")
}
