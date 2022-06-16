package common

import "net/url"

func ParseURL(urlStr string) *url.URL {
	link, _ := url.Parse(urlStr)

	return link
}
