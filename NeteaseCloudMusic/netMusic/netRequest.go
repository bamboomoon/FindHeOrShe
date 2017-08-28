package netMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type commentMusic struct {
	HotComments []user
	Comments    []user
	Total       uint64
}

type user struct {
	User    userName
	Content string
}
type userName struct {
	Nickname string
}

//RequestParam post request body
type RequestParam struct {
	Offset    uint32 `json:"offset"`
	Rid       string `json:"rid"`
	Limit     int    `json:"limit"`
	CsrfToken string `json:"csrf_token"`
}

var findUseName = "在流浪路上越走越远"
var IsContinue = true
var Wg = sync.WaitGroup{}
var total uint64

func proxyClient() *http.Client {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://61.136.163.246:80")
	}
	//Proxy: proxy,
	t := &http.Transport{Proxy: proxy}
	c := http.Client{Timeout: time.Duration(1) * time.Minute, Transport: t}
	return &c
}

func randomUserAgent() string {
	userAgentList := []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
		"Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 5.1.1; Nexus 6 Build/LYZ28E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Mobile/14F89;GameHelper",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.1.1 Safari/603.2.4",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:46.0) Gecko/20100101 Firefox/46.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:46.0) Gecko/20100101 Firefox/46.0",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",
		"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)",
		"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Win64; x64; Trident/6.0)",
		"Mozilla/5.0 (Windows NT 6.3; Win64, x64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/13.10586",
		"Mozilla/5.0 (iPad; CPU OS 10_0 like Mac OS X) AppleWebKit/602.1.38 (KHTML, like Gecko) Version/10.0 Mobile/14A300 Safari/602.1"}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pos := r.Intn(len(userAgentList))
	return userAgentList[pos]
}
func clientRequest(page uint32, rid string) (*http.Response, error) {
	parsms := RequestParam{Offset: page * 100, Limit: 100, Rid: rid, CsrfToken: ""} //uint64(page) * 100,
	body, err := Encrypt(&parsms)
	catchError(err)

	urlStr := fmt.Sprintf("http://music.163.com/weapi/v1/resource/comments/R_SO_4_%s/?csrf_token=", rid)
	v := url.Values{}
	v.Set("params", body.Params)
	v.Add("encSecKey", body.EncSecKey)

	request, err := http.NewRequest("POST", urlStr, strings.NewReader(v.Encode()))
	userAgent := randomUserAgent()
	header := map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,gl;q=0.6,zh-TW;q=0.4",
		"Connection":      "keep-alive",
		"Content-Type":    "application/x-www-form-urlencoded",
		"Referer":         "http://music.163.com",
		"Host":            "music.163.com",
		"Cookie":          "",
		"User-Agent":      userAgent,
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}

	return proxyClient().Do(request)
}

//GetComments get comment model
func GetComments(ch chan uint32) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}

		// fmt.Println("Done")
		Wg.Done()
	}()
	page := <-ch //read
	if total != 0 && uint64(page*100) > total {
		fmt.Println("uint64(page*100) > total：", page, total)
		IsContinue = false
		ch <- page + 1
		return
	}

	ch <- page + 1
	resp, err := clientRequest(page, "436514312")
	catchError(err)
	if resp == nil {
		return
	}

	p, err := ioutil.ReadAll(resp.Body)
	catchError(err)

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if len(p) == 0 {
		fmt.Println("数据获取失败")
		os.Exit(-1)
	}

	var comment commentMusic
	err = json.Unmarshal(p, &comment)
	if err != nil {
		panic(err)
	}
	findComment(&comment)

}

func findComment(comment *commentMusic) {
	total = comment.Total //comment total
	for _, userCommtent := range comment.Comments {
		// fmt.Println(userCommtent.User.Nickname, userCommtent.Content)
		if userCommtent.User.Nickname == findUseName {
			fmt.Println("找到了", userCommtent.Content)
			IsContinue = false
		}
	}
}

//error handle
func catchError(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
