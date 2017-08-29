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

	//page
	ch := make(chan uint32, 40)
	ch <- uint32(0)

	//线程数量
	count := 0
	httPIP := okIP[0]
	fmt.Println("okip", httPIP)
	for netMusic.IsContinue {

		if count == 40 {
			netMusic.Wg.Wait()
			count = 0
		}
		count++
		netMusic.Wg.Add(1)
		go netMusic.GetComments(ch, httPIP)
	}
	fmt.Println("循环结束了")
	netMusic.Wg.Wait()
	fmt.Println("查找完毕！！！")

}
