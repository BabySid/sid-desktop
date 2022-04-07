package common

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/sahilm/fuzzy"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	// BuiltInHttpRequestHeader zero-value means calculated when request is sent
	builtInHttpRequestHeader = []*HttpHeader{
		{Key: "Accept", Value: "*/*"},
		{Key: "Accept-Encoding", Value: "gzip, deflate, br"},
		{Key: "Connection", Value: "keep-alive"},
		{Key: "Content-Length", Value: "<calculated when request is sent>"},
		{Key: "User-Agent", Value: "Sid Desktop"},
		{Key: "Content-Type", Value: "application/json"},
	}

	AuthHeader = HttpHeader{
		Key:   "Authorization",
		Value: "<calculated when request is sent>",
	}
	httpHeaderName []string

	HttpMethod = []string{
		"POST",
		"GET",
	}
)

type HttpHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func BuiltInHttpHeaderName() []string {
	if httpHeaderName != nil {
		return httpHeaderName
	}

	httpHeaderName = make([]string, 0)
	for _, header := range builtInHttpRequestHeader {
		httpHeaderName = append(httpHeaderName, header.Key)
	}

	return httpHeaderName
}

func NewHttpHeader() *HttpHeader {
	return &HttpHeader{
		Key:   "",
		Value: "",
	}
}

func NewBuiltInHttpHeader() []interface{} {
	rs := make([]interface{}, 0)
	for _, header := range builtInHttpRequestHeader {
		header := &HttpHeader{
			Key:   header.Key,
			Value: header.Value,
		}
		rs = append(rs, header)
	}

	return rs
}

type HttpRequest struct {
	ID     int64  `json:"-"`
	Method string `json:"method"`
	Url    string `json:"url"`

	ReqHeader   []HttpHeader `json:"req_header"`
	ReqBody     []byte       `json:"req_body"`
	ReqBodyType string       `json:"req_body_type"`

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

func (s *HttpRequestList) Append(d HttpRequest) {
	s.requests = append(s.requests, d)
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
		fmt.Println(req.ID, req.Method, req.Url, req.ReqHeader, req.ReqBody, req.ReqBodyType, req.CreateTime, req.AccessTime)
	}
}

func DoHttpRequest(method string, url string, reqBody string, header map[string]string) (string, http.Header, []byte, error) {
	client := &http.Client{}

	var req *http.Request
	var err error

	switch method {
	case "POST":
		req, err = http.NewRequest(method, url, strings.NewReader(reqBody))
	case "GET":
		req, err = http.NewRequest(method, url, nil)
	default:
		gobase.AssertHere()
	}

	if err != nil {
		return "", nil, nil, err
	}

	for k, v := range header {
		if k != "" && k != "Content-Length" {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return resp.Status, resp.Header, body, err
}
