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

var once1 sync.Once
var once2 sync.Once

//IsContinue 是否继续获取评论
var IsContinue = true

//WgRquest 请求获取评论main同步
var WgRquest = sync.WaitGroup{}
var WgDealComment = sync.WaitGroup{}
var total uint64

//验证用户名和ip是否有提供
func doCheck() {
	if FindUseName == "" || SongID == "" {
		fmt.Println("\n未提供需要查找的用户名或者歌曲ID")
		os.Exit(-1)
	}
}

//GetComments get comment model
func GetComments(ch chan uint32, proxyIP string) {
	once1.Do(doCheck)
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

	resp, err := clientRequest(page, SongID, proxyIP)
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
	once2.Do(func() {
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
