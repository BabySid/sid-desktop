package common

import (
	"fmt"
	"github.com/sahilm/fuzzy"
	"strings"
	"time"
)

var (
	// BuiltInHttpRequestHeader zero-value means calculated when request is sent
	//builtInHttpRequestHeader = map[string]interface{}{
	//	"Accept":          "*/*",
	//	"Accept-Encoding": "gzip, deflate, br",
	//	"Connection":      "keep-alive",
	//	"Content-Length":  0,
	//	"User-Agent":      "Sid Desktop",
	//	"Content-Type":    "application/json",
	//}
	builtInHttpRequestHeader = map[string]string{
		"0": "0",
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
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
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewHttpHeader() *HttpHeader {
	return &HttpHeader{
		Key:   "",
		Value: "",
	}
}

func NewBuiltInHttpHeader() []interface{} {
	rs := make([]interface{}, 0)
	for k, v := range builtInHttpRequestHeader {
		header := &HttpHeader{
			Key:   k,
			Value: v,
		}
		rs = append(rs, header)
	}

	return rs
}

type HttpRequest struct {
	ID     int64  `json:"-"`
	Method string `json:"method"`
	Url    string `json:"url"`

	ReqHeader []HttpHeader `json:"req_header"`

	CreateTime int64 `json:"create_time"`
	AccessTime int64 `json:"access_time"`
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
