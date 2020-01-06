package scrape

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Response defines the payload of scraper result.
type Response struct {
	R *http.Response

	doc *goquery.Document
	err error
}

// Get downloads the response from url.
func Get(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return &Response{R: resp}, nil
}

// Doc creates *goquery.Document for content scraping.
func (r *Response) Doc() (*goquery.Document, error) {
	if r.doc == nil {
		doc, err := goquery.NewDocumentFromReader(r.R.Body)
		r.R.Body.Close()
		r.doc, r.err = doc, err
	}
	return r.doc, r.err
}

// Card returns the sharing card data.
func (r *Response) Card() (*Card, error) {
	doc, err := r.Doc()
	if err != nil {
		return nil, err
	}
	return NewCardFromDoc(doc)
}
