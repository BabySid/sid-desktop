package common

import (
	"encoding/base64"
	"github.com/BabySid/gobase"
	"github.com/sahilm/fuzzy"
	"io/ioutil"
	"net/http"
	"strings"
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

	ReqHeader []HttpHeader `json:"req_header"`
	ReqBody   []byte       `json:"req_body"`

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
	return s.requests[i].Method + " " + s.requests[i].Url
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

func (s *HttpRequestList) Get(method, url string) (HttpRequest, bool) {
	for _, req := range s.requests {
		if req.Method == method && req.Url == url {
			return req, true
		}
	}

	return HttpRequest{}, false
}

func (s *HttpRequestList) Append(d HttpRequest) {
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

func EncodeBasicAuth(user, pass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
}

func DecodeBasicAuth(code string) (string, string) {
	if len(code) <= len("Basic ") {
		return "", ""
	}

	auth, err := base64.StdEncoding.DecodeString(code[len("Basic "):])
	if err != nil {
		return "", ""
	}

	arr := strings.Split(string(auth), ":")
	if len(arr) != 2 {
		return "", ""
	}
	return arr[0], arr[1]
}

func DoHttpRequest(method string, url string, reqBody string, headers []HttpHeader) (int, string, http.Header, []byte, error) {
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
		return 0, "", nil, nil, err
	}

	for _, header := range headers {
		if header.Key != "" && header.Key != "Content-Length" {
			req.Header.Add(header.Key, header.Value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", nil, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, resp.Status, resp.Header, body, err
}
