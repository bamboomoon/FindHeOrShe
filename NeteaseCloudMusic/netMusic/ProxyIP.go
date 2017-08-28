package netMusic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type proxyIP struct {
	Proxies []ip //`json:"proxies"`
	Code    int  //`json:"code"`
}

type ip struct {
	HTTP string //`json:"http"`
}

//GetProxyIP 获取代理 IP
func GetProxyIP() (*proxyIP, error) {
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
	var proxy proxyIP
	err = json.Unmarshal(respBody, &proxy)
	if err != nil {
		return nil, err
	}
	return &proxy, nil
}
