package base

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func GetWebPageTitle(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", genWebPageError(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", genWebPageError(fmt.Errorf("status code is not 200: %d", resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", genWebPageError(err)
	}

	title := doc.Find("title").Text()
	return title, nil
}

func genWebPageError(err error) error {
	return fmt.Errorf("process webpage failed. %s", err)
}
