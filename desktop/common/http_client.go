package common

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/sahilm/fuzzy"
	"strings"
	"time"
)

var (
	// BuiltInHttpRequestHeader zero-value means calculated when request is sent
	builtInHttpRequestHeader = map[string]interface{}{
		"Accept":          "*/*",
		"Accept-Encoding": "gzip, deflate, br",
		"Connection":      "keep-alive",
		"Content-Length":  0,
		"User-Agent":      "Sid Desktop",
		"Content-Type":    "application/json",
	}
	defMethod = "POST"

	HttpHeaderName = []string{
		"Accept",
		"Accept-Encoding",
		"Connection",
		"Content-Length",
		"User-Agent",
		"Content-Type",
	}

	HttpMethod = []string{
		"POST",
		"GET",
	}
)

type HttpHeader struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type HttpRequest struct {
	ID     int64  `json:"-"`
	Method string `json:"method"`
	Url    string `json:"url"`

	ReqHeader []HttpHeader `json:"req_header"`

	CreateTime int64 `json:"create_time"`
	AccessTime int64 `json:"access_time"`
}

func InitHttpRequest(req *HttpRequest) {
	gobase.True(req != nil)

	req.Method = defMethod
	req.ReqHeader = make([]HttpHeader, 0)

	for k, v := range builtInHttpRequestHeader {
		header := HttpHeader{
			Key:   k,
			Value: v,
		}
		req.ReqHeader = append(req.ReqHeader, header)
	}

	req.CreateTime = time.Now().Unix()
	req.AccessTime = time.Now().Unix()
}

func (s *HttpRequest) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.ReqHeader), len(s.ReqHeader))
	for i := range s.ReqHeader {
		rs[i] = s.ReqHeader[i]
	}
	return rs
}

type HttpRequestList struct {
	requests []HttpRequest
}

func NewHttpRequestList() *HttpRequestList {
	return &HttpRequestList{
		requests: make([]HttpRequest, 0),
	}
}

func (s *HttpRequestList) Find(name string) *HttpRequestList {
	matches := fuzzy.FindFrom(name, s)

	rs := NewHttpRequestList()
	for _, match := range matches {
		rs.requests = append(rs.requests, s.requests[match.Index])
	}

	return rs
}

func (s *HttpRequestList) String(i int) string {
	return s.requests[i].Url
}

func (s *HttpRequestList) Len() int {
	return len(s.requests)
}

func (s *HttpRequestList) Set(d []HttpRequest) {
	if d == nil {
		return
	}
	s.requests = d
}

func (s *HttpRequestList) Upsert(d HttpRequest) {
	d.Method = strings.ToUpper(d.Method)

	for i, req := range s.requests {
		if req.Method == d.Method && req.Url == d.Url {
			d.ID = req.ID
			d.AccessTime = time.Now().Unix()
			req = d
			s.requests[i] = req
			return
		}
	}
	s.requests = append(s.requests, d)
}

func (s *HttpRequestList) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.requests), len(s.requests))
	for i := range s.requests {
		rs[i] = s.requests[i]
	}
	return rs
}

func (s *HttpRequestList) GetHttpRequest() []HttpRequest {
	return s.requests
}

func (s *HttpRequestList) Debug() {
	for _, req := range s.requests {
		fmt.Println(req.ID, req.Method, req.Url, req.ReqHeader, req.CreateTime, req.AccessTime)
	}
}
