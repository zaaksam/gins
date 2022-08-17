package test

import (
	"net/url"
)

type Request struct {
	header map[string]string
	query  url.Values
	json   map[string]any
}

func NewRequest() *Request {
	return &Request{}
}

func (req *Request) SetHeader(key, val string) {
	if req.header == nil {
		req.header = make(map[string]string)
	}

	req.header[key] = val
}

func (req *Request) AddQuery(key, val string) {
	if req.query == nil {
		req.query = make(url.Values)
	}

	req.query.Add(key, val)
}

func (req *Request) AddJSON(key string, val any) {
	if req.json == nil {
		req.json = make(map[string]any)
	}

	req.json[key] = val
}
