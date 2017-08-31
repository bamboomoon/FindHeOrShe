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

//FindUseName 查找的用户名
var FindUseName string

//SongID 查找的歌曲ID
var SongID string

var once sync.Once

//IsContinue 是否继续获取评论
var IsContinue = true

//WgRquest 请求获取评论main同步
var WgRquest = sync.WaitGroup{}
var WgDealComment = sync.WaitGroup{}
var total uint64

//GetComments get comment model
func GetComments(ch chan uint32, proxyIP string) {
	page := <-ch //read

	defer func() {
		if err := recover(); err != nil {

		}
		WgRquest.Done()
	}()
	getCommentsCount := uint64(page * 100)
	if total != 0 && getCommentsCount > total && getCommentsCount-total > 100 {
		// fmt.Println("uint64(page*100) > total：", page, total)
		IsContinue = false
		ch <- page + 1
		return
	}

	ch <- page + 1

	comments, err := sendRequest(page, proxyIP)
	if err != nil {
		fmt.Printf("第%d页没有获取到\n", page)
		return
	}
	findComment(comments, page)
}

func sendRequest(page uint32, proxyIP string) (*commentMusic, error) {
	resp, err := clientRequest(page, SongID, proxyIP)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}
	return &comment, nil
}

func findComment(comment *commentMusic, page uint32) {
	once.Do(func() {
		total = comment.Total //comment total
		fmt.Println("总评论数:", total)
	})
	WgDealComment.Add(1)
	go func(comment *commentMusic) {
		defer WgDealComment.Done()
		for _, userCommtent := range comment.Comments {
			if userCommtent.User.Nickname == FindUseName {
				fmt.Println("找到了---", userCommtent.Content)
			}
		}
	}(comment)
}

//error handle
func catchError(err error, line int) {
	if err != nil {
		fmt.Println("error:", err, line)

		panic(err)
	}
}
