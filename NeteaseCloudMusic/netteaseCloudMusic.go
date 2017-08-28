package main

import (
	"NeteaseCloudMusic/netMusic"
	"fmt"
	"time"
	// "./netMusic"
)

func main() {
	ip, err := netMusic.GetProxyIP()
	if err != nil {
		fmt.Println("获取代理ip 出错了", err)
		return
	}
	fmt.Println(ip, ip.Proxies[0].HTTP)
	return
	ch := make(chan uint32, 10)

	ch <- uint32(0)
	i := 0
	page := 0
	for netMusic.IsContinue {
		page++
		for i == 20 {
			time.Sleep(time.Second * time.Duration(5))
			i = 0
		}
		i++
		netMusic.Wg.Add(1)
		go netMusic.GetComments(ch)
	}
	netMusic.Wg.Wait()
	fmt.Println("main over")
}
