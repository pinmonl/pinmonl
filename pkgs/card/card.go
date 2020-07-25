package card

import (
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Card scrapes the information of social media card.
type Card struct {
	HtmlTitle       string
	HtmlDescription string

	FacebookTitle       string
	FacebookDescription string
	FacebookImageURL    string
	FacebookSiteName    string

	TwitterTitle       string
	TwitterDescription string
	TwitterImageURL    string

	Response *http.Response
	Document *goquery.Document
}

// NewCard downloads card information from the url.
func NewCard(rawurl string) (*Card, error) {
	res, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return &Card{
		HtmlTitle:       doc.Find("head title").Text(),
		HtmlDescription: doc.Find("head meta[name='description']").AttrOr("content", ""),

		FacebookTitle:       doc.Find("head meta[property='og:title']").AttrOr("content", ""),
		FacebookDescription: doc.Find("head meta[property='og:description']").AttrOr("content", ""),
		FacebookImageURL:    doc.Find("head meta[property='og:image']").AttrOr("content", ""),
		FacebookSiteName:    doc.Find("head meta[property='og:site_name']").AttrOr("content", ""),

		TwitterTitle:       doc.Find("head meta[name='twitter:title']").AttrOr("content", ""),
		TwitterDescription: doc.Find("head meta[name='twitter:description']").AttrOr("content", ""),
		TwitterImageURL:    doc.Find("head meta[name='twitter:image:src']").AttrOr("content", ""),

		Response: res,
		Document: doc,
	}, nil
}

// Title retrieves the suitable card title.
func (c *Card) Title() string {
	if c.FacebookTitle != "" {
		return c.FacebookTitle
	}
	if c.TwitterTitle != "" {
		return c.TwitterTitle
	}
	return c.HtmlTitle
}

// Description retrieves the suitable card description.
func (c *Card) Description() string {
	if c.FacebookDescription != "" {
		return c.FacebookDescription
	}
	if c.TwitterDescription != "" {
		return c.TwitterDescription
	}
	return c.HtmlDescription
}

// ImageURL retrieves the suitable card image url.
func (c *Card) ImageURL() string {
	if c.FacebookImageURL != "" {
		return c.FacebookImageURL
	}
	if c.TwitterImageURL != "" {
		return c.TwitterImageURL
	}
	return ""
}

// Image returns the image content at ImageURL.
func (c *Card) Image() ([]byte, error) {
	url := c.ImageURL()
	if url == "" {
		return nil, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	img, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}
