package main

import (
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	govList := home()
	govList = append(govList, "gov-.htm")
	govAll := make([]string, 0)
	for _, gov := range govList {
		govs := getGovList("http://www.5566.net/" + gov)
		govAll = append(govAll, govs...)
	}
	fmt.Println(govAll)
	fmt.Println(len(govAll))
}

func home() []string {
	bodyStr := tryGetBody("http://www.5566.net/gov-.htm", 3)
	compile := regexp.MustCompile(`<a class="p12" href="(.*?)" target="_self">.*?</a>`)
	allString := compile.FindAllStringSubmatch(bodyStr, -1)
	govList := make([]string, len(allString))
	for i, match := range allString {
		govList[i] = match[1]
	}
	return govList
}

//get gov.cn list
func getGovList(url string) []string {
	bodyStr := tryGetBody(url, 3)
	limiterCompile := regexp.MustCompile(`.wscckey=`)
	islimiter := limiterCompile.MatchString(bodyStr)
	if islimiter {
		url = getLimiterUrl(bodyStr)
		return getGovList(url)
	} else {
		compile := regexp.MustCompile(`,.*.gov.cn`)
		allString := compile.FindAllStringSubmatch(bodyStr, -1)
		var govStr string
		for _, match := range allString {
			govStr = match[0]
		}
		if len(govStr) != 0 {

			return strings.Split(govStr, ",")
		}
		return nil
	}
}
func getLimiterUrl(body string) string {
	urlReg := `<HTML><HEAD><script>window.location.href="(.*?)";</script></HEAD><BODY>`
	limiterCompile := regexp.MustCompile(urlReg)
	limiterUrl := limiterCompile.FindAllStringSubmatch(body, -1)
	return limiterUrl[0][1]
}

func tryGetBody(url string, num int) string {
	response := getBody(url)
	if response.StatusCode != 200 && num > 0 {
		for {
			num--
			fmt.Printf("[Error] %d 休眠 30 秒重试 \n", response.StatusCode)
			time.Sleep(time.Duration(30) * time.Second)
			return tryGetBody(url, num)
		}
	} else {
		body := response.Body
		return readBody(body)
	}
}

func getBody(url string) *http.Response {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("fatal error ", err.Error())
		os.Exit(0)
	}
	request.Header.Add("Accept-Language", "")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	request.Header.Add("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	response, err := client.Do(request)
	if err != nil || response == nil {
		fmt.Println("fatal error")
		panic(err)
	}
	return response
}

// 读取body
func readBody(body io.ReadCloser) string {
	byte2, _ := ioutil.ReadAll(body)
	defer body.Close()
	env := mahonia.NewDecoder("GBK")
	out := env.ConvertString(string(byte2))
	return out
}
