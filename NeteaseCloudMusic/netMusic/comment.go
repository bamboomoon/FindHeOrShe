package netMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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

var findUseName = "一叶倾城"

//IsContinue 是否继续获取评论
var IsContinue = true

// var sn *sync.Mutex
var IsChangeIP bool

//Wg main同步
var Wg = sync.WaitGroup{}
var total uint64

//GetComments get comment model
func GetComments(ch chan uint32, proxyIP string) {

	page := <-ch //read

	defer func() {
		if err := recover(); err != nil {

			// sn.Lock()
			IsChangeIP = true
			// sn.Unlock()
		}
		Wg.Done()
	}()
	getCommentsCount := uint64(page * 100)
	if total != 0 && getCommentsCount > total && getCommentsCount-total > 100 {
		fmt.Println("uint64(page*100) > total：", page, total)
		IsContinue = false
		ch <- page + 1
		return
	}

	ch <- page + 1

	resp, err := clientRequest(page, "65538", proxyIP)
	catchError(err, 67)
	if resp == nil || resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}

	p, err := ioutil.ReadAll(resp.Body)
	catchError(err, 73)

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
	catchError(err, 88)
	findComment(&comment, page)
}

func findComment(comment *commentMusic, page uint32) {
	total = comment.Total //comment total
	for _, userCommtent := range comment.Comments {
		if userCommtent.User.Nickname == findUseName {
			fmt.Println("找到了---", userCommtent)
		}
	}
}

//error handle
func catchError(err error, line int) {
	if err != nil {
		fmt.Println("error:", err, line)

		panic(err)
	}
}
