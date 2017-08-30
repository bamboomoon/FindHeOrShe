package netMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

// ProxyIP  代理IP
type ProxyIP struct {
	Proxies []Ip //`json:"proxies"`
	Code    int  //`json:"code"`
}

type Ip struct {
	HTTP string //`json:"http"`
}

//GetProxyIP 获取代理 IP
func GetProxyIP() (*ProxyIP, error) {
	resp, err := http.Get("http://lab.crossincode.com/proxy/get/?num=20")
	if err != nil {
		return nil, err
	}
	defer func() {
		resp.Body.Close()
	}()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var proxy ProxyIP
	err = json.Unmarshal(respBody, &proxy)
	if err != nil {
		return nil, err
	}
	return &proxy, nil
}

func GetOkProxyIP() []string {
	fmt.Println("正在获取代理IP:")
	ip, err := GetProxyIP()
	if err != nil || ip.Code != 1 {
		fmt.Println("获取代理IP error:", err)
		os.Exit(-1)
	}
	var okIP []string
	var sn sync.Mutex
	var wg sync.WaitGroup
	for _, v := range ip.Proxies {

		wg.Add(1)
		go func(http string) {
			defer wg.Done()

			c := proxyClient(http, 10)
			resp, err := c.Get("http://music.163.com")
			if err != nil {
				return
			}
			if resp.StatusCode == 200 && resp.ContentLength != 0 {
				fmt.Println(http)
				sn.Lock()
				okIP = append(okIP, http)
				sn.Unlock()
			}
		}(v.HTTP)
	}
	wg.Wait()
	if len(okIP) == 0 {
		fmt.Println("无可用代理IP")
		os.Exit(-1)
	}
	return okIP
}
