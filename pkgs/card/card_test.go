package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	tests := []struct {
		rawurl  string
		fbtitle string
		fbdesc  string
		err     error
	}{
		{
			rawurl:  "https://github.com/ahshum/card",
			fbtitle: "ahshum/card",
			fbdesc:  "card card card card card card card card card card card - ahshum/card",
			err:     nil,
		},
		{
			rawurl:  "https://github.com/ahshum/not-existed",
			fbtitle: "Build software better, together",
			fbdesc:  "GitHub is where people build software. More than 50 million people use GitHub to discover, fork, and contribute to over 100 million projects.",
			err:     nil,
		},
		{
			rawurl:  "https://gitlab.com/ahshum/card",
			fbtitle: "Ah Shum / Card",
			fbdesc:  "card card card card card card card card card card card",
			err:     nil,
		},
		{
			rawurl:  "https://bitbucket.org/ahshum/card",
			fbtitle: "",
			fbdesc:  "",
			err:     nil,
		},
	}

	for _, test := range tests {
		card, err := NewCard(test.rawurl)
		t.Errorf("%v\n", card)
		assert.Equal(t, test.err, err)
		if err == nil {
			assert.Equal(t, test.fbtitle, card.FacebookTitle)
			assert.Equal(t, test.fbdesc, card.FacebookDescription)
		}
	}
}
