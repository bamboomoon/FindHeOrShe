package netMusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ProxyIP  代理IP
type ProxyIP struct {
	Proxies []ip //`json:"proxies"`
	Code    int  //`json:"code"`
}

type ip struct {
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
	ip, err := GetProxyIP()
	if err != nil || ip.Code != 1 {
		fmt.Println("获取代理IP error:", err)
		os.Exit(-1)
	}
	var okIP []string
	for _, v := range ip.Proxies {
		c := proxyClient(v.HTTP, 10)
		fmt.Println(v.HTTP)
		resp, err := c.Get("http://music.163.com")
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 {
			okIP = append(okIP, v.HTTP)
		}
	}
	return okIP
}
