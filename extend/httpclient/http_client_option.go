package httpclient

import (
	"net/http"
	"net/url"
	"time"
)

// Option http请求客户端参数
type Option struct {
	URL         string
	Method      string // 默认 POST
	ContentType string // 默认 application/json
	Header      http.Header
	Query       url.Values
	Body        []byte
	Timeout     time.Duration // 请求超时
}
