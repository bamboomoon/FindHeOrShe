package netMusic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

var findedSync sync.Mutex
var searchCommens []string //已经找到的

//GetComments 内部获取页数
func getPageThenBegin(ch chan uint32, proxyIP string) {
	page := <-ch //read

	defer func() {
		if err := recover(); err != nil {

		}
		wgRquest.Done()
	}()
	getCommentsCount := uint64(page * 100)
	if total != 0 && getCommentsCount > total && getCommentsCount-total > 100 {
		// fmt.Println("uint64(page*100) > total：", page, total)
		isContinue = false
		ch <- page + 1
		return
	}
	ch <- page + 1

	getComments(page, proxyIP, false)
}

//外部传入页数
//处理错误页使用
func getComments(page uint32, proxyIP string, isDealErr bool) {
	if isDealErr {
		defer wgDealErros.Done()
	}
	comments, err := sendRequest(page, proxyIP)
	if err != nil {
		fmt.Printf("第%d页没有获取到\n", page)
		rw.Lock()
		errosPages = append(errosPages, page)
		rw.Unlock()
		return
	}
	findComment(comments, page)
}

func sendRequest(page uint32, proxyIP string) (*commentMusic, error) {
	resp, err := commentReq(page, songID, proxyIP)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
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

	var comment commentMusic
	err = json.Unmarshal(p, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func findComment(comment *commentMusic, page uint32) {
	wgDealComment.Add(1)
	go func(comment *commentMusic) {
		defer wgDealComment.Done()
		for _, userCommtent := range comment.Comments {
			if userCommtent.User.Nickname == findUseName {
				findedSync.Lock()
				searchCommens = append(searchCommens, userCommtent.Content)
				findedSync.Unlock()
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
