package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"./netMusic"
)

func main() {
	begin()
}
func stdinUserNameAndID() {
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
		netMusic.FindUseName = os.Args[2]
	} else {
		secondParam := os.Args[2]
		_, err := strconv.ParseUint(secondParam, 10, 64)
		if err != nil {
			panic(errors.New(fmt.Sprint("提供的歌曲ID格式错误!!!"+"--", err)))
		}
		netMusic.SongID = secondParam
	}

	if os.Args[3] == "-ID" {
		secondParam := os.Args[4]
		_, err := strconv.ParseUint(secondParam, 10, 64)
		if err != nil {
			panic(errors.New(fmt.Sprint("提供的歌曲ID格式错误!!!"+"--", err)))
		}
		netMusic.SongID = secondParam
	} else {
		netMusic.FindUseName = os.Args[4]
	}
	if netMusic.SongID == "" || netMusic.FindUseName == "" {
		panic(errors.New("歌曲ID获取或者用户名不能为空"))
	}
	fmt.Println("你的输入为-", "查找用户名:", netMusic.FindUseName, "查找歌曲ID:", netMusic.SongID)
	fmt.Println("请确认(y/n)")
	var correct string
	fmt.Scanln(&correct)
	if correct == "y" || correct == "Y" {
		begin()
	} else {
		os.Exit(-1)
	}
}
func begin() {
	stdinUserNameAndID()

	//代理IP
	okIP := netMusic.GetOkProxyIP()
	ipCount := len(okIP)
	ipIndex := 0
	var httpIP string
	//page
	ch := make(chan uint32, 40)
	ch <- uint32(0)

	//goroutine数量
	count := 0
	allCount := ipCount * 10
	fmt.Printf("开始查找『%s』在「%s」下的评论:\n", netMusic.FindUseName, netMusic.SongID)
	for netMusic.IsContinue {
		if ipIndex > ipCount-1 {
			ipIndex = 0
		}
		httpIP = okIP[ipIndex]
		if count == allCount { //防止发送请求过于频繁被封
			netMusic.WgRquest.Wait()
			count = 0
		}
		count++
		ipIndex++
		netMusic.WgRquest.Add(1)
		go netMusic.GetComments(ch, httpIP)
	}
	netMusic.WgRquest.Wait()
	netMusic.WgDealComment.Wait()
	fmt.Println("查找完毕！！！")

}
