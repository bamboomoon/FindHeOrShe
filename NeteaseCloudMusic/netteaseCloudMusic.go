package main

import (
	"fmt"

	"./netMusic"
)

func main() {
	begin()

}

func begin() {
	okIP := netMusic.GetOkProxyIP()
	ipCount := len(okIP)
	ipIndex := 0
	//page
	ch := make(chan uint32, 40)
	ch <- uint32(0)

	//go数量
	count := 0
	allCount := ipCount * 10

	var httPIP string

	for netMusic.IsContinue {
		if ipIndex > ipCount-1 {
			ipIndex = 0
		}
		httPIP = okIP[ipIndex]
		if count == allCount {
			netMusic.Wg.Wait()
			count = 0
		}
		count++
		ipIndex++
		netMusic.Wg.Add(1)
		go netMusic.GetComments(ch, httPIP)
	}
	netMusic.Wg.Wait()
	fmt.Println("查找完毕！！！")

}
