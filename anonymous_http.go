package libtools

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/proxy"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

var socks5Client *http.Client

func init() {
	// Create a socks5 dialer
	proxyConf, _ := config.String("anonymous_proxy")
	dialer, err := proxy.SOCKS5("tcp", proxyConf, nil, proxy.Direct)
	if err != nil {
		logs.Error("[SimpleAnonymousHttpClient] proxy dialer err: %v, proxyConf: %s", err, proxyConf)
	}

	// Setup HTTP transport
	tr := &http.Transport{
		Dial: dialer.Dial,
	}
	socks5Client = &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}
}

// 简单的http客户端,支持POST表单域,但不支持上传文件
func SimpleAnonymousHttpClient(reqMethod string, reqUrl string, reqHeaders map[string]string, reqBody string) ([]byte, int, error) {
	var httpStatusCode int
	var emptyBody []byte

	req, err := http.NewRequest(reqMethod, reqUrl, strings.NewReader(reqBody))
	if err != nil {
		logs.Error("[SimpleAnonymousHttpClient] http.NewRequest fail, reqUrl:", reqUrl)
		return emptyBody, httpStatusCode, err
	}

	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	resp, err := socks5Client.Do(req)
	if err != nil {
		logs.Error("[SimpleAnonymousHttpClient] do request fail, reqUrl:", reqUrl, ", err:", err)
		return emptyBody, httpStatusCode, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("[SimpleAnonymousHttpClient] read request fail, reqUrl:", reqUrl, ", err:", err)
		return emptyBody, httpStatusCode, err
	}

	return body, resp.StatusCode, err
}
