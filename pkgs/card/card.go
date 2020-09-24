package card

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// Card scrapes the information of social media card.
type Card struct {
	HtmlTitle       string
	MetaTitle       string
	MetaDescription string

	FacebookTitle       string
	FacebookDescription string
	FacebookImageURL    string
	FacebookSiteName    string
	FacebookURL         string

	TwitterTitle       string
	TwitterDescription string
	TwitterImageURL    string
	TwitterURL         string

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
		HtmlTitle:       doc.Find("title").Text(),
		MetaTitle:       doc.Find("meta[name='title']").AttrOr("content", ""),
		MetaDescription: doc.Find("meta[name='description']").AttrOr("content", ""),

		FacebookTitle:       doc.Find("meta[property='og:title']").AttrOr("content", ""),
		FacebookDescription: doc.Find("meta[property='og:description']").AttrOr("content", ""),
		FacebookImageURL:    doc.Find("meta[property='og:image']").AttrOr("content", ""),
		FacebookSiteName:    doc.Find("meta[property='og:site_name']").AttrOr("content", ""),
		FacebookURL:         doc.Find("meta[property='og:url']").AttrOr("content", ""),

		TwitterTitle:       doc.Find("meta[name='twitter:title']").AttrOr("content", ""),
		TwitterDescription: doc.Find("meta[name='twitter:description']").AttrOr("content", ""),
		TwitterImageURL:    doc.Find("meta[name='twitter:image:src']").AttrOr("content", ""),
		TwitterURL:         doc.Find("meta[name='twitter:url']").AttrOr("content", ""),

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
	if c.MetaTitle != "" {
		return c.MetaTitle
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
	return c.MetaDescription
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
	imgurl := c.ImageURL()
	if imgurl == "" {
		return nil, nil
	}

	u, err := url.Parse(imgurl)
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		u.Host = c.Response.Request.URL.Host
	}
	if u.Scheme == "" {
		u.Scheme = c.Response.Request.URL.Scheme
	}

	res, err := http.Get(u.String())
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

func (c *Card) URL() string {
	if c.FacebookURL != "" {
		return c.FacebookURL
	}
	if c.TwitterURL != "" {
		return c.TwitterURL
	}
	return ""
}
