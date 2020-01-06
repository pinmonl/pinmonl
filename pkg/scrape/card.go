package scrape

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Card defines the data for sharing card.
type Card struct {
	title       string
	description string
	image       []byte
	imageURL    string
}

// NewCardFromDoc creates Card from *goquery.Document.
func NewCardFromDoc(doc *goquery.Document) (*Card, error) {
	h := doc.Find("head")
	if h == nil {
		return nil, fmt.Errorf("scrape: head is empty")
	}

	ct, _ := h.Find("meta[property='og:title']").Attr("content")
	if ct == "" {
		ct, _ = h.Find("meta[name='twitter:title']").Attr("content")
	}

	cd, _ := h.Find("meta[property='og:description']").Attr("content")
	if cd == "" {
		cd, _ = h.Find("meta[name='twitter:description']").Attr("content")
	}

	ci, _ := h.Find("meta[property='og:image']").Attr("content")
	if ci == "" {
		ci, _ = h.Find("meta[name='twitter:image:src']").Attr("content")
	}

	return &Card{
		title:       ct,
		description: cd,
		imageURL:    ci,
	}, nil
}

// Title returns the card title.
func (c *Card) Title() string {
	return c.title
}

// Description returns the card description.
func (c *Card) Description() string {
	return c.description
}

// Image downloads image from the image url.
func (c *Card) Image() ([]byte, error) {
	if c.image == nil {
		res, err := http.Get(c.imageURL)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		c.image = bs
	}
	return c.image, nil
}
