package utils

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Scrape(url string) (*http.Response, error) {
	res, err := http.Get(url)
	return res, err
}

func ScrapeWithDoc(url string) (*goquery.Document, *http.Response, error) {
	res, err := Scrape(url)
	if err != nil {
		return nil, res, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc, res, err
}
