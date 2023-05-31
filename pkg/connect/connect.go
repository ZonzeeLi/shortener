package connect

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

// client is a http client with a short timeout.
var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

// Get 判断url是否能请求通
func Get(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorw("connect client.Get failed", logx.Field("err", err.Error()))
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK // 返回重定向也不认为ping通
}
