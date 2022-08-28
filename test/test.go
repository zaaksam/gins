package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/zaaksam/gins/constant"
)

func Get(t *testing.T, req *Request, urlStr string) {
	do(t, req, "GET", urlStr)
}

func Post(t *testing.T, req *Request, urlStr string) {
	do(t, req, "POST", urlStr)
}

func Do(t *testing.T, req *Request, method, urlStr string) {
	do(t, req, method, urlStr)
}

func do(t *testing.T, req *Request, method, urlStr string) {
	onceInit()

	method = strings.ToUpper(method)
	var (
		httpReq *http.Request
		body    []byte
		err     error
	)

	if req.query != nil {
		if strings.Contains(urlStr, "?") {
			urlStr += "&"
		} else {
			urlStr += "?"
		}
		urlStr += req.query.Encode()

	}

	// POST 请求时，json 为空时必须传 {}
	if method == "POST" && req.json == nil {
		req.json = make(map[string]any)
	}

	if req.json != nil {
		body, err = json.Marshal(req.json)
		if err != nil {
			t.Fatalf("序列化 json 数据时出错：%s", err)
		}

		httpReq, err = http.NewRequest(method, urlStr, bytes.NewReader(body))
	} else {
		httpReq, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		t.Fatalf("创建请求对象时出错：%s", err)
	}

	if method == "POST" {
		httpReq.Header.Set("Content-Type", "application/json")
	} else {
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	for k, v := range req.header {
		httpReq.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	instance.Engine().ServeHTTP(w, httpReq)

	if err != nil {
		t.Fatalf("处理请求时出错：%s", err)
	}

	httpRes := w.Result()

	if httpRes.StatusCode != http.StatusOK {
		t.Fatalf("请求响应错误，http code：%d", httpRes.StatusCode)
	}

	res, err := NewResponse(w.Body.Bytes())
	if err != nil {
		t.Fatalf("响应数据解析时出错：%s", err)
	}

	if res.Code != constant.API_OK {
		t.Fatalf("响应 code 不为 200，响应内容：\n\n%s\n\n", res.content)
	}

	t.Log(res.content)
}
