package monl

import "github.com/PuerkitoBio/goquery"

type Card interface {
	Title() string
	Description() string
	ImageUrl() string
}

type SimpleCard struct {
	title       string
	description string
	imageUrl    string
}

func NewCard(title, description, imageUrl string) *SimpleCard {
	return &SimpleCard{
		title:       title,
		description: description,
		imageUrl:    imageUrl,
	}
}

func NewCardFromDocument(doc *goquery.Document) *SimpleCard {
	return NewCard("", "", "")
}

func (s *SimpleCard) Title() string { return s.title }

func (s *SimpleCard) Description() string { return s.description }

func (s *SimpleCard) ImageUrl() string { return s.imageUrl }
