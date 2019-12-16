package database

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseMigration(t *testing.T) {
	testCases := []struct {
		name  string
		file  string
		wants Migration
	}{
		{
			name: "standard one",
			file: "./testdata/source/standard/0001_one.sql",
			wants: Migration{
				Up:   []string{"up one;\n"},
				Down: []string{"down one;\n"},
			},
		},
		{
			name: "standard two",
			file: "./testdata/source/standard/0002_two.sql",
			wants: Migration{
				Up:   []string{"up two;\n"},
				Down: []string{"down two;\n"},
			},
		},
		{
			name: "discardline one",
			file: "./testdata/source/discardline/0001_one.sql",
			wants: Migration{
				Up:   []string{"up one;\n"},
				Down: []string{"down one;\n"},
			},
		},
		{
			name: "dicardline two",
			file: "./testdata/source/discardline/0002_two.sql",
			wants: Migration{
				Up:   []string{"up two;\n"},
				Down: []string{"down two;\n"},
			},
		},
	}

	for _, tc := range testCases {
		content, _ := ioutil.ReadFile(tc.file)
		got, _ := parseMigration("", bytes.NewBuffer(content))
		if wants := tc.wants; !reflect.DeepEqual(got, wants) {
			t.Errorf("case %q fails", tc.name)
		}
	}
}
