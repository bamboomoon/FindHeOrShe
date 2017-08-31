package netMusic

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

//FindUseName 查找的用户名
var findUseName string

//SongID 查找的歌曲ID
var songID string

var once sync.Once

//IsContinue 是否继续获取评论
var isContinue = true

//WgRquest 请求获取评论main同步
var wgRquest = sync.WaitGroup{}
var wgDealComment = sync.WaitGroup{}
var total uint64

//错误页 重新获取
var sn sync.Mutex
var errosPages []uint32
var wgDealErros = sync.WaitGroup{}

//可用的代理IP
var okIPs []string

func stdinUserNameAndID() bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println(len(os.Args))
	if len(os.Args) != 5 {
		panic(errors.New("命令有误"))
	}

	if os.Args[1] == "-name" {
		findUseName = os.Args[2]
	} else {
		secondParam := os.Args[2]
		_, err := strconv.ParseUint(secondParam, 10, 64)
		if err != nil {
			panic(errors.New(fmt.Sprint("提供的歌曲ID格式错误!!!"+"--", err)))
		}
		songID = secondParam
	}

	if os.Args[3] == "-ID" {
		secondParam := os.Args[4]
		_, err := strconv.ParseUint(secondParam, 10, 64)
		if err != nil {
			panic(errors.New(fmt.Sprint("提供的歌曲ID格式错误!!!"+"--", err)))
		}
		songID = secondParam
	} else {
		findUseName = os.Args[4]
	}
	if songID == "" || findUseName == "" {
		panic(errors.New("歌曲ID获取或者用户名不能为空"))
	}
	fmt.Println("您输入的查找的用户名为:", findUseName, "查找的歌曲ID:", songID)
	fmt.Println("请确认(y/n)")
	var correct string
	fmt.Scanln(&correct)
	if correct == "y" || correct == "Y" {
		return true
	}
	return false
}

//Begin 程序入口
func Begin() {
	isBegin := stdinUserNameAndID()
	if isBegin == false {
		os.Exit(-1)
	}

	//代理IP
	okIPs = getOkProxyIP()
	ipCount := len(okIPs)
	ipIndex := 0
	var httpIP string
	//page
	ch := make(chan uint32, 40)
	ch <- uint32(0)

	//goroutine数量
	count := 0
	allCount := ipCount * 10

	fmt.Printf("开始查找『%s』在「%s」下的评论:\n", findUseName, songID)
	fmt.Println("总评论数:", total)
	for isContinue {
		if ipIndex > ipCount-1 {
			ipIndex = 0
		}
		httpIP = okIPs[ipIndex]
		if count == allCount { //防止发送请求过于频繁被封
			wgRquest.Wait()
			count = 0
		}
		count++
		ipIndex++
		wgRquest.Add(1)
		go getPageThenBegin(ch, httpIP)
	}
	wgRquest.Wait()
	wgDealComment.Wait()

	//处理错误页
	timer := time.NewTimer(time.Duration(5) * time.Minute)
	go func() {
		<-timer.C
		fmt.Println("处理错误超时")
		os.Exit(-1)
	}()

	for len(errosPages) > 0 {
		dealErrorPage()
	}
	fmt.Println("查找完毕！！！")

}

//处理没有获出错了的页面的评论
func dealErrorPage() {
	ipCount := len(okIPs)
	ipIndex := 0
	var newErrorPages []uint32
	for i, v := range errosPages {
		if ipIndex >= ipCount {
			ipIndex = 0
		}
		wgDealErros.Add(1)
		fmt.Printf("正在重新获取第%d页\n", v)
		go getComments(v, okIPs[ipIndex], true)
		ipIndex++
		newErrorPages = append(errosPages[:i], errosPages[i+1:]...) //移除使用了的
	}
	errosPages = newErrorPages
	wgDealErros.Wait()
}
