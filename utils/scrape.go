package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Scrape gets httpResponse
func Scrape(url string) (*http.Response, error) {
	res, err := http.Get(url)
	return res, err
}

// ScrapeWithDocument gets httpResponse and GoQuery for deeper handling
func ScrapeWithDocument(url string) (*goquery.Document, *http.Response, error) {
	res, err := Scrape(url)
	if err != nil {
		return nil, res, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	return doc, res, err
}
