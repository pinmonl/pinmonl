package scrape

import "testing"

func TestScrapeGet(t *testing.T) {
	res, _ := Get("https://github.com/vuejs/vue")
	doc, _ := res.Doc()

	h := doc.Find("head")

	t.Error("-- facebook")
	t.Error(h.Find("meta[property='og:title']").Attr("content"))
	t.Error(h.Find("meta[property='og:description']").Attr("content"))
	t.Error(h.Find("meta[property='og:image']").Attr("content"))

	t.Error("-- twitter")
	t.Error(h.Find("meta[name='twitter:title']").Attr("content"))
	t.Error(h.Find("meta[name='twitter:description']").Attr("content"))
	t.Error(h.Find("meta[name='twitter:image:src']").Attr("content"))

	card, _ := res.Card()
	t.Error(card.Image())
}
