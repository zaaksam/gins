package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
)

// instance http请求客户端
var instance client

type client struct{}

// Get 发起 http get 请求
func Get(urlStr string, query url.Values, headerOpt ...http.Header) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	return instance.Get(urlStr, query, headerOpt...)
}

// Get 发起 http get 请求
func (c *client) Get(urlStr string, query url.Values, headerOpt ...http.Header) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	opt := &Option{
		URL:         urlStr,
		Method:      "GET",
		ContentType: "application/x-www-form-urlencoded",
		Query:       query,
	}

	if len(headerOpt) == 1 {
		opt.Header = headerOpt[0]
	}

	return c.Request(opt)
}

// PostForm 发起 http post 表单请求
func PostForm(urlStr string, queryBody url.Values, headerOpt ...http.Header) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	return instance.PostForm(urlStr, queryBody, headerOpt...)
}

// PostForm 发起 http post 表单请求
func (c *client) PostForm(urlStr string, queryBody url.Values, headerOpt ...http.Header) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	opt := &Option{
		URL:         urlStr,
		ContentType: "application/x-www-form-urlencoded",
		Body:        []byte(queryBody.Encode()),
	}

	if len(headerOpt) == 1 {
		opt.Header = headerOpt[0]
	}

	return c.Request(opt)
}

// PostJSON 发起 http post json 请求
func PostJSON(urlStr string, jsonBody []byte, headerOpt ...http.Header) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	return instance.PostJSON(urlStr, jsonBody, headerOpt...)
}

// PostJSON 发起 http post json 请求
func (c *client) PostJSON(urlStr string, jsonBody []byte, headerOpt ...http.Header) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	opt := &Option{
		URL:  urlStr,
		Body: jsonBody,
	}

	if len(headerOpt) == 1 {
		opt.Header = headerOpt[0]
	}

	return c.Request(opt)
}

// Request 发起 http 请求
func Request(opt *Option) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	return instance.Request(opt)
}

// Request 发起 http 请求
func (*client) Request(opt *Option) (resStatus int, resHeader http.Header, resBody []byte, err error) {
	client := http.Client{
		Transport: cleanhttp.DefaultTransport(),
	}
	opt.Method = strings.ToUpper(opt.Method)

	if opt.Method == "" {
		opt.Method = "POST"
	}

	if opt.ContentType == "" {
		opt.ContentType = "application/json"
	}

	if opt.Timeout.Seconds() > 0 {
		client.Timeout = opt.Timeout
	}

	urlStr := opt.URL
	if len(opt.Query) > 0 {
		if strings.Contains(urlStr, "?") {
			urlStr += "&"
		} else {
			urlStr += "?"
		}
		urlStr += opt.Query.Encode()
	}

	var (
		res *http.Response
		req *http.Request
	)

	if len(opt.Body) == 0 {
		req, err = http.NewRequest(opt.Method, urlStr, nil)
	} else {
		req, err = http.NewRequest(opt.Method, opt.URL, bytes.NewReader(opt.Body))
	}
	if err != nil {
		return
	}

	if len(opt.Header) > 0 {
		req.Header = opt.Header
	}
	req.Header.Set("Content-Type", opt.ContentType)

	res, err = client.Do(req)
	if err != nil {
		return
	}

	resStatus = res.StatusCode
	resHeader = res.Header

	//读取响应结果
	resBody, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	return
}
