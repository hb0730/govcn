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
)

func main() {
	govList := home()
	for _, gov := range govList {
		govs := getGovList("http://www.5566.net/" + gov)
		fmt.Println(govs)
	}
}

func home() []string {
	bodyStr := getBody("http://www.5566.net/gov-.htm")
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
	bodyStr := getBody(url)
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

func getBody(url string) string {
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
	body := response.Body
	return readBody(body)
}

// 读取body
func readBody(body io.ReadCloser) string {
	byte2, _ := ioutil.ReadAll(body)
	defer body.Close()
	env := mahonia.NewDecoder("GBK")
	out := env.ConvertString(string(byte2))
	return out
}
